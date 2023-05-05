package parser

import (
	"fmt"
	"strings"
)

const (
	debug = false

	traceIdentPlaceholder = "\t"
)

var traceLevel = 0

func trace(msg string) string {
	traceLevel++
	tracePrint("BEGIN " + msg)
	return msg
}

func untrace(msg string) {
	tracePrint("END " + msg)
	traceLevel--
}

func tracePrint(msg string) {
	if debug {
		_, _ = fmt.Printf("%s%s\n", indent(), msg)
	}
}

func indent() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel-1)
}
