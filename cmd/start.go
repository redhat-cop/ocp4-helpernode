package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Args: func(cmd *cobra.Command, args []string) error {
		validateArgs(args)
		return nil
	},
	Use:   "start",
	Short: "Starts HelperNode containers based on the provided manifest.",
	Long: `This will start the containers needed for the HelperNode to run.
It will run the services depending on what manifest is passed.
Examples:
	helpernodectl start --config=helpernode.yaml
	
	cp helpernode.yaml ~/.helpernode.yaml
	helpernodectl start
This manifest should have all the information the services need to start
up successfully.`,
	Run: func(cmd *cobra.Command, args []string) {
		verifyConfig()
		skippreflight, _ := cmd.Flags().GetBool("skip-preflight")
		if skippreflight {
			logrus.Info("Skipping Preflightchecks\n======================\n")
		} else {
			preflightCmd.Run(cmd, []string{})
			logrus.Info("Starting Containers\n======================\n")
		}

		runContainers()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("skip-preflight", "s", false, "Skips preflight checks and tries to start the containers")
	//TODO right now this will ALL images.
	// Need to update to only whats in comma separated list if that is passed
	//TODO add --disable?

}

func runContainers() {
	reconcileImageList(imageList)
	if logrus.GetLevel().String() == "debug" {
		for _, name := range imageList {
			logrus.Info("Starting: " + name)
		}
	}
	for name, image := range images {
		if isImageRunning("helpernode-" + name) {
			logrus.Info("SKIPPING: Container helpernode-" + name + " already running.")
		} else {
			startImage(image, getEncodedConfuration(), name)
		}
	}
}
