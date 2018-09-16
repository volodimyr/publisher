package publisher

import (
	"github.com/volodimyr/publisher/pkg/models"
	"github.com/volodimyr/publisher/pkg/response"
	"github.com/volodimyr/publisher/pkg/storage"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	published = "Published"
	postOnly  = "POST method only"

	errorNotRegistered = "Event wasn't registered"
)

//Handlers handles /publish endpoints
//It also holds essential dependencies to be using
type Handlers struct {
	logger *log.Logger
}

//SetupRoutes setups all initial endpoints for listener handlers
func (h *Handlers) SetupRoutes(sm *http.ServeMux) {
	sm.HandleFunc("/publish/", h.Logger(h.publish))
}

func (h *Handlers) publish(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		defer r.Body.Close()
		eventNames := strings.Split(r.URL.Path, "/publish/")
		if len(eventNames) < 2 || eventNames[1] == "" {
			h.logger.Println("server: Event name has been empty")
			http.Error(w, "Event name must be specified", http.StatusBadRequest)
			return
		}
		s := storage.New(h.logger)
		if _, ok := s.Events[eventNames[1]]; !ok {
			h.logger.Println("server: Couldn't publish to non-existing event")
			http.Error(w, errorNotRegistered, http.StatusNotFound)
			return
		}
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			h.logger.Println("server: Invalid body")
			http.Error(w, "Cannot read body", http.StatusBadRequest)
			return
		}
		done := make(chan struct{})
		s.Broadcast <- storage.Publish{Done: done, PublishMessage: models.PublishMessage{Event: eventNames[1], Body: bs}}
		<-done
		resp.OK(w, published)
		return
	}
	h.logger.Printf("server: method [%s] not available for publish endpoint\n", r.Method)
	http.Error(w, postOnly, http.StatusMethodNotAllowed)
}

//Logger is a middleware for the publish handlers
func (h *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer h.logger.Printf("request processed in [%s]\n", time.Now().Sub(start))
		next(w, r)
	}
}

//NewHandlers create Publish Handlers and establishes all dependencies
//Important SetupRoutes needs to be called before it can be used
//if logger == nil, default will be taken
func NewHandlers(logger *log.Logger) *Handlers {
	if logger == nil {
		logger = log.New(os.Stdout, "server: ", log.LstdFlags|log.Lshortfile)
	}
	return &Handlers{logger: logger}
}
