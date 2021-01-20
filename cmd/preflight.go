package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
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
