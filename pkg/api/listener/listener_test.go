package listener

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var body = strings.NewReader(`{"event":"event","name":"test_1","address":"http://localhost:8090/test"}`)

func TestRegister(t *testing.T) {
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
			l := NewHandlers(nil)
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

func TestUnregister(t *testing.T) {
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
