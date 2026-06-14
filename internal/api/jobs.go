package api

import (
	"encoding/json"
	"net/http"

	"github.com/Modd1e/durarun/internal/postgres/dbgen"
)

type createJobRequest struct {
	Queue   int32  `json:"queue"`
	Payload string `json:"payload"`
}

func (a *API) createJob(w http.ResponseWriter, r *http.Request) {
	var request createJobRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	job, err := a.queries.CreateJob(r.Context(), dbgen.CreateJobParams{
		Queue:   &request.Queue,
		Payload: &request.Payload,
	})
	if err != nil {
		a.log.Error("create job", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(job); err != nil {
		a.log.Error("encode job response", "error", err)
	}
}

func (a *API) jobsCount(w http.ResponseWriter, r *http.Request) {
	count, err := a.queries.CountJobs(r.Context())
	if err != nil {
		a.log.Error("count jobs", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]int64{
		"count": count,
	})
}
