package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	/* Robert says no
	"github.com/spf13/viper"
	*/
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"net"
)

type Config struct {
	Services []string `yaml:"disableservice"`
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}


// checking if service is running
func isServiceRunning(servicename string) bool {
	// check if the service is active
	activestate, err := exec.Command("systemctl", "show", "-p", "ActiveState", servicename).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", activestate, err)
		os.Exit(53)
	}
	// return the status
	as := strings.TrimSuffix(strings.Split(string(activestate), "=")[1], "\n")
	return as == "active"
}

// checking if service is running
func isServiceEnabled(servicename string) bool {
	// check if the service is active
	enabledstate, err := exec.Command("systemctl", "is-enabled", servicename).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", enabledstate, err)
		os.Exit(53)
	}
	// return the status
	es := strings.TrimSuffix(string(enabledstate), "\n")
	return es == "enabled"
}

// stopping service
func stopService(servicename string) {

	// stop the service only if it's running
	if isServiceRunning(servicename) {
		logrus.Info("Stopping service: " + servicename)
		//Stop the service with systemd
		cmd, err := exec.Command("systemctl", "stop", servicename).Output()
		// Check to see if the stop was successful
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
			os.Exit(53)
		}
	}
}

// stopping service
func startService(servicename string) {

	// start the service only if it isn't running
	if !isServiceRunning(servicename) {
		logrus.Info("Starting service: " + servicename)
		//Start the service with systemd
		cmd, err := exec.Command("systemctl", "start", servicename).Output()
		// Check to see if the start was successful
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
			os.Exit(53)
		}
	}
}

// disable service
func disableService(servicename string) {

	// Disable only if it needs to be
	if isServiceEnabled(servicename) {
		logrus.Info("Disabling service: " + servicename)
		//Stop the service with systemd
		cmd, err := exec.Command("systemctl", "disable", servicename).Output()
		// Check to see if the stop was successful
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
			os.Exit(53)
		}
	}
}

// enable service
func enableService(servicename string) {

	// Enable only if it needs to be
	if !isServiceEnabled(servicename) {
		logrus.Info("Enabling service: " + servicename)
		//enable the service with systemd
		cmd, err := exec.Command("systemctl", "enable", servicename).Output()
		// Check to see if the enable was successful
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
			os.Exit(53)
		}
	}
}

// get current firewalld rules and return as a slice of string
func getCurrentFirewallRules() []string {

	// get list of ports currently configured
	cmd, err := exec.Command("firewall-cmd", "--list-ports").Output()

	// check for error of command
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
		os.Exit(253)
	}

	// create a slice of string based on the output, trimming the newline first and splitting on " " (space)
	s := strings.Split(strings.TrimSuffix(string(cmd), "\n"), " ")

	// get the list of services currenly configured
	scmd, err := exec.Command("firewall-cmd", "--list-services").Output()

	// check for error
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", scmd, err)
		os.Exit(253)
	}

	// create a slice of string based on the output, trimming the newline first and splitting on " " (space)
	svc := strings.Split(strings.TrimSuffix(string(scmd), "\n"), " ")

	// create a new array based on this new svc array. We will be converting service names to port output
	// simiar to what we got with: firewall-cmd --list--ports
	var ns = []string{}

	// range over the service, find out it's port and append it to the array we just created
	for _, v := range svc {
		lc, err := exec.Command("firewall-cmd", "--service", v, "--get-ports", "--permanent").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", lc, err)
			os.Exit(253)
		}
		nv := strings.TrimSuffix(string(lc), "\n")
		if strings.Contains(nv, " ") {
			ls := strings.Split(nv, " ")
			for _, l := range ls {
				ns = append(ns, l)
			}
		} else {
			ns = append(ns, nv)
		}
	}

	// append this new array of string into the original
	for _, v := range ns {
		s = append(s, v)
	}

	// Let's return this slice of string
	return s
}

func 	openPort(port string) {

	// Open Ports using the port number
	cmd, err := exec.Command("firewall-cmd", "--add-port", port, "--permanent", "-q").Output()

	// check for error of command
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running add-port command %s: %s\n", cmd, err)
		os.Exit(253)
	}

	// Reload the firewall to get the most up to date table
	rcmd, err := exec.Command("firewall-cmd", "--reload").Output()

	// check for error of command
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running reload command %s: %s\n", rcmd, err)
		os.Exit(253)
	}
}

func verifyContainerRuntime() {
	_, err := exec.LookPath("podman")
	if err != nil {
		logrus.Fatal("Podman not found. Please install")
	}

}

func verifyFirewallCommand() {

	_, err := exec.LookPath("firewall-cmd")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error looking for firewall-cmd: %s\n", err)
		os.Exit(1)
	}
}

//This reconciles a list of images to start or stop
//defaults to all images unless specifically
func reconcileImageList(list []string) {

	//TODO change this to read from helpernodectl viper configuration
	/* Robert told me to do it this way
	disabledServices := viper.GetStringSlice("disabledServices")
	*/
	disabledServices := helperConfig.GetStringSlice("disabledServices")

	//all is implied so need to remove disabledServices
	if list[0] == "all" {
		//lets remove any disabled images
		for name := range disabledServices {
			delete(images, disabledServices[name])
		}
	} else {
		//create a new list from our args
		var subsetOfServices = make(map[string]string)

		for _, name := range list {
			subsetOfServices[name] = images[name]
		}
		images = subsetOfServices

		if logrus.GetLevel().String() == "debug" {
			for name, image := range subsetOfServices {
				logrus.Debug("Subset: " + name + ":" + image)
			}
		}
	}
	//TODO add plugable images
}

//TODO need to update this to use helperconfig
func getEncodedConfuration() string {
	//Not sure if this will stay here but lets do some validation on the configuration
	if( !validateConfiguration() ){
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

func validateConfiguration()  bool{
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
	if(!found){
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

		//TODO make sure plugable images is added to images var
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
func verifyConfig(){
	if !helpernodectlConfig.IsSet("configFile") &&  !rootCmd.PersistentFlags().Changed("config")  {
		logrus.Fatal("Config file was not passed or has no previous save")
	}else{
		logrus.Info("Found a configuration")
	}

}
