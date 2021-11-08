package server

import "strings"

func imageTagChangedOrNoChange(was, next string) bool {
	if was == next {
		return true
	}
	wasParts := strings.SplitN(was, ":", 2)
	nextParts := strings.SplitN(next, ":", 2)

	if wasParts[0] == nextParts[0] {
		return true
	}
	return false
}
