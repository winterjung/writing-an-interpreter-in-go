package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"go-interpreter/repl"
)

func main() {
	kst, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic(err)
	}
	now := time.Now().In(kst).Format("2006-01-02 15:04:05")
	_, _ = fmt.Printf("Unnamed Programming Language (main, %s) on %s(%s)\n", now, runtime.GOOS, runtime.GOARCH)
	repl.Start(os.Stdin, os.Stdout)
}
