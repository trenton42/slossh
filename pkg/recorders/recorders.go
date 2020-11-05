package recorders

import (
	"github.com/trenton42/slossh/pkg/session"
)

// Recorder is the interface for recording ssh attempts
type Recorder interface {
	Record(session.SlosshSession)
	Close()
}

var recorders []Recorder

// Register adds a recorder to the registry
func Register(r Recorder) {
	recorders = append(recorders, r)
}

// Recorders returns a list of recorders
func Recorders() []Recorder {
	return recorders
}
