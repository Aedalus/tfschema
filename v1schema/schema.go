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
	Type            string              `json:"type"`
	Description     string              `json: "description"`
	Attributes      []TFScalarAttribute `json:"attributes"`
	AttributeBlocks []TFBlockAttribute  `json:"attribute_blocks"`
}

type TFScalarAttribute struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Computed    bool   `json:"computed"`
	Description string `json:"description"`
}

type TFBlockAttribute struct {
	Name            string              `json:"name"`
	Required        bool                `json:"required"`
	Computed        bool                `json:"computed"`
	Description     bool                `json:"description"`
	Attributes      []TFScalarAttribute `json:"attributes"`
	AttributeBlocks []TFBlockAttribute  `json:"attribute_blocks"`
}

// Easily find and return a resource
func (s TFProviderSchema) GetResource(resourceType string) *TFResourceSchema {
	for _, resource := range s.Resources {
		if resource.Type == resourceType {
			return &resource
		}
	}
	return nil
}

// Sort function for deterministic results
type SortResourceByName []TFResourceSchema

func (a SortResourceByName) Len() int           { return len(a) }
func (a SortResourceByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortResourceByName) Less(i, j int) bool { return strings.Compare(a[i].Type, a[j].Type) < 1 }

func (s TFProviderSchema) SortResourcesByName() []TFResourceSchema {
	sort.Sort(SortResourceByName(s.Resources))
	return s.Resources
}

// Sort function for deterministic results
type SortAttrByName []TFScalarAttribute

func (a SortAttrByName) Len() int           { return len(a) }
func (a SortAttrByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortAttrByName) Less(i, j int) bool { return strings.Compare(a[i].Name, a[j].Name) < 1 }

func (r TFResourceSchema) SortAttributesByName() []TFScalarAttribute {
	sort.Sort(SortAttrByName(r.Attributes))
	return r.Attributes
}
