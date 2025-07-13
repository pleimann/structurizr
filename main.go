package main

import (
	"embed"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/ysmood/gson"
)

//go:embed all:frontend
var assets embed.FS

func main() {
	url, err := serve()
	if err != nil {
		log.Fatal(err)
	}

	debug := flag.Bool("debug", false, "Turns off Chromium headless mode for debugging")
	watch := flag.Bool("watch", false, "Watch for changes and re-export all views")
	outDir := flag.String("outdir", "./", "Specify an alternative directory to store rendered views")
	flag.Parse()

	var workspaceFileName = "workspace.json"
	if flag.NArg() >= 1 {
		workspaceFileName = flag.Arg(0)
	}

	exportAllViews(outDir, workspaceFileName, debug, url)

	if *watch {
		watchWorkspaceFile(workspaceFileName, outDir, debug, url)

		// Block main goroutine forever.
		<-make(chan struct{})
	}
}

func watchWorkspaceFile(workspaceFileName string, outDir *string, debug *bool, url string) {
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
					exportAllViews(outDir, workspaceFileName, debug, url)
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

func exportAllViews(outDir *string, workspaceFileName string, debug *bool, url string) {
	if _, err := os.Stat(*outDir); err != nil && errors.Is(err, fs.ErrNotExist) {
		os.MkdirAll(*outDir, os.ModePerm)
	}

	fmt.Printf("Loading workspace `%s`...", workspaceFileName)
	workspaceContentBytes, err := os.ReadFile(workspaceFileName)
	if err != nil {
		log.Fatal(err)
	}

	workspaceContent := string(workspaceContentBytes)

	l := launcher.New()
	defer l.Cleanup()

	l.Leakless(true)

	if *debug {
		fmt.Println("\nDebug mode -- headless is disabled")
		l.Headless(false).Devtools(true)
	}

	browser := rod.New().ControlURL(l.MustLaunch()).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(url).MustWaitStable()

	viewRenderCount := 0
	page.MustExpose("savePng", func(g gson.JSON) (interface{}, error) {
		saveDiagram(*outDir, g.Get("viewKey").Str(), g.Get("png").Str(), &viewRenderCount)

		return nil, nil
	})

	// Wait for structurizr to be ready
	page.Wait(&rod.EvalOptions{JS: "structurizr.scripting && structurizr.scripting.isDiagramRendered() === true"})

	// Load workspace into structurizr
	views := page.MustEval("(workspaceContent) => load(workspaceContent)", workspaceContent).Arr()

	fmt.Println(" DONE")

	for _, view := range views {
		page.MustEval("async (viewKey) => { await render(viewKey).then(png => savePng({ viewKey, png })); }", view.Get("key"))
	}

	fmt.Printf("Exported %d of %d diagram(s)\n", viewRenderCount, len(views))
}

func saveDiagram(outDir string, diagramName string, b64data string, viewRenderCount *int) string {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64data))

	fileName := filepath.Join(outDir, fmt.Sprint(diagramName, ".png"))
	file, err := os.Create(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fmt.Printf("Exporting `%s`...", fileName)

	_, err = io.Copy(file, reader)

	(*viewRenderCount)++

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(" DONE")

	return fileName
}

func serve() (string, error) {
	frontendAssets, err := fs.Sub(assets, "frontend")
	if err != nil {
		return "", err
	}

	http.Handle("/", http.FileServer(http.FS(frontendAssets)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	go func() {
		err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := fmt.Sprintf("http://localhost:%s", port)

	return url, nil
}
