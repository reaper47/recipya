package config

import "errors"

var ErrWaitNegative = errors.New("config error: 'wait' must be >= 1")

var ErrIndexIntervalInvalid = errors.New("config error: 'indexInterval' is invalid. Please check the legal values in config.yaml")
