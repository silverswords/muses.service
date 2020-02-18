package eventbus

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/go-redis/redis/v7"
)

// EventBusRegister defines subscription-related bus behavior
type Register interface {
	Register(channel string, fn interface{}) error
	RegisterAsync(channel string, fn interface{}, transactional bool) error
	RegisterOnce(channel string, fn interface{}) error
	RegisterOnceAsync(channel string, fn interface{}) error
	UnRegister(channel string, handler interface{}) error
}

type Publisher interface {
	Publish(channel string, msg ...interface{})
}

type Mux interface {
}

// A Handler responds to an HTTP request.
//
// ServeHTTP should write reply headers and data to the ResponseWriter
// and then return. Returning signals that the request is finished; it
// is not valid to use the ResponseWriter or read from the
// Request.Body after or concurrently with the completion of the
// ServeHTTP call.
//
// Depending on the HTTP client software, HTTP protocol version, and
// any intermediaries between the client and the Go server, it may not
// be possible to read from the Request.Body after writing to the
// ResponseWriter. Cautious handlers should read the Request.Body
// first, and then reply.
//
// Except for reading the body, handlers should not modify the
// provided Request.
//
// If ServeHTTP panics, the server (the caller of ServeHTTP) assumes
// that the effect of the panic was isolated to the active request.
// It recovers the panic, logs a stack trace to the server error log,
// and either closes the network connection or sends an HTTP/2
// RST_STREAM, depending on the HTTP protocol. To abort a handler so
// the client sees an interrupted response but the server doesn't log
// an error, panic with the value ErrAbortHandler.
type Handler interface {
	ServeHTTP(http.ResponseWriter, *Request)
}

// http.ResponseWriter
// Header() Header map[string][]string
// Write([]byte) (int, error)
// WriteHeader(statusCode int)

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(http.ResponseWriter, *Request)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *Request) {
	f(w, r)
}

type Request struct {
	*redis.Message
}

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux { return new(ServeMux) }

type ServeMux struct {
	mu    sync.RWMutex
	m     map[string]muxEntry
	es    []muxEntry // slice of entries sorted from longest to shortest.
	hosts bool       // whether any patterns contain hostnames
}

type muxEntry struct {
	h       Handler
	pattern string
}

// DefaultServeMux is the default ServeMux used by Serve.
var DefaultEventBus = &defaultEventBus

var defaultEventBus EventBus

// HandleFunc registers the handler function for the given pattern.
func (mux *ServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	mux.Handle(pattern, HandlerFunc(handler))
}

// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (mux *ServeMux) Handle(pattern string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if _, exist := mux.m[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}
	e := muxEntry{h: handler, pattern: pattern}
	mux.m[pattern] = e
	if pattern[len(pattern)-1] == '/' {
		mux.es = appendSorted(mux.es, e)
	}

	if pattern[0] != '/' {
		mux.hosts = true
	}
}
func appendSorted(es []muxEntry, e muxEntry) []muxEntry {
	n := len(es)
	i := sort.Search(n, func(i int) bool {
		return len(es[i].pattern) < len(e.pattern)
	})
	if i == n {
		return append(es, e)
	}
	// we now know that i points at where we want to insert
	es = append(es, muxEntry{}) // try to grow the slice in place, any entry works.
	copy(es[i+1:], es[i:])      // Move shorter entries down
	es[i] = e
	return es
}

func (b *EventBus) HandleFunc(pattern string, handler func(http.ResponseWriter, *Request)) {
	DefaultEventBus.mux.HandleFunc(pattern, handler)
}

// EventBus to get subscription and do HandleFunc.
type EventBus struct {
	redisC *redis.Client
	// use for unsub target channel
	subC map[string]*redis.PubSub

	wg  sync.WaitGroup
	mux ServeMux
}

//------------------------------------------------------------------------------
func init() {
	os.Setenv("REDIS_URL", "redis://localhost:6379/0")
}

