package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os/exec"
)

// getClientsCmd represents the getClients command
var getClientsCmd = &cobra.Command{
	Use:     "get-clients",
	Aliases: []string{"getclients"},
	Short:   "Gets the needed clients from the holding container.",
	Long: `This will get the  needed clients from the holding container.
Which is, by default, the http container. It saves this in your current
working directory. Example:

	helpernodectl get-clients

This command must be run on the node to be the helpernode.`,
	Run: func(cmd *cobra.Command, args []string) {
		getTheClients("http")
	},
}

func init() {
	rootCmd.AddCommand(getClientsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getClientsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getClientsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getTheClients(image string) {
	// By default the client path is an apache path
	clientpath := "/var/www/html/artifacts/"

	// If image is NOT running, it must be started
	if !isImageRunning("helpernode-" + image) {
		/* Right now, just start the image the way it is with a dummy value.
		TODO: Start the container with `sleep infinity`. Maybe build it into the startup.sh file?
		*/
		logrus.Info("Image helpernode-" + image + " is NOT running...starting temporarily")
		startImage(images[image], "bm90OiAidXNlZCIK", "http")
		for _, v := range clients {
			logrus.Info("Getting file " + v)
			// get the artifact - should probably make a put/get function later
			cmd := exec.Command("podman", "cp", "helpernode-"+image+":"+clientpath+v, ".")
			runCmd(cmd)
		}
		// Assume the user wants it stopped since it wasn't running
		stopImage("http")
	} else {
		for _, v := range clients {
			logrus.Info("Getting file " + v)
			// get the artifact
			cmd := exec.Command("podman", "cp", "helpernode-"+image+":"+clientpath+v, ".")
			runCmd(cmd)
		}
	}
}
