package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	graphql2 "github.com/graph-gophers/graphql-go"
	"github.com/jackc/pgx/v5/pgtype"
	errorMiddleware "github.com/wisdom-oss/microservice-middlewares/v5/error"

	"microservice/graphql"
	"microservice/types"
)

func Measurements(w http.ResponseWriter, r *http.Request) {
	errorHandler := r.Context().Value(errorMiddleware.ChannelName).(chan<- interface{})
	var from, until time.Time
	var station string
	var err error
	if r.URL.Query().Has("from") {
		rawFromTimestamp := r.URL.Query().Get("from")
		from, err = time.Parse(time.RFC3339, rawFromTimestamp)
		if err != nil {
			errorHandler <- err
		}
	}

	if r.URL.Query().Has("until") {
		rawFromTimestamp := r.URL.Query().Get("until")
		until, err = time.Parse(time.RFC3339, rawFromTimestamp)
		if err != nil {
			errorHandler <- err
		}
	}

	if r.URL.Query().Has("station") {
		station = r.URL.Query().Get("station")
	}

	measurements, err := graphql.Query{}.Measurements(&graphql.MeasurementArguments{
		From:    &graphql2.Time{Time: from},
		Until:   &graphql2.Time{Time: until},
		Station: &station,
	})
	if err != nil {
		errorHandler <- err
		return
	}

	var output []types.Measurement
	for _, measurement := range measurements {
		pgStation := pgtype.Text{}
		_ = pgStation.Scan(measurement.Station)

		pgDate := pgtype.Date{}
		_ = pgDate.Scan(measurement.Date.Time)

		pgClassification := pgtype.Text{}
		if measurement.Classification != nil {
			_ = pgClassification.Scan(*measurement.Classification)
		}

		pgNHN := pgtype.Numeric{}
		if measurement.WaterLevelNHN != nil {
			_ = pgNHN.Scan(fmt.Sprintf("%f", *measurement.WaterLevelNHN))

		}

		pgGOK := pgtype.Numeric{}
		if measurement.WaterLevelGOK != nil {
			_ = pgGOK.Scan(fmt.Sprintf("%f", *measurement.WaterLevelGOK))

		}

		output = append(output, types.Measurement{
			Station:        pgStation,
			Date:           pgDate,
			Classification: pgClassification,
			WaterLevelNHN:  pgNHN,
			WaterLevelGOK:  pgGOK,
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		errorHandler <- fmt.Errorf("error while encoding measurement stations: %w", err)
		return
	}
}
