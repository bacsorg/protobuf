package generator

import (
    "fmt"
    "log"
    "os"
    "path"
    "path/filepath"
)

func getGoPath() ([]string, error) {
    gopath_env := os.Getenv("GOPATH")
    if gopath_env == "" {
        return nil, fmt.Errorf("GOPATH is not set")
    }
    return filepath.SplitList(gopath_env), nil
}

func findProjectInPath(import_path string) (string, error) {
    gopaths, err := getGoPath()
    if err != nil {
        return "", err
    }
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

func findImportPathForProject(project_path string) (string, error) {
    gopaths, err := getGoPath()
    if err != nil {
        return "", err
    }
    for _, gopath := range gopaths {
        srcpath := path.Join(gopath, "src")
        if rel, err := filepath.Rel(srcpath, project_path); err == nil {
            log.Printf("Found import path \"%s\" for project \"%s\"", rel, project_path)
            return rel, nil
        }
    }
    return "", fmt.Errorf("Unable to find import path for project \"%s\"", project_path)
}
