package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/trenton42/slossh"
	_ "github.com/trenton42/slossh/pkg/filerecorder"
)

func main() {
	port := pflag.IntP("port", "p", 2022, "Port to listen on")
	pflag.Parse()
	s, err := slossh.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s.Serve(*port)
}
