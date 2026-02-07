package dbutils

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// errorFilters is a list of functions that attempt to classify a database error
// into a known application-level error type. Each filter returns:
//   - a boolean indicating whether it matches the given error
//   - a mapped error to return if it matches
var errorFilters = []func(err error) (bool, error){
	filterDuplicationType,
	filterRecordNotFound,
}

// CatchDBErr converts low-level database errors into normalized application errors.
//
// If err is nil, it returns nil.
// Otherwise, it passes the error through a list of filters. The first filter
// that matches will return its mapped error. If no filters match, the original
// error is returned unchanged.
func CatchDBErr(err error) error {
	if err == nil {
		return nil
	}
	for _, filter := range errorFilters {
		check, newErr := filter(err)
		if check {
			return newErr
		}
	}
	return err
}

var (
	// ErrDuplicationType indicates a database duplication/unique constraint error.
	// This is typically used when an insert or update violates a unique constraint.
	ErrDuplicationType = errors.New("duplicate type")

	// ErrNotFoundType indicates that a requested record was not found in the database.
	// This usually maps to gorm.ErrRecordNotFound.
	ErrNotFoundType = errors.New("not found type")
)

// filterDuplicationType checks whether the given error represents a unique
// constraint / duplication error (e.g., duplicate key).
//
// It returns true and ErrDuplicationType if the error message suggests a
// uniqueness constraint violation; otherwise it returns false and ErrDuplicationType.
func filterDuplicationType(err error) (bool, error) {
	return strings.Contains(strings.ToLower(err.Error()), "unique constraint"), ErrDuplicationType
}

// filterRecordNotFound checks whether the given error is a "record not found" error
// from GORM.
//
// It returns true and ErrNotFoundType if the error matches gorm.ErrRecordNotFound;
// otherwise it returns false and ErrNotFoundType.
func filterRecordNotFound(err error) (bool, error) {
	return errors.Is(err, gorm.ErrRecordNotFound), ErrNotFoundType
}
