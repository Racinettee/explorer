package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	
	ui "github.com/VladimirMarkelov/clui"
	term "github.com/nsf/termbox-go"
)

type FileListBox struct {
	*ui.ListBox
	itemClickHandler func(ui.Event)
}

func CreateFileListBox(parent ui.Control, w, h, s int) *FileListBox {
	fileListBox := new(FileListBox)
	fileListBox.ListBox = ui.CreateListBox(nil, w, h, s)

	return fileListBox
}

func (fileListBox *FileListBox) ProcessEvent(event ui.Event) bool {
	fileListBox.ListBox.ProcessEvent(event)
	if event.Type == ui.EventClick {
		if event.Y <= fileListBox.ItemCount() {
			fileListBox.itemClickHandler(event)
			return true
		}
	}
	if event.Type == ui.EventKey && event.Key == term.KeyEnter {
		fileListBox.itemClickHandler(event)
	}
	return true
}

type FileExplorer struct {
	*ui.Window
	currentDir string
	flistBox *FileListBox
}

func CreateFileExplorer(x, y, w, h int, initialDir string) (*FileExplorer, error) {
	explorer := new(FileExplorer)
	explorer.Window = ui.AddWindow(x, y, w, h, "Explorer")
	explorer.flistBox = CreateFileListBox(explorer, w, h, 1)
	explorer.flistBox.itemClickHandler = fileItemClickedHandler(explorer)
	explorer.AddChild(explorer.flistBox)
	explorer.LoadDir(initialDir)

	return explorer, nil
}

func (explorer *FileExplorer) LoadDir(path string) error {
	dir, err := os.Open(path)
	
	if err != nil {
		return err
	}
	dirInfo, err := dir.Stat()
	
	if err != nil {
		return err
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("FileExplorer requires a directory name")
	}
	
	subFileInfos, err := dir.Readdir(0)
	dir.Close()

	if err != nil {
		return err
	}
	explorer.currentDir = path
	explorer.flistBox.Clear()
	if path != root() {
		explorer.flistBox.AddItem("..")
	}
	
	for _, fileInfo := range subFileInfos {
		if fileInfo.IsDir() {
			explorer.flistBox.AddItem(fileInfo.Name() + string(os.PathSeparator))
		} else {
			explorer.flistBox.AddItem(fileInfo.Name())
		}
	}
	return nil
}

func main() {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	CreateFileExplorer(0, 0, 30, 20, "./")
	
	ui.MainLoop()
}

func root() string {
   s := os.TempDir()
   return s[:strings.IndexRune(s, os.PathSeparator) + 1]
}

func fileItemClickedHandler(explorer *FileExplorer) func(ui.Event) {
	return func(ev ui.Event) {
		path := explorer.currentDir+string(os.PathSeparator)+
			explorer.flistBox.SelectedItemText()
		file, err := os.Open(path)

		if err != nil {
			log.Println(err)
			return
		}
		fileStats, err := file.Stat()
		if err != nil {
			log.Println(err)
			return
		}
		if fileStats.IsDir() {
			explorer.LoadDir(path)
		}
	}
}
