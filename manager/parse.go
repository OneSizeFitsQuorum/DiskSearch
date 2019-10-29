package manager

func (m *Manager) parseFileContent(filePath string) {
	return
}

func (m *Manager) parseFileName(name, filePath string) {
	// last := strings.LastIndex(name, ".")
	// if last != -1 {
	// 	name = name[:last]
	// }
	results := m.Cut(name)
	for _, result := range results {
		m.invertedIndex[result] = append(m.invertedIndex[result], filePath)
	}
}

func (m *Manager) Cut(text string) []string {
	hmm := seg.Cut(text, true)
	return hmm
}
