package errors

import "fmt"

// Все возможные ошибки при работе с данными.
var (
	ErrInvalidDataFormat    = fmt.Errorf("invalid data format")
	ErrInvalidDate          = fmt.Errorf("invalid date or date format")
	ErrInvalidRepeatRule    = fmt.Errorf("invalid repeat rule")
	ErrInvalidMonth         = fmt.Errorf("invalid month")
	ErrDublicateMonth       = fmt.Errorf("dublicate month")
	ErrInvalidDay           = fmt.Errorf("invalid day")
	ErrDublicateDay         = fmt.Errorf("dublicate day")
	ErrMissingId            = fmt.Errorf("missing identifier")
	ErrInvalidId            = fmt.Errorf("invalid identifier")
	ErrTaskNotFound         = fmt.Errorf("task not found")
	ErrInvalidPassword      = fmt.Errorf("invalid password")
	ErrInvalidSigningMethod = fmt.Errorf("invalid token signing method")
	ErrInvalidClaims        = fmt.Errorf("invalid claims")
	ErrExpiredToken         = fmt.Errorf("expired token")
	ErrInvalidToken         = fmt.Errorf("invalid token")
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
