package models

import "fmt"

//Events represents entity of named events
//Value of the first map is Listener in format name: address
//One event obviously can contain multiple listeners
type Events map[string]map[string]string

//Listener represents entity of servers who are looking for new messages
//None of these fields can be empty
//Address should be in the format http://domain.com/endpoint, but it isn't restricted
type Listener struct {
	Event   string `json:"event"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

//IsEmpty checks whether fields are not nil
func (listn *Listener) IsEmpty() error {
	if listn.Name == "" {
		return fmt.Errorf("empty 'Name' field. Validation error [%v]", listn)
	}
	if listn.Event == "" {
		return fmt.Errorf("empty 'Event' field. Validation error [%v]", listn)
	}
	if listn.Address == "" {
		return fmt.Errorf("empty 'Address' field. Validation error [%v]", listn)
	}
	return nil
}

//PublishMessage defines event and therefore listeners where messsage should be published
type PublishMessage struct {
	Event string
	Body  []byte
}

//IsEmpty checks whether fields are not nil
func (e *PublishMessage) IsEmpty() error {
	if e.Event == "" {
		return fmt.Errorf("empty 'Name' field. Validation error [%v]", e)
	}
	if e.Body == nil || len(e.Body) == 0 {
		return fmt.Errorf("empty 'Body' field. Validation error [%v]", e)
	}

	return nil
}
