package renderer

import (
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/ysmood/gson"
)

//go:embed all:frontend
var assets embed.FS

type Renderer struct {
	url               string
	browserPath       string
	debug             bool
	workspaceFileName string
}

func New(workspaceFileName string, browserPath *string, debug bool) Renderer {
	url, err := serve()

	if err != nil {
		log.Fatalf("Error serving static Structurizr UI site: %v", err)
	}

	path, _ := launcher.LookPath()
	if browserPath != nil && *browserPath != "" {
		path = *browserPath
		fmt.Printf("Browser path (from command line): `%s`\n", path)

	} else {
		fmt.Printf("Browser path (from lookup): `%s`\n", path)
	}

	r := Renderer{
		workspaceFileName: workspaceFileName,
		url:               url,
		browserPath:       path,
		debug:             debug,
	}

	return r
}

func (r *Renderer) ExportAllViews(outDir *string) {
	if _, err := os.Stat(*outDir); err != nil && errors.Is(err, fs.ErrNotExist) {
		os.MkdirAll(*outDir, os.ModePerm)
	}

	fmt.Printf("Loading workspace `%s`...", r.workspaceFileName)
	workspaceContentBytes, err := os.ReadFile(r.workspaceFileName)
	if err != nil {
		log.Fatal(err)
	}

	workspaceContent := string(workspaceContentBytes)

	l := launcher.New()
	defer l.Cleanup()

	// Find locally installed broweser
	l.NoSandbox(true).Leakless(true).Bin(r.browserPath)

	if r.debug {
		fmt.Println("\nDebug mode -- headless is disabled")
		l.Headless(false).Devtools(true)
	}

	browser := rod.New().ControlURL(l.MustLaunch()).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(r.url).MustWaitStable()

	viewRenderCount := 0
	page.MustExpose("savePng", func(g gson.JSON) (interface{}, error) {
		viewKey := g.Get("viewKey").Str()
		png := g.Get("png").Str()

		r.saveDiagram(*outDir, viewKey, png, &viewRenderCount)

		return nil, nil
	})

	// Wait for structurizr to be ready
	page.Wait(&rod.EvalOptions{JS: "structurizr.scripting && structurizr.scripting.isDiagramRendered() === true"})

	// Load workspace into structurizr
	views := page.MustEval("(workspaceFilename, workspaceContent) => load(workspaceFilename, workspaceContent)",
		r.workspaceFileName, workspaceContent).Arr()

	fmt.Println(" DONE")

	for _, view := range views {
		viewKey := view.Get("key").Str()
		viewType := view.Get("type").Str()

		if viewType == "Image" {
			contentUrl := view.Get("content").Str()
			contentType := view.Get("contentType").Str()

			if strings.Contains(contentType, "image/png") {
				resp, err := http.Get(contentUrl)
				if err != nil {
					fmt.Printf("failed to make HTTP request for image (%s): %s\n", contentUrl, err.Error())
				}
				defer resp.Body.Close() // Close the response body when the function exits

				// Check if the request was successful (status code 200 OK)
				if resp.StatusCode != http.StatusOK {
					fmt.Printf("failed to make HTTP request for image (%s): bad status code: %s\n", contentUrl, resp.Status)
				}

				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("Error reading response body: %s\n", err.Error())
				}

				png := base64.StdEncoding.EncodeToString(bodyBytes)

				r.saveDiagram(*outDir, viewKey, png, &viewRenderCount)
			}

		} else {
			page.MustEval("async (viewKey) => { await render(viewKey).then(png => savePng({ viewKey, png })); }", viewKey)
		}
	}

	fmt.Printf("Exported %d of %d diagram(s)\n", viewRenderCount, len(views))
}

func (r *Renderer) saveDiagram(outDir string, diagramName string, b64data string, viewRenderCount *int) string {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64data))

	fileName := filepath.Join(outDir, fmt.Sprint(diagramName, ".png"))
	file, err := os.Create(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fmt.Printf("Saving `%s`...", fileName)

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
