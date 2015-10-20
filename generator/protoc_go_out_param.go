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

func (o *protocGoOutParam) addParam(param string) {
	if !o.first {
		o.params += ","
	}
	o.first = false
	o.params += param
}

func (o *protocGoOutParam) setPath(path string) {
	o.path = path
}

func (o *protocGoOutParam) String() string {
	if o.first {
		return "--go_out=" + o.path
	} else {
		return "--go_out=" + o.params + ":" + o.path
	}
}
