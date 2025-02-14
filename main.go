package main

import (
	"embed"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/ysmood/gson"
)

//go:embed all:frontend
var assets embed.FS

var viewRenderCount = 0

func main() {
	port, err := serve()
	if err != nil {
		log.Fatal(err)
	}

	debug := flag.Bool("debug", false, "Turns off Chromium headless mode for debugging")
	flag.Parse()

	if *debug {
		fmt.Println("Debug mode -- headless is disabled")
	}

	var workspaceFileName = "workspace.json"
	if flag.NArg() >= 1 {
		workspaceFileName = flag.Arg(0)
	}

	fmt.Printf("Loading workspace `%s`...", workspaceFileName)
	workspaceContent := loadWorkspace(workspaceFileName)

	l := launcher.New()
	defer l.Cleanup()

	if *debug {
		l.Headless(false).Devtools(true)
	}

	browser := rod.New().ControlURL(l.MustLaunch()).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(fmt.Sprintf("http://localhost:%s", port)).MustWaitStable()

	page.MustExpose("savePng", func(g gson.JSON) (interface{}, error) {
		saveDiagram("./", g.Get("viewKey").Str(), g.Get("png").Str())

		return nil, nil
	})

	page.MustExpose("log", func(g gson.JSON) (interface{}, error) {
		fmt.Printf("JS Log: %s", g.Str())

		return nil, nil
	})

	// Wait for structurizr to be ready
	page.Wait(&rod.EvalOptions{JS: "structurizr.scripting && structurizr.scripting.isDiagramRendered() === true"})

	// Load workspace into structurizr
	views := page.MustEval("(workspaceContent) => load(workspaceContent)", workspaceContent).Arr()

	fmt.Println(" DONE")

	for _, view := range views {
		page.MustEval("(viewKey) => render(viewKey)", view.Get("key"))
	}

	time.Sleep(time.Duration(2) * time.Second)

	fmt.Printf("Exported %d diagrams of %d expected\n", viewRenderCount, len(views))
}

func loadWorkspace(workspaceFileName string) string {

	workspaceContentBytes, err := os.ReadFile(workspaceFileName)
	if err != nil {
		fmt.Println("")
		log.Fatal(err)
	}
	workspaceContent := string(workspaceContentBytes)

	return workspaceContent
}

func saveDiagram(localDir string, diagramName string, dataURI string) string {
	b64data := dataURI[strings.IndexByte(dataURI, ',')+1:]
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64data))

	fileName := filepath.Join(localDir, fmt.Sprint(diagramName, ".png"))
	file, err := os.Create(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fmt.Printf("Exporting %s...", fileName)

	_, err = io.Copy(file, reader)

	viewRenderCount++

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

	return port, nil
}
