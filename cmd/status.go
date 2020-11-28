package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Shows the status of running containers",
	Long: `This shows the status of the running containers
on the host. Example:

	helpernodectl status

This simply passes the information provided by the contianer
runtime.`,
	Run: func(cmd *cobra.Command, args []string) {
		getStatus()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getStatus() {
	// Call `podman ps` filtering by containers that start with "helpernode" and only print out what is needed
	cmd, err := exec.Command("podman", "ps", "-a", "--filter", "name=^helpernode", "--format", fmt.Sprintf(`table {{.Names}} {{.Status}} {{.Image}}`)).Output()

	// check if there's an error
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running podman-cp command %s: %s\n", cmd, err)
		os.Exit(253)
	}

	//print whatever it gave me
	fmt.Print(string(cmd))
}