func NewOptions(opts *redis.Options) *EventBus {
	return &EventBus{redisC: redis.NewClient(opts),
		subC: make(map[string]*redis.PubSub)}
}

func NewDefault() *EventBus {
	return NewOptions(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

// redisUrl = os.Getenv("REDIS_URL")
func New() (b *EventBus, err error) {
	return NewURL(os.Getenv("REDIS_URL"))
}

// redisUrl, _ := url.Parse("redis://localhost:6379")
func NewURL(s string) (b *EventBus, err error) {
	redisPassword := ""
	redisURL, err := url.Parse(s)
	if redisURL.User != nil {
		if password, ok := redisURL.User.Password(); ok {
			redisPassword = password
		}
	}
	db := 0
	if len(redisURL.Path) > 1 {
		db, err = strconv.Atoi(strings.TrimPrefix(redisURL.Path, "/"))
		if err != nil {
			return
		}
	}

	b = NewOptions(&redis.Options{
		Addr:     redisURL.Host,
		Password: redisPassword,
		DB:       db, // use default DB
	})
	return
}

// Close closes the client, releasing any open resources.
//
// It is rare to Close a Client, as the Client is meant to be
// long-lived and shared between many goroutines.
func (b *EventBus) Close() error {
	// don't need because of pubsub is in pool of redisC
	// for _, subC := range b.subC {
	// 	subC.Close()
	// }
	return b.redisC.Close()
}

//------------------------------------------------------------------------------
// Publish posts the message to the channel.
// http://redisdoc.com/pubsub/pubsub.html
func (b *EventBus) Broadcast(msg ...interface{}) {
	b.PPublish("", msg...)
}

func (b *EventBus) PPublish(pattern string, msg ...interface{}) {
	channels := b.redisC.PubSubChannels(pattern).Val()
	for _, channel := range channels {
		b.redisC.Publish(channel, fmt.Sprintln(msg...))
	}
}

func (b *EventBus) Publish(channel string, msg ...interface{}) {
	b.redisC.Publish(channel, fmt.Sprint(msg...))
}

//------------------------------------------------------------------------------
// sub message
// usage
// for msg := range ch {
// 	fmt.Println("send: ", msg.Channel, msg.Pattern, msg.Payload)
// 	err = ws.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
// }
func (b *EventBus) Subscribe(channels ...string) <-chan *redis.Message {
	pubsub := b.redisC.Subscribe(channels...)
	iface, err := pubsub.Receive()
	if err != nil {
		return nil
	}
	switch iface.(type) {
	case *redis.Subscription:
		// subscribe succeeded
	case *redis.Message:
		// received first message
	case *redis.Pong:
		// pong received
	default:
		// handle error
	}

	for _, channel := range channels {
		b.subC[channel] = pubsub
	}
	return pubsub.Channel()
}

//http://redisdoc.com/pubsub/psubscribe.html
func (b *EventBus) PSubscribe(pattern ...string) <-chan *redis.Message {
	pubsub := b.redisC.PSubscribe(pattern...)
	iface, err := pubsub.Receive()
	if err != nil {
		return nil
	}
	switch iface.(type) {
	case *redis.Subscription:
		// subscribe succeeded
	case *redis.Message:
		// received first message
	case *redis.Pong:
		// pong received
	default:
		// handle error
	}

	for _, channel := range pattern {
		b.subC[channel] = pubsub
	}
	return pubsub.Channel()
}

func (b *EventBus) UnSubscribe(channels ...string) error {
	for _, channel := range channels {
		if b.subC[channel] == nil {
			continue
		}
		err := b.subC[channel].Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//------------------------------------------------------------------------------
// evt callback note that msg only be string. so recommend JSON form

func (b *EventBus) Register(channel string, fn interface{}) {

	// select channel message and do handler
	go func() {

	}()
}

func (b *EventBus) UnRegister(channel string) {}

func (b *EventBus) Event(channel string, msg ...interface{}) {}