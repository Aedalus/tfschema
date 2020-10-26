package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
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
		binary, lookErr := exec.LookPath("terraform")
		if lookErr != nil {
			fmt.Println("terraform not found on the system")
			os.Exit(1)
		}

		tfArgs := []string{"providers", "schema", "-json"}

		out, err := exec.Command(binary, tfArgs...).Output()
		if err != nil {
			log.Fatal(err)
		}
		outVal := string(out)
		if outVal == "{\"format_version\":\"0.1\"}\n" {
			fmt.Println("No providers found. You might have initialized in an empty directory.")
			os.Exit(1)
		} else {
			fmt.Print(outVal)
		}
	},
}
