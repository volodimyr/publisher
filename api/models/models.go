package models

import (
	"fmt"
)

type PublishedEvent struct {
	Name string
	Body []byte
}

type WorkPublisher struct {
	Done  chan struct{}
	Event PublishedEvent
}

type WorkNew struct {
	Done     chan struct{}
	Listener Listener
}

type WorkDiscard struct {
	Name string
	Done chan struct{}
}

type Listener struct {
	Event   string `json:"event"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (listn *Listener) EmptyStrings() error {
	if listn.Name == "" {
		return fmt.Errorf("Empty 'Name' field. Validation error. [%v]\n", listn)
	}
	if listn.Event == "" {
		return fmt.Errorf("Empty 'Event' field. Validation error. [%v]\n", listn)
	}
	if listn.Address == "" {
		return fmt.Errorf("Empty 'Address' field. Validation error. [%v]\n", listn)
	}
	return nil
}
