package manager
import (
	"os/exec"
	"fmt"
	"strings"
)

var tikaSuffix = []string{".html", "xml", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"}

func (m *Manager) parseFileContent(filePath string) {
	meet := false
	for _, suffix := range tikaSuffix {
		if strings.HasSuffix(filePath, suffix) {
			meet = true
			break
		}
	}
	var cmd *exec.Cmd
	if (meet) {
		cmd = exec.Command("tika", "--text", filePath)
	} else {
		cmd = exec.Command("cat", filePath)
	}
	buf, err := cmd.Output()
	
	// fmt.Printf("File: " + filePath)
	// fmt.Printf("output: %s\n", buf)
	// fmt.Printf("err: %v\n", err)
	if (err == nil) {
		results := m.Cut(string(buf))
		fmt.Printf("Length of results of file %s is %d\n", filePath, len(results))
		m.mutex.Lock()
		for _, result := range results {
			m.invertedIndex[result] = append(m.invertedIndex[result], filePath)
		}
		m.mutex.Unlock()
	}
	return
}

func (m *Manager) parseFileName(name, filePath string) {
	// last := strings.LastIndex(name, ".")
	// if last != -1 {
	// 	name = name[:last]
	// }
	results := m.Cut(name)
	m.mutex.Lock()
	for _, result := range results {
		m.invertedIndex[result] = append(m.invertedIndex[result], filePath)
	}
	m.mutex.Unlock()
}

func (m *Manager) Cut(text string) []string {
	hmm := seg.Cut(text, true)
	return hmm
}
