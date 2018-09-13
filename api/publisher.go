package main

import (
	"encoding/json"
	"github.com/volodimyr/event_publisher/api/client"
	"github.com/volodimyr/event_publisher/api/config"
	"github.com/volodimyr/event_publisher/api/models"
	"github.com/volodimyr/event_publisher/api/response"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Publisher struct {
	Events    map[string]map[string]string
	New       chan models.WorkNew
	Discard   chan models.WorkDiscard
	Broadcast chan models.WorkPublisher
	Stop      chan struct{}
}

func main() {
	publisher := &Publisher{
		Events:    make(map[string]map[string]string, 10),
		New:       make(chan models.WorkNew, 10),
		Discard:   make(chan models.WorkDiscard, 10),
		Broadcast: make(chan models.WorkPublisher, 10),
		Stop:      make(chan struct{}),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/listener", func(w http.ResponseWriter, req *http.Request) {
		register(publisher, w, req)
	})
	log.Println("server: The next POST endpoint '/listener' could be used for register new listener")
	mux.HandleFunc("/listener/", func(w http.ResponseWriter, req *http.Request) {
		unregister(publisher, w, req)
	})
	log.Println("server: The next DELETE endpoint '/listener/:name' could be used for unregister a listener")
	mux.HandleFunc("/publish/", func(w http.ResponseWriter, req *http.Request) {
		publish(publisher, w, req)
	})
	log.Println("server: The next POST endpoint '/publish/:event_name' could be used for publish message to registered listeners")

	go publisher.Service()

	log.Printf("server: Sevrer starts working at http://localhost%s/\n", config.Port)
	err := http.ListenAndServe(config.Port, mux)
	if err != nil {
		publisher.Stop <- struct{}{}
		log.Fatalln("server: ", err)
	}
	publisher.Stop <- struct{}{}
}

func (p *Publisher) Service() {
	log.Println("server: Publisher service is online")
	for {
		select {
		case wn := <-p.New:
			//register new listener into existing event
			if reg, ok := p.Events[wn.Listener.Event]; ok {
				reg[wn.Listener.Name] = wn.Listener.Address
				wn.Done <- struct{}{}
				log.Printf("server: Registered new listener [%v] into existing event [%s]\n", wn.Listener, p.Events[wn.Listener.Event])
				continue
			}
			//create new event and add new listener
			p.Events[wn.Listener.Event] = map[string]string{wn.Listener.Name: wn.Listener.Address}
			log.Printf("server: Created new event [%s] and registered new listener [%s]\n", wn.Listener.Event, wn.Listener.Name)
			wn.Done <- struct{}{}
		case wd := <-p.Discard:
			for _, Listeners := range p.Events {
				if _, ok := Listeners[wd.Name]; ok {
					delete(Listeners, wd.Name)
				}
			}
			log.Printf("server: Discard executed for the next listeners [%s]\n", wd.Name)
			wd.Done <- struct{}{}
		case wp := <-p.Broadcast:
			for name, url := range p.Events[wp.Event.Name] {
				log.Printf("server: Sending event to the next listener: [%s] at [%s]\n", name, url)
				client.DoPOST(url, wp.Event.Body)
			}
			log.Printf("server: Broadcasted message for the event [%s]\n", wp.Event.Name)
			wp.Done <- struct{}{}
		case <-p.Stop:
			log.Println("server: Publisher service is offline")
			return
		}
	}
}

func register(p *Publisher, w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		defer req.Body.Close()
		dec := json.NewDecoder(req.Body)
		l := models.Listener{}
		err := dec.Decode(&l)
		if err != nil {
			log.Println("server: Invalid body")
			http.Error(w, "Couldn't parse body", http.StatusBadRequest)
			return
		}
		if err := l.EmptyStrings(); err != nil {
			log.Printf("server: Listener should containe valid non-empty fields [%v]\n", l)
			http.Error(w, "Body contains invalid values", http.StatusBadRequest)
			return
		}
		done := make(chan struct{})
		p.New <- models.WorkNew{Listener: l, Done: done}
		<-done
		response.Created(w)
		return
	}
	log.Printf("server: method [%s] not available for register endpoint\n", req.Method)
	http.Error(w, "POST method only", http.StatusMethodNotAllowed)
}

func unregister(p *Publisher, w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodDelete {
		lNames := strings.Split(req.URL.Path, "/listener/")
		if len(lNames) < 2 || lNames[1] != "" {
			done := make(chan struct{})
			p.Discard <- models.WorkDiscard{Name: lNames[1], Done: done}
			<-done
			response.OK(w)
			return
		}
		log.Printf("server: Listener name has been empty")
		http.Error(w, "Listener name must be specified", http.StatusBadRequest)
		return
	}
	log.Printf("server: method [%s] not available for unregister endpoint\n", req.Method)
	http.Error(w, "DELETE method only", http.StatusMethodNotAllowed)
}

func publish(p *Publisher, w http.ResponseWriter, req *http.Request) {
	defer latencyTrack(time.Now(), "Publishing")
	if req.Method == http.MethodPost {
		defer req.Body.Close()
		eventNames := strings.Split(req.URL.Path, "/publish/")
		if len(eventNames) < 2 || eventNames[1] == "" {
			log.Println("server: Event name has been empty")
			http.Error(w, "Event name must be specified", http.StatusBadRequest)
			return
		}
		if _, ok := p.Events[eventNames[1]]; !ok {
			log.Println("server: Couldn't publish to non-existing event")
			http.Error(w, "Event wasn't registered", http.StatusNotFound)
			return
		}
		bs, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println("server: Invalid body")
			http.Error(w, "Cannot read body", http.StatusBadRequest)
			return
		}
		done := make(chan struct{})
		p.Broadcast <- models.WorkPublisher{Done: done, Event: models.PublishedEvent{Body: bs, Name: eventNames[1]}}
		<-done
		response.OK(w)
		return
	}
	log.Printf("server: method [%s] not available for publish endpoint\n", req.Method)
	http.Error(w, "POST method only", http.StatusMethodNotAllowed)
}

func latencyTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
