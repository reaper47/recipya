package consts

import "errors"

// Error codes returned by invalid values in the configuration file.
var (
	ErrWaitNegative         = errors.New("config error: 'wait' must be >= 1")
	ErrIndexIntervalInvalid = errors.New(
		"config error: 'indexInterval' is invalid. Please check the legal values in config.yaml",
	)
)

var (
	ErrEntryNotFound = errors.New("entry not found")
)
