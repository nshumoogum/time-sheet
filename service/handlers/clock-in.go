package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	errs "github.com/time-sheet/apierrors"
	"github.com/time-sheet/common"
	"github.com/time-sheet/service/models"
	"github.com/time-sheet/store"
)

// Store provides a backend for timesheets
type Store struct {
	store.Storer
}

// PostStart ...
func (api *Store) PostStart(w http.ResponseWriter, r *http.Request) {
	defer common.DrainBody(r)
	ctx := r.Context()

	time := time.Now().UTC()

	if r.Body == nil {
		fmt.Println("got here")
	}

	var clockInRequest models.ClockInRequest
	if err := models.UnmarshalClockInRequest(r.Body, time, &clockInRequest); err != nil {
		handleError(ctx, w, err)
		return
	}

	newTimeRecord := &models.TimeRecord{
		Start: &models.TimeObject{
			Hour:   time.Hour(),
			Minute: time.Minute(),
			Second: time.Second(),
		},
	}

	currentTimesheet, err := api.GetTimesheet(time.Day(), time.Month(), time.Year())
	if err != nil {
		if err != errs.ErrTimesheetNotFound {
			handleError(ctx, w, err)
			return
		}
	}

	var b []byte
	if currentTimesheet != nil {
		var newTimeRecords []*models.TimeRecord

		// Check not already clocked in
		for _, timeRecord := range currentTimesheet.Time {
			if timeRecord.End == nil {
				handleError(ctx, w, errs.ErrStillClockedIn)
				return
			}

			newTimeRecords = append(newTimeRecords, timeRecord)
		}

		// Append new time to previous time list
		newTimeRecords = append(newTimeRecords, newTimeRecord)

		updateTimesheet := &models.Timesheet{
			ID:   currentTimesheet.ID,
			Time: newTimeRecords,
		}

		// Update timesheet with new clock-in time
		if err = api.UpdateTimesheet(updateTimesheet); err != nil {
			handleError(ctx, w, err)
			return
		}

		currentTimesheet.Time = newTimeRecords

		b, err = json.Marshal(currentTimesheet)
		if err != nil {
			handleError(ctx, w, err)
			return
		}

	} else {
		timesheet := &models.Timesheet{
			Assignment:     clockInRequest.Assignment,
			CompletedHours: 0,
			ExpectedHours:  clockInRequest.ExpectedHours,
			ID:             &clockInRequest.ID,
			Note:           clockInRequest.Note,
		}

		var timeRecords []*models.TimeRecord
		timeRecords = append(timeRecords, newTimeRecord)
		timesheet.Time = timeRecords

		// Add new timesheet
		if err = api.AddTimesheet(timesheet); err != nil {
			handleError(ctx, w, err)
			return
		}

		b, err = json.Marshal(timesheet)
		if err != nil {
			handleError(ctx, w, err)
			return
		}
	}

	writeBody(ctx, w, b)
}

func handleError(ctx context.Context, w http.ResponseWriter, err error) {
	var status int
	switch {
	case errs.NotFoundMap[err]:
		status = http.StatusNotFound
	case errs.BadRequestMap[err]:
		status = http.StatusBadRequest
	case errs.ForbiddenRequestMap[err]:
		status = http.StatusForbidden
	case errs.ConflictRequestMap[err]:
		status = http.StatusConflict
	default:
		status = http.StatusInternalServerError
	}

	fmt.Printf("[%v]: request failed: [%v]", status, err)
	http.Error(w, err.Error(), status)
}

func writeBody(ctx context.Context, w http.ResponseWriter, b []byte) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(b); err != nil {
		fmt.Printf("failed to write body: [%v]", err)
	}
}
