package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const callerSkip = 1

func assert(ok bool, args ...any) {
	if !ok {
		var caller string
		_, file, line, hasCaller := runtime.Caller(callerSkip)
		if hasCaller {
			caller = " " + callerPath(file) + ":" + strconv.Itoa(line) + ":"
		}

		_, _ = fmt.Fprintf(os.Stderr, "Fatal:%s %s", caller, fmt.Sprint(args...))
		os.Exit(1)
	}
}

var callerPrefix string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if ok {
		callerPrefix = filepath.ToSlash(filepath.Dir(filepath.Dir(file))) + "/"
	}
}

func callerPath(path string) string {
	return strings.TrimPrefix(filepath.ToSlash(path), callerPrefix)
}
