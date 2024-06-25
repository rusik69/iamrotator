package args

import (
	"fmt"
	"os"
)

// Parse parses the command line arguments
func Parse() Args {
	if len(os.Args) < 3 {
		fmt.Println("Usage: iamrotator <action> <configpath>")
		os.Exit(1)
	}
	action := os.Args[1]
	configPath := os.Args[2]
	return Args{Action: action, ConfigPath: configPath}
}
