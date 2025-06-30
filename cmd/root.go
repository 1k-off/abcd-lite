package cmd

import (
	"github.com/1k-off/abcd-lite/internal/embeddata"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "abcd-lite",
		Short: "ABCD Lite",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	embeddedData embeddata.EmbedData
)

// SetEmbeddedData allows main package to provide embedded FS data to commands
func SetEmbeddedData(data embeddata.EmbedData) {
	embeddedData = data
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(configCmd)
}
