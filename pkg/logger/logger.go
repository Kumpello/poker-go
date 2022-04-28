package logger

import "github.com/sirupsen/logrus"

type Logger logrus.FieldLogger

func NewLogger() Logger {
	l := logrus.New()
	l.SetReportCaller(true)
	return l
}
