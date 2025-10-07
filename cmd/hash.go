package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Think-iT-Labs/dirhash/lib"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ignoredPaths []string
var outputFormat string

func init() {
	rootCmd.AddCommand(hashCmd)
	hashCmd.Flags().StringSliceVarP(&ignoredPaths, "ignore", "x", nil, "ignored glob paths")
	hashCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "output format: text|json")
}

var hashCmd = &cobra.Command{
	Use:   "sha256",
	Short: "Print the sha256 hash of a directory",
	Long:  "Compute the SHA-256 hash of a directory. You can ignore files or folders using one or more -x/--ignore glob patterns (supports **).",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var directory = args[0]
		log.Debug("directory: ", directory)
		log.Debug("ignore: ", ignoredPaths)
		switch outputFormat {
		case "text":
			fmt.Println(lib.DirHash(directory, ignoredPaths))
		case "json":
			overall, details := lib.DirHashDetails(directory, ignoredPaths)
			payload := struct {
				Hash   string         `json:"hash"`
				Files  []lib.FileHash `json:"files"`
				Algo   string         `json:"algo"`
				Format string         `json:"format"`
			}{Hash: overall, Files: details, Algo: "sha256", Format: "json"}
			b, err := json.MarshalIndent(payload, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(b))
		default:
			log.Fatalf("unsupported output format: %s", outputFormat)
		}
	},
	Example: "dirhash sha256 -x node_modules/** -x '**/*.log' <directory>",
}
