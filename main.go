package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Debt is the thing you need to do.
type Debt struct {
	ID       int64     `json:"id"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Payments int64     `json:"payments"`
}

var debts map[int64]Debt
var debtMutex sync.RWMutex

func main() {
	debts = make(map[int64]Debt)

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/debts", listDebts).Methods("GET")
	r.HandleFunc("/debts", makeDebt).Methods("PUT")
	r.HandleFunc("/debts/{debtid}", showDebt).Methods("GET")
	r.HandleFunc("/debts/{debtid}", payDebt).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, debtor.")
}

func listDebts(w http.ResponseWriter, r *http.Request) {
	debtMutex.RLock()
	defer debtMutex.RUnlock()

	out, err := json.Marshal(debts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func makeDebt(w http.ResponseWriter, r *http.Request) {
	id := rand.Int63()
	debtMutex.Lock()
	defer debtMutex.Unlock()

	debt := Debt{
		ID:      id,
		Created: time.Now(),
		Updated: time.Now(),
	}

	debts[id] = debt
	out, err := json.Marshal(debt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func payDebt(w http.ResponseWriter, r *http.Request) {
	debtMutex.Lock()
	defer debtMutex.Unlock()

	did := mux.Vars(r)["debtid"]
	if did == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing a debt id")
		return
	}

	id, err := strconv.ParseInt(did, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid debt id %s", did)
		return
	}

	d := debts[id]
	if d.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	d.Updated = time.Now()
	d.Payments++
	debts[id] = d

	out, err := json.Marshal(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func showDebt(w http.ResponseWriter, r *http.Request) {
	debtMutex.RLock()
	defer debtMutex.RUnlock()

	did := mux.Vars(r)["debtid"]
	if did == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing a debt id")
		return
	}

	id, err := strconv.ParseInt(did, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid debt id %s", did)
		return
	}

	d := debts[id]
	if d.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)

	out, err := json.Marshal(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
