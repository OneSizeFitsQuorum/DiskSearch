package manager

import "strings"

func (m *Manager) Cut(text string) []string {
	hmm := seg.Cut(text, true)
	return hmm
}

func (m *Manager) Meet(suffix []string, name string) bool {
	for _, suffix := range suffix {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}
