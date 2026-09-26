package service

import "errors"

// ErrKafkaUnavailable — публикация события не удалась.
var ErrKafkaUnavailable = errors.New("kafka unavailable")
