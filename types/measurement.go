package types

import "github.com/jackc/pgx/v5/pgtype"

type Measurement struct {
	// Station contains the Station.WebsiteID of the water measurement station at
	// which the measurement has been taken
	Station pgtype.Text `json:"station" db:"station"`

	// Date contains the date at which the measurement has been taken
	Date pgtype.Date `json:"date" db:"date"`

	// Classification contains the textual representation of the water levels
	// classification.
	// The values are documented here:
	// https://www.grundwasserstandonline.nlwkn.niedersachsen.de/Hinweis#einstufungGrundwasserstandsklassen
	Classification pgtype.Text `json:"classification" db:"classification"`

	// WaterLevelNHN refers to the water level in reference to the sea level in
	// Germany
	WaterLevelNHN pgtype.Numeric `json:"waterLevelNHN" db:"water_level_nhn"`

	// WaterLevelGOK refers to the water level in reference to the terrain
	// height around the measurement station
	WaterLevelGOK pgtype.Numeric `json:"waterLevelGOK" db:"water_level_gok"`
}
