package metadata

import (
	"encoding/json"

	"github.com/tidwall/gjson"
)

type Metadata struct {
	Name        string
	Description string
	Image       string
	ExternalURL string
	Attributes  []map[string]string
	Extra       map[string]gjson.Result
}

func (md *Metadata) UnmarshalJSON(b []byte) error {
	fields := gjson.ParseBytes(b).Map()

	attrs := fields["attributes"].Array()

	md.Name = fields["name"].String()
	md.Description = fields["description"].String()
	md.Image = fields["image"].String()
	md.ExternalURL = fields["external_url"].String()
	md.Extra = make(map[string]gjson.Result)
	md.Attributes = make([]map[string]string, len(attrs))

	for i, attr := range attrs {
		md.Attributes[i] = make(map[string]string)

		for k, v := range attr.Map() {
			md.Attributes[i][k] = v.String()
		}
	}

	delete(fields, "name")
	delete(fields, "description")
	delete(fields, "image")
	delete(fields, "external_url")
	delete(fields, "attributes")

	md.Extra = make(map[string]gjson.Result, len(fields))
	for k, v := range fields {
		md.Extra[k] = v
	}

	return nil
}

func (md Metadata) MarshalJSON() ([]byte, error) {
	raw := map[string]interface{}{
		"name":         md.Name,
		"description":  md.Description,
		"image":        md.Image,
		"external_url": md.ExternalURL,
		"attributes":   md.Attributes,
	}

	for k, v := range md.Extra {
		raw[k] = json.RawMessage(v.Raw)
	}

	return json.Marshal(raw)
}
