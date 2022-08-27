package repl

import (
	"errors"
	"flag"
	"os"
)

func parseFlags(args []string) {
	os.Args = []string{""}
	os.Args = append(os.Args, args...)
	flag.Parse()
}

func resetFlags() {
	path = ""
	cid = -1
	sid = -1
	session = false
}

func nextArg(args []string) ([]string, error) {
	if len(args) <= 0 {
		return args, errors.New("args is null")
	}
	args = args[1:]
	return args, nil
}
