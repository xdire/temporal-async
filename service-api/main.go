package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/xdire/temporal-async/messaging"
	"github.com/xdire/temporal-async/util"
	"github.com/xdire/temporal-async/workflows"
	"go.temporal.io/sdk/client"
	"log"
	"net/http"
)

var tmprl client.Client

func main() {
	// set up the worker
	var err error
	tmprl, err = client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer tmprl.Close()

	// w := worker.New(c, "saga-process", worker.Options{})

	mux := http.NewServeMux()
	// curl -X POST http://localhost:8093/process/create?name=xxx&desc=testprocess
	mux.HandleFunc("/process/create", CreateProcess)
	// curl -X PATCH http://localhost:8093/process/stage?id=proc_uuid&stage=execute
	mux.HandleFunc("/process/stage", UpdateProcessStage)
	// curl -X PATCH http://localhost:8093/process/confirm?id=proc_uuid
	mux.HandleFunc("/process/confirm", UpdateProcessConfirmation)
	server := &http.Server{Addr: ":8093", Handler: mux}

	log.Fatal(server.ListenAndServe())
}

func CreateProcess(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	desc := r.URL.Query().Get("desc")
	_, err := tmprl.ExecuteWorkflow(r.Context(), client.StartWorkflowOptions{
		ID:        "proc_" + uuid.New().String(),
		TaskQueue: "saga-process",
	}, workflows.StartProcess, &messaging.CreateProcess{
		Name: name,
		Desc: desc,
	})
	if err != nil {
		http.Error(w, fmt.Errorf("unable to start workflow, error %w", err).Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateProcessStage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	stage := r.URL.Query().Get("stage")

	dbc := &util.DB{}
	err := dbc.Connect()
	if err != nil {
		http.Error(w, "no database connection", http.StatusInternalServerError)
		return
	}

	proc := &messaging.Process{}
	records := dbc.Conn().Where("uuid = ?", id).First(proc)
	if records.Error != nil {
		http.Error(w, fmt.Errorf("not found: %w", records.Error).Error(), http.StatusInternalServerError)
		return
	}

	err = tmprl.SignalWorkflow(r.Context(), id, "", "processStageUpdates", stage)
	if err != nil {
		http.Error(w, "failed to update process stage", http.StatusInternalServerError)
		return
	}
}

func UpdateProcessConfirmation(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	dbc := &util.DB{}
	err := dbc.Connect()
	if err != nil {
		http.Error(w, "no database connection", http.StatusInternalServerError)
		return
	}

	proc := &messaging.Process{}
	records := dbc.Conn().Where("uuid = ?", id).First(proc)
	if records.Error != nil {
		http.Error(w, fmt.Errorf("not found: %w", records.Error).Error(), http.StatusInternalServerError)
		return
	}

	err = tmprl.SignalWorkflow(r.Context(), id, "", "processUpdates", "confirmed")
	if err != nil {
		http.Error(w, fmt.Errorf("failed to update process, %w", err).Error(), http.StatusInternalServerError)
		return
	}
}
