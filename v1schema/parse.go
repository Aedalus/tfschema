package v1schema

import (
	json "encoding/json"
	"fmt"
)

func ParseV1Schema(jsonFile []byte) (*[]TFProviderSchema, error) {
	var fullSchema map[string]interface{}
	if err := json.Unmarshal(jsonFile, &fullSchema); err != nil {
		return nil, err
	}

	providerSchemas := fullSchema["provider_schemas"].(map[string]interface{})

	tfschema := TFProviderSchema{
		Provider:    "",
		Resources:   []TFResourceSchema{},
		DataSources: []TFResourceSchema{},
	}

	for provider, providerSchemaRaw := range providerSchemas {
		tfschema.Provider = provider
		providerSchema := providerSchemaRaw.(map[string]interface{})

		// Compute resource schemas
		resourceSchemas := providerSchema["resource_schemas"].(map[string]interface{})
		for resourceType, content := range resourceSchemas {
			r := TFResourceSchema{
				Type:       resourceType,
				Attributes: []TFScalarAttribute{},
			}
			jsonContent := content.(map[string]interface{})
			block := jsonContent["block"].(map[string]interface{})
			attributes := block["attributes"].(map[string]interface{})

			for attrName, attrRaw := range attributes {
				attr := attrRaw.(map[string]interface{})
				a := TFScalarAttribute{
					Name:     attrName,
					Type:     []string{},
					Optional: false,
					Computed: false,
				}
				if val, ok := attr["type"]; ok {
					switch v := val.(type) {
					case string:
						a.Type = []string{v}
					case []interface{}:
						for _, str := range v {
							a.Type = append(a.Type, str.(string))
						}
					default:
						a.Type = []string{"UNKNOWN_TYPE"}
						fmt.Printf("Unknown type: %T\n", val)
					}
				}
				if val, ok := attr["optional"]; ok {
					a.Optional = val.(bool)
				}
				if val, ok := attr["required"]; ok {
					a.Optional = !val.(bool)
				}
				if val, ok := attr["computed"]; ok {
					a.Computed = !val.(bool)
				}

				r.Attributes = append(r.Attributes, a)
			}

			// Add resource to the schema
			tfschema.Resources = append(tfschema.Resources, r)
		}

		// dataSchemas := providerSchema["data_source_schemas"].(map[string]interface{})

	}

	return &[]TFProviderSchema{
		tfschema,
	}, nil
}
