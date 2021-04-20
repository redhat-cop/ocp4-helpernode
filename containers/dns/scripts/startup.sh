#!/bin/bash
#
## This is the startup script for the DNS container

#
## Variables for BIND
namedConfig=/etc/named.conf
zonefileConfig=/var/named/zonefile.db
reverseConfig=/var/named/reverse.db
namedPidFile=/run/named/named.pid
namedConfigTemplate=/usr/local/src/named.conf.j2
zonefileConfigTemplate=/usr/local/src/zonefile.j2
reverseConfigTemplate=/usr/local/src/reverse.j2
helperPodYaml=/usr/local/src/helperpod.yaml
ansibleLog=/var/log/helperpod_ansible_run.log
extraVars="-e setup_registry=false -e serialnumber=$(date +%s)"

#
## Make sure the HELPERPOD_CONFIG_YAML env var has size.
[[ ${#HELPERPOD_CONFIG_YAML} -eq 0 ]] && echo "FATAL: HELPERPOD_CONFIG_YAML env var not set!!!" && exit 254

#
## Take the HELPERPOD_CONFIG_YAML env variable and write out the YAML file.
echo ${HELPERPOD_CONFIG_YAML} | base64 -d > ${helperPodYaml}

#
## Create BIND files based on the template and yaml passed in.
ansible localhost -c local -e @${helperPodYaml} -m template -a "src=${namedConfigTemplate} dest=${namedConfig} mode='0640'" >> ${ansibleLog} 2>&1
ansible localhost -c local -e @${helperPodYaml} ${extraVars} -m template -a "src=${zonefileConfigTemplate} dest=${zonefileConfig} mode='0644'" >> ${ansibleLog} 2>&1
ansible localhost -c local -e @${helperPodYaml} ${extraVars} -m template -a "src=${reverseConfigTemplate} dest=${reverseConfig} mode='0644'" >> ${ansibleLog} 2>&1

#
## Test for the validity of the config file. Run the BIND process if it passes
if ! /usr/sbin/named-checkconf -z ${namedConfig} ; then
	echo "=========================="
	echo "FATAL: Invalid BIND config"
	echo "=========================="
	exit 254
else
	echo "========================"
	echo "Starting BIND service..."
	echo "========================"
	rm -f ${namedPidFile}
	/usr/sbin/named -u named -c ${namedConfig} -f
fi
##
##