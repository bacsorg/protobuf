package generator

import (
    "flag"
    "fmt"
    "github.com/bacsorg/protobuf/generator/config"
    "log"
    "os"
    "os/exec"
    "path"
)

var protoc = flag.String("protoc", "protoc", "protoc location")
var protoc_gen_go = flag.String(
    "protoc-gen-go", "protoc-gen-go", "protoc-gen-go location")

func Generate() error {
    current_project, err := os.Getwd()
    if err != nil {
        return err
    }
    cfg, err := config.ParseProject(current_project)
    if err != nil {
        return err
    }
    protoc_gen_go_path, err := exec.LookPath(*protoc_gen_go)
    if err != nil {
        return err
    }
    current_root := path.Join(current_project, cfg.Local.SourcePrefix)
    protoc_path_args := []string{
        "--proto_path=" + current_root,
    }
    protoc_gen_go_param := "--plugin=protoc-gen-go=" + protoc_gen_go_path
    protoc_go_out_param := "--go_out="
    // TODO add Mimport/file.proto=prefix/import/file.proto mappings
    protoc_go_out_param = protoc_go_out_param + "."
    return walkProtoPackages(current_root,
        func(root string, prefix string, protos []string) error {
            protoc_args := append([]string{
                protoc_gen_go_param,
                protoc_go_out_param,
            }, protoc_path_args...)
            for _, proto := range protos {
                protoc_args = append(protoc_args, path.Join(root, prefix, proto))
            }
            log.Printf("Generating protobufs for [%s] %s", root, prefix)
            cmd := exec.Command(*protoc, protoc_args...)
            output, err := cmd.CombinedOutput()
            if err != nil {
                return fmt.Errorf("protoc %v: %s", err, output)
            }
            return nil
        })
}
