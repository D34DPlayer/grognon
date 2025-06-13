package backend

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"

	inertia "github.com/romsar/gonertia"
)

func initInertia(ssrHost string) *inertia.Inertia {
	viteHotFile := "./public/hot"
	rootViewFile := "./index.html"

	// check if laravel-vite-plugin is running in dev mode (it puts a "hot" file in the public folder)
	_, err := os.Stat(viteHotFile)
	if err == nil {
		i, err := inertia.NewFromFile(
			rootViewFile,
			inertia.WithSSR(ssrHost),
		)
		if err != nil {
			slog.Error(err.Error())
		}
		i.ShareTemplateFunc("vite", func(entry string) (string, error) {
			content, err := os.ReadFile(viteHotFile)
			if err != nil {
				return "", err
			}
			url := strings.TrimSpace(string(content))
			if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
				url = url[strings.Index(url, ":")+1:]
			} else {
				url = "//localhost:8080"
			}
			if entry != "" && !strings.HasPrefix(entry, "/") {
				entry = "/" + entry
			}
			return url + entry, nil
		})
		return i
	}

	// laravel-vite-plugin not running in dev mode, use build manifest file
	manifestPath := "./public/build/manifest.json"

	// check if the manifest file exists, if not, rename it
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		// move the manifest from ./public/build/.vite/manifest.json to ./public/build/manifest.json
		// so that the vite function can find it
		err := os.Rename("./public/build/.vite/manifest.json", "./public/build/manifest.json")
		if err != nil {
			return nil
		}
	}

	i, err := inertia.NewFromFile(
		rootViewFile,
		inertia.WithVersionFromFile(manifestPath),
		inertia.WithSSR(ssrHost),
	)
	if err != nil {
		slog.Error("Failed initializing Inertia", "error", err)
	}

	i.ShareTemplateFunc("vite", vite(manifestPath, "/build/"))

	return i
}

func vite(manifestPath, buildDir string) func(path string) (string, error) {
	f, err := os.Open(manifestPath)
	if err != nil {
		slog.Error("cannot open provided vite manifest file", slog.Any("error", err))
		return nil
	}
	defer f.Close()

	viteAssets := make(map[string]*struct {
		File   string `json:"file"`
		Source string `json:"src"`
	})
	err = json.NewDecoder(f).Decode(&viteAssets)
	// print content of viteAssets
	for k, v := range viteAssets {
		slog.Info("viteAssets", slog.Any(k, v.File))
	}

	if err != nil {
		slog.Error("cannot unmarshal vite manifest file to json", slog.Any("error", err))
		return nil
	}

	return func(p string) (string, error) {
		if val, ok := viteAssets[p]; ok {
			return path.Join("/", buildDir, val.File), nil
		}
		return "", fmt.Errorf("asset %q not found", p)
	}
}

func handleServerErr(w http.ResponseWriter, err error) {
	slog.Error("http error", slog.Any("error", err))
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("server error"))
}

func Render(w http.ResponseWriter, r *http.Request, i *inertia.Inertia, name string, props inertia.Props) {
	err := i.Render(w, r, name, props)
	if err != nil {
		handleServerErr(w, err)
	}
}
