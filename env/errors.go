package env

import "errors"

var (
	ErrorMetricsPortNotSet  = errors.New("metrics port not set")
	ErrorMetricsPortInvalid = errors.New("metrics port invalid")
)
