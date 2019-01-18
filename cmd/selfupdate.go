package cmd

import (
	"log"

	"github.com/mpppk/connect-to-gce-win/lib"
	"github.com/spf13/cobra"
)

var selfUpdateCmd = &cobra.Command{
	Use:   "selfupdate",
	Short: "update connect-to-gce-win",
	//Long: `Update connect-to-gce-win`,
	Run: func(cmd *cobra.Command, args []string) {
		updated, err := lib.DoSelfUpdate()
		if err != nil {
			log.Println("Binary update failed:", err)
			return
		}
		if updated {
			log.Println("Current binary is the latest version", lib.Version)
		} else {
			log.Println("Successfully updated to version", lib.Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(selfUpdateCmd)
}
