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
var protoc_plugin_grpc = flag.Bool("protoc-plugin-grpc", true, "Enable gRPC")

func Generate() error {
    current_project, err := os.Getwd()
    if err != nil {
        return err
    }
    cfg, err := config.ParseProject(current_project)
    if err != nil {
        return err
    }

    import_paths := make([]string, 0, 128)
    import_map := make(map[string]string)
    err = walkProtoProjects(
        func(root string, import_path string, cfg *config.Config) error {
            local_proto_root := path.Join(root, cfg.Local.ProtoPrefix)
            import_paths = append(import_paths, local_proto_root)
            import_walker := func(proto_root, prefix string, protos []string) error {
                full_import := path.Join(import_path, cfg.Local.GoPrefix, prefix)
                for _, proto := range protos {
                    full_proto := path.Join(prefix, proto)
                    import_map[full_proto] = full_import
                }
                return nil
            }
            err := walkProtoPackages(local_proto_root, import_walker)
            if err != nil {
                return err
            }
            for _, imp := range cfg.Local.Import {
                err = walkProtoPackagesCustom(imp.Prefix, imp.Path, import_walker)
                if err != nil {
                    return err
                }
            }
            return nil
        })
    if err != nil {
        return err
    }

    protoc_gen_go_path, err := exec.LookPath(*protoc_gen_go)
    if err != nil {
        return err
    }
    proto_root := path.Join(current_project, cfg.Local.ProtoPrefix)
    go_root := path.Join(current_project, cfg.Local.GoPrefix)
    err = os.MkdirAll(go_root, 0777)
    if err != nil {
        return err
    }
    protoc_path_args := []string{
        "--proto_path=" + proto_root,
    }
    for _, import_path := range import_paths {
        protoc_path_args = append(protoc_path_args, "--proto_path="+import_path)
    }
    protoc_gen_go_param := "--plugin=protoc-gen-go=" + protoc_gen_go_path
    protoc_go_out_param := newProtocGoOutParam()
    if *protoc_plugin_grpc {
        protoc_go_out_param.addParam("plugins=grpc")
    }
    for key, value := range import_map {
        protoc_go_out_param.addParam("M" + key + "=" + value)
    }
    protoc_go_out_param.setPath(go_root)
    proto_generator := func(root string, prefix string, protos []string) error {
        protoc_args := append([]string{
            protoc_gen_go_param,
            protoc_go_out_param.String(),
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
    }
    err = walkProtoPackages(proto_root, proto_generator)
    if err != nil {
        return err
    }
    for _, imp := range cfg.Local.Import {
        err = walkProtoPackagesCustom(imp.Prefix, imp.Path, proto_generator)
        if err != nil {
            return err
        }
    }
    return nil
}
