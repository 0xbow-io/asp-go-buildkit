package chainalysis

import (
	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
)

type Config struct {
	ChainalysisAPIKey string `env:"CHAINALYSIS_API_KEY"`
}

func (c Config) GetLogger() *logrus.Logger {
	logger := logrus.New()
	lvl, err := logrus.ParseLevel(logrus.InfoLevel.String())
	if err != nil {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(lvl)
	}

	logger.SetReportCaller(true)
	if true {
		logger.SetFormatter(joonix.NewFormatter())
	}
	return logger
}
