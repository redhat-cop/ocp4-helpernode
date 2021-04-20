#!/bin/bash
#
## This is the startup script for the LoadBalancer container

#
## Variables for DHCP
dhcpConfig=/etc/dhcp/dhcpd.conf
dhcpConfigTemplate=/usr/local/src/dhcpd.conf.j2
helperPodYaml=/usr/local/src/helperpod.yaml
ansibleLog=/var/log/helperpod_ansible_run.log

#
## Make sure the HELPERPOD_CONFIG_YAML env var has size
[[ ${#HELPERPOD_CONFIG_YAML} -eq 0 ]] && echo "FATAL: HELPERPOD_CONFIG_YAML env var not set!!!" && exit 254

#
## Take the HELPERPOD_CONFIG_YAML env variable and write out the YAML file.
echo ${HELPERPOD_CONFIG_YAML} | base64 -d > ${helperPodYaml}

#
## Create dhcpd.conf based on the template and yaml passed in.
ansible localhost -c local -e @${helperPodYaml} -m template -a "src=${dhcpConfigTemplate} dest=${dhcpConfig}" > ${ansibleLog} 2>&1

#
## Test for the validity of the config file. Run the DHCP process if it passes
if ! /usr/sbin/dhcpd -t -cf ${dhcpConfig} ; then
	echo "=========================="
	echo "FATAL: Invalid DHCP config"
	echo "=========================="
	exit 254
else
	echo "========================"
	echo "Starting DHCP service..."
	echo "========================"
	/usr/sbin/dhcpd -f -cf ${dhcpConfig} -user dhcpd -group dhcpd --no-pid
fi
##
##