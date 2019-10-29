package main

import (
	"DiskSearch/manager"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	path = kingpin.Flag("path", "the path of directory").Default(".").String()
)

func main() {
	kingpin.Parse()
	d, err := filepath.Abs(filepath.Dir(*path))
	if err != nil {
		logrus.WithError(err).Error("Open path failed")
		return
	}
	_ = manager.NewManager(d)
}
