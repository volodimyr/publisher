package storage

import (
	"github.com/volodimyr/publisher/pkg/client"
	"github.com/volodimyr/publisher/pkg/models"
	"log"
	"sync"
)

var (
	once sync.Once

	instance *storage
)

type storage struct {
	models.Events
	New       chan Add
	Discard   chan Discard
	Broadcast chan Publish
	Stop      chan struct{}
}

//New establishes new persistent and concurrency safe storage
//Returns *storage which is a singleton instance
//logger must be specified by the first call of the function
func New(logger *log.Logger) *storage {
	once.Do(func() {
		instance = &storage{
			Events:    make(map[string]map[string]string, 10),
			New:       make(chan Add, 10),
			Discard:   make(chan Discard, 10),
			Broadcast: make(chan Publish, 10),
			Stop:      make(chan struct{}),
		}
		go instance.service(logger)
	})
	return instance
}

//Publish is a type of work for broadcasting message between whole event: []listeners
//PublishMessage defines event and therefore listeners where messsage should be published
//Done uses for notifying caller everything is done
type Publish struct {
	Done chan struct{}
	models.PublishMessage
}

//Add is a type of work using to add new event: listeners{}
//Done uses for notifying caller everything is done
type Add struct {
	Done chan struct{}
	models.Listener
}

//Discard is a type of work to remove a specific listener
//Name is the listener name
//Done uses for notifying caller everything is done
type Discard struct {
	Name string
	Done chan struct{}
}

func (s *storage) service(logger *log.Logger) {
	logger.Println("Publisher service is online")
	for {
		select {
		case n := <-s.New:
			//register new listener into existing event
			if reg, ok := s.Events[n.Listener.Event]; ok {
				reg[n.Listener.Name] = n.Listener.Address
				n.Done <- struct{}{}
				logger.Printf("Registered new listener [%v] into existing event [%s]\n", n.Listener, s.Events[n.Listener.Event])
				continue
			}
			//create new event and add new listener
			s.Events[n.Listener.Event] = map[string]string{n.Listener.Name: n.Listener.Address}
			logger.Printf("Created new event [%s] and registered new listener [%s]\n", n.Listener.Event, n.Listener.Name)
			n.Done <- struct{}{}
		case d := <-s.Discard:
			for _, Listeners := range s.Events {
				if _, ok := Listeners[d.Name]; ok {
					delete(Listeners, d.Name)
				}
			}
			logger.Printf("Discard executed for the next listeners [%s]\n", d.Name)
			d.Done <- struct{}{}
		case b := <-s.Broadcast:
			for name, url := range s.Events[b.PublishMessage.Event] {
				logger.Printf("Sending event to the next listener: [%s] at [%s]\n", name, url)
				client.DoPOST(url, b.PublishMessage.Body, logger)
			}
			logger.Printf("Broadcasted message for the event [%s]\n", b.PublishMessage.Event)
			b.Done <- struct{}{}
		case <-s.Stop:
			logger.Println("Publisher service is offline")
			return
		}
	}
}
