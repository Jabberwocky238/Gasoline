package device

import "github.com/sirupsen/logrus"

var log = logrus.New()

func init() {
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}
