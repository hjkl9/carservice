package conv

import (
	"strconv"
	"strings"
)

func ToStringWithSep(sep byte, d ...uint8) string {
	l := len(d)
	var builder strings.Builder
	// 1,2,3,4,5
	builder.Grow(l*2 - 1)
	for k, v := range d {
		builder.WriteString(strconv.FormatUint(uint64(v), 10))
		if k != l-1 {
			builder.WriteByte(sep)
		}
	}
	return builder.String()
}
