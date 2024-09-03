package plan

import "github.com/invopop/jsonschema"

func GenSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             true,
		RequiredFromJSONSchemaTags: true,
	}

	schema := r.Reflect(&planBody{})
	schema.AdditionalProperties = jsonschema.TrueSchema // to allow anchors at the top level

	return schema
}
