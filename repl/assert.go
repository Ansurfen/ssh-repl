package repl

import "github.com/sirupsen/logrus"

func Panic(e error) {
	if e != nil {
		panic(e)
	}
}

func PanicWithLogger(err error, msg string) {
	if err != nil {
		logrus.Fatalf("%s error: %v", msg, err)
	}
}
