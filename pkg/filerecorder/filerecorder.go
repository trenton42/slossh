package filerecorder

import (
	"encoding/json"
	"fmt"
	"os"

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
	fp *os.File
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
	f.fp.Close()
}
