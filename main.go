package main

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/colinmarc/hdfs"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	conf = kingpin.Flag("conf", "TOML config file path.").Required().ExistingFile()
)

// config Main program config struct.
type config struct {
	// HDFS server hostname.
	Host string

	// HDFS source directory path to sync from.
	Src string

	// Local storage directory path to sync to.
	Dst string
}

func exitOnErr(desc string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", desc, err)
	}
}

func main() {
	kingpin.Parse()
	var c config

	_, err := toml.DecodeFile(*conf, &c)
	exitOnErr("TOML error", err)

	client, err := hdfs.New(c.Host)
	exitOnErr("HDFS new", err)

	srcStat, err := client.Stat(c.Src)
	exitOnErr("HDFS stat", err)

	if !srcStat.IsDir() {
		log.Fatalf("Given source path '%s' is not a directory!", c.Src)
	}
}
