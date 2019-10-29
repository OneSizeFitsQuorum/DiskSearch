package manager

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"sync"

	"github.com/go-ego/gse"
	"github.com/sirupsen/logrus"
)

var (
	fileSuffix = []string{".txt", ".html", ".h", ".c", ".cc", ".cxx", ".cpp", ".hpp", ".java", ".go", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"}
	seg        gse.Segmenter
)

type Manager struct {
	mutex         *sync.RWMutex
	rootPath      string
	invertedIndex map[string][]string
}

func NewManager(root string) *Manager {
	m := &Manager{
		mutex:         new(sync.RWMutex),
		rootPath:      root,
		invertedIndex: make(map[string][]string),
	}
	seg.LoadDict()
	m.scanner(m.rootPath)
	fmt.Println(m.invertedIndex)
	return m
}

func (m *Manager) scanner(curPath string) {
	files, err := ioutil.ReadDir(curPath)
	if err != nil {
		logrus.WithError(err).Error("Open scanner directory failed")
		return
	}
	for _, file := range files {
		filePath := path.Join(curPath, file.Name())
		if file.IsDir() {
			m.scanner(filePath)
		} else {
			meet := false
			for _, suffix := range fileSuffix {
				if strings.HasSuffix(file.Name(), suffix) {
					m.parseFileName(file.Name(), filePath)
					m.parseFileContent(filePath)
					meet = true
					break
				}
			}
			if !meet {
				m.parseFileName(file.Name(), filePath)
			}
		}
	}
}

func (m *Manager) parseFileContent(filePath string) {
	return
}

func (m *Manager) parseFileName(name, filePath string) {
	last := strings.LastIndex(name, ".")
	if last != -1 {
		name = name[:last]
	}
	results := m.Cut(name)
	for _, result := range results {
		m.invertedIndex[result] = append(m.invertedIndex[result], filePath)
	}
}

func (m *Manager) Cut(text string) []string {
	hmm := seg.Cut(text, true)
	return hmm
}
