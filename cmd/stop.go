/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
