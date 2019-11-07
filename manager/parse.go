package manager

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

func (m *Manager) addFileContent(filePath string) {
	var cmd *exec.Cmd
	if m.Meet(tikaSuffix, filePath) {
		cmd = exec.Command("tika", "--text", filePath)
	} else {
		cmd = exec.Command("cat", filePath)
	}
	buf, err := cmd.Output()
	if err != nil {
		logrus.WithError(err).WithField("file", filePath).Error("scanner file failed")
	} else {
		results := m.Cut(string(buf))
		m.mutex.Lock()
		defer m.mutex.Unlock()
		for _, result := range results {
			if m.file2word[filePath] == nil {
				m.file2word[filePath] = NewSet()
			}
			m.file2word[filePath].Add(result)
			if m.word2file[result] == nil {
				m.word2file[result] = NewSet()
			}
			m.word2file[result].Add(filePath)
		}
	}
}

func (m *Manager) addFileName(name, filePath string) {
	results := m.Cut(name)
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, result := range results {
		if m.file2word[filePath] == nil {
			m.file2word[filePath] = NewSet()
		}
		m.file2word[filePath].Add(result)
		if m.word2file[result] == nil {
			m.word2file[result] = NewSet()
		}
		m.word2file[result].Add(filePath)
	}
}

func (m *Manager) removeFile(filePath string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	set, ok := m.file2word[filePath]
	if ok {
		for _, value := range set.Values() {
			m.word2file[value].Remove(filePath)
		}
	}
	delete(m.file2word, filePath)
}

func (m *Manager) updateFile(name, filePath string) {
	m.removeFile(filePath)
	m.addFileName(name, filePath)
	if m.Meet(fileSuffix, name) {
		m.addFileContent(filePath)
	}
}
