#!/bin/bash
#
## This is the startup script for the LoadBalancer container

#
## Variables for PXE
tftpBootDir=/var/lib/tftpboot
pxeConfig=${tftpBootDir}/pxelinux.cfg
rhcosDir=${tftpBootDir}/rhcos
bootstrapPxeTemplate=/usr/local/src/pxe-bootstrap.j2
masterPxeTemplate=/usr/local/src/pxe-master.j2
workerPxeTemplate=/usr/local/src/pxe-worker.j2
helperPodYaml=/usr/local/src/helperpod.yaml
ansibleLog=/var/log/helperpod_ansible_run.log

#
## Make sure the HELPERPOD_CONFIG_YAML env var has size
[[ ${#HELPERPOD_CONFIG_YAML} -eq 0 ]] && echo "FATAL: HELPERPOD_CONFIG_YAML env var not set!!!" && exit 254

#
## Take the HELPERPOD_CONFIG_YAML env variable and write out the YAML file.
echo ${HELPERPOD_CONFIG_YAML} | base64 -d > ${helperPodYaml}

#
## Create pxe/tftp files based on the template and yaml passed in.

# First create the bootstrap, only if it's provided
if [[ ! $(yq -r .bootstrap ${helperPodYaml}) == "null" ]]; then
	ansible localhost -c local -e @${helperPodYaml} -e "http_port=${HELPERNODE_HTTP_PORT}" -e disk=$(yq -r .bootstrap.disk ${helperPodYaml} | tr [:upper:] [:lower:]) -m template -a "src=${bootstrapPxeTemplate} dest=${pxeConfig}/01-$(yq -r .bootstrap.macaddr ${helperPodYaml} | tr [:upper:] [:lower:] | sed 's~:~-~g') mode=0555" >> ${ansibleLog} 2>&1
fi

# For the masters, we need to loop
masters=($(yq -r '.masters[] | @base64'  < ${helperPodYaml}))
for master in ${!masters[@]}
do
	ansible localhost -c local -e @${helperPodYaml} -e "http_port=${HELPERNODE_HTTP_PORT}" -e disk=$(yq -r .masters[${master}].disk ${helperPodYaml} | tr [:upper:] [:lower:]) -m template -a "src=${masterPxeTemplate} dest=${pxeConfig}/01-$(yq -r .masters[${master}].macaddr ${helperPodYaml} | tr [:upper:] [:lower:] | sed 's~:~-~g') mode=0555" >> ${ansibleLog} 2>&1
done

# Only loop through the workers if there is any (i.e. "compact cluster" mode)
if [[ ! $(yq -r .workers ${helperPodYaml}) == "null" ]]; then
	workers=($(yq -r '.workers[] | @base64'  < ${helperPodYaml}))
	for worker in ${!workers[@]}
	do
		ansible localhost -c local -e @${helperPodYaml} -e "http_port=${HELPERNODE_HTTP_PORT}" -e disk=$(yq -r .workers[${worker}].disk ${helperPodYaml} | tr [:upper:] [:lower:]) -m template -a "src=${workerPxeTemplate} dest=${pxeConfig}/01-$(yq -r .workers[${worker}].macaddr ${helperPodYaml} | tr [:upper:] [:lower:] | sed 's~:~-~g') mode=0555" >> ${ansibleLog} 2>&1
	done
fi

#
## PXE is a "best effort" service that is kind of "old". So putting this here as a placeholder until someone has time to write a "checker"
if false ; then
	echo "=============================="
	echo "FATAL: Invalid PXE/TFTP config"
	echo "=============================="
	exit 254
else
	echo "============================"
	echo "Starting PXE/TFTP service..."
	echo "============================"
	/usr/sbin/in.tftpd -L --verbosity 4  --secure ${tftpBootDir}
fi
##
##
