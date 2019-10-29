package manager

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/go-ego/gse"
	"github.com/sirupsen/logrus"
)

var (
	fileSuffix    = []string{".txt", ".html", ".h", ".c", ".cc", ".cxx", ".cpp", ".hpp", ".java", ".go", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"}
	seg           gse.Segmenter
	semaphoreChan = make(chan struct{}, runtime.GOMAXPROCS(runtime.NumCPU()))
)

type Manager struct {
	mutex         *sync.RWMutex
	wg            *sync.WaitGroup
	rootPath      string
	invertedIndex map[string][]string
}

func NewManager(root string) *Manager {
	m := &Manager{
		mutex:         new(sync.RWMutex),
		wg:            new(sync.WaitGroup),
		rootPath:      root,
		invertedIndex: make(map[string][]string),
	}
	seg.LoadDict()
	m.wg.Add(1)
	m.scanner(m.rootPath)
	m.wg.Wait()
	go m.monitor()
	return m
}

func (m *Manager) Repl() {
	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Please enter your query : ")
		input, err := inputReader.ReadString('\n')
		if err != nil {
			logrus.WithError(err).Error("Illegal input")
		}
		key := strings.TrimSpace(input[:len(input)-1])
		m.mutex.RLock()
		results, ok := m.invertedIndex[key]
		m.mutex.RUnlock()
		if !ok || len(results) == 0 {
			fmt.Println("No Result.")
			continue
		}
		for _, result := range results {
			fmt.Println(result)
		}
	}
}

func (m *Manager) scanner(curPath string) {
	semaphoreChan <- struct{}{}
	defer func() {
		<-semaphoreChan
		m.wg.Done()
	}()
	files, err := ioutil.ReadDir(curPath)
	if err != nil {
		logrus.WithError(err).Error("Open scanner directory failed")
		return
	}
	for _, file := range files {
		filePath := path.Join(curPath, file.Name())
		if file.IsDir() {
			m.wg.Add(1)
			go m.scanner(filePath)
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
