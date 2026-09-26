package domain

import "errors"

// ErrOrderNotFound возвращается, когда заказ не найден.
var ErrOrderNotFound = errors.New("order not found")
