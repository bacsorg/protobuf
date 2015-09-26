package generator

import (
    "fmt"
    "github.com/bacsorg/protobuf/generator/config"
    "os"
)

type walkProjectFunc func(root string, import_path string, cfg *config.Config) error

const (
    not_done = iota
    doing    = iota
    done     = iota
)

func walkProtoProjectsCustom(
    root string, import_path string, state map[string]int,
    walkFn walkProjectFunc) error {

    if s, ok := state[import_path]; ok {
        switch s {
        case not_done:
            state[import_path] = doing
        case doing:
            return fmt.Errorf("cyclic dependency detected at \"%s\"", import_path)
        case done:
            return nil
        }
    } else {
        state[import_path] = doing
    }
    defer func() { state[import_path] = done }()

    cfg, err := config.ParseProject(root)
    if err != nil {
        return err
    }
    for _, dep := range cfg.Dependencies {
        dep_path, err := findProjectInPath(dep)
        if err != nil {
            return err
        }
        err = walkProtoProjectsCustom(dep_path, dep, state, walkFn)
        if err != nil {
            return err
        }
    }
    return walkFn(root, import_path, cfg)
}

func walkProtoProjects(walkFn walkProjectFunc) error {
    root, err := os.Getwd()
    if err != nil {
        return err
    }
    import_path, err := findImportPathForProject(root)
    if err != nil {
        return err
    }
    state := make(map[string]int)
    return walkProtoProjectsCustom(root, import_path, state, walkFn)
}
