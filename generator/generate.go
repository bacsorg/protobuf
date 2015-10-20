package generator

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/bunsanorg/protoutils/generator/config"
)

var protoc = flag.String("protoc", "protoc", "protoc location")
var protocGenGo = flag.String(
	"protoc-gen-go", "protoc-gen-go", "protoc-gen-go location")
var protocPluginGrpc = flag.Bool("protoc-plugin-grpc", true, "Enable gRPC")

type generator struct {
	protocGenGoPath string
	pathContext     ProtocPathContext
	pathArgs        []string
}

func newGenerator(pathContext ProtocPathContext) (*generator, error) {
	protocGenGoPath, err := exec.LookPath(*protocGenGo)
	if err != nil {
		return nil, err
	}
	return &generator{
		protocGenGoPath: protocGenGoPath,
		pathContext:     pathContext,
	}, nil
}

func (g *generator) walkImportProjects(
	root, importPath string, cfg *config.Config) error {

	localProtoRoot := path.Join(root, cfg.Local.ProtoPrefix)
	g.pathContext.AddImportPath(localProtoRoot)
	importWalker := func(proto_root, prefix string, protos []string) error {
		fullImport := path.Join(importPath, cfg.Local.GoPrefix, prefix)
		for _, proto := range protos {
			fullProto := path.Join(prefix, proto)
			g.pathContext.RegisterProto(fullProto, fullImport)
		}
		return nil
	}
	err := walkProtoPackages(localProtoRoot, importWalker)
	if err != nil {
		return err
	}
	for _, imp := range cfg.Local.Import {
		err = walkProtoPackagesCustom(imp.Prefix, imp.Path, importWalker)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *generator) makeGoPluginParam() string {
	return "--plugin=protoc-gen-go=" + g.protocGenGoPath
}

func (g *generator) walkGenProtos(root, prefix string, protos []string) error {
	if g.pathArgs == nil {
		g.pathArgs = g.pathContext.MakePathArgs()
	}
	protocArgs := append([]string{
		g.makeGoPluginParam(),
		g.pathContext.MakeGoOutParam(),
	}, g.pathArgs...)
	for _, proto := range protos {
		protocArgs = append(protocArgs, path.Join(root, prefix, proto))
	}
	log.Printf("Generating protobufs for [%s] %s", root, prefix)
	cmd := exec.Command(*protoc, protocArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("protoc %v: %s", err, output)
	}
	return nil
}

func Generate() error {
	currentProject, err := os.Getwd()
	if err != nil {
		return err
	}
	cfg, err := config.ParseProject(currentProject)
	if err != nil {
		return err
	}
	pathContext, err := newProtocPathContext(currentProject, cfg)
	if err != nil {
		return err
	}

	generator, err := newGenerator(pathContext)
	if err != nil {
		return err
	}

	err = walkProtoProjectsFromCurrent(generator.walkImportProjects)
	if err != nil {
		return err
	}

	err = walkProtoPackages(pathContext.ProtoRoot(), generator.walkGenProtos)
	if err != nil {
		return err
	}
	for _, imp := range cfg.Local.Import {
		err = walkProtoPackagesCustom(imp.Prefix, imp.Path, generator.walkGenProtos)
		if err != nil {
			return err
		}
	}
	return nil
}
