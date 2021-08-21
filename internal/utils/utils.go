package utils

import (
	"strconv"
	"strings"
)

const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
)

func ParseSize(size string) (bytes int64) {
	delimiter := len(size) - 1
	usize, err := strconv.ParseInt(size[:delimiter], 10, 64)
	if err != nil {
		panic(err)
	}
	switch suffix := strings.ToUpper(size[delimiter:]); suffix {
	case "G":
		bytes = GB * usize
	case "M":
		bytes = MB * usize
	case "K":
		bytes = KB * usize
	default:
		bytes = usize
	}
	return
}
