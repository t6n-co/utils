package env

import "os"

const PsmUndefined string = "psm"

func PSM() string {
	if psm := os.Getenv(KeyPSM); psm != "" {
		return psm
	}
	return PsmUndefined
}
