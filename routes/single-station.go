package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	errorMiddleware "github.com/wisdom-oss/microservice-middlewares/v5/error"

	"microservice/graphql"
)

func SingleStation(w http.ResponseWriter, r *http.Request) {
	errorHandler := r.Context().Value(errorMiddleware.ChannelName).(chan<- interface{})
	stationID := chi.URLParam(r, "stationID")
	measurementStation, err := graphql.Query{}.Station(struct{ WebsiteID string }{WebsiteID: stationID})

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(measurementStation)
	if err != nil {
		errorHandler <- fmt.Errorf("error while encoding measurement stations: %w", err)
		return
	}
}
