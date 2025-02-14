package errors

import (
	"errors"
)

// Все возможные ошибки при работе с данными.
var (
	ErrInvalidDataFormat    = errors.New("invalid data format")
	ErrInvalidDate          = errors.New("invalid date or date format")
	ErrInvalidRepeatRule    = errors.New("invalid repeat rule")
	ErrInvalidMonth         = errors.New("invalid month")
	ErrDublicateMonth       = errors.New("dublicate month")
	ErrInvalidDay           = errors.New("invalid day")
	ErrDublicateDay         = errors.New("dublicate day")
	ErrMissingId            = errors.New("missing identifier")
	ErrInvalidId            = errors.New("invalid identifier")
	ErrTaskNotFound         = errors.New("task not found")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrInvalidSigningMethod = errors.New("invalid token signing method")
	ErrInvalidClaims        = errors.New("invalid claims")
	ErrExpiredToken         = errors.New("expired token")
	ErrInvalidToken         = errors.New("invalid token")
)

// Для быстрого поиска нужной ошибки сохраним их в мапу.
var ErrList map[error]interface{} = map[error]interface{}{
	ErrInvalidDataFormat: new(interface{}),
	ErrInvalidDate:       new(interface{}),
	ErrInvalidRepeatRule: new(interface{}),
	ErrInvalidMonth:      new(interface{}),
	ErrDublicateMonth:    new(interface{}),
	ErrInvalidDay:        new(interface{}),
	ErrDublicateDay:      new(interface{}),
	ErrMissingId:         new(interface{}),
	ErrInvalidId:         new(interface{}),
	ErrTaskNotFound:      new(interface{}),
	ErrInvalidPassword:   new(interface{}),
}
