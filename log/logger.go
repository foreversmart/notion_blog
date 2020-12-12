package log

import (
	formatter "github.com/banzaicloud/logrus-runtime-formatter"
	log "github.com/sirupsen/logrus"
	"os"
)

var Logger *log.Entry

func init() {
	// Log as JSON instead of the default ASCII formatter.
	runtimeFormatter := &formatter.Formatter{
		ChildFormatter: &log.TextFormatter{},
		Line:           true,
		File:           true,
	}
	log.SetFormatter(runtimeFormatter)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
	Logger = log.WithFields(log.Fields{
		"notion": "export",
	})
}
