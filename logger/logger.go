package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

var ErrUnsupportedDestination = fmt.Errorf("unsupported destination for log output")

// New creates a new logger with the specified target and returns a pointer to it.
// The target parameter specifies the destination for the log output.
// Supported targets are:
//   - "console": Output will be written to the standard output.
//   - "file": Output will be written to the "pipeline.log" file.
//   - "none": Output will be discarded.
//
// If the target is not one of the supported values, an error will be returned.
func New(target string) (*log.Logger, error) {
	logger := log.Default()
	logger.SetPrefix("pipeline: ")

	if target == "console" {
		logger.SetOutput(os.Stdout)
	} else if target == "file" {
		f, err := os.Create("pipeline.log")
		if err != nil {
			return nil, err
		}
		logger.SetOutput(f)
	} else if target == "none" {
		log.SetOutput(io.Discard)
	} else {
		return nil, ErrUnsupportedDestination
	}

	return logger, nil
}
