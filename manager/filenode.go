package manager

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

type FileNode struct {
	Name      string
	Path      string
	IsDir     bool
	FileNodes map[string]*FileNode
	WordSet   *Set
}

func NewFileNode(name, path string, isDir bool) *FileNode {
	f := &FileNode{
		Name:      name,
		Path:      path,
		IsDir:     isDir,
		FileNodes: make(map[string]*FileNode),
		WordSet:   NewSet(),
	}
	return f
}

func (f *FileNode) Build(m *Manager) {
	f.addFileName(m)
	if Meet(RoleFileSuffix, f.Name) {
		f.addFileContent(m)
	}
}

func (f *FileNode) GetFatherNode() string {
	return f.Path[:len(f.Path)-len(f.Name)-1]
}

func (f *FileNode) addFileContent(m *Manager) {
	var cmd *exec.Cmd
	if Meet(RoleTikaSuffix, f.Path) {
		cmd = exec.Command("tika", "--text", f.Path)
	} else {
		cmd = exec.Command("cat", f.Path)
	}
	buf, err := cmd.Output()
	if err != nil {
		logrus.WithError(err).WithField("file", f.Path).Error("scanner file failed")
	} else {
		results := Cut(string(buf))
		for _, result := range results {
			f.WordSet.Add(result)
			if m.word2fileNode[result] == nil {
				m.word2fileNode[result] = NewSet()
			}
			m.word2fileNode[result].Add(f)
		}
	}
}

func (f *FileNode) addFileName(m *Manager) {
	results := Cut(f.Name)
	for _, result := range results {
		f.WordSet.Add(result)
		if m.word2fileNode[result] == nil {
			m.word2fileNode[result] = NewSet()
		}
		m.word2fileNode[result].Add(f)
	}
}
