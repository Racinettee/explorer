package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Racinettee/explorer/pkg/explorer"
	ui "github.com/VladimirMarkelov/clui"
)

func main() {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	initPath, err := resolveInitialPath()

	if err != nil {
		log.Println(err)
	}
	log.Printf("Starting explorer @ %v\n", initPath)
	explorer.CreateFileExplorer(0, 0, 30, 20, initPath)

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
