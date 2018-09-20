package listener

import (
	"encoding/json"
	"github.com/volodimyr/publisher/pkg/models"
	"github.com/volodimyr/publisher/pkg/persistence"
	"github.com/volodimyr/publisher/pkg/response"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	registered   = "Registered"
	unregistered = "Removed"

	deleteOnly = "DELETE method only"
	postOnly   = "POST method only"

	invalidBody = "Body contains invalid values"
)

//Handlers handles /listener endpoints
//It also holds essential dependencies to be using
type Handlers struct {
	logger *log.Logger
	s      *persistence.Storage
}

//SetupRoutes setups all initial endpoints for listener handlers
func (h *Handlers) SetupRoutes(sm *http.ServeMux) {
	sm.HandleFunc("/listener", h.Logger(h.register))
	sm.HandleFunc("/listener/", h.Logger(h.unregister))
}

func (h *Handlers) register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		defer r.Body.Close()
		dec := json.NewDecoder(r.Body)
		l := models.Listener{}
		err := dec.Decode(&l)
		if err != nil {
			h.logger.Println("server: Invalid body")
			http.Error(w, invalidBody, http.StatusBadRequest)
			return
		}
		if err := l.IsEmpty(); err != nil {
			h.logger.Printf("server: Listener should containe valid non-empty fields [%v]\n", l)
			http.Error(w, invalidBody, http.StatusBadRequest)
			return
		}
		done := make(chan struct{})
		h.s.New <- persistence.Add{Listener: l, Done: done}
		<-done
		resp.Created(w, registered)
		return
	}
	h.logger.Printf("server: method [%s] not available for register endpoint\n", r.Method)
	http.Error(w, postOnly, http.StatusMethodNotAllowed)
}

func (h *Handlers) unregister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		if r.Body != nil {
			defer r.Body.Close()
		}
		lNames := strings.Split(r.URL.Path, "/listener/")
		if len(lNames) < 2 || lNames[1] != "" {
			done := make(chan struct{})
			h.s.Discard <- persistence.Discard{Name: lNames[1], Done: done}
			<-done
			resp.OK(w, unregistered)
			return
		}
		h.logger.Println("server: Listener name has been empty")
		http.Error(w, "Listener name must be specified", http.StatusBadRequest)
		return
	}
	h.logger.Printf("server: method [%s] not available for unregister endpoint\n", r.Method)
	http.Error(w, deleteOnly, http.StatusMethodNotAllowed)
}

//Logger is a middleware for the listener handlers
func (h *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer h.logger.Printf("request processed in [%s]\n", time.Now().Sub(start))
		next(w, r)
	}
}

//NewHandlers create Listener Handlers and establishes all dependencies
//Important SetupRoutes needs to be called before it can be used
//if logger == nil, default will be taken
func NewHandlers(logger *log.Logger, storage *persistence.Storage) *Handlers {
	if logger == nil {
		logger = log.New(os.Stdout, "server: ", log.LstdFlags|log.Lshortfile)
	}
	return &Handlers{logger: logger, s: storage}
}
