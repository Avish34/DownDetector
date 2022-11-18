package urlchecker

import (
	_ "fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// root cmd
var rootCmd = &cobra.Command{
	Use:   "url-monitor",
	Short: "url-monitor",
	Long:  `url-monitor`,
}

// This is the function which run when packet get's executed
func init() {
	rootCmd.AddCommand(checkStatusCmd)
	checkStatusCmd.Flags().BoolP("statistics", "s", false, "Show statistics")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
