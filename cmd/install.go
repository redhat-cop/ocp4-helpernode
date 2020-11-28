package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Args: func(cmd *cobra.Command, args []string) error {
		validateArgs(args)
		return nil
	},
	Use:   "install",
	Short: "Install creates a helpernode configuration",
	Long:  `Install creates pulls images and sets up initial ~/.helpernodectl.yaml config file`,
	Run: func(cmd *cobra.Command, args []string) {


		logrus.Info("called install")
		validateConfiguration()
		if logrus.GetLevel().String() == "debug" {
			for _,name := range imageList {
				logrus.Debug("Starting: " + name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
