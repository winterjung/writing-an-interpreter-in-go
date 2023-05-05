package parser

import (
	"fmt"
	"strings"
	"sync/atomic"
)

const (
	debug = false

	traceIdentPlaceholder = "\t"
)

var traceLevel int32

func trace(msg string) string {
	atomic.AddInt32(&traceLevel, 1)
	tracePrint("BEGIN " + msg)
	return msg
}

func untrace(msg string) {
	tracePrint("END " + msg)
	atomic.AddInt32(&traceLevel, -1)
}

func tracePrint(msg string) {
	if debug {
		_, _ = fmt.Printf("%s%s\n", indent(), msg)
	}
}

func indent() string {
	return strings.Repeat(traceIdentPlaceholder, int(traceLevel-1))
}
