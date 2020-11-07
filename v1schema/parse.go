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
				Type:        resourceType,
				Attributes:  []TFScalarAttribute{},
				Description: "",
			}
			jsonContent := content.(map[string]interface{})
			block := jsonContent["block"].(map[string]interface{})
			attributes := block["attributes"].(map[string]interface{})

			if val, ok := block["description"]; ok {
				r.Description = val.(string)
			}

			for attrName, attrRaw := range attributes {
				attr := attrRaw.(map[string]interface{})
				a := TFScalarAttribute{
					Name:        attrName,
					Type:        "",
					Required:    false,
					Computed:    false,
					Description: "",
				}
				if val, ok := attr["type"]; ok {
					switch v := val.(type) {
					case string:
						a.Type = v
					case []interface{}:
						// Get the outer type
						if outerType, ok := v[0].(string); ok {
							if innerType, ok := v[1].(string); ok {
								a.Type = fmt.Sprintf("%s(%s)", outerType, innerType)
							} else {
								a.Type = fmt.Sprintf("%s(%s", outerType, "UNKNOWN")
							}
						} else {
							return nil, fmt.Errorf("could not parse 'type' field of resource attribute %s:%s. type found: [](%+v)", r.Type, attrName, v)
						}
					default:
						a.Type = "UNKNOWN_TYPE"
					}
				}
				if val, ok := attr["optional"]; ok {
					a.Required = !val.(bool)
				}
				if val, ok := attr["required"]; ok {
					a.Required = val.(bool)
				}
				if val, ok := attr["computed"]; ok {
					a.Computed = val.(bool)
				}
				if val, ok := attr["description"]; ok {
					a.Description = val.(string)
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
