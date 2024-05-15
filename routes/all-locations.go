package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"microservice/graphql"
)
import errorMiddleware "github.com/wisdom-oss/microservice-middlewares/v5/error"

func AllLocations(w http.ResponseWriter, r *http.Request) {
	errorHandler := r.Context().Value(errorMiddleware.ChannelName).(chan<- interface{})

	var q graphql.Query
	stations, err := q.Stations()
	if err != nil {
		errorHandler <- err
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(stations)
	if err != nil {
		errorHandler <- fmt.Errorf("error while encoding measurement stations: %w", err)
		return
	}
}
