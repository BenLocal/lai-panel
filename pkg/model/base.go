package model

type MetadataBase struct {
	Name       string            `json:"name,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type Metadata struct {
	MetadataBase `json:",inline"`
}

func (m *Metadata) ToMap() map[string]string {
	return m.Properties
}

func ToMetadataMap(ms []*Metadata, name string) (map[string]string, bool) {
	for _, m := range ms {
		if m.Name == name {
			return m.Properties, true
		}
	}
	return nil, false
}
