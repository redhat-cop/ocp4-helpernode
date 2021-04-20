#!/bin/bash
#
## This is the startup script for the LoadBalancer container

#
## Variables for HTTPD
httpConfig=/etc/httpd/conf/httpd.conf
httpConfigTemplate=/usr/local/src/httpd.conf.j2
helperPodYaml=/usr/local/src/helperpod.yaml
ansibleLog=/var/log/helperpod_ansible_run.log

#
## Make sure the HELPERPOD_CONFIG_YAML env var has size
[[ ${#HELPERPOD_CONFIG_YAML} -eq 0 ]] && echo "FATAL: HELPERPOD_CONFIG_YAML env var not set!!!" && exit 254

#
## Take the HELPERPOD_CONFIG_YAML env variable and write out the YAML file.
echo ${HELPERPOD_CONFIG_YAML} | base64 -d > ${helperPodYaml}

#
## Create httpd.conf based on the template and yaml passed in.
ansible localhost -c local -e @${helperPodYaml} -e "http_port=${HELPERNODE_HTTP_PORT}" -m template -a "src=${httpConfigTemplate} dest=${httpConfig}" > ${ansibleLog} 2>&1

#
## Test for the validity of the config file. Run the HTTPD process if it passes
if ! /usr/sbin/httpd -f ${httpConfig} -t; then
	echo "=========================="
	echo "FATAL: Invalid HTTP config"
	echo "=========================="
	exit 254
else
	echo "========================"
	echo "Starting HTTP service..."
	echo "========================"
	/usr/sbin/httpd -f ${httpConfig} -D FOREGROUND
fi
##
##
