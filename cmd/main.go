package main

import (
	"DiskSearch/manager"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	path = kingpin.Flag("path", "the path of directory").Default(".").String()
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	kingpin.Parse()
	d, err := filepath.Abs(filepath.Dir(*path))
	if err != nil {
		logrus.WithError(err).Error("Wrong path format")
		return
	}
	_, err = os.Stat(d)
	if err != nil {
		logrus.WithError(err).Error("Path not exist")
		return
	}
	m := manager.NewManager(d)
	m.Repl()
}
