package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/falcosecurity/falcosidekick/types"
	"golang.org/x/net/websocket"
)

type eventPayload struct {
	UUID    string             `json:"uuid,omitempty"`
	Event   types.FalcoPayload `json:"event,omitempty"`
	Stats   map[string]int64   `json:"stats,omitempty"`
	Outputs []string           `json:"outputs,omitempty"`
}

type eventStore struct {
	statsByUUID map[string]map[string]int64
	Events      []types.FalcoPayload `json:"events,omitempty"`
	Stats       map[string]int64     `json:"stats,omitempty"`
	Outputs     []string             `json:"outputs,omitempty"`
	Retention   int                  `json:"retention,omitempty"`
}

var (
	broadcast chan eventStore
	store     eventStore
	retention int
	mutex     *sync.RWMutex
)

func init() {
	broadcast = make(chan eventStore, 20)
	store.statsByUUID = make(map[string]map[string]int64)
	store.Stats = make(map[string]int64)
	mutex = &sync.RWMutex{}
}

func main() {
	a := flag.String("a", "0.0.0.0", "Listen Address")
	p := flag.Int("p", 2802, "Listen Port")
	r := flag.Int("r", 200, "Number of events to keep in retention")
	flag.Parse()

	if ip := net.ParseIP(*a); ip == nil {
		log.Fatalf("[ERROR] : Failed to parse Address")
	}

	store.Retention = *r

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("/events", eventsHandler)
	http.Handle("/ui", http.StripPrefix("/ui", http.FileServer(http.Dir("./static"))))
	http.Handle("/ws", websocket.Handler(socket))

	log.Printf("[INFO]  : Falco Sidekick Web UI is up and listening on %s:%d\n", *a, *p)
	log.Printf("[INFO]  : Retention is %d last events\n", store.Retention)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", *a, *p), nil); err != nil {
		log.Fatalf("[ERROR] : %v\n", err.Error())
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please send a valid request body", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Please send with post http method", http.StatusBadRequest)
		return
	}

	d := json.NewDecoder(r.Body)
	d.UseNumber()

	var e eventPayload
	err := d.Decode(&e)
	if err != nil {
		http.Error(w, "Please send a valid request body", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	store.Outputs = e.Outputs
	if len(store.Events) >= store.Retention {
		store.Events = append(store.Events[1:len(store.Events)-1], e.Event)
	} else {
		store.Events = append(store.Events, e.Event)
	}
	store.statsByUUID[e.UUID] = e.Stats
	mutex.Unlock()

	temp := make(map[string]int64)
	mutex.RLock()
	for _, i := range store.statsByUUID {
		for j, k := range i {
			temp[j] += k
			temp["Total"] += k
		}
	}
	mutex.RUnlock()
	mutex.Lock()
	store.Stats = temp
	mutex.Unlock()

	broadcast <- store
}

// healthHandler is a simple handler to test if daemon is UP.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	// nolint: errcheck
	w.Write([]byte(`{"status": "ok"}`))
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	s, _ := json.Marshal(store)

	// nolint: errcheck
	w.Write(s)
}

func socket(ws *websocket.Conn) {
	log.Printf("[INFO]  : A Websocket connection to WebUI has been established\n")
	for {
		events := <-broadcast
		if err := websocket.JSON.Send(ws, events); err != nil {
			break
		}
	}
}
