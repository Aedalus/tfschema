package v1schema

import (
	"sort"
	"strings"
)

type TFProviderSchema struct {
	Provider    string             `json:"provider"`
	Resources   []TFResourceSchema `json:"resources"`
	DataSources []TFResourceSchema `json:"data_sources"`
}

type TFResourceSchema struct {
	Type        string              `json:"type"`
	Description string              `json: "description"`
	Attributes  []TFScalarAttribute `json:"attributes"`
}

type TFScalarAttribute struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Computed    bool   `json:"computed"`
	Description string `json:"description"`
}

func (s TFProviderSchema) GetResource(resourceType string) *TFResourceSchema {
	for _, resource := range s.Resources {
		if resource.Type == resourceType {
			return &resource
		}
	}

	return nil
}

type SortAttrByName []TFScalarAttribute

func (a SortAttrByName) Len() int           { return len(a) }
func (a SortAttrByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortAttrByName) Less(i, j int) bool { return strings.Compare(a[i].Name, a[j].Name) < 1 }

func (r TFResourceSchema) SortAttributesByName() []TFScalarAttribute {
	sort.Sort(SortAttrByName(r.Attributes))
	return r.Attributes
}
