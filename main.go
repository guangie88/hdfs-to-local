package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/colinmarc/hdfs"
	"github.com/evalphobia/logrus_fluent"
	joonix "github.com/joonix/log"
	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	conf = kingpin.Flag("conf", "TOML config file path.").Required().ExistingFile()
)

type fluentd struct {
	Host string
	Port int
	Tag  string
}

// config Main program config struct.
type config struct {
	// HDFS server hostname.
	Host string

	// HDFS source directory path to sync from.
	Src string

	// Local storage directory path to sync to.
	Dst string

	// Regex filters accepting the source files to copy from.
	Filters []string

	//
	UseFluentd bool

	//
	Fluentd fluentd
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

func isMatchingFilters(srcPath string, filters []*regexp.Regexp) bool {
	for _, r := range filters {
		if r.MatchString(srcPath) {
			return true
		}
	}

	return false
}

func isSimilarFile(srcPath string, dstPath string, client *hdfs.Client) (bool, error) {
	srcData, err := client.ReadFile(srcPath)

	if err != nil {
		return false, err
	}

	// allow for dst file not to exist
	dstData, err := ioutil.ReadFile(dstPath)

	if err != nil {
		return false, nil
	}

	return md5.Sum(srcData) == md5.Sum(dstData), nil
}

func exitOnErr(desc string, err error) {
	if err != nil {
		log.WithFields(log.Fields{
			"error": fmt.Sprintf("%v", err),
		}).Error(desc)

		os.Exit(1)
	}
}

func initLog(c config) error {
	if c.UseFluentd {
		hook, err := logrus_fluent.NewWithConfig(logrus_fluent.Config{
			Host:          c.Fluentd.Host,
			Port:          c.Fluentd.Port,
			MarshalAsJSON: true,
		})

		if err != nil {
			return err
		}

		hook.SetLevels([]log.Level{
			log.InfoLevel,
			log.ErrorLevel,
		})

		hook.SetTag(c.Fluentd.Tag)

		log.SetFormatter(&joonix.FluentdFormatter{})
		log.AddHook(hook)
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}

	return nil
}

func main() {
	kingpin.Parse()

	var c config
	_, err := toml.DecodeFile(*conf, &c)
	exitOnErr("TOML DecodeFile", err)

	err = initLog(c)
	exitOnErr("initLog", err)

	filters := make([]*regexp.Regexp, len(c.Filters))

	for i, f := range c.Filters {
		filter, err := regexp.Compile(f)
		exitOnErr("regexp.Compile", err)
		filters[i] = filter
	}

	client, err := hdfs.New(c.Host)
	exitOnErr("HDFS New", err)

	srcStat, err := client.Stat(c.Src)
	exitOnErr("HDFS Stat", err)

	if !srcStat.IsDir() {
		log.Fatalf("HDFS src: Given source path '%s' is not a directory!", c.Src)
	}

	// recursive portion
	walkDir("", c.Src, c.Dst, client, func(srcPath string, dstPath string, client *hdfs.Client, f os.FileInfo) {
		if !f.IsDir() && isMatchingFilters(srcPath, filters) {
			err := os.MkdirAll(path.Dir(dstPath), 0755)
			exitOnErr("os.MkdirAll", err)

			isSimilar, err := isSimilarFile(srcPath, dstPath, client)
			exitOnErr("isSimilarFile", err)

			if isSimilar {
				log.Printf("SIMILAR %s AND %s, not copying...", srcPath, dstPath)
			} else {
				log.Printf("COPY %s -> %s", srcPath, dstPath)
				err = client.CopyToLocal(srcPath, dstPath)
				exitOnErr("hdfs.Client.CopyToLocal", err)
			}
		}
	})
}
