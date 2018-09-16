package models

import "testing"

func TestListener_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		listener Listener
	}{
		{name: "Empty event", listener: Listener{Event: "", Name: "Default", Address: "Default"}},
		{name: "Empty name", listener: Listener{Event: "Default", Name: "", Address: "Default"}},
		{name: "Empty address", listener: Listener{Event: "Default", Name: "Default", Address: ""}},
		{name: "Empty", listener: Listener{Event: "", Name: "", Address: ""}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test := test
			if err := test.listener.IsEmpty(); err == nil {
				t.Logf("[%s] - Test failure. IsEmpty func works incorrect.", test.name)
				t.Fail()
			}
		})
	}
}

func TestPublishMessage_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		pm   PublishMessage
	}{
		{name: "Empty", pm: PublishMessage{Event: "", Body: []byte("")}},
		{name: "Empty event", pm: PublishMessage{Event: "", Body: []byte(`{data": "random"}`)}},
		{name: "Empty body", pm: PublishMessage{Event: "Default", Body: []byte("")}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.pm.IsEmpty(); err == nil {
				t.Logf("[%s] - Test failure. IsEmpty func works incorrect.", test.name)
				t.Fail()
			}
		})
	}
}
