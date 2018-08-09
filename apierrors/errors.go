package apierrors

import "errors"

// A list of error messages for Dataset API
var (
	ErrFailToParseRequestBody     = errors.New("failed to parse request body")
	ErrFailToUnmarshalRequestBody = errors.New("failed to unmarshal request body")
	ErrStillClockedIn             = errors.New("still clocked in, you must clock out before clocking back in")
	ErrTimesheetNotFound          = errors.New("timesheet not found")

	NotFoundMap = map[error]bool{
		ErrTimesheetNotFound: true,
	}

	BadRequestMap = map[error]bool{
		ErrFailToParseRequestBody:     true,
		ErrFailToUnmarshalRequestBody: true,
	}

	ForbiddenRequestMap = map[error]bool{}

	ConflictRequestMap = map[error]bool{
		ErrStillClockedIn: true,
	}
)
