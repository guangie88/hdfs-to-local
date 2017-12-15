package main

import (
	"log"
	"path"

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

// Function literal type to take a HDFS src path and perform an action
type pathAct func(string)

func walkDir(dirname string, basePath string, client *hdfs.Client, act pathAct) {
	dirPath := path.Join(basePath, dirname)
	fileInfo, err := client.ReadDir(dirPath)

	exitOnErr("HDFS ReadDir", err)

	for _, f := range fileInfo {
		filePath := path.Join(dirPath, f.Name())
		act(filePath)

		if f.IsDir() {
			walkDir(f.Name(), dirPath, client, act)
		}
	}
}

func main() {
	kingpin.Parse()
	var c config

	_, err := toml.DecodeFile(*conf, &c)
	exitOnErr("TOML DecodeFile", err)

	client, err := hdfs.New(c.Host)
	exitOnErr("HDFS New", err)

	srcStat, err := client.Stat(c.Src)
	exitOnErr("HDFS Stat", err)

	if !srcStat.IsDir() {
		log.Fatalf("HDFS src: Given source path '%s' is not a directory!", c.Src)
	}

	// recursive portion
	walkDir(c.Src, "", client, func(srcPath string) {
		log.Printf("File info: %s", srcPath)
	})
}
