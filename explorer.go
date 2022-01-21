package main

import (
	"archive/zip"
	"path/filepath"
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
	explorer.SetTitle(fmt.Sprintf("Exploring: %v", dirInfo.Name()))
	return nil
}

func (explorer *FileExplorer) LoadZip(path string) error {
	zipf, err := zip.OpenReader(path)

	if err != nil {
		return err
	}
	defer zipf.Close()
	explorer.currentDir = path
	explorer.flistBox.Clear()

	explorer.flistBox.AddItem("..")
	for _, file := range zipf.File {
		explorer.flistBox.AddItem(file.Name)
	}

	return nil
}

func main() {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	initPath, err := resolveInitialPath()

	if err != nil {
		log.Println(err)
	}
	log.Printf("Starting explorer @ %v\n", initPath)
	CreateFileExplorer(0, 0, 30, 20, initPath)
	
	ui.MainLoop()
}

func resolveInitialPath() (result string, err error) {
	result, err = os.Getwd()
	if err != nil {
		return ".", err
	}
	switch {
		case len(os.Args) == 1:
			return
		case len(os.Args) > 1:
			if os.Args[1] == "." {
				return result, nil
			}
			if os.Args[1] == ".." {
				result, err = filepath.Abs(os.Args[1])
			}
	}
	return
}

func root() string {
   s := os.TempDir()
   return s[:strings.IndexRune(s, os.PathSeparator) + 1]
}

func fileItemClickedHandler(explorer *FileExplorer) func(ui.Event) {
	return func(ev ui.Event) {
		path, _ := filepath.Abs(explorer.currentDir+string(os.PathSeparator)+
			explorer.flistBox.SelectedItemText())
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
		log.Printf("Naving to: %v", path)
		if strings.HasSuffix(fileStats.Name(), ".zip") {
			explorer.LoadZip(path)
			return
		}
		if fileStats.IsDir() {
			explorer.LoadDir(path)
		}
	}
}
