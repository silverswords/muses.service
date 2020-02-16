package service

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v7"
)

type eventbus struct {
	redisC *redis.Client
	subC   map[string]*redis.PubSub
}

func init() {
	os.Setenv("REDIS_URL", "redis://localhost:6379/0")
}

func NewEventBus() *eventbus {
	return &eventbus{redisC: redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})}
}

// redisUrl, _ := url.Parse("redis://localhost:6379")
func NewEventBusURL(s string) (b *eventbus, err error) {
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

	b = &eventbus{redisC: redis.NewClient(&redis.Options{
		Addr:     redisURL.Host,
		Password: redisPassword,
		DB:       db, // use default DB
	})}
	return
}

// redisUrl = os.Getenv("REDIS_URL")
func NewEventURLenv() (b *eventbus, err error) {
	return NewEventBusURL(os.Getenv("REDIS_URL"))
}

func NewEventBusOptions(opts *redis.Options) *eventbus {
	return &eventbus{redisC: redis.NewClient(opts)}
}

//------------------------------------------------------------------------------
// Publish posts the message to the channel.
// http://redisdoc.com/pubsub/pubsub.html
func (b *eventbus) Broadcast(msg ...interface{}) {
	b.PPublish("", msg...)
}

func (b *eventbus) PPublish(pattern string, msg ...interface{}) {
	channels := b.redisC.PubSubChannels(pattern).Val()
	for _, channel := range channels {
		b.redisC.Publish(channel, fmt.Sprintln(msg...))
	}
}

func (b *eventbus) Publish(topic string, msg ...interface{}) {
	b.redisC.Publish(topic, fmt.Sprint(msg...))
}

//------------------------------------------------------------------------------
// sub message
// usage
// for msg := range ch {
// 	fmt.Println("send: ", msg.Channel, msg.Pattern, msg.Payload)
// 	err = ws.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
// }
func (b *eventbus) Subscribe(topic ...string) <-chan *redis.Message {
	pubsub := b.redisC.Subscribe(topic...)
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

	return pubsub.Channel()
}

//http://redisdoc.com/pubsub/psubscribe.html
func (b *eventbus) PSubscribe(pattern ...string) <-chan *redis.Message {
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

	return pubsub.Channel()
}

func (b *eventbus) UnSubscribe(topic string) error {
	return b.subC[topic].Close()
}

func (b *eventbus) BUnSubscribe(topics ...string) error {
	for _, topic := range topics {
		err := b.subC[topic].Close()
		if err != nil {
			return err
		}
	}
	return nil
}
