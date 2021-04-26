/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Args: func(cmd *cobra.Command, args []string) error {
		validateArgs(args)
		return nil
	},
	Use:   "restart",
	Short: "Restart containers",
	Long: `Restart a set of containers
The bare command will restart all running containers
Optionally you can pass a comma of delimited list of containers
i.e. helpernodectl restart dns,pxe`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Restarting services")
		reconcileImageList(imageList)
		for name, image := range images {
			if isImageRunning("helpernode-" + name) {
				logrus.Info("RESTARTING: Container helpernode-" + name)
				logrus.Info("STOPPING: Container helpernode-" + name)
				stopImage(name)
				logrus.Info("STARTING: Container helpernode-" + name)
				startImage(image, getEncodedConfuration(), name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
