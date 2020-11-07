package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
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
	resourcesCmd.Flags().BoolVarP(&resourcesVerbose, "verbose", "v", false, "prints additional information like computed attributes")
}

// Flags
var resourcesTFFile string
var resourcesVerbose bool

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
		fmt.Println(color.HiCyanString(strings.Repeat("-", 80)))
		fmt.Printf(color.HiYellowString(" %s - %s\n"), foundProvider, foundResource.Type)
		fmt.Println(color.HiCyanString(strings.Repeat("-", 80)))
		if foundResource.Description != "" {
			wrappedDesc := wrapperLong(foundResource.Description)
			fmt.Printf("%s\n", color.HiBlackString(wrappedDesc))
			fmt.Println(strings.Repeat("-", 80))
		}
		fmt.Printf("--- %s -----------------------------------------------------------------\n", color.HiYellowString("Attributes"))
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

		sortedAttr := foundResource.SortAttributesByName()
		// Write the non-computed first
		for _, attr := range sortedAttr {
			if attr.Computed == false {
				formatResourceAttr(w, attr)
			}
		}

		if resourcesVerbose {
			ln := fmt.Sprintf("--- %s -----------------------------------------------------------------\n", color.HiYellowString("Computed"))
			w.Write([]byte(ln))
			// Then the computed
			for _, attr := range sortedAttr {
				if attr.Computed == true {
					formatResourceAttr(w, attr)
				}
			}
		}

		w.Flush()
	} else {
		log.Fatalf("Resource %s not found", resourceType)
	}
}

func formatResourceAttr(w *tabwriter.Writer, attr v1schema.TFScalarAttribute) {
	// Optional Text

	rqText := " "
	if attr.Required {
		rqText = "*"
	}

	fmt.Fprintf(w, "\t%s\t%s\t%s\n", attr.Name, rqText, color.HiBlueString(attr.Type))
	if attr.Description != "" {
		wrappedDesc := wrapper(attr.Description)
		for _, ln := range strings.Split(wrappedDesc, "\n") {
			fmt.Fprintf(w, "\t\t%s\n", color.HiBlackString(ln))
		}
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
