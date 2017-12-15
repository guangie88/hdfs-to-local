package main

import (
	"log"

	"github.com/colinmarc/hdfs"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	host = kingpin.Flag("host", "HDFS server hostname.").Required().String()
	root = kingpin.Flag("root", "Root path to data").Required().String()
)

func exitOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	kingpin.Parse()

	client, err := hdfs.New(*host)
	exitOnErr(err)

	rootFile, err := client.Stat(*root)
	exitOnErr(err)

	if !rootFile.IsDir() {
		log.Fatalf("Given root path '%s' is not a directory!", *root)
	}
}
