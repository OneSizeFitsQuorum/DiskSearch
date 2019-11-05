package manager

import (
	"os/exec"
	"fmt"
	"io"
	"bufio"
	"strings"
)

func (m *Manager) monitor() {
	cmd := exec.Command("fswatch", "-x", "--event=Created", "--event=Updated", "--event=Removed", "--event=Renamed", m.rootPath)
	stdout, err := cmd.StdoutPipe()
    if err != nil {
		fmt.Println("[MONITOR] Error")
        fmt.Println(err)
	} else {
		cmd.Start()
		reader := bufio.NewReader(stdout)
		fmt.Println("[MONITOR] Started")
		//实时循环读取输出流中的一行内容
		for {
			line, err2 := reader.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				break
			}
			// fmt.Println(line)
			data := strings.TrimSpace(line)
			file_paths := strings.Split(data, "\n")
			for _, file_path := range file_paths {
				// if strings.HasSuffix(file_path, " Updated") {
				// 	fmt.Println("[Monitor] FileChanged: " + file_path[:-7])
				// }
				fmt.Println("[Monitor] Event: " + file_path)
			}
		}
		cmd.Wait()
	}
	fmt.Println("[Monitor] Exited")
}
