package manager

import (
	"bufio"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

//新建文件： path + Created
//刚刚新建文件名字修改：oldpath + Created Renamed + '\n' + newpath + Renamed
//修改文件内容： path + Updated
//修改文件名字： oldpath + Renamed + '\n' + newpath + Renamed
//删除文件：path Removed || path Renamed || path Created Renamed
//文件移出监控区：path Renamed
//文件移入监控区：path Renamed
func (m *Manager) monitor() {
	cmd := exec.Command("fswatch", "-x", "--event=Created", "--event=Updated", "--event=Removed", "--event=Renamed", "--batch-marker", m.rootFileNode.Path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logrus.WithError(err).Error("[MONITOR] Error")
	} else {
		cmd.Start()
		monitorChan <- struct{}{}
		for {
			reader := bufio.NewReader(stdout)
			for {
				changes := make([][]string, 0)
				for true {
					line, err := reader.ReadString('\n')
					if err != nil || strings.TrimSpace(line) == "NoOp" {
						break
					}
					item := strings.Split(strings.TrimSpace(line), " ")
					if !strings.Contains(item[0], "/.") {
						changes = append(changes, item)
					}
				}
				for i := 0; i < len(changes); i++ {
					item := changes[i]
					path := item[0]
					if item[len(item)-1] == "Created" {
						m.CreateNode(path)
					} else if item[len(item)-1] == "Updated" {
						m.UpdateNode(path, path)
					} else if item[len(item)-1] == "Removed" || (item[len(item)-1] == "Renamed" && item[len(item)-2] == "Removed") {
						m.RemoveNode(path)
					} else if item[len(item)-1] == "Renamed" && len(changes) == 2 {
						item1 := changes[0]
						path1 := item1[0]
						item2 := changes[1]
						path2 := item2[0]
						if item1[len(item1)-1] == "Renamed" && item2[len(item2)-1] == "Renamed" {
							m.UpdateNode(path1, path2)
						}
						break
					} else if item[len(item)-1] == "Renamed" && len(changes) == 1 {
						_, err := os.Stat(path)
						if err != nil {
							m.RemoveNode(path)
						}
						m.CreateNode(path)
					}
				}
			}
		}
	}
}
