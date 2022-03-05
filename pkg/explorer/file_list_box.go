package explorer

import (
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
