package cmd

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

//used by all commands to print output meaningfully
//logs as Info
//TODO maybe see if can bundle up a list of cmds to run
func runCmd(cmd *exec.Cmd) {
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {

		logrus.Fatal(err)
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		m := scanner.Text()
		logrus.Info(m)
	}
	cmd.Wait()
}

//going to covert this to use the podman module in the future
func pullImage(image string, version string) {

	cmd := exec.Command(containerRuntime, "pull", image)
	runCmd(cmd)
}

//going to covert this to use the podman module in the future
//TODO clean up this to just take one string. build the string elsewhere
//TODO we need to adjust startImage to account for pluggable container that won't take the encoded file or need --net-host
//TODO lets add a generic label to make stop all easier.  filter on helpernode=true or something similar
func startImage(image string, encodedyaml string, containername string) {
	startOptions := []string{containerRuntime, "run", "--rm", "-d", "--label=helpernode-" + containername + "=" + VERSION, "--net=host", "--name=helpernode-" + containername}
	coreImageOptions := "--env=HELPERPOD_CONFIG_YAML=" + encodedyaml

	pluggableServices := helperConfig.GetStringMapString("pluggableServices")
	if _, ok := pluggableServices[containername]; ok {
		pContainerStartOptions := helperConfig.GetString("pluggableServices." + containername + ".startupOptions")
		parsedOptions := strings.Fields(pContainerStartOptions)
		startOptions = append(startOptions, parsedOptions...)
	} else {
		startOptions = append(startOptions, coreImageOptions)
	}
	logrus.Info("Starting helpernode-" + containername)
	startOptions = append(startOptions, image)
	cmd := exec.Command(containerRuntime)
	cmd.Args = startOptions
	runCmd(cmd)

}

//going to covert this to use the podman module in the future
func stopImage(containername string) {

	logrus.Info("Stopping: helpernode-" + containername)
	//TODO check if image is runnign here rather than in start.go
	cmd := exec.Command(containerRuntime, "stop", "helpernode-"+containername)
	runCmd(cmd)

}

//check if an image is running. Return true if it is
func isImageRunning(containername string) bool {
	out, err := exec.Command("podman", "ps", "--format", "{{.Names}}", "--filter=label="+containername+"="+VERSION).Output()
	name := strings.TrimSuffix(string(out), "\n")
	if err != nil {
		logrus.Debug(err)
	} else if name == containername {
		logrus.Debugf("%s is running", name)
		return true
	}
	return false
}

func createImageList() {
	//TODO probably move this to setupCtlConfig
	helpernodectlConfig.SetEnvPrefix("helpernode")
	helpernodectlConfig.BindEnv("image_prefix")
	if helpernodectlConfig.GetString("image_prefix") == "" {
		logrus.Debug("HELPERNODE_IMAGE_PREFIX not found")
		helpernodectlConfig.Set("image_prefix", "quay.io")
	} else {
		logrus.Debug("Using quay.io as the registry")
	}

	helpernodectlConfig.AutomaticEnv() // read in environment variables that match

	registry = helpernodectlConfig.GetString("image_prefix")

	for _, name := range coreImageNames {
		images[name] = registry + "/" + repository + "/" + name + ":" + VERSION
	}
	//TODO Add pluggable images here
	pluggableServices := helperConfig.GetStringMapString("pluggableServices")
	for pluggableImageName := range pluggableServices {
		pImageName := helperConfig.GetString("pluggableServices." + pluggableImageName + ".image")
		logrus.Debugf("image value is %s\n", pImageName)
		images[pluggableImageName] = pImageName

		//lets get ports
		ports := helperConfig.GetStringSlice("pluggableServices." + pluggableImageName + ".ports")
		for _, v := range ports {
			portvalue := strings.Split(v, "/")
			portlist[portvalue[0]] = append(portlist[portvalue[0]], portvalue[1])
		}
	}

	//Just some logic to print if in debug
	if logrus.GetLevel().String() == "debug" {
		logrus.Debug("Using registry : " + registry)
		for name, image := range images {
			logrus.Debug(name + ":" + image)
		}
	}
}

func verifyContainerRuntime() {
	_, err := exec.LookPath("podman")
	if err != nil {
		logrus.Fatal("Podman not found. Please install")
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
