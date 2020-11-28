package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// copyIgnCmd represents the copyIgn command
var copyIgnCmd = &cobra.Command{
	Use:   "copy-ign",
	Aliases:   []string{"copy-ignition", "cp-ign", "copyign", "cpign", "copyignition"},
	Short: "Copies ignition configs from given direcotry into the http container.",
	Long: `This command takes ignition configurations from the given directory,
and copies those files into the http contianer. For example:

	helpernodectl copy-ign --dir=/path/to/install/dir

This command must be run on the host that is to be the
helpernode. There is no support for copying the ignition
files to an external webserver.`,
	Run: func(cmd *cobra.Command, args []string) {
		// take what what passed into --dir
		dir, _ := cmd.Flags().GetString("dir")

		//check if it exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.Fatal("Please specify your install directory where the ignition files are.")
		} else {
			copyIgnToWebServer(dir)
		}

	},
}

func init() {
	rootCmd.AddCommand(copyIgnCmd)
	copyIgnCmd.PersistentFlags().String("dir", "", "the directory where the ignition files are located")

	//make the --dir flag required
	copyIgnCmd.MarkPersistentFlagRequired("dir")
}

func copyIgnToWebServer(dir string) {

	//if the webserver isn't running, don't bother running
	if !isImageRunning("helpernode-http") {
		logrus.Fatal("ERROR: helpernode-http isn't running\n")
	}

	// get the list of the ignition files from the directory
	dirglob := dir + "/*.ign"
	ignfiles, err := filepath.Glob(dirglob)

	// check for error of command
	if err != nil {
		logrus.Fatalf("Error getting ignition files %s: %s\n", ignfiles, err)
	}

	// Check to see if there's ANY igniton files there
	if len(ignfiles) == 0 {
		logrus.Fatalf("ERROR: No ignition files found in: " + dir + "\n")
	}

	// itterate over them
	clientpath := "/var/www/html/ignition/"
	for _, v := range ignfiles {
		logrus.Info("Copying over " + v + " to http container")
		// copy the igniton file to the http server - I should create a generic get/put function
		cmd := exec.Command("podman", "cp", v, "helpernode-http:" + clientpath)
		runCmd(cmd)
	}

	// Now we must fix permissions
	fixcmd := exec.Command("podman", "exec", "-it", "helpernode-http", "chmod", "+r", "-R", clientpath)
	runCmd(fixcmd)
}
