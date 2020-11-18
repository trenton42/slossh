package filerecorder

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/trenton42/slossh/pkg/recorders"
	"github.com/trenton42/slossh/pkg/session"
)

func init() {
	fr := FileRecorder{}
	fp, err := os.OpenFile("attempts.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	fr.fp = fp
	recorders.Register(&fr)
}

// FileRecorder records ssh attempts into a file as json objects
type FileRecorder struct {
	fp       *os.File
	filePath string
}

// Name returns this name of this recorder
func (f *FileRecorder) Name() string {
	return "file"
}

// Options returns the flags that can configure this recorder
func (f *FileRecorder) Options() *pflag.FlagSet {
	flags := pflag.NewFlagSet("file", pflag.ContinueOnError)
	flags.StringVar(&f.filePath, "file-path", "", "Path to json file to store results")
	return flags
}

// Open prepares this recorder to start working
func (f *FileRecorder) Open() error {
	if f.filePath == "" {
		f.filePath = "attempts.json"
	}
	fp, err := os.OpenFile(f.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	f.fp = fp
	return nil
}

// Record records sessions
func (f *FileRecorder) Record(s session.SlosshSession) {
	data, err := json.Marshal(s)
	if err != nil {
		return
	}
	f.fp.Write(data)
}

// Close closes this recorder and any associated resources
func (f *FileRecorder) Close() {
	if f.fp == nil {
		return
	}
	f.fp.Close()
}
