package graphql

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/graph-gophers/graphql-go"
	"github.com/jackc/pgx/v5/pgtype"

	"microservice/globals"
	"microservice/types"
)

type MeasurementArguments struct {
	From    *graphql.Time
	Until   *graphql.Time
	Station *string
}

type Measurement struct {
	// Station contains the Station.WebsiteID of the water measurement station at
	// which the measurement has been taken
	Station string `json:"station" db:"station"`

	// Date contains the date at which the measurement has been taken
	Date graphql.Time `json:"date" db:"date"`

	// Classification contains the textual representation of the water levels
	// classification.
	// The values are documented here:
	// https://www.grundwasserstandonline.nlwkn.niedersachsen.de/Hinweis#einstufungGrundwasserstandsklassen
	Classification *string `json:"classification" db:"classification"`

	// WaterLevelNHN refers to the water level in reference to the sea level in
	// Germany
	WaterLevelNHN *float64 `json:"waterLevelNHN" db:"water_level_nhn"`

	// WaterLevelGOK refers to the water level in reference to the terrain
	// height around the measurement station
	WaterLevelGOK *float64 `json:"waterLevelGOK" db:"water_level_gok"`
}

func (Query) Measurements(args *MeasurementArguments) ([]Measurement, error) {
	var rawQuery string
	var err error
	var queryArgs []interface{}
	switch {
	case args.Until == nil && args.From == nil && args.Station == nil:
		rawQuery, err = globals.SqlQueries.Raw("get-measurements")

	case args.Until == nil && args.From == nil && args.Station != nil:
		rawQuery, err = globals.SqlQueries.Raw("get-measurements-for-station")
		queryArgs = append(queryArgs, *args.Station)

	case (args.Until != nil || args.From != nil) && (args.Station == nil || (args.Station != nil && strings.TrimSpace(*args.Station) == "")):
		rawQuery, err = globals.SqlQueries.Raw("get-measurements-in-range")
		if args.Until == nil || args.Until.IsZero() {
			args.Until = &graphql.Time{Time: time.Now()}
		}
		from := pgtype.Date{}
		err := from.Scan((*args.From).Time)
		if err != nil {
			return nil, err
		}

		until := pgtype.Date{}
		err = until.Scan((*args.Until).Time)
		if err != nil {
			return nil, err
		}

		queryArgs = append(queryArgs, from, until)

	case (args.Until != nil || args.From != nil) && args.Station != nil:
		rawQuery, err = globals.SqlQueries.Raw("get-measurements-for-station-in-range")
		if args.Until == nil || args.Until.IsZero() {
			args.Until = &graphql.Time{Time: time.Now()}
		}
		from := pgtype.Date{}
		err := from.Scan((*args.From).Time)
		if err != nil {
			return nil, err
		}

		until := pgtype.Date{}
		err = until.Scan((*args.Until).Time)
		if err != nil {
			return nil, err
		}

		queryArgs = append(queryArgs, from, until, *args.Station)
	}
	if err != nil {
		return nil, err
	}

	var measurements []types.Measurement
	err = pgxscan.Select(context.Background(), globals.Db, &measurements, rawQuery, queryArgs...)
	if err != nil {
		return nil, err
	}

	var output []Measurement
	for _, m := range measurements {
		var nhn, gok []byte
		var nhnf64, gokf64 *float64

		nhn, _ = m.WaterLevelNHN.MarshalJSON()
		gok, _ = m.WaterLevelGOK.MarshalJSON()
		err = json.Unmarshal(nhn, &nhnf64)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(gok, &gokf64)
		if err != nil {
			return nil, err
		}

		output = append(output, Measurement{
			Station:        m.Station.String,
			Date:           graphql.Time{Time: m.Date.Time},
			Classification: &m.Classification.String,
			WaterLevelNHN:  nhnf64,
			WaterLevelGOK:  gokf64,
		})
	}
	return output, err

}
