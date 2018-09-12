package main

import (
	"github.com/volodimyr/event_publisher/api/models"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTearDown(t *testing.T) (func(t *testing.T), *Publisher) {
	t.Log("Setup Publisher for testing endpoints")
	p := &Publisher{
		Events:    make(map[string]map[string]string, 10),
		New:       make(chan models.WorkNew, 10),
		Discard:   make(chan models.WorkDiscard, 10),
		Broadcast: make(chan models.WorkPublisher, 10),
		Stop:      make(chan struct{}),
	}
	go p.Service()
	return func(t *testing.T) {
		p.Stop <- struct{}{}
		t.Log("Tear down. Test is off.")
	}, p
}

func TestRegister(t *testing.T) {
	tearDown, publisher := setupTearDown(t)
	defer tearDown(t)
	cases := []struct {
		method string
		body   io.Reader
		status int
		path   string
	}{{
		"PUT",
		strings.NewReader(`{"event":"PUT","name":"test_1","address":"http://localhost:8090/test"}`),
		405,
		"https://domain.com/listener",
	}, {
		"DELETE",
		nil,
		405,
		"https://domain.com/listener",
	}, {
		"GET",
		nil,
		405,
		"https://domain.com/listener",
	}, {
		"POST",
		strings.NewReader(`{"event":"event_1","name":"test_1","address":"http://localhost:8090/test"}`),
		201,
		"https://domain.com/listener",
	}}
	for _, c := range cases {
		t.Run("Test register: "+c.method, func(t *testing.T) {
			r, err := http.NewRequest(c.method, c.path, c.body)
			w := httptest.NewRecorder()
			if err != nil {
				t.Error(err)
			}
			register(publisher, w, r)
			res := w.Result()
			if res.StatusCode != c.status {
				t.Errorf("It looks like endpoint '/listener' works unexpectedly wrong."+
					"Expected status code [%d] actual status code [%d] \n", c.status, res.StatusCode)
			}
		})
	}
	if v, _ := publisher.Events["event_1"]["test_1"]; v != "http://localhost:8090/test" {
		t.Errorf("Event [event_1] and listener [test_1]'http://localhost:8090/test' should be created\n")
	}
}

func TestUnregister(t *testing.T) {
	event := "delete_event"
	listener := "delete_listener"
	tearDown, publisher := setupTearDown(t)
	defer tearDown(t)
	publisher.Events[event] = map[string]string{listener: "local"}
	cases := []struct {
		method string
		body   io.Reader
		status int
		path   string
	}{{
		"POST",
		strings.NewReader(`{"event":"POST","name":"test_1","address":"http://localhost:8090/test"}`),
		405,
		"https://domain.com/listener/" + listener,
	}, {
		"PUT",
		strings.NewReader(`{"event":"PUT","name":"test_1","address":"http://localhost:8090/test"}`),
		405,
		"https://domain.com/listener/" + listener,
	}, {
		"DELETE",
		nil,
		200,
		"https://domain.com/listener/" + listener,
	}, {
		"GET",
		nil,
		405,
		"https://domain.com/listener/" + listener,
	}}

	for _, c := range cases {
		t.Run("Test unregister: "+c.method, func(t *testing.T) {
			r, err := http.NewRequest(c.method, c.path, c.body)
			w := httptest.NewRecorder()
			if err != nil {
				t.Error(err)
			}
			unregister(publisher, w, r)
			res := w.Result()
			if res.StatusCode != c.status {
				t.Errorf("It looks like endpoint '/listener/:name' works unexpectedly wrong."+
					"Expected status code [%d] actual status code [%d] \n", c.status, res.StatusCode)
			}
		})
	}
	if _, ok := publisher.Events[event]; !ok {
		t.Errorf("Delete event should be persisted")
	}
	if _, ok := publisher.Events[event][listener]; ok {
		t.Errorf("Delete listener should not be registered")
	}
}

func TestPublish(t *testing.T) {
	event := "event_publish"
	listener := "publish_listener"
	tearDown, publisher := setupTearDown(t)
	defer tearDown(t)

	publisher.Events[event] = map[string]string{listener: "https://localhost:8080"}
	cases := []struct {
		method string
		body   io.Reader
		status int
		path   string
	}{{
		"POST",
		strings.NewReader(`{"data": "random"}`),
		404,
		"https://domain.com/publish/event_1",
	}, {
		"PUT",
		strings.NewReader(`{"data": "random"}`),
		405,
		"https://domain.com/publish/event_1",
	}, {
		"DELETE",
		nil,
		405,
		"https://domain.com/publish/event_1",
	}, {
		"POST",
		strings.NewReader(`{"data": "random"}`),
		200,
		"https://domain.com/publish/event_publish",
	}, {
		"GET",
		nil,
		405,
		"https://domain.com/publish/event_1",
	}}
	for _, c := range cases {
		t.Run("Test publish: "+c.method, func(t *testing.T) {
			r, err := http.NewRequest(c.method, c.path, c.body)
			w := httptest.NewRecorder()
			if err != nil {
				t.Error(err)
			}
			publish(publisher, w, r)
			res := w.Result()
			if res.StatusCode != c.status {
				t.Errorf("It looks like endpoint '/publish/:event' works unexpectedly wrong."+
					"Expected status code [%d] actual status code [%d] \n", c.status, res.StatusCode)
			}
		})
	}
}

func EndToEnd(t *testing.T) {

}
