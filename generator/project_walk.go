package generator

import (
    "flag"
    "fmt"
    "os"

    "github.com/bunsanorg/buildutils"
    "github.com/bunsanorg/protoutils/generator/config"
)

type walkProjectFunc func(root, importPath string, cfg *config.Config) error

const (
    notDone = iota
    doing   = iota
    done    = iota
)

var bootstrapProject = flag.String(
    "bootstrap-project", "github.com/bunsanorg/protoutils", "Project containing base types")
var dependOnBootstrapProject = flag.Bool(
    "depend-on-bootstrap-project", true, "Implicitly depend on bootstrap project")

func walkProtoProjectsCustom(
    root, importPath string, state map[string]int, walkFn walkProjectFunc) error {

    if s, ok := state[importPath]; ok {
        switch s {
        case notDone:
            state[importPath] = doing
        case doing:
            return fmt.Errorf("cyclic dependency detected at %q", importPath)
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
        depPath, err := buildutils.SrcDir(dep)
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
    importPath, err := buildutils.ImportPath(root)
    if err != nil {
        return err
    }
    state := make(map[string]int)
    if *dependOnBootstrapProject && importPath != *bootstrapProject {
        bootstrapRoot, err := buildutils.SrcDir(*bootstrapProject)
        if err != nil {
            return err
        }
        err = walkProtoProjectsCustom(bootstrapRoot, *bootstrapProject, state, walkFn)
        if err != nil {
            return err
        }
    }
    return walkProtoProjectsCustom(root, importPath, state, walkFn)
}
