package dbutils

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var errorFilters = []func(err error) (bool, error){
	filterDuplicationType,
	filterRecordNotFound,
}

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
	ErrDuplicationType = errors.New("duplicate type")
	ErrNotFoundType    = errors.New("not found type")
)

func filterDuplicationType(err error) (bool, error) {
	return strings.Contains(strings.ToLower(err.Error()), "unique constraint"), ErrDuplicationType
}

func filterRecordNotFound(err error) (bool, error) {
	return errors.Is(err, gorm.ErrRecordNotFound), ErrNotFoundType
}
