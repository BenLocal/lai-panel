package model

type MetadataBase struct {
	Name       string            `json:"name,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type Metadata struct {
	MetadataBase `json:",inline"`
}
