package manager

import (
	"os"
	"path"
)

func (m *Manager) removeFile(node *FileNode) {
	values := node.WordSet.Values()
	for _, value := range values {
		m.word2fileNode[value.(string)].Remove(node)
	}
	delete(m.filePath2fileNode, node.Path)
	fatherNode, ok := m.filePath2fileNode[node.GetFatherNode()]
	if ok {
		delete(fatherNode.FileNodes, node.Name)
	}
}

func (m *Manager) RemoveNode(path string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	node, ok := m.filePath2fileNode[path]
	if ok {
		if node.IsDir {
			q := NewQueue()
			q.Add(node)
			for q.Length() != 0 {
				size := q.Length()
				for i := 0; i < size; i++ {
					root := q.Remove().(*FileNode)
					for _, value := range root.FileNodes {
						q.Add(value)
					}
					m.removeFile(root)
				}
			}
		} else {
			m.removeFile(node)
		}
	}
}

func (m *Manager) CreateNode(path string) {
	f, err := os.Stat(path)
	if err != nil {
		return
	}
	node := NewFileNode(GetFileNameFromFilePath(path), path, f.IsDir())
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.filePath2fileNode[node.Path] = node
	node.Build(m)
	fatherNode, ok := m.filePath2fileNode[node.GetFatherNode()]
	if ok {
		fatherNode.FileNodes[node.Name] = node
	}
	if f.IsDir() {
		m.wg.Add(1)
		m.scanner(node, false)
		m.wg.Wait()
	}
}

func (m *Manager) UpdateNode(oldpath, newpath string) {
	f, err := os.Stat(newpath)
	if err != nil {
		return
	}
	if f.IsDir() {
		m.mutex.Lock()
		defer m.mutex.Unlock()
		node, ok := m.filePath2fileNode[oldpath]
		if ok {
			node.Path = newpath
			node.Name = f.Name()
			q := NewQueue()
			q.Add(node)
			for q.Length() != 0 {
				size := q.Length()
				for i := 0; i < size; i++ {
					root := q.Remove().(*FileNode)
					for _, value := range root.FileNodes {
						value.Path = path.Join(root.Path, value.Name)
						q.Add(value)
					}
				}
			}
		}
	} else {
		m.RemoveNode(oldpath)
		m.CreateNode(newpath)
	}

}
