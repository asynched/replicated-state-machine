package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/asynched/replicated-state-machine/config"
	"github.com/asynched/replicated-state-machine/fsm"
	"github.com/asynched/replicated-state-machine/replication"
	"github.com/hashicorp/raft"
)

var (
	configFlag = flag.String("config", "config.toml", "path to config file")
)

func main() {
	flag.Parse()

	config, err := config.ParseConfig(*configFlag)

	if err != nil {
		panic(err)
	}

	ra, err := replication.GetRaft(config)

	if err != nil {
		panic(err)
	}

	router := http.NewServeMux()

	router.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Value int `json:"value"`
		}

		raftAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)

		if ra.Leader() != raft.ServerAddress(raftAddr) {
			http.Error(w, "Not a leader", http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		changeset := fsm.Changeset{
			Value: data.Value,
		}

		bytes, err := json.Marshal(changeset)

		if err != nil {
			panic(err)
		}

		if err := ra.Apply(bytes, time.Second*10).Error(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Value is applied"))
	})

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port+1)

	log.Printf("Starting server on http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
