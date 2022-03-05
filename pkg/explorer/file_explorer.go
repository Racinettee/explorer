package explorer

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	ui "github.com/VladimirMarkelov/clui"
)

type FileExplorer struct {
	*ui.Window
	currentDir string
	flistBox   *FileListBox
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

func fileItemClickedHandler(explorer *FileExplorer) func(ui.Event) {
	return func(ev ui.Event) {
		path, _ := filepath.Abs(explorer.currentDir + string(os.PathSeparator) +
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

func root() string {
	s := os.TempDir()
	return s[:strings.IndexRune(s, os.PathSeparator)+1]
}
