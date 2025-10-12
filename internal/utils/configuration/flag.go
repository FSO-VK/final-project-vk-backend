package configuration

import (
	"flag"
	"fmt"
)

// ReadConfigPathFlag reads path to file (-f, --file flags) from command line.
// It takes default path as an argument.
func ReadConfigPathFlag(defaultPath string) (string, error) {
	fFlagPtr := flag.String("f", defaultPath, "path to config file (shortcut)")
	fileFlagPtr := flag.String("file", defaultPath, "path to config file")

	flag.Parse()

	filePath := defaultPath
	// set priority to --flag.
	if *fileFlagPtr != defaultPath {
		filePath = *fileFlagPtr
	} else if *fFlagPtr != defaultPath {
		filePath = *fFlagPtr
	}

	if filePath == "" {
		return "", fmt.Errorf("must specify path if flag -f or --file are using")
	}

	return filePath, nil
}
