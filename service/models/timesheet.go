package models

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// Timesheet ...
type Timesheet struct {
	Assignment     string        `json:"assignment"      bson:"assignment"`
	CompletedHours float64       `json:"completed_hours" bson:"completed_hours"`
	ExpectedHours  float64       `json:"expected_hours"  bson:"expected_hours"`
	ID             *CompoundID   `json:"id"              bson:"id"`
	LastUpdated    string        `json:"last_updated"    bson:"last_updated"`
	Note           string        `json:"note,omitempty"  bson:"note,omitempty"`
	Time           []*TimeRecord `json:"time"            bson:"time"`
}

// CompoundID ...
type CompoundID struct {
	Day   int `json:"day"   bson:"day"`
	Month int `json:"month" bson:"month"`
	Year  int `json:"year"  bson:"year"`
}

// TimeRecord ...
type TimeRecord struct {
	Start *TimeObject `json:"start,omitempty"   bson:"start,omitempty"`
	End   *TimeObject `json:"end,omitempty"     bson:"end,omitempty"`
}

// TimeObject ...
type TimeObject struct {
	Hour   int `json:"hour"   bson:"hour"`
	Minute int `json:"minute" bson:"minute"`
	Second int `json:"second" bson:"second"`
}

// UnmarshalTimeSheet ...
func UnmarshalTimeSheet(reader io.Reader, model Timesheet) (err error) {
	var b []byte
	b, err = ioutil.ReadAll(reader)
	if err != nil {
		return
	}

	if err = json.Unmarshal(b, &model); err != nil {
		return
	}

	return
}
