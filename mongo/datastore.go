package mongo

import (
	"errors"
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	errs "github.com/time-sheet/apierrors"
	"github.com/time-sheet/service/models"
)

// Mongo represents a simplistic MongoDB configuration.
type Mongo struct {
	Collection     string
	Database       string
	lastPingTime   time.Time
	lastPingResult error
	Session        *mgo.Session
	URI            string
}

// Init creates a new mgo.Session with a strong consistency and a write mode of "majority".
func (m *Mongo) Init() (session *mgo.Session, err error) {
	if session != nil {
		return nil, errors.New("session already exists")
	}

	if session, err = mgo.Dial(m.URI); err != nil {
		return nil, err
	}

	session.EnsureSafe(&mgo.Safe{WMode: "majority"})
	session.SetMode(mgo.Strong, true)
	return session, nil
}

// AddTimesheet ...
func (m *Mongo) AddTimesheet(timesheet *models.Timesheet) error {
	s := m.Session.Copy()
	defer s.Close()

	timesheet.LastUpdated = time.Now().UTC().Format(time.RFC3339Nano)
	if err := s.DB(m.Database).C(m.Collection).Insert(&timesheet); err != nil {
		return err
	}

	return nil
}

// GetTimesheet ...
func (m *Mongo) GetTimesheet(day int, month time.Month, year int) (*models.Timesheet, error) {
	s := m.Session.Copy()
	defer s.Close()

	selector := bson.M{
		"id.day":   day,
		"id.month": month,
		"id.year":  year,
	}

	var timesheet models.Timesheet
	err := s.DB(m.Database).C(m.Collection).Find(selector).One(&timesheet)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, errs.ErrTimesheetNotFound
		}
		return nil, err
	}

	return &timesheet, nil
}

// UpdateTimesheet ...
func (m *Mongo) UpdateTimesheet(timesheet *models.Timesheet) error {
	s := m.Session.Copy()
	defer s.Close()

	selector := bson.M{
		"id": timesheet.ID,
	}

	update := bson.M{"$set": updateTimesheet(timesheet)}

	if err := s.DB(m.Database).C(m.Collection).Update(selector, update); err != nil {
		if err == mgo.ErrNotFound {
			return errs.ErrTimesheetNotFound
		}
		return err
	}

	return nil
}

func updateTimesheet(timesheet *models.Timesheet) bson.M {
	updates := make(bson.M)

	fmt.Printf("building update query for timesheet resource: [%v]", timesheet)

	if timesheet.Assignment != "" {
		updates["assignment"] = timesheet.Assignment
	}

	if timesheet.CompletedHours != 0 {
		updates["completed_hours"] = timesheet.CompletedHours
	}

	if timesheet.ExpectedHours != 0 {
		updates["expected_hours"] = timesheet.ExpectedHours
	}

	updates["last_updated"] = time.Now().UTC().Format(time.RFC3339Nano)

	if timesheet.Note != "" {
		updates["note"] = timesheet.Note
	}

	if timesheet.Time != nil {
		updates["time"] = timesheet.Time
	}

	return updates
}
