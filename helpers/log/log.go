package log

import "github.com/sirupsen/logrus"

func Error(args ...interface{}) {
	logrus.Error(args...)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

func Printf(format string, args ...interface{}) {
	logrus.Printf(format, args...)
}

func Println(args ...interface{}) {
	logrus.Println(args...)
}
