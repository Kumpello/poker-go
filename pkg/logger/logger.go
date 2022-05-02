package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"time"
)

type Logger logrus.FieldLogger

func NewLogger() Logger {
	l := logrus.New()
	l.SetReportCaller(true)
	return l
}

func MakeEchoLogEntry(logger Logger, c echo.Context) Logger {
	if c == nil {
		return logger.WithFields(logrus.Fields{
			"at": time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	return logger.WithFields(logrus.Fields{
		"at":     time.Now().Format("2006-01-02 15:04:05"),
		"method": c.Request().Method,
		"uri":    c.Request().URL.String(),
		"ip":     c.Request().RemoteAddr,
	})
}
