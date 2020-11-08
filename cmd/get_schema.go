package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
)

func init() {
	RootCmd.AddCommand(getSchemaCmd)

	// Bind flags
	getSchemaCmd.Flags().StringVarP(&getSchemaOutFile, "out", "o", "", "print the schema to a file")
}

// Flags
var getSchemaOutFile string

var getSchemaCmd = &cobra.Command{
	Use:   "get-schema",
	Short: "Executes a 'terraform providers schema'",
	Args:  cobra.RangeArgs(0, 0),
	Run: func(cmd *cobra.Command, args []string) {

		schema, err := getLocalProvidersSchema()
		if err != nil {
			log.Fatalf("error retrieving local providers: %v", err)
		}
		fmt.Println(schema)
	},
}

func getLocalProvidersSchema() (string, error) {
	binary, lookErr := exec.LookPath("terraform")
	if lookErr != nil {
		return "", fmt.Errorf("terraform not found on the system")
	}

	tfArgs := []string{"providers", "schema", "-json"}

	out, err := exec.Command(binary, tfArgs...).Output()
	if err != nil {
		return "", fmt.Errorf("error running terraform: %v", err)
	}
	outVal := string(out)
	if outVal == "{\"format_version\":\"0.1\"}\n" {
		return "", errors.New("no providers found. you might have initialized in an empty directory")
	}

	return outVal, nil
}
