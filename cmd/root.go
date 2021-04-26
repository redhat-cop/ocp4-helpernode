package cmd

import (
	"bufio"
	"encoding/base64"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"strings"
)

var cfgFile string

//This is the configuration passed into each container
var helperConfig = viper.New()

//Used to store tool related configurations like what services to start/stop by default
//most importantly it stores the location of the helpernodeConfig file if set
var helpernodectlConfig = viper.New()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "helpernodectl",
	Short: "Utility for the HelperNode",
	Long: `This cli utility is used to stop/start the HelperNode
on the host it's ran from. You need to provide a helpernode.yaml file
with information about your helper config. A simple example to start
your HelperNode is:

helpernodectl start --config=helpernode.yaml`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	//	setUpLogging()
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	if runtime.GOOS != "linux" {
		logrus.Fatal("helpernodectl only runs on linux")
	} else {
		verifyContainerRuntime()
		verifyFirewallCommand()
	}

	//TODO lets move --config to subcommands that need it. that way we can set it required or not.
	// Need it in start and pull
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.helpernodectl.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level (e.g. \"debug | info | warn | error\")")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	setUpLogging()
	setupCtlConfig()
	setupHelperConfig()
	createImageList()
	//TODO addDefaults()?

}

func setupCtlConfig() {
	//Used to store tool related configurations like what services to start/stop by default
	//most importantly it stores the location of the helpernodeConfig file if set

	home, err := homedir.Dir()
	if err != nil {
		logrus.Fatal(err)
	}
	helpernodectlConfig.AddConfigPath(home)
	helpernodectlConfig.SetConfigName(".helpernodectl")
	helpernodectlConfig.SetConfigType("yaml")

	helpernodectlConfig.SetDefault("services", map[string]bool{"dns": true, "dhcp": true, "http": true, "loadbalancer": true, "pxe": true})
	//	helpernodectlConfig.SetDefault("configFile", home + "/.helpernodectl.yaml")
	if err := helpernodectlConfig.ReadInConfig(); err != nil {
		//we got an error trying to read the config
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//our error was a ConfigFileNotFound error. Lets try to create it
			if emptyFile, err := os.Create(home + "/.helpernodectl.yaml"); err != nil {
				//we could not create the file. Lets log that but keep going using defaults
				logrus.Debug("Could not create $HOME/.helpernodectl.yaml file. Using defaults")
			} else {
				logrus.Debug("Created file $HOME/.helpernodectl.yaml")
				emptyFile.Close()
			}
		} else {
			//we could not read the file but it wasn't a ConfigNotFound error
			logrus.Debugf("We got the following error trying to read in file %s", err)
		}
	} else {
		//we read in the file...do we need to do anything else?
		logrus.Debug("Using file located at $HOME/.helpernodectl.yaml")
	}

	if err := helpernodectlConfig.WriteConfig(); err != nil {
		logrus.Error("Error writing config file")
	} else {
		logrus.Trace("Writing to: $HOME/.helpernodectl.yaml")
	}
}
func setupHelperConfig() {
	//steps that need to happen
	//if passed via command line use it
	//if not use what is defined in ~/.helpernodectl.yaml "configFile" field
	//    would need to test if it exists
	if cfgFile != "" {
		// Use config file from the flag.
		helperConfig.SetConfigFile(cfgFile)
		logrus.Info("Using --config file:" + cfgFile)
	} else {
		//TODO this will change to read from helpernodectl viper configuration
		// Find home directory.
		/* CH comment out for now
		logrus.Info("Found saved configuration in ~/.helper.yaml")
		*/
		home, err := homedir.Dir()
		if err != nil {
			logrus.Fatal(err)
		}

		// Search config in home directory with name ".helpernodectl" (without extension).
		helperConfig.AddConfigPath(home)
		helperConfig.SetConfigName(".helper")
		helperConfig.SetConfigType("yaml")
	}
	if err := helperConfig.ReadInConfig(); err != nil {
		//we got an error trying to read the config
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//our error was a ConfigFileNotFound error. Lets try to create it
			logrus.Debug("Could not find the configuration file")
		} else {
			//we could not read the file but it wasn't a ConfigNotFound error
			logrus.Debugf("We got the following error trying to read in file %s", err)
		}
	}
	logrus.Debug("Reading configuration file")
}

func setUpLogging() {
	//TODO set up something to read log-level from .helpernodectl.yaml
	logrus.SetOutput(os.Stdout)
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "invalid log-level"))
	}
	logrus.SetLevel(level)
	//	logrus.SetReportCaller(true)
}

//TODO need to update this to use helperconfig
func getEncodedConfuration() string {
	//Not sure if this will stay here but lets do some validation on the configuration
	if !validateConfiguration() {
		logrus.Fatal("Error in configuration file!!!")
	}
	// Check to see if file exists
	logrus.Trace("Config file used: " + helperConfig.ConfigFileUsed())
	var encoded string
	configurationFile := helperConfig.ConfigFileUsed()
	if _, err := os.Stat(helperConfig.ConfigFileUsed()); os.IsNotExist(err) {
		logrus.Error("File " + configurationFile + " does not exist")
	} else {
		// Open file on disk
		f, _ := os.Open(configurationFile)
		// Read file into a byte slice
		reader := bufio.NewReader(f)
		content, _ := ioutil.ReadAll(reader)
		//Encode to base64
		encoded = base64.StdEncoding.EncodeToString(content)
	}
	return encoded
}

func validateConfiguration() bool {
	//initially lets check that the interface name matches something on this NIC
	logrus.Infof("Validationg configuration in %s", helperConfig.ConfigFileUsed())
	interfaces, _ := net.Interfaces()
	found := false
	configInterface := helperConfig.GetString("helper.networkifacename")
	for _, ifc := range interfaces {
		if ifc.Name == configInterface {
			found = true
			break
		}
	}
	if !found {
		logrus.Errorf("Could not find %s in interface list of this machine", configInterface)
	}
	return found
}

func validateArgs(args []string) {
	imageCount := len(args)

	//if bare start command assume "all"
	if imageCount == 0 {
		logrus.Debug("Setting target images to all")
		imageList = []string{"all"}
	} else if imageCount == 1 {

		logrus.Debug("starting: " + args[0])
		//parse image list
		imageList = strings.Split(args[0], ",")
		logrus.Info(imageList)

		//TODO make sure pluggable images is added to images var
		//Lets make sure its in our list of images (should include pluggable images)
		for _, name := range imageList {
			if _, exists := images[name]; exists {
				continue
			} else {
				logrus.Fatal("Listed service is not part of image list ")
			}

		}
	} else {
		logrus.Fatal("Wrong number of arguments passed. Must be comma separated list")
	}

}
func verifyConfig() {
	if !helpernodectlConfig.IsSet("configFile") && !rootCmd.PersistentFlags().Changed("config") {
		logrus.Fatal("Config file was not passed or has no previous save")
	} else {
		logrus.Info("Found a configuration")
	}

}
