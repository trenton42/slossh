package recorders

import (
	"github.com/spf13/pflag"
	"github.com/trenton42/slossh/pkg/session"
)

// Recorder is the interface for recording ssh attempts
type Recorder interface {
	Record(session.SlosshSession)
	Name() string
	Close()
	Options() *pflag.FlagSet
}

// Opener is an interface for recorders that need to do work before starting
type Opener interface {
	Open() error
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

// RecorderMap returns a map of recorders, keyed by their name
func RecorderMap() map[string]Recorder {
	out := make(map[string]Recorder)
	for _, rec := range recorders {
		out[rec.Name()] = rec
	}
	return out
}
