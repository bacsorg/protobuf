package generator

import (
    "fmt"
    "github.com/bacsorg/protobuf/generator/config"
    "os"
)

type walkProjectFunc func(root, importPath string, cfg *config.Config) error

const (
    notDone = iota
    doing   = iota
    done    = iota
)

func walkProtoProjectsCustom(
    root, importPath string, state map[string]int, walkFn walkProjectFunc) error {

    if s, ok := state[importPath]; ok {
        switch s {
        case notDone:
            state[importPath] = doing
        case doing:
            return fmt.Errorf("cyclic dependency detected at \"%s\"", importPath)
        case done:
            return nil
        }
    } else {
        state[importPath] = doing
    }
    defer func() { state[importPath] = done }()

    cfg, err := config.ParseProject(root)
    if err != nil {
        return err
    }
    for _, dep := range cfg.Dependencies {
        depPath, err := findProjectInPath(dep)
        if err != nil {
            return err
        }
        err = walkProtoProjectsCustom(depPath, dep, state, walkFn)
        if err != nil {
            return err
        }
    }
    return walkFn(root, importPath, cfg)
}

func walkProtoProjectsFromCurrent(walkFn walkProjectFunc) error {
    root, err := os.Getwd()
    if err != nil {
        return err
    }
    importPath, err := findImportPathForProject(root)
    if err != nil {
        return err
    }
    state := make(map[string]int)
    return walkProtoProjectsCustom(root, importPath, state, walkFn)
}
