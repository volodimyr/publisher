package publisher

import (
	"fmt"
	"github.com/volodimyr/publisher/pkg/models"
	"github.com/volodimyr/publisher/pkg/persistence"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const (
	event        = "test_event"
	publishedMsg = `{"data":"default"}`
)

var (
	storage *persistence.Storage
	logger  *log.Logger
)

func setup() {
	logger = log.New(os.Stdout, "test: ", log.LstdFlags|log.Lshortfile)
	if storage == nil {
		storage = persistence.New(logger)
	}
}

func shutdown() {
	if storage != nil {
		storage.Stop <- struct{}{}
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestPublish(t *testing.T) {
	tests := []struct {
		name           string
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
		expectedBody   string
	}{
		{name: "POST", in: httptest.NewRequest("POST", "/publish/event",
			strings.NewReader(publishedMsg)),
			out: httptest.NewRecorder(), expectedStatus: http.StatusNotFound, expectedBody: errorNotRegistered + "\n"},
		{name: "GET", in: httptest.NewRequest("GET", "/publish/event", nil),
			out: httptest.NewRecorder(), expectedStatus: http.StatusMethodNotAllowed, expectedBody: postOnly + "\n"},
		{name: "PUT", in: httptest.NewRequest("PUT", "/publish/event", strings.NewReader(`{"data":"random"}`)),
			out: httptest.NewRecorder(), expectedStatus: http.StatusMethodNotAllowed, expectedBody: postOnly + "\n"},
		{name: "PUT_NIL_BODY", in: httptest.NewRequest("PUT", "/publish/event", nil),
			out: httptest.NewRecorder(), expectedStatus: http.StatusMethodNotAllowed, expectedBody: postOnly + "\n"},
		{name: "DELETE", in: httptest.NewRequest("DELETE", "/publish/event", strings.NewReader(`{"data":"random"}`)),
			out: httptest.NewRecorder(), expectedStatus: http.StatusMethodNotAllowed, expectedBody: postOnly + "\n"},
		{name: "DELETE_NIL_BODY", in: httptest.NewRequest("DELETE", "/publish/event", nil),
			out: httptest.NewRecorder(), expectedStatus: http.StatusMethodNotAllowed, expectedBody: postOnly + "\n"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			l := NewHandlers(logger, storage)
			l.publish(test.out, test.in)
			if test.out.Code != test.expectedStatus {
				t.Logf("Expected [%d], but got [%d]", test.expectedStatus, test.out.Code)
				t.Fail()
			}

			body := test.out.Body.String()
			if body != test.expectedBody {
				t.Logf("Expected [%s], but got [%s]", test.expectedBody, body)
				t.Fail()
			}
		})
	}
}

func fakeServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Test work incorrect [%v] \n", err)
		}
		actual := string(bs)
		if actual != publishedMsg {
			t.Logf("Expected published event [%s], but got [%s]\n", event, actual)
			t.Fail()
		}
	}))
}

func setupFakeServer(t *testing.T) {
	fake := fakeServer(t)
	done := make(chan struct{})
	storage.New <- persistence.Add{Done: done, Listener: models.Listener{Event: event, Name: fake.URL, Address: fake.URL}}
	<-done
}

func TestPublishWithFakeServer(t *testing.T) {
	setupFakeServer(t)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", fmt.Sprintf("/publish/%s", event), strings.NewReader(publishedMsg))
	p := NewHandlers(logger, storage)

	p.publish(w, r)

	//redundant
	if w.Code != http.StatusOK {
		t.Logf("Expected status [%d], but got [%d]", http.StatusOK, w.Code)
		t.Fail()
	}
}

func TestNewHandlers(t *testing.T) {
	p := NewHandlers(nil, storage)

	if p.logger == nil {
		t.Log("Logger cannot be nil")
		t.Fail()
	}
}
