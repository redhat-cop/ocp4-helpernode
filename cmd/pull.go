package cmd

import (
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Args: func(cmd *cobra.Command, args []string) error {
		validateArgs(args)
		return nil
	},
	Use:   "pull",
	Short: "Pulls images into your node",
	Long: `This will pull the  core helpernode images onto your local host. These images are used to
start all the services needed for the HelperNode. These images are:

quay.io/helpernode/pxe
quay.io/helpernode/http
quay.io/helpernode/loadbalancer
quay.io/helpernode/dns
quay.io/helpernode/dhcp`,
	Run: func(cmd *cobra.Command, args []string) {
		pullImages()
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}

// Loop through images and pull them
func pullImages() {
	reconcileImageList(imageList)
	for _, image := range images {
		//TODO need to update this to version aftertesting
			pullImage(image, VERSION)
	}
}
