package logging

import (
	"mtgpoolservice/setting"

	log "github.com/sirupsen/logrus"
)

// Setup initialize the log instance
func Setup(settings setting.App) {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
}
