package springcfg

import (
	"strings"
)

func buildPath(path, delimiter string) []string {
	parts := strings.Split(path, delimiter)
	rslt := parts[:0]
	for _, p := range parts {
		if p != "" {
			rslt = append(rslt, p)
		}
	}
	return rslt
}
