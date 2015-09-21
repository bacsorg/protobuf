package generator

import (
    "fmt"
    "os"
    "path"
    "path/filepath"
)

func findProjectInPath(import_path string) (string, error) {
    gopath_env := os.Getenv("GOPATH")
    if gopath_env == "" {
        return "", fmt.Errorf("GOPATH is not set")
    }
    gopaths := filepath.SplitList(gopath_env)
    for _, gopath := range gopaths {
        project_path := path.Join(gopath, "src", import_path)
        info, err := os.Stat(project_path)
        if err != nil {
            continue
        }
        if !info.IsDir() {
            continue
        }
        return project_path, nil
    }
    return "", fmt.Errorf("Unable to find \"%s\" in GOPATH", import_path)
}
