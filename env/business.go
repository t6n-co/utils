package env

import "os"

const PsmUndefined string = "UNDEFINED_PSM"

func PSM() string {
	if psm := os.Getenv(KeyPSM); psm != "" {
		return psm
	}
	return PsmUndefined
}
