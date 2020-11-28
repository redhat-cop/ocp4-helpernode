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
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var file string
// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save is used to save a configuration file -f is required",
	Long: `Save is used to save a configuration file -f is required`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := homedir.Dir()
		if err != nil {
			logrus.Fatal(err)
		}
		helpernodectlConfig.Set("configFile", home + "/.helper.yaml")
		helpernodectlConfig.WriteConfig()

		helperConfig.SetConfigFile(file)
		helperConfig.ReadInConfig()

		helperConfig.AddConfigPath(home)
		helperConfig.SetConfigName(".helper")
		helperConfig.SetConfigType("yaml")

		if helperFile, err := os.Create(home + "/.helper.yaml"); err != nil {
			//we could not create the file. Lets log that but keep going using defaults
			logrus.Debug("Could not create $HOME/.helper.yaml file. Using defaults")
		} else {
			logrus.Debug("Created file $HOME/.helper.yaml")
			helperFile.Close()
		}



		if err := helperConfig.WriteConfig(); err != nil {
			logrus.Debugf("Error writing config file %s", err)
		} else {
			logrus.Trace("Writing to: $HOME/.helper.yaml")
		}
		logrus.Info("Saving " + file + " to ~/.helper.yaml")



	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
	saveCmd.Flags().StringVarP(&file, "file", "f", "config file (saved to $HOME/.helper.yaml)", "")
	saveCmd.MarkFlagRequired("file")

}
