package v1schema

type TFProviderSchema struct {
	Provider    string             `json:"provider"`
	Resources   []TFResourceSchema `json:"resources"`
	DataSources []TFResourceSchema `json:"data_sources"`
}

type TFResourceSchema struct {
	Type       string              `json:"type"`
	Attributes []TFScalarAttribute `json:"attributes"`
}

type TFScalarAttribute struct {
	Name     string   `json:"name"`
	Type     []string `json:"type"`
	Optional bool     `json:"optional"`
	Computed bool     `json:"computed"`
}

func (s TFProviderSchema) GetResource(resourceType string) *TFResourceSchema {
	for _, resource := range s.Resources {
		if resource.Type == resourceType {
			return &resource
		}
	}

	return nil
}
