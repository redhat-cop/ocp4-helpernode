#!/bin/bash
nfsnamespace=nfs-provisioner
rbac=/usr/local/src/nfs-provisioner-rbac.yaml
deploy=/usr/local/src/nfs-provisioner-deployment.yaml
sc=/usr/local/src/nfs-provisioner-sc.yaml
#
export PATH=/usr/local/bin:$PATH
#
## Check openshift connection
if ! oc get project default -o jsonpath={.metadata.name} > /dev/null 2>&1 ; then
	echo "ERROR: Cannot connect to OpenShift. Are you sure you exported your KUBECONFIG path and are admin?"
	echo ""
	echo "...remember this is a POST INSTALL step."
	exit 254
fi
#
## Check to see if the namespace exists
if [ "$(oc get project default -o jsonpath={.metadata.name})" = "${nfsnamespace}" ]; then
	echo "ERROR: Seems like NFS provisioner is already deployed"
	exit 254
fi
#
## Check to see if important files are there
for file in ${rbac} ${deploy} ${sc}
do
	[[ ! -f ${file} ]] && echo "FATAL: File ${file} does not exist" && exit 254
done
#
## Check if the project is already there
if oc get project ${nfsnamespace} -o jsonpath={.metadata.name} > /dev/null 2>&1 ; then
	echo "ERROR: Looks like you've already deployed the nfs-provisioner"
	exit 254
fi
#
## If we are here; I can try and deploy
oc new-project ${nfsnamespace}
oc project ${nfsnamespace}
oc create -f ${rbac}
oc adm policy add-scc-to-user hostmount-anyuid system:serviceaccount:${nfsnamespace}:nfs-client-provisioner
oc create -f ${deploy} -n ${nfsnamespace}
oc create -f ${sc}
oc annotate storageclass nfs-storage-provisioner storageclass.kubernetes.io/is-default-class="true"
oc project default
oc rollout status deployment nfs-client-provisioner -n ${nfsnamespace}
#
## Show some info
cat <<EOF

Deployment started; you should monitor it with "oc get pods -n ${nfsnamespace}"

EOF
##
##
