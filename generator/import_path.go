package generator

import (
    "fmt"
    "log"
    "os"
    "path"
    "path/filepath"
)

func getGoPath() ([]string, error) {
    gopathEnv := os.Getenv("GOPATH")
    if gopathEnv == "" {
        return nil, fmt.Errorf("GOPATH is not set")
    }
    return filepath.SplitList(gopathEnv), nil
}

func findProjectInPath(importPath string) (string, error) {
    gopaths, err := getGoPath()
    if err != nil {
        return "", err
    }
    for _, gopath := range gopaths {
        project_path := path.Join(gopath, "src", importPath)
        info, err := os.Stat(project_path)
        if err != nil {
            continue
        }
        if !info.IsDir() {
            continue
        }
        return project_path, nil
    }
    return "", fmt.Errorf("unable to find %q in GOPATH", importPath)
}

func findImportPathForProject(project_path string) (string, error) {
    gopaths, err := getGoPath()
    if err != nil {
        return "", err
    }
    for _, gopath := range gopaths {
        srcpath := path.Join(gopath, "src")
        if rel, err := filepath.Rel(srcpath, project_path); err == nil {
            log.Printf("Found import path %q for project %q", rel, project_path)
            return rel, nil
        }
    }
    return "", fmt.Errorf("unable to find import path for project %q", project_path)
}
