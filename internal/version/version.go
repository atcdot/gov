package version

import (
	"regexp"
)

func ExtractMinor(raw string) string {
	return regexp.MustCompile(`\d+(\.\d+)?`).FindString(raw)
}

func ExtractFull(raw string) string {
	return regexp.MustCompile(`\d+((\.\d+)+)?`).FindString(raw)
}
