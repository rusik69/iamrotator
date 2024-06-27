package args

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Parse parses the command line arguments
func Parse() Args {
	if len(os.Args) < 3 {
		logrus.Panic("Usage: iamrotator <action> <configpath>")
	}
	action := os.Args[1]
	configPath := os.Args[2]
	return Args{Action: action, ConfigPath: configPath}
}
