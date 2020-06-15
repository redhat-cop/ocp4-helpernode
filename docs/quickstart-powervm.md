# Helper Node Quickstart Install

This quickstart will get you up and running for HMC managed PowerVM.  You will need knowledge of the HMC to create the LPARs, here are docs for [HMC](https://www.ibm.com/support/knowledgecenter/en/9009-22A/p9eh6/p9eh6_kickoff.htm).

> **NOTE** For now static IP is not supported for PowerVM.

To start, login to your HMC GUI to perform required operation or ssh to HMC host to use CLI to do some operations. But most of the operations need to be done with HMC GUI.

## Create Virtual Network

We can use exist network, or create a new virtual network for cluster to use only.

List exist virtual network for a managed system:
```
lshwres -m <managed_system_name> -r virtualio --rsubtype vnetwork 
```
To create or view the virtual network in HMC for a managed system:
1. Click the managed system you want to list/create the network under `All Systems` view to get the `Partitions` view for the managed system
2. On left side of view under `PowerVM section` click the  `Virtual Networks`, all exist networks will be list under `Virtial Networks` section, if you want to use exist network, you can stop here, otherwise you can contine to to create a new network
3. Click the `Add Virtual Network` to open the wizard to create a new virtual network
> **NOTE** If the cluster cross mutiple managed systems, you need to create the virtual network with same vlan_id on all of the managed systems.

## Create the LPAR

 The disk should havs the OS image on it, so it can be set as the first boot device. So Clone the RHEL7/RHEL8 image disk is the best way to create the new disk for helper, and all cluster nodes.

__Create LPAR__

Here are the commands to create the LPAR:
```
mksyscfg -r lpar -m <managed_system> -i name=<lpar_name>, profile_name=default_profile, lpar_env=aixlinux, shared_proc_pool_util_auth=1, \
  min_mem=4096, desired_mem=4096, max_mem=8192, proc_mode=shared, min_proc_units=0.2, desired_proc_units=2.0, \
  max_proc_units=4.0, min_procs=1, desired_procs=4, max_procs=4, sharing_mode=uncap, uncap_weight=128, max_virtual_slots=64, \
  boot_mode=norm, conn_monitoring=1, shared_proc_pool_util_auth=1
```
In HMC, on the selected managed system's `Partitions` list view by click the `Create Partition ...` to open the `Create Partition(s)` dialog, and using the values from above `mksyscfg` to fill out the dialog for `Partition Name`, `Maximum Virtal Adapters`, and Processor and Memory configuration.

__Add Network__

1. Under managed system's partition list view, double click the new created the LPAR to open the LPAR `General` view
2. On left `Virtual I/O` section, click the `Virtual Networks` to open the `Virtual Networks` view
3. Click the `Attach Virtual Network` to open the pop up dialog to select virtual networks to use
4. Click `OK` to complete the operation

__Add Storage__

You can use any HMC supported storage from your environment. Below steps are just for using Fibre Channel storage(NPIV):
1. Under managed system's partition list view, double click the LPAR to open the LPAR `General` view
2. On left `Virtual I/O` section, click the `Virtual Storage` to open the `Virtual Storage` view
3. Click the `Virtual Fibre Channel` tab to show table of the defined `virtual Fibere Channel Devices`
4. Click the `Add Virtual Fibre Channel Device` to popup dialog for `Add Virtual Fibre Channel`
5. Check box to select the availabe  `Fibre Channel Port`
6. Click `OK` to complete the operation
> **NOTE** You need to setup disk on SAN storage to map to the host(LPAR) to let this LPAR to use it

## Boot Up LPAR with Bootp

After the LPAR is created and configured, now it is time to boot it up. There are two ways to do it:
* Use HMC to boot the LPAR 
1. Check the LPAR on `Partitions` list view
2. Click the `Actions` to select `Activate...` from popup menu to open the `Activate` dialog
3. For `Operation Type`, check the `Normal`, and  for Boot Mode` select `System Management
4. Using `vtmenu` on HMC host select managed system and LPAR to open the console for that LPAR
5. In SMS, select `Select Boot Options`->`Select Install/Boot Device`->`Network`->`Bootp`->`Select Device` to select the network adapter to boot from
6. Continue to start the bootp
* using HMC CLI to directly boot up LPAR with bootp:
```
lpar_netboot  -t ent -m <macaddr> -s auto -d auto <lpar_name> <profile_name> <managed_system>
```
> **NOTE** The `<macaddr>` format is `fad38e3ca520`, which does not contain `:`. 
Before `lpar_netboot`, the console to the LPAR has to be closed, otherwise it will fail.
After `lpar_netboot` completed, we can open the console to check the boot progress. 
If the disk is not set as first boot device, the bootp will loop on. To solve it, stop the boot to SMS after bootp complete with write the OS image to disk and reboot. Then set the disk as first boot device in SMS, and continue as normal boot.

## Create the Helper Node

Create helper LPAR using the steps descripted in `Create the LPAR` section with proper processor and memory configuration.

After helper is up and running, configured it with correct network configurations based on your network:
* IP - <helper_ip>
* NetMask - 255.255.255.0
* Default Gateway - <default_gateway>
* DNS Server - <default_DNS>


## Create Cluster Nodes

Create 6 LPARs using same steps descripted in above `Create the LPAR` section. Please follow the [min requirements](https://docs.openshift.com/container-platform/4.3/installing/installing_ibm_power/installing-ibm-power.html#minimum-resource-requirements_installing-ibm-power) for these LPARs.

__Bootstrap__

Create bootstrap LPAR

```
mksyscfg -r lpar -m <managed_system> -i name=ocp4-bootstrap, profile_name=default_profile, lpar_env=aixlinux, shared_proc_pool_util_auth=1, \
  min_mem=8192, desired_mem=16384, max_mem=16384, proc_mode=shared, min_proc_units=0.2, desired_proc_units=2.0, \
  max_proc_units=4.0, min_procs=1, desired_procs=4, max_procs=4, sharing_mode=uncap, uncap_weight=128, max_virtual_slots=64, \
  boot_mode=norm, conn_monitoring=1, shared_proc_pool_util_auth=1
```

__Masters__

Create the master LPARs

```
for i in master{0..2}
do
  mksyscfg -r lpar -m <managed_system> -i name="ocp4-${i}", profile_name=default_profile, lpar_env=aixlinux, shared_proc_pool_util_auth=1, \
    min_mem=4096, desired_mem=8192, max_mem=8192, proc_mode=shared, min_proc_units=0.2, desired_proc_units=2.0, \
    max_proc_units=4.0, min_procs=1, desired_procs=4, max_procs=4, sharing_mode=uncap, uncap_weight=128, max_virtual_slots=64, \
    boot_mode=norm, conn_monitoring=1, shared_proc_pool_util_auth=1
done
```

__Workers__

Create the worker LPARs

```
for i in worker{0..1}
do
  mksyscfg -r lpar -m <managed_system> -i name="ocp4-${i}", profile_name=default_profile, lpar_env=aixlinux, shared_proc_pool_util_auth=1, \
    min_mem=8192, desired_mem=8192, max_mem=16384, proc_mode=shared, min_proc_units=0.2, desired_proc_units=2.0, \
    max_proc_units=4.0, min_procs=1, desired_procs=4, max_procs=4, sharing_mode=uncap, uncap_weight=128, max_virtual_slots=64, \
    boot_mode=norm, conn_monitoring=1, shared_proc_pool_util_auth=1
done
```
> **NOTE** Make sure attached them to the network your choice, and add storage to them based your storage configuration after seccessfully created LPAR.

## Prepare the Helper Node

After the helper node is installed; login to it

```
ssh <helper_user>@<helper_ip>
```

> **NOTE** If using RHEL 7 - you need to enable the `rhel-7-server-rpms` and the `rhel-7-server-extras-rpms` repos. If you're using RHEL 8 you will need to enable `rhel-8-for-ppc64le-baseos-rpms`, `rhel-8-for-ppc64le-appstream-rpms`, and `ansible-2.9-for-rhel-8-ppc64le-rpms`

Install EPEL

```
yum -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-$(rpm -E %rhel).noarch.rpm
```

Install `ansible` and `git` and clone this repo

```
yum -y install ansible git
git clone https://github.com/RedHatOfficial/ocp4-helpernode
cd ocp4-helpernode
```

Get the Mac addresses with this command running from your HMC host:

```
for i in <managed_systems>
do
  lshwres -m $i -r virtualio --rsubtype eth --level lpar -F lpar_name,mac_addr
done
```
Or if used SRIOV's logical port for network:
```
for i in <managed_systems>
do
  lshwres -m $i -r sriov --rsubtype logport --level eth -F lpar_name,mac_addr
done
```
Edit the [vars.yaml](examples/vars-ppc64le.yaml) file with the mac addresses of the cluster node LPARs.

```
cp docs/examples/vars-pppc64le.yaml vars.yaml
```

> **NOTE** See the `vars.yaml` [documentation page](vars-doc.md) for more info about what it does.

## Run the playbook

Run the playbook to setup your helper node

```
ansible-playbook -e @vars.yaml tasks/main.yml
```

After it is done run the following to get info about your environment and some install help


```
/usr/local/bin/helpernodecheck
```

## Create Ignition Configs

Now you can start the installation process. Create an install dir.

```
mkdir ~/ocp4
cd ~/ocp4
```

Create a place to store your pull-secret

```
mkdir -p ~/.openshift
```

Visit [try.openshift.com](https://cloud.redhat.com/openshift/install) and select "Bare Metal". Download your pull secret and save it under `~/.openshift/pull-secret`

```shell
# ls -1 ~/.openshift/pull-secret
/root/.openshift/pull-secret
```

This playbook creates an sshkey for you; it's under `~/.ssh/helper_rsa`. You can use this key or create/user another one if you wish.

```shell
# ls -1 ~/.ssh/helper_rsa
/root/.ssh/helper_rsa
```

> :warning: If you want you use your own sshkey, please modify `~/.ssh/config` to reference your key instead of the one deployed by the playbook

Next, create an `install-config.yaml` file.

> :warning: Make sure you update if your filenames or paths are different.

```
cat <<EOF > install-config.yaml
apiVersion: v1
baseDomain: example.com
compute:
- hyperthreading: Enabled
  name: worker
  replicas: 0
controlPlane:
  hyperthreading: Enabled
  name: master
  replicas: 3
metadata:
  name: ocp4
networking:
  clusterNetworks:
  - cidr: 10.254.0.0/16
    hostPrefix: 24
  networkType: OpenShiftSDN
  serviceNetwork:
  - 172.30.0.0/16
platform:
  none: {}
pullSecret: '$(< ~/.openshift/pull-secret)'
sshKey: '$(< ~/.ssh/helper_rsa.pub)'
EOF
```

Create the installation manifests

```
openshift-install create manifests
```

Edit the `manifests/cluster-scheduler-02-config.yml` Kubernetes manifest file to prevent Pods from being scheduled on the control plane machines by setting `mastersSchedulable` to `false`.

```shell
$ sed -i 's/mastersSchedulable: true/mastersSchedulable: false/g' manifests/cluster-scheduler-02-config.yml
```

It should look something like this after you edit it.

```shell
$ cat manifests/cluster-scheduler-02-config.yml
apiVersion: config.openshift.io/v1
kind: Scheduler
metadata:
  creationTimestamp: null
  name: cluster
spec:
  mastersSchedulable: false
  policy:
    name: ""
status: {}
```

Next, generate the ignition configs

```
openshift-install create ignition-configs
```

Finally, copy the ignition files in the `ignition` directory for the websever

```
cp ~/ocp4/*.ign /var/www/html/ignition/
restorecon -vR /var/www/html/
chmod o+r /var/www/html/ignition/*.ign
```

## Install LPARs

Boot up the LPARs using steps described in `Boot Up with Bootp` section.

Boot/install the LPARs in the following order

* Bootstrap
* Masters
* Workers

On your laptop/workstation visit the status page

```
firefox http://<helper_ip>:9000
```

You'll see the bootstrap turn "green" and then the masters turn "green", then the bootstrap turn "red". This is your indication that you can continue.

Also you can check all cluster node LPAR status in HMC's partition list view.

## Wait for install

The boostrap LPAR actually does the install for you; you can track it with the following command.

```
openshift-install wait-for bootstrap-complete --log-level debug
```

Once you see this message below...

```
DEBUG OpenShift Installer v4.2.0-201905212232-dirty
DEBUG Built from commit 71d8978039726046929729ad15302973e3da18ce
INFO Waiting up to 30m0s for the Kubernetes API at https://api.ocp4.example.com:6443...
INFO API v1.13.4+838b4fa up
INFO Waiting up to 30m0s for bootstrapping to complete...
DEBUG Bootstrap status: complete
INFO It is now safe to remove the bootstrap resources
```

...you can continue....at this point you can delete the bootstrap server.


## Finish Install

First, login to your cluster

```
export KUBECONFIG=/root/ocp4/auth/kubeconfig
```

Set the registry for your cluster

First, you have to set the `managementState` to `Managed` for your cluster

```
oc patch configs.imageregistry.operator.openshift.io cluster --type merge --patch '{"spec":{"managementState":"Managed"}}'
```

For PoCs, using `emptyDir` is ok (to use PVs follow [this](https://docs.openshift.com/container-platform/latest/installing/installing_bare_metal/installing-bare-metal.html#registry-configuring-storage-baremetal_installing-bare-metal) doc)

```
oc patch configs.imageregistry.operator.openshift.io cluster --type merge --patch '{"spec":{"storage":{"emptyDir":{}}}}'
```

If you need to expose the registry, run this command

```
oc patch configs.imageregistry.operator.openshift.io/cluster --type merge -p '{"spec":{"defaultRoute":true}}'
```

> Note: You can watch the operators running with `oc get clusteroperators`

Watch your CSRs. These can take some time; go get come coffee or grab some lunch. You'll see your nodes' CSRs in "Pending" (unless they were "auto approved", if so, you can jump to the `wait-for install-complete` step)

```
watch oc get csr
```

To approve them all in one shot...

```
oc get csr --no-headers | awk '{print $1}' | xargs oc adm certificate approve
```

Check for the approval status (it should say "Approved,Issued")

```
oc get csr | grep 'system:node'
```

Once Approved; finish up the install process

```
openshift-install wait-for install-complete
```

## Login to the web console

The OpenShift 4 web console will be running at `https://console-openshift-console.apps.{{ dns.clusterid }}.{{ dns.domain }}` (e.g. `https://console-openshift-console.apps.ocp4.example.com`)

* Username: kubeadmin
* Password: the output of `cat /root/ocp4/auth/kubeadmin-password`

## Upgrade

If you didn't install the latest 4.3.Z release then just run the following.

```
oc adm upgrade --to-latest
```

If you're having issues upgrading you can try adding `--force` to the upgrade command.

```
oc adm upgrade --to-latest --force
```

See [issue #46](https://github.com/RedHatOfficial/ocp4-helpernode/issues/46) to understand why the `--force` is necessary and an alternative to using it.


Scale the router if you need to

```
oc patch --namespace=openshift-ingress-operator --patch='{"spec": {"replicas": 3}}' --type=merge ingresscontroller/default
```

## DONE

Your install should be done! You're a UPI master!
