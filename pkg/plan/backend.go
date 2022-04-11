package plan

type Backend interface {
	Import(p *Plan) error
	Export(p *Plan) error
}

const LocalScheme = "fs://"

var Backends = map[string]Backend{
	"s3://":     &BackendS3{},
	LocalScheme: &BackendLocal{},
}
