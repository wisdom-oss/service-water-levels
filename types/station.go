package types

import (
	"encoding/json"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
)

type Station struct {
	WebsiteID *string `db:"website_id"`
	PublicID  *string `db:"public_id"`
	Name      *string `db:"name"`
	Operator  *string `db:"operator"`
	Location  *geom.T `db:"location"`
}

// MarshalJSON is used to control the generation of the GeoJSON output for the
// station's location.
// This is required since the [geom.T] interface does not implement the
// conversion to GeoJSON directly
func (s Station) MarshalJSON() ([]byte, error) {
	loc, err := geojson.Marshal(*s.Location)
	if err != nil {
		return nil, err
	}
	outputStation := struct {
		WebsiteID *string         `json:"websiteID,omitempty"`
		PublicID  *string         `json:"publicID,omitempty"`
		Name      *string         `json:"name,omitempty"`
		Operator  *string         `json:"operator,omitempty"`
		Location  json.RawMessage `json:"location,omitempty"`
	}{
		WebsiteID: s.WebsiteID,
		PublicID:  s.PublicID,
		Name:      s.Name,
		Operator:  s.Operator,
		Location:  loc,
	}
	return json.Marshal(outputStation)
}
