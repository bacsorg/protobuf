package generator

import (
    "os"
    "path"

    "github.com/bunsanorg/protoutils/generator/config"
)

type ProtocPathContext interface {
    AddImportPath(path string)
    RegisterProto(proto, goPackage string)
    ProtoRoot() string
    GoRoot() string
    MakeGoOutParam() string
    MakePathArgs() []string
}

type protocPathContext struct {
    currentProject string
    protoRoot      string
    goRoot         string
    importPaths    []string
    importMap      map[string]string
}

func newProtocPathContext(
    currentProject string, cfg *config.Config) (ProtocPathContext, error) {

    pathContext := &protocPathContext{
        currentProject: currentProject,
        protoRoot:      path.Join(currentProject, cfg.Local.ProtoPrefix),
        goRoot:         path.Join(currentProject, cfg.Local.GoPrefix),
        importPaths:    make([]string, 0, 128),
        importMap:      make(map[string]string),
    }
    err := os.MkdirAll(pathContext.goRoot, 0777)
    if err != nil {
        return nil, err
    }
    return pathContext, nil
}

func (c *protocPathContext) AddImportPath(path string) {
    c.importPaths = append(c.importPaths, path)
}

func (c *protocPathContext) RegisterProto(proto, goPackage string) {
    c.importMap[proto] = goPackage
}

func (c *protocPathContext) ProtoRoot() string {
    return c.protoRoot
}
func (c *protocPathContext) GoRoot() string {
    return c.goRoot
}

func (c *protocPathContext) MakeGoOutParam() string {
    protoc_go_out_param := newProtocGoOutParam()
    if *protoc_plugin_grpc {
        protoc_go_out_param.addParam("plugins=grpc")
    }
    for key, value := range c.importMap {
        protoc_go_out_param.addParam("M" + key + "=" + value)
    }
    protoc_go_out_param.setPath(c.goRoot)
    return protoc_go_out_param.String()
}

func (c *protocPathContext) MakePathArgs() []string {
    protoc_path_args := []string{
        "--proto_path=" + c.protoRoot,
    }
    for _, importPath := range c.importPaths {
        protoc_path_args = append(protoc_path_args, "--proto_path="+importPath)
    }
    return protoc_path_args
}
