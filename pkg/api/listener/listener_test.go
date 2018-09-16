package listener

import (
	"fmt"
	"github.com/volodimyr/publisher/pkg/storage"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

const (
	lAddr = "http://localhost:8080/test"
	lName = "listener001"
)

func TestRegister(t *testing.T) {
	var body = strings.NewReader(fmt.Sprintf(`{"event":"event_001","name":"%s","address":"%s"}`, lName, lAddr))
	l := NewHandlers(nil)
	tests := []struct {
		name           string
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
		expectedBody   string
	}{
		{name: "GET", in: httptest.NewRequest("GET", "/listener", nil),
			out: httptest.NewRecorder(), expectedStatus: http.StatusMethodNotAllowed, expectedBody: postOnly + "\n"},
		{name: "PUT", in: httptest.NewRequest("PUT", "/listener", body),
			out: httptest.NewRecorder(), expectedStatus: http.StatusMethodNotAllowed, expectedBody: postOnly + "\n"},
		{name: "PUT_NIL_BODY", in: httptest.NewRequest("PUT", "/listener", nil),
			out: httptest.NewRecorder(), expectedStatus: http.StatusMethodNotAllowed, expectedBody: postOnly + "\n"},
		{name: "DELETE", in: httptest.NewRequest("DELETE", "/listener", body),
			out: httptest.NewRecorder(), expectedStatus: http.StatusMethodNotAllowed, expectedBody: postOnly + "\n"},
		{name: "DELETE_NIL_BODY", in: httptest.NewRequest("DELETE", "/listener", nil),
			out: httptest.NewRecorder(), expectedStatus: http.StatusMethodNotAllowed, expectedBody: postOnly + "\n"},
		{name: "POST", in: httptest.NewRequest("POST", "/listener", body),
			out: httptest.NewRecorder(), expectedStatus: http.StatusCreated, expectedBody: registered},
		{name: "POST_NIL_BODY", in: httptest.NewRequest("POST", "/listener", nil),
			out: httptest.NewRecorder(), expectedStatus: http.StatusBadRequest, expectedBody: invalidBody + "\n"},
		{name: "POST_INVALID_BODY", in: httptest.NewRequest("POST", "/listener", strings.NewReader("{absolutely epic}")),
			out: httptest.NewRecorder(), expectedStatus: http.StatusBadRequest, expectedBody: invalidBody + "\n"},
		{name: "POST_NIL_EVENT", in: httptest.NewRequest("POST", "/listener", strings.NewReader(`{"event":"","name":"test_1","address":"http://localhost:8090/test"}`)),
			out: httptest.NewRecorder(), expectedStatus: http.StatusBadRequest, expectedBody: invalidBody + "\n"},
		{name: "POST_NIL_ADDR", in: httptest.NewRequest("POST", "/listener", strings.NewReader(`{"event":"event","name":"test_1"}`)),
			out: httptest.NewRecorder(), expectedStatus: http.StatusBadRequest, expectedBody: invalidBody + "\n"},
		{name: "POST_NIL_NAME", in: httptest.NewRequest("POST", "/listener", strings.NewReader(`{"event":"event","addr":"localhost:8080"}`)),
			out: httptest.NewRecorder(), expectedStatus: http.StatusBadRequest, expectedBody: invalidBody + "\n"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			l.register(test.out, test.in)
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

func setupEvents(t *testing.T) map[string]map[string]string {
	events := make(map[string]map[string]string)
	const max = 15
	random := func(len int) string {
		bytes := make([]byte, len)
		for i := 0; i < len; i++ {
			bytes[i] = byte(65 + rand.Intn(25))
		}
		return string(bytes)
	}

	for i := 0; i < 10; i++ {
		eName := random(max)
		events[eName] = map[string]string{lName: lAddr}
	}

	w := httptest.NewRecorder()
	l := NewHandlers(nil)

	for k, _ := range events {
		body := strings.NewReader(fmt.Sprintf(`{"event":"%s","name":"%s","address":"%s"}`, k, lName, lAddr))
		r := httptest.NewRequest("POST", "/listener", body)
		l.register(w, r)
		if w.Code != http.StatusCreated {
			t.Logf("Expected [%d], but got [%d]", http.StatusCreated, w.Code)
			t.Fail()
		}
	}

	return events
}

func TestRegisterAndCheckStorage(t *testing.T) {
	setupEvents := setupEvents(t)

	s := storage.New(nil)
	for k, v := range setupEvents {
		if _, ok := s.Events[k]; !ok {
			t.Logf("Expected key [%s], but actual hasn't got it", k)
			t.Fail()
		}
		if !reflect.DeepEqual(v, s.Events[k]) {
			t.Logf("Expected key [%s], but actual hasn't got it", k)
			t.Fail()
		}
	}
}

func TestUnregisterAndCheckStorage(t *testing.T) {
	setupEvents(t)

	w := httptest.NewRecorder()
	l := NewHandlers(nil)
	s := storage.New(l.logger)
	r := httptest.NewRequest("DELETE", fmt.Sprintf("/listener/%s", lName), nil)

	l.unregister(w, r)

	for _, value := range s.Events {
		if len(value) != 0 {
			t.Logf("Expected len to be [0], but got [%d]", len(value))
			t.Fail()
		}
	}
}

func TestNewHandlers(t *testing.T) {
	l := NewHandlers(nil)

	if l.logger == nil {
		t.Log("Logger cannot be nil")
		t.Fail()
	}
}

func TestUnregister(t *testing.T) {
	var body = strings.NewReader(fmt.Sprintf(`{"event":"event_001","name":"%s","address":"%s"}`, lName, lAddr))
	tests := []struct {
		name           string
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
		expectedBody   string
	}{
		{name: "GET", in: httptest.NewRequest("GET", "/listener/event_1", nil), out: httptest.NewRecorder(),
			expectedStatus: http.StatusMethodNotAllowed, expectedBody: deleteOnly + "\n"},
		{name: "POST", in: httptest.NewRequest("POST", "/listener/event_1", body), out: httptest.NewRecorder(),
			expectedStatus: http.StatusMethodNotAllowed, expectedBody: deleteOnly + "\n"},
		{name: "POST_NIL_BODY", in: httptest.NewRequest("POST", "/listener/event_1", nil), out: httptest.NewRecorder(),
			expectedStatus: http.StatusMethodNotAllowed, expectedBody: deleteOnly + "\n"},
		{name: "PUT", in: httptest.NewRequest("PUT", "/listener/event_1", body), out: httptest.NewRecorder(),
			expectedStatus: http.StatusMethodNotAllowed, expectedBody: deleteOnly + "\n"},
		{name: "PUT_NIL_BODY", in: httptest.NewRequest("PUT", "/listener/event_1", nil), out: httptest.NewRecorder(),
			expectedStatus: http.StatusMethodNotAllowed, expectedBody: deleteOnly + "\n"},
		{name: "DELETE", in: httptest.NewRequest("DELETE", "/listener/event_1", nil), out: httptest.NewRecorder(),
			expectedStatus: http.StatusOK, expectedBody: unregistered},
		{name: "DELETE_WITH_ BODY", in: httptest.NewRequest("DELETE", "/listener/event_1", body), out: httptest.NewRecorder(),
			expectedStatus: http.StatusOK, expectedBody: unregistered},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			l := NewHandlers(nil)
			l.unregister(test.out, test.in)
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
