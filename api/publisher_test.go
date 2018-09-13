package main

import (
	"bytes"
	"fmt"
	"github.com/volodimyr/event_publisher/api/models"
	"io"
	"io/ioutil"
	"log"
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
				t.Fatalf("Test work incorrect [%v] \n", err)
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
				t.Fatalf("Test work incorrect [%v] \n", err)
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
				t.Fatalf("Test work incorrect [%v] \n", err)
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

func TestEndToEnd(t *testing.T) {
	tearDown, publisher := setupTearDown(t)
	defer tearDown(t)

	event := "{'fake':'event'}"
	//create fake server with body checking
	fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Test work incorrect [%v] \n", err)
		}
		actual := string(bs)
		if actual != event {
			t.Errorf("Expected published event [%s], but got [%s]\n", event, actual)
		}
	}))

	eventName := "publish_event"
	//register listener
	r, err := http.NewRequest("POST", "https://domain.com/listener",
		strings.NewReader(fmt.Sprintf(`{"event":"%s","name":"test_1","address":"%s"}`, eventName, fake.URL)))

	w := httptest.NewRecorder()
	if err != nil {
		t.Fatalf("Test work incorrect [%v] \n", err)
	}
	register(publisher, w, r)
	res := w.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("It looks like endpoint '/listener' works unexpectedly wrong."+
			"Expected status code [%d] actual status code [%d] \n", http.StatusCreated, res.StatusCode)
	}

	//publish event to registered listener
	r, err = http.NewRequest("POST", "https://domain.com/publish/"+eventName, strings.NewReader(event))
	w = httptest.NewRecorder()
	if err != nil {
		t.Fatalf("Test work incorrect [%v] \n", err)
	}
	publish(publisher, w, r)
	res = w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("It looks like endpoint '/publish/:event' works unexpectedly wrong."+
			"Expected status code [%d] actual status code [%d] \n", http.StatusOK, res.StatusCode)
	}
}

func benchSetupTearDown(b *testing.B) (func(b *testing.B), *Publisher) {
	//log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
	b.Log("Setup Publisher for benchmark test")
	p := &Publisher{
		Events:    make(map[string]map[string]string, 10),
		New:       make(chan models.WorkNew, 10),
		Discard:   make(chan models.WorkDiscard, 10),
		Broadcast: make(chan models.WorkPublisher, 10),
		Stop:      make(chan struct{}),
	}
	go p.Service()

	return func(b *testing.B) {
		p.Stop <- struct{}{}
		b.Log("Tear down benchmark test")
	}, p
}

func setupRequest(b *testing.B, method string, path string, body io.Reader) (*http.Request, *httptest.ResponseRecorder) {
	r, err := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	if err != nil {
		b.Fatalf("Benchmark incorrect setup [%v] \n", err)
	}

	return r, w
}

func BenchmarkRegister(b *testing.B) {
	tearDown, publisher := benchSetupTearDown(b)
	defer tearDown(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r, w := setupRequest(b, "POST", "https://domain.com/listener",
			bytes.NewBuffer([]byte(`{"event":"event_1","name":"test_1","address":"http://localhost:8090/test"}`)))
		register(publisher, w, r)
		res := w.Result()
		if res.StatusCode != http.StatusCreated {
			b.Errorf("It looks like endpoint '/listener' works unexpectedly wrong."+
				"Expected status code [%d] actual status code [%d] \n", http.StatusCreated, res.StatusCode)
		}
	}
}

func BenchmarkPublish(b *testing.B) {
	tearDown, publisher := benchSetupTearDown(b)
	defer tearDown(b)
	event := "event_publish"
	listener := "publish_listener"
	publisher.Events[event] = map[string]string{listener: "https://localhost:8080"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r, w := setupRequest(b, "POST", "https://domain.com/publish/event_publish",
			bytes.NewBuffer([]byte(`{"data": "random"}`)))
		publish(publisher, w, r)
		res := w.Result()
		if res.StatusCode != http.StatusOK {
			b.Errorf("It looks like endpoint '/publish/:name' works unexpectedly wrong."+
				"Expected status code [%d] actual status code [%d] \n", http.StatusOK, res.StatusCode)
		}
	}
}
