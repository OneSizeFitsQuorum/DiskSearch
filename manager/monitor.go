package manager

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"errors"

	"github.com/sirupsen/logrus"
)

func (m *Manager) monitor() {
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
				fmt.Println("[Monitor Debug] Event: " + change)
				// There may be many status words, so need a loop structure
				newOperator := false
				deleteOperator := false
				notFinished := true
				for ; notFinished; {
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
				filepath := change
				if newOperator {
					f, err := os.Stat(filepath)
					if err != nil {
						if !(strings.HasSuffix(err.Error(), "no such file or directory") && deleteOperator) {  // Because there are two logs of `Renamed` event, both origin and new names will show
							logrus.WithError(err).WithField("file", filepath).Error("[MONITOR] Error")
						}
					} else {
						m.updateFile(f.Name(), filepath)
						fmt.Println("[Monitor] FileChanged: " + filepath)
					}
				}
				// Here don't use `else` because we need to get the origin name of renamed file
				if deleteOperator {
					_, err := os.Stat(filepath)
					if err != nil && strings.HasSuffix(err.Error(), "no such file or directory") {
						m.removeFile(filepath)
						fmt.Println("[Monitor] FileRemoved: " + filepath)
					} else if !(err == nil && newOperator) {
						logrus.WithError(errors.New("Error when delete")).WithField("file", filepath).Error("[MONITOR] Error")
					}
				}
			}
		}
		cmd.Wait()
	}
	logrus.Debug("[MONITOR] Exited")
}
