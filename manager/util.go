package manager

import (
	"os"
	"strings"
)

func Cut(text string) []string {
	hmm := seg.Cut(text, true)
	return hmm
}

func Meet(mode int, name string) bool {
	var suffix *[]string
	if mode == RoleFileSuffix {
		suffix = &fileSuffix
	} else {
		suffix = &tikaSuffix
	}
	for _, suffix := range *suffix {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}

func GetFileNameFromFilePath(path string) string {
	f, _ := os.Stat(path)
	return f.Name()
}
