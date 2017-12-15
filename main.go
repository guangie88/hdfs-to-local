package main

import (
	"log"
	"os"
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

// Function literal type to take a HDFS src path, local dst path, and HDFS client
type pathAct func(string, string, *hdfs.Client, os.FileInfo)

func walkDir(dirname string, src string, dst string, client *hdfs.Client, act pathAct) {
	srcDirPath := path.Join(src, dirname)
	dstDirPath := path.Join(dst, dirname)

	fileInfo, err := client.ReadDir(srcDirPath)
	exitOnErr("HDFS ReadDir", err)

	for _, f := range fileInfo {
		srcPath := path.Join(srcDirPath, f.Name())
		dstPath := path.Join(dstDirPath, f.Name())

		act(srcPath, dstPath, client, f)

		if f.IsDir() {
			walkDir(f.Name(), srcDirPath, dstDirPath, client, act)
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
	walkDir("", c.Src, c.Dst, client, func(srcPath string, dstPath string, client *hdfs.Client, f os.FileInfo) {
		if f.IsDir() {
			err := os.MkdirAll(dstPath, f.Mode())
			exitOnErr("os.MkdirAll", err)
		} else {
			log.Printf("%s -> %s", srcPath, dstPath)
			err := client.CopyToLocal(srcPath, dstPath)
			exitOnErr("hdfs.Client.CopyToLocal", err)
		}
	})
}
