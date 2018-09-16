package publisher

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPublish(t *testing.T) {
	tests := []struct {
		name           string
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
		expectedBody   string
	}{
		{name: "POST", in: httptest.NewRequest("POST", "/publish/event",
			strings.NewReader(`{"data":"random"}`)),
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
			l := NewHandlers(nil)
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

//TODO: publish OK!
