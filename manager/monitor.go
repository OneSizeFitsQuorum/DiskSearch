package manager

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func (m *Manager) monitor() {
	for {
		cmd := exec.Command("fswatch", "-x", "--event=Created", "--event=Updated", "--event=Removed", "--event=Renamed", m.rootPath)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			logrus.WithError(err).Error("[MONITOR] Error")
		} else {
			cmd.Start()
			reader := bufio.NewReader(stdout)
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					break
				}
				changes := strings.Split(strings.TrimSpace(line), "\n")
				for _, change := range changes {
					//fmt.Println("[Monitor Debug] Event: " + change)
					// There may be many status words, so need a loop structure
					newOperator := false
					deleteOperator := false
					notFinished := true
					for notFinished {
						if strings.HasSuffix(change, " Renamed") {
							newOperator = true
							deleteOperator = true
							change = change[:len(change)-8]
						} else if strings.HasSuffix(change, " Updated") || strings.HasSuffix(change, " Created") {
							newOperator = true
							change = change[:len(change)-8]
						} else if strings.HasSuffix(change, " Removed") {
							deleteOperator = true
							change = change[:len(change)-8]
						} else {
							notFinished = false
						}
					}
					filePath := change
					if newOperator {
						f, err := os.Stat(filePath)
						if err != nil {
							if !(strings.HasSuffix(err.Error(), "no such file or directory") && deleteOperator) { // Because there are two logs of `Renamed` event, both origin and new names will show
								logrus.WithError(err).WithField("file", filePath).Error("[MONITOR] Error")
							}
						} else {
							m.updateFile(f.Name(), filePath)
							//fmt.Println("[Monitor] FileChanged: " + filePath)
						}
					}
					// Here don't use `else` because we need to get the origin name of renamed file
					if deleteOperator {
						_, err := os.Stat(filePath)
						if err != nil && strings.HasSuffix(err.Error(), "no such file or directory") {
							m.removeFile(filePath)
							//fmt.Println("[Monitor] FileRemoved: " + filePath)
						} else if !(err == nil && newOperator) {
							logrus.WithError(errors.New("Error when delete")).WithField("file", filePath).Error("[MONITOR] Error")
						}
					}
				}
			}
			cmd.Wait()
		}
	}
}
