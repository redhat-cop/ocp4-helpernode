package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Args: func(cmd *cobra.Command, args []string) error {
		validateArgs(args)
		return nil
	},
	Use:   "stop",
	Short: "Stops helpernode containres",
	Long:  "Stops helpernode containres",
	Run: func(cmd *cobra.Command, args []string) {
		stopContainers()
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func stopContainers() {
	reconcileImageList(imageList)
	for name, _ := range images {
		if !isImageRunning("helpernode-" + name) {
			logrus.Info("SKIPPING: Container helpernode-" + name + " already stopped.")
		} else {
			stopImage(name)
		}
	}
}
