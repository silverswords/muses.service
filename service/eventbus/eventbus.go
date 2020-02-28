package eventbus

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

//https://github.com/grpc/grpc-go/blob/master/stream.go#L122:2
//stream read write
type Endpoint interface {
	Descriptor() string
	Receive() (Event, error)
	Detach()
	Probe() bool
}

type EventBus interface {
	Attach(topic string) Endpoint
	Send(topic string, e Event)
}

type Event interface {
	Type() string
	Payload() []byte
}

type Pong struct {
	pID     string
	payload []byte
}

func NewMessage(pid string, payload []byte) Pong { return Pong{pID: pid, payload: payload} }
func (p Pong) Type() string                      { return p.pID }
func (p Pong) Payload() []byte                   { return p.payload }

type Redbus struct {
	epcount int
	eps     map[string]map[int]chan Event
	rm      sync.RWMutex
}

func New() *Redbus {
	return &Redbus{
		epcount: 0,
		eps:     make(map[string]map[int]chan Event),
		rm:      sync.RWMutex{},
	}
}

func (eb *Redbus) Attach(topic string) Endpoint {
	eb.rm.Lock()
	ep := eb.Newendpoint(topic + " " + strconv.Itoa(eb.epcount))
	if _, found := eb.eps[topic]; !found {
		eb.eps[topic] = make(map[int]chan Event)
	}
	eb.eps[topic][eb.epcount] = ep.Ch
	eb.epcount++
	eb.rm.Unlock()
	return ep
}

func (eb *Redbus) Send(topic string, data Event) {
	eb.rm.RLock()
	if eps, found := eb.eps[topic]; found {
		go func(data Event, eps map[int]chan Event) {
			for _, ep := range eps {
				ep <- data.(Event)
			}
		}(data, eps)
	}
	eb.rm.RUnlock()
}

type endpoint struct {
	eb         Redbus
	descriptor string
	Ch         chan Event
}

func (eb *Redbus) Newendpoint(descriptor string) *endpoint {
	return &endpoint{eb: *eb, descriptor: descriptor, Ch: make(chan Event, 100)}
}

func (ep *endpoint) Send(topic string, e Event) (err error) {
	//add time out control
	ep.eb.Send(topic, e)
	return nil
}

// a block method until receive a event from the eventbus
func (ep *endpoint) Receive() (data Event, err error) {
	data, haddata := <-ep.Ch
	if !haddata {
		return data, errors.New("channel no data and closed")
	}
	return data, nil
}

func (ep *endpoint) Detach() { //	close channel on EventBus
	desc := strings.Split(ep.descriptor, " ")
	epcount, _ := strconv.Atoi(desc[1])
	delete(ep.eb.eps[desc[0]], epcount)
	close(ep.Ch)
}

func (ep *endpoint) Probe() bool {
	desc := strings.Split(ep.descriptor, " ")
	epcount, _ := strconv.Atoi(desc[1])
	if _, found := ep.eb.eps[desc[0]][epcount]; !found {
		return false
	}
	return true
}

func (ep *endpoint) Descriptor() string { return ep.descriptor }
