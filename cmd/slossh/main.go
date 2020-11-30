package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/trenton42/slossh"
	_ "github.com/trenton42/slossh/pkg/filerecorder"
	_ "github.com/trenton42/slossh/pkg/httprecorder"
	"github.com/trenton42/slossh/pkg/recorders"
)

var (
	version = "dev"
	date    = "unknown"
)

func main() {
	recs := recorders.RecorderMap()
	names := []string{}
	for key := range recs {
		names = append(names, key)
	}
	recorder := pflag.StringArrayP("recorder", "r", nil, "recorder to use (can be specified multiple times). Available recorders: "+strings.Join(names, ", "))
	port := pflag.IntP("port", "p", 22, "Port to listen on")
	showVersion := pflag.BoolP("version", "v", false, "version information")
	for _, rec := range recs {
		pflag.CommandLine.AddFlagSet(rec.Options())
	}
	pflag.Parse()
	if *showVersion {
		fmt.Printf("Slossh version %s, build date: %s\n", version, date)
		os.Exit(0)
	}
	useRecs := make([]recorders.Recorder, 0)
	for _, recName := range *recorder {
		if rec, ok := recs[recName]; !ok {
			fmt.Printf("Recorder %s is not available. Available recorders: %s\n", recName, strings.Join(names, ", "))
			os.Exit(1)
		} else {
			useRecs = append(useRecs, rec)
			if openRec, ok := rec.(recorders.Opener); ok {
				err := openRec.Open()
				if err != nil {
					fmt.Printf("Error in the %s recorder: %s\n", rec.Name(), err)
					os.Exit(1)
				}
			}
		}
	}
	s, err := slossh.New(useRecs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s.Serve(*port)
}
