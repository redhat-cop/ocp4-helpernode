package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
	"os"
	"os/exec"
	"strings"
)

// preflight error counter
var preflightErrorCount int = 0

// preflightCmd represents the preflight command
var preflightCmd = &cobra.Command{
	Use:   "preflight",
	Short: "Checks for any conflicts on the host.",
	Long: `This checks for conflicts on the host and can optionally fix
errors it finds. For example:
	
	helpernodectl preflight

	helpernodectl preflight --fix-all


This checks for port conflicts, systemd conflicts, and also checks any 
firewall rules. It will optionally fix systemd and firewall rules by
passing the --fix-all option (EXPERIMENTAL).`,
	Run: func(cmd *cobra.Command, args []string) {

		fixall, _ := cmd.Flags().GetBool("fix-all")
		logrus.Info("RUNNING PREFLIGHT TASKS")
		if fixall {
			logrus.Info("==========================BESTEFFORT IN FIXING ERRORS============================\n")
		}
		//fix-all defaults to false unless passed on the command line
		//		systemdCheck(fixall)
		//		portCheck()
		//		firewallRulesCheck(fixall)

		logrus.WithFields(logrus.Fields{
			"SystemdCheck": systemdCheck(fixall),
			"PortCheck":    portCheck(),
			"FWRules":      firewallRulesCheck(fixall),
		}).Info("Preflight Summary")
		if preflightErrorCount == 0 {
			logrus.Infof("No preflight errors found")
		} else {
			if !fixall {
				logrus.Fatal("Cannot Start, preflight errors found")
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(preflightCmd)
	preflightCmd.Flags().BoolP("fix-all", "x", false, "Does the needful and fixes errors it finds - EXPERIMENTAL")

}

func portCheck() int {
	logrus.Info("Starting Port Checks")
	// set the error count to 0
	porterrorcount := 0

	for port, protocolArray := range portlist {
		for _, protocol := range protocolArray {
			logrus.Debugf("Testing port %s on protocol %s", port, protocol)
			//check if you can listen on this port on TCP
			if protocol == "tcp" {
				if t, err := net.Listen(protocol, ":"+port); err == nil {
					// If this returns an error, then something else is listening on this port
					if err != nil {
						if logrus.GetLevel().String() == "debug" {
							logrus.Warnf("Port check  %s/%s is in use", port, protocol)
						}
						porterrorcount += 1
					}
					t.Close()

				}
			} else if protocol == "udp" {
				if u, err := net.ListenPacket(protocol, ":"+port); err == nil {
					// If this returns an error, then something else is listening on this port
					if err != nil {
						if logrus.GetLevel().String() == "debug" {
							logrus.Warnf("Port check  %s/%s is in use", port, protocol)
						}
						porterrorcount += 1
					}
					u.Close()

				}
			}
		}
	}

	// Display that no errors were found
	if porterrorcount > 0 {
		preflightErrorCount += 1
	}
	logrus.WithFields(logrus.Fields{"Port Issues": porterrorcount}).Info("Preflight checks for Ports")
	return porterrorcount
}

func systemdCheck(fix bool) int {
	// set the error count to 0
	svcerrorcount := 0
	logrus.Info("Starting Systemd Checks")

	for _, s := range systemdsvc {
		if isServiceRunning(s) {
			logrus.Debug("Service " + s + " is running")
			svcerrorcount += 1
			if fix {
				logrus.Info("STOPPING/DISABLING SERVICE: " + s)
				stopService(s)
				disableService(s)
			}
		}
	}
	// Display that no errors were found
	if svcerrorcount > 0 {
		preflightErrorCount += 1
	}
	logrus.WithFields(logrus.Fields{"Systemd Issues": svcerrorcount}).Info("Preflight checks for Systemd")
	return svcerrorcount

}

func firewallRulesCheck(fix bool) int {
	// set the error count to 0
	fwerrorcount := 0
	fwfixCount := 0

	logrus.Info("Running firewall checks")
	// Check if firewalld service is running
	if !isServiceRunning("firewalld.service") {
		//		fwerrorcount += 1
		logrus.Debug("Service firewalld.service is NOT running")
		if fix {
			startService("firewalld.service")
			enableService("firewalld.service")
		}
	}

	// get the current firewall rules on the host and set it to "s"
	s := getCurrentFirewallRules()
	// loop through each firewall rule:
	// If there's a match, that means the rule is there and nothing needs to be done.
	// If it's NOT there, it needs to be enabled (if requested)
	for port, protocolArray := range portlist {
		for _, protocol := range protocolArray {
			_, found := find(s, port+"/"+protocol)
			if !found {
				if logrus.GetLevel().String() == "debug" {
					//this is a bit weird but only want to log these in debug mode.
					//BUT using WARN so they show up yellow
					logrus.Warnf("Firewall rule %s not found", port+"/"+protocol)
				}
				fwerrorcount += 1
				if fix {
					logrus.Info("OPENING PORT: " + port + "/" + protocol)
					openPort(port + "/" + protocol)
					fwfixCount++
				}
			}
		}
	}

	// Display that no errors were found
	if fwerrorcount > 0 {
		preflightErrorCount += 1
	}
	if fix {
		logrus.WithFields(logrus.Fields{"Firewall Issues": fwerrorcount, "Firewall rules added": fwfixCount}).Info("Preflight checks for Firewall")
	} else {
		logrus.WithFields(logrus.Fields{"Firewall Issues": fwerrorcount}).Info("Preflight checks for Firewall")
	}
	return fwerrorcount
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

func openPort(port string) {

	// Open Ports using the port number
	cmd := exec.Command("firewall-cmd", "--add-port", port, "--permanent", "-q")
	runCmd(cmd)

	// check for error of command
	/*	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running add-port command %s: %s\n", cmd, err)
		os.Exit(253)
	}*/

	// Reload the firewall to get the most up to date table
	cmd = exec.Command("firewall-cmd", "--reload")
	runCmd(cmd)
}

func verifyFirewallCommand() {

	_, err := exec.LookPath("firewall-cmd")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error looking for firewall-cmd: %s\n", err)
		os.Exit(1)
	}
}

