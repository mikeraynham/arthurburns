package pathfmt

import (
	"strings"
)

func ToTitle(s string) string {
	return strings.Title(strings.Replace(s, "-", " ", -1))
}
