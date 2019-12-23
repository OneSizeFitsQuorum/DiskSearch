package manager

import (
	"os"
	"path"
)

func (m *Manager) RemoveNode(path string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	node, ok := m.filePath2fileNode[path]
	if ok {
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
}

func (m *Manager) CreateNode(path string) {
	f, err := os.Stat(path)
	if err != nil {
		return
	}
	node := NewFileNode(GetFileNameFromFilePath(path), path, f.IsDir())
	m.mutex.Lock()
	m.filePath2fileNode[node.Path] = node
	node.Build(m)
	fatherNode, ok := m.filePath2fileNode[node.GetFatherNode()]
	if ok {
		fatherNode.FileNodes[node.Name] = node
	}
	m.mutex.Unlock()
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
			delete(m.filePath2fileNode, oldpath)
			m.filePath2fileNode[newpath] = node
			father, ok := m.filePath2fileNode[node.GetFatherNode()]
			if ok {
				delete(father.FileNodes, node.Name)
			}
			node.Path = newpath
			node.Name = f.Name()
			q := NewQueue()
			q.Add(node)
			for q.Length() != 0 {
				size := q.Length()
				for i := 0; i < size; i++ {
					root := q.Remove().(*FileNode)
					for _, value := range root.FileNodes {
						delete(m.filePath2fileNode, value.Path)
						value.Path = path.Join(root.Path, value.Name)
						m.filePath2fileNode[value.Path] = value
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
