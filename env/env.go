package env

import "os"

const PsmUndefined string = "UNDEFINED_PSM"

func PSM() string {
	if psm := os.Getenv("PSM"); psm != "" {
		return psm
	}
	return PsmUndefined
}
