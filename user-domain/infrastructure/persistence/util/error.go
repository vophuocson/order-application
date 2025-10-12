package persistenceutil

import (
	"errors"
	"fmt"
	domainerror "user-domain/internal/domain/error"

	"gorm.io/gorm"
)

func MapErrorToHTTPStatus(err error) error {

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return domainerror.ErrCodeNotFound
	case errors.Is(err, gorm.ErrMissingWhereClause),
		errors.Is(err, gorm.ErrInvalidField),
		errors.Is(err, gorm.ErrInvalidValue),
		errors.Is(err, gorm.ErrInvalidValueOfLength),
		errors.Is(err, gorm.ErrPrimaryKeyRequired),
		errors.Is(err, gorm.ErrModelValueRequired),
		errors.Is(err, gorm.ErrModelAccessibleFieldsRequired),
		errors.Is(err, gorm.ErrForeignKeyViolated),
		errors.Is(err, gorm.ErrCheckConstraintViolated):
		return domainerror.ErrCodeInvalidInput
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return domainerror.ErrCodeConflict

	default:
		return domainerror.ErrCodeInternal
	}

}

func Wrap(msg string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}
