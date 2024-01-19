package config

import (
	"flag"
	"os"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Verbose bool
	ErrorLogFile string
	OutDir string
	Website string
}

func GetConfig() Config {
	// Setup the config
	config := &Config{}

	// Check if the config file exists
	_, err := os.Stat("config.toml")
	if (err == nil) {
		tree, err := toml.LoadFile("config.toml")
		if err != nil {
			panic(err)
		}

		err = tree.Unmarshal(config)
		if err != nil {
			panic(err)
		}
	}

	// Flags
	verbosePtr := flag.Bool("v", false, "Enable verbose output")
	errorLogFilePtr := flag.String("elog", "", "Path to the error log file")
	websitePtr := flag.String("website", "", "Website to crawl")
	outDirPtr := flag.String("outdir", "archive", "Output directory")
	flag.Parse()

	// Map the flags onto the config.
	config.Verbose = *verbosePtr
	config.ErrorLogFile = *errorLogFilePtr
	config.Website = *websitePtr
	config.OutDir = *outDirPtr

	return config
}