package generator

type protocGoOutParam struct {
    first  bool
    params string
    path   string
}

func newProtocGoOutParam() *protocGoOutParam {
    return &protocGoOutParam{
        first: true,
    }
}

func (go_out *protocGoOutParam) addParam(param string) {
    if !go_out.first {
        go_out.params += ","
    }
    go_out.first = false
    go_out.params += param
}

func (go_out *protocGoOutParam) setPath(path string) {
    go_out.path = path
}

func (go_out *protocGoOutParam) String() string {
    if go_out.first {
        return "--go_out=" + go_out.path
    } else {
        return "--go_out=" + go_out.params + ":" + go_out.path
    }
}
