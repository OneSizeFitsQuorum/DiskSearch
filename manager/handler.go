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

const (
	RoleFileSuffix = iota
	RoleTikaSuffix
)

var (
	fileSuffix  = []string{".txt", ".html", ".h", ".c", ".cc", ".cxx", ".cpp", ".hpp", ".java", ".go", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"}
	tikaSuffix  = []string{".html", ".xml", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"}
	scannerChan = make(chan struct{}, runtime.GOMAXPROCS(runtime.NumCPU()))
	monitorChan = make(chan struct{}, 0)
	seg         gse.Segmenter
)

type Manager struct {
	mutex             *sync.RWMutex
	wg                *sync.WaitGroup
	rootFileNode      *FileNode
	filePath2fileNode map[string]*FileNode
	word2fileNode     map[string]*Set
}

func NewManager(root string) *Manager {
	m := &Manager{
		mutex:             new(sync.RWMutex),
		wg:                new(sync.WaitGroup),
		rootFileNode:      NewFileNode(GetFileNameFromFilePath(root), root, true),
		filePath2fileNode: make(map[string]*FileNode),
		word2fileNode:     make(map[string]*Set),
	}
	go m.monitor()
	seg.LoadDict()
	m.rootFileNode.Build(m)
	m.wg.Add(1)
	m.scanner(m.rootFileNode, true)
	m.wg.Wait()
	<-monitorChan
	return m
}

func (m *Manager) scanner(node *FileNode, log bool) {
	scannerChan <- struct{}{}
	defer func() {
		<-scannerChan
		m.wg.Done()
	}()
	files, err := ioutil.ReadDir(node.Path)
	if err != nil {
		logrus.WithError(err).Error("Open scanner directory failed")
		return
	}
	if log {
		logrus.WithFields(logrus.Fields{"dir": node.Path, "childItemsNum": len(files)}).Debug("scanning dir...")
	}
	for _, file := range files {
		if file.Name()[0] != '.' {
			filePath := path.Join(node.Path, file.Name())
			child := NewFileNode(file.Name(), filePath, file.IsDir())
			if file.IsDir() {
				m.wg.Add(1)
				go m.scanner(child, true)
			}
			m.mutex.Lock()
			child.Build(m)
			m.filePath2fileNode[filePath] = child
			node.FileNodes[file.Name()] = child
			m.mutex.Unlock()
		}
	}
}

func (m *Manager) Repl() {
	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Please enter your query : ")
		input, err := inputReader.ReadString('\n')
		if err != nil {
			logrus.WithError(err).Error("Illegal input")
		}
		keys := Cut(strings.TrimSpace(input))
		results := NewSet()
		m.mutex.RLock()
		for _, key := range keys {
			result, ok := m.word2fileNode[key]
			if ok {
				for _, value := range result.Values() {
					results.Add((*(value.(*FileNode))).Path)
				}
			}
		}
		m.mutex.RUnlock()
		values := results.Values()
		if len(values) == 0 {
			fmt.Println("No Result.")
		} else {
			fmt.Println(len(values))
			for _, result := range values {
				fmt.Println(result)
			}
		}
	}
}
