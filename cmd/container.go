package cmd

import (
	"bufio"
	"fmt"
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

	fmt.Println("Pulling: " + image + ":" + version)
	cmd := exec.Command(containerRuntime, "pull", image+":"+version)
	runCmd(cmd)
}

//going to covert this to use the podman module in the future
//TODO clean up this to just take one string. build the string elsewhere
//TODO we need to adjust startImage to account for pluggable container that won't take the encoded file or need --net-host
//TODO lets add a generic label to make stop all easier.  filter on helpernode=true or something similar
func startImage(image string, encodedyaml string, containername string) {
	logrus.Info("Starting helpernode-" + containername)
	cmd := exec.Command(containerRuntime, "run", "--rm", "-d", "--env=HELPERPOD_CONFIG_YAML="+encodedyaml, "--label=helpernode-"+containername+"="+VERSION, "--net=host", "--name=helpernode-"+containername, image)
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
