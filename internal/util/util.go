package util

import "strconv"

func FormatPort(port int) string {
	return ":" + strconv.Itoa(port)
}
