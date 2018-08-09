package models

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	errs "github.com/time-sheet/apierrors"
)

// ClockInRequest ...
type ClockInRequest struct {
	Assignment    string     `json:"assignment"`
	ExpectedHours float64    `json:"expected_hours"`
	ID            CompoundID `json:"id"`
	Note          string     `json:"note,omitempty"`
}

const defaultExpectedHours = 7.5

// UnmarshalClockInRequest ...
func UnmarshalClockInRequest(reader io.Reader, time time.Time, clockInRequest *ClockInRequest) error {

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("unable to read request body: [%v]", err)
		return errs.ErrFailToParseRequestBody
	}

	// if body is empty do not try to unmarshal
	if len(b) != 0 {
		if err = json.Unmarshal(b, &clockInRequest); err != nil {
			fmt.Printf("unable to unmarshal request body: [%v]", err)
			return errs.ErrFailToUnmarshalRequestBody
		}
	}

	if clockInRequest.Assignment == "" {
		clockInRequest.Assignment = "not-applicable"
	}

	if clockInRequest.ExpectedHours == 0 {
		clockInRequest.ExpectedHours = defaultExpectedHours
	}

	clockInRequest.ID = CompoundID{
		Day:   time.Day(),
		Month: int(time.Month()),
		Year:  time.Year(),
	}

	return nil
}
