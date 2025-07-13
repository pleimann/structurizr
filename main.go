package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/pleimann/structurizr-renderer/renderer"
)

func main() {
	debug := flag.Bool("debug", false, "Turns off Chromium headless mode for debugging")
	watch := flag.Bool("watch", false, "Watch for changes and re-export all views")
	outDir := flag.String("outdir", "./", "Specify an alternative directory to store rendered views")
	flag.Parse()

	var workspaceFileName = "workspace.json"
	if flag.NArg() >= 1 {
		workspaceFileName = flag.Arg(0)
	}

	r := renderer.New(workspaceFileName, *debug)

	r.ExportAllViews(outDir)

	if *watch {
		watchWorkspaceFile(workspaceFileName, r, outDir)

		// Block main goroutine forever.
		<-make(chan struct{})
	}
}

func watchWorkspaceFile(workspaceFileName string, r renderer.Renderer, outDir *string) {
	fmt.Println("\nWatching for changes...")

	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Write) && event.Name == workspaceFileName {
					fmt.Println("modified file:", event.Name)
					r.ExportAllViews(outDir)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(workspaceFileName)
	if err != nil {
		log.Fatal(err)
	}
}
