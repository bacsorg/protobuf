package generator

import (
    "io/ioutil"
    "path"
    "path/filepath"
    "strings"
)

type walkProtoFunc func(root string, prefix string, protos []string) error

func walkProtoPackagesCustom(root string, prefix string, walkFn walkProtoFunc) error {
    protos := make([]string, 0, 64)
    files, err := ioutil.ReadDir(path.Join(root, prefix))
    if err != nil {
        return err
    }
    for _, f := range files {
        if f.IsDir() {
            nprefix := path.Join(prefix, f.Name())
            err = walkProtoPackagesCustom(root, nprefix, walkFn)
            if err != nil {
                return err
            }
        } else {
            if strings.HasSuffix(f.Name(), ".proto") {
                protos = append(protos, f.Name())
            }
        }
    }
    if len(protos) > 0 {
        return walkFn(root, prefix, protos)
    }
    return nil
}

func walkProtoPackages(root string, walkFn walkProtoFunc) error {
    abs_root, err := filepath.Abs(root)
    if err != nil {
        return err
    }
    return walkProtoPackagesCustom(abs_root, "", walkFn)
}
