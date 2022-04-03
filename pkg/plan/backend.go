package plan

type Backend interface {
	Import(p *Plan) error
	Export(p *Plan) error
}

var Backends = map[string]Backend{
	"s3://":   &BackendS3{},
	"file://": &BackendLocal{},
}
