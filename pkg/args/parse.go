package args

import "flag"

// Parse parses the command line arguments
func Parse() Args {
	configPath := flag.String("config", "config.yaml", "path to the configuration file")
	flag.Parse()
	return Args{ConfigPath: *configPath}
}
