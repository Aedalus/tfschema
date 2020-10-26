package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"tfschema/v1schema"
)

func init() {
	RootCmd.AddCommand(resourcesCmd)

	// Bind flags
	resourcesCmd.Flags().StringVarP(&resourcesTFFile, "file", "f", "", "schema.json to read from")
}

// Flags
var resourcesTFFile string

var resourcesCmd = &cobra.Command{
	Use:   "resources [specific resource]",
	Short: "Prints information about resources",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var providerSchemas *[]v1schema.TFProviderSchema
		if resourcesTFFile != "" {
			jsonBytes, err := ioutil.ReadFile(resourcesTFFile)
			if err != nil {
				log.Fatalf("error reading schema file: %v", err)
			}

			providerSchemas, err = v1schema.ParseV1Schema(jsonBytes)
			if err != nil {
				log.Fatalf("error parsing schema file: %v", err)
			}
		}

		if len(args) == 0 {
			listAllResources(providerSchemas)
		} else {
			resourceType := args[0]
			getResourceDetails(providerSchemas, resourceType)
		}
	},
}

func getResourceDetails(providerSchemas *[]v1schema.TFProviderSchema, resourceType string) {
	var foundProvider string
	var foundResource *v1schema.TFResourceSchema

	for _, provider := range *providerSchemas {
		foundProvider = provider.Provider
		foundResource = provider.GetResource(resourceType)
		if foundResource != nil {
			break
		}
	}

	// If foundResource is not nil, we found a valid resource
	if foundResource != nil {
		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("%s - %s\n", foundProvider, foundResource.Type)
		fmt.Println(strings.Repeat("-", 80))
		fmt.Println("Attributes:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		for _, attr := range foundResource.Attributes {
			fmt.Fprintf(w, "\t%s\t%s\n", attr.Name, strings.Join(attr.Type, ","))
		}
		w.Flush()
	} else {
		log.Fatalf("Resource %s not found", resourceType)
	}
}
func listAllResources(schemas *[]v1schema.TFProviderSchema) {
	// They specified a schema file

	for _, schema := range *schemas {
		providerShort := schema.Provider[strings.LastIndex(schema.Provider, "/")+1:]
		fmt.Printf("%s (%s)\n", providerShort, schema.Provider)
		for _, resource := range schema.Resources {
			fmt.Printf("  %s\n", resource.Type)
		}
	}

}
