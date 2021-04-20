#!/bin/bash
#
## This is the startup script for the LoadBalancer container

#
## Variables for HAProxy
haproxyConfig=/etc/haproxy/haproxy.cfg
haproxyPidFile=/run/haproxy.pid
haproxyConfigTemplate=/usr/local/src/haproxy.cfg.j2
helperPodYaml=/usr/local/src/helperpod.yaml
ansibleLog=/var/log/helperpod_ansible_run.log

#
## Make sure the HELPERPOD_CONFIG_YAML env var has size
[[ ${#HELPERPOD_CONFIG_YAML} -eq 0 ]] && echo "FATAL: HELPERPOD_CONFIG_YAML env var not set!!!" && exit 254

#
## Take the HELPERPOD_CONFIG_YAML env variable and write out the YAML file.
echo ${HELPERPOD_CONFIG_YAML} | base64 -d > ${helperPodYaml}

#
## Create haproxy.cfg based on the template and yaml passed in.
ansible localhost -c local -e @${helperPodYaml} -m template -a "src=${haproxyConfigTemplate} dest=${haproxyConfig}" > ${ansibleLog} 2>&1

#
## Test for the validity of the config file. Run the HAProxy process if it passes
if ! /usr/sbin/haproxy -f ${haproxyConfig} -c -q ; then
	echo "============================="
	echo "FATAL: Invalid HAProxy config"
	echo "============================="
	exit 254
else
	echo "==========================="
	echo "Starting HAproxy service..."
	echo "==========================="
	rm -f ${haproxyPidFile}
	/usr/sbin/haproxy -Ws -f ${haproxyConfig} -p ${haproxyPidFile}
fi
##
##