package manager

import (
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
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
	if meet {
		cmd = exec.Command("tika", "--text", filePath)
	} else {
		cmd = exec.Command("cat", filePath)
	}
	buf, err := cmd.Output()
	if err == nil {
		results := m.Cut(string(buf))
		logrus.WithFields(logrus.Fields{"file": filePath, "fileSize": len(string(buf))}).Debug("scanning file...")
		m.mutex.Lock()
		for _, result := range results {
			if m.invertedIndex[result] == nil {
				m.invertedIndex[result] = NewSet()
			}
			m.invertedIndex[result].Add(filePath)
		}
		m.mutex.Unlock()
	}
	return
}

func (m *Manager) parseFileName(name, filePath string) {
	results := m.Cut(name)
	m.mutex.Lock()
	for _, result := range results {
		if m.invertedIndex[result] == nil {
			m.invertedIndex[result] = NewSet()
		}
		m.invertedIndex[result].Add(filePath)
	}
	m.mutex.Unlock()
}

func (m *Manager) Cut(text string) []string {
	hmm := seg.Cut(text, true)
	return hmm
}
