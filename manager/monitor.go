package manager

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
				if strings.HasSuffix(change, " Updated") || strings.HasSuffix(change, " Created") {
					filepath := change[:len(change)-8]
					f, err := os.Stat(filepath)
					if err != nil {
						logrus.WithError(err).WithField("file", filepath).Error("[MONITOR] Error")
					}
					m.updateFile(f.Name(), filepath)
					fmt.Println("[Monitor] FileChanged: " + filepath)
				}
				fmt.Println("[Monitor] Event: " + change)
			}
		}
		cmd.Wait()
	}
	logrus.Debug("[MONITOR] Exited")
}
