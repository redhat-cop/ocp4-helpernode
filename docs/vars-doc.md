# How to use vars.yaml

This page gives you an explanation of the variables found in the [vars.yaml](examples/vars.yaml) example given in this repo to help you formulate your own/edit the provided example.

There are [examples provided](#example-vars-file) in this page for various workflows.

## Disk to install RHCOS

> **OPTIONAL** if doing static ips

In the first section, you'll see that it's asking for a disk

```
disk: vda
```

This needs to be set to the disk where you are installing RHCOS on the masters/workers. This will be set in the boot options for the [pxe server](../templates/default.j2).

**NOTE**: This will be the same for ALL masters and workers. Support for "mixed disk" (i.e. if your masters use `sda` and your workers are `vda`) is not supported at this time

You can, however edit the `/var/lib/tftpboot/pxelinux.cfg/default` file by hand after the install.


## Helper Section

This section sets the variables for the helpernode

```
helper:
  name: "helper"
  ipaddr: "192.168.7.77"
  networkifacename: "eth0"
```

This is how it breaks down

* `helper.name` - *REQUIRED*: This needs to be set to the hostname you want your helper to be (some people leave it as "helper" others change it to "bastion")
* `helper.ipaddr` - *REQUIRED* Set this to the current IP address of the helper. In case of high availability cluster, set this to the virtual IP address of the helpernodes. This is used to set up the [reverse dns definition](../templates/named.conf.j2#L65)
* `helper.networkifacename` - *OPTIONAL*: By default the playbook uses `{{ ansible_default_ipv4.interface }}` for the interface of the helper or helpernodes (In case of high availability). This option can be set to override the interface used for the helper or helpernodes (if, for example, you're on a dual homed network or your helper has more than one interface).

**NOTE**: The `helper.networkifacename` is the ACTUAL name of the interface, NOT the NetworkManager name (you should _NEVER_ need to set it to something like `System eth0`. Set it to what you see in `ip addr`)


## DNS Section

This section sets up your DNS server.

```
dns:
  domain: "example.com"
  clusterid: "ocp4"
  forwarder1: "8.8.8.8"
  forwarder2: "8.8.4.4"
  lb_ipaddr: "{{ helper.ipaddr }}"
```

Explanation of the DNS variables:

* `dns.domain` - This is what domain the installed DNS server will have. This needs to match what you will put for the [baseDomain](examples/install-config-example.yaml#L2) inside the `install-config.yaml` file.
* `dns.clusterid` - This is what your clusterid will be named and needs to match what you will for [metadata.name](examples/install-config-example.yaml#L12) inside the `install-config.yaml` file.
* `dns.forwarder1` - This will be set up as the DNS forwarder. This is usually one of the corporate (or "upstream") DNS servers.
* `dns.forwarder2` - This will be set up as the second DNS forwarder. This is usually one of the corporate (or "upstream") DNS servers.
* `lb_ipaddr` - This is the load balancer IP, it is optional, the default value is `helper.ipaddr`.

The DNS server will be set up using `dns.clusterid` + `dns.domain` as the domain it's serving. In the above example, the helper will be setup to be the SOA for `ocp4.example.com`. The helper will also be setup as it's [own DNS server](../templates/resolv.conf.j2)

**NOTE**: Although you _CAN_ use the helper as your dns server. It's best to have your DNS server delegate the `dns.clusterid` + `dns.domain` domain to the helper (i.e. Delegate `ocp4.example.com` to the helper)

## DHCP Section

> **OPTIONAL** if doing static ips

This section sets up the DHCP server.

```
dhcp:
  router: "192.168.7.1"
  bcast: "192.168.7.255"
  netmask: "255.255.255.0"
  dns: "{{ helper.ipaddr }}"
  poolstart: "192.168.7.10"
  poolend: "192.168.7.30"
  ipid: "192.168.7.0"
  netmaskid: "255.255.255.0"
```

Explanation of the options you can set:

* `dhcp.router` - This is the default gateway of your network you're going to assign to the masters/workers
* `dhcp.bcast` - This is the broadcast address for your network
* `dhcp.netmask` - This is the netmask that gets assigned to your masters/workers
* `dhcp.dns` - This is the domain name server, it is optional, the default value is set to `helper.ipaddr`
* `dhcp.poolstart` - This is the first address in your dhcp address pool
* `dhcp.poolend` - This is the last address in your dhcp address pool
* `dhcp.ipid` - This is the ip network id for the range
* `dhcp.netmaskid` - This is the networkmask id for the range.

These variables are used to set up the [dhcp config file](../templates/dhcpd.conf.j2)

## Bootstrap Node Section

This section defines the bootstrap node configuration

```
bootstrap:
  name: "bootstrap"
  ipaddr: "192.168.7.20"
  macaddr: "52:54:00:60:72:67"
```

The options are:

* `bootstrap.name` - The hostname (**__WITHOUT__** the fqdn) of the bootstrap node you want to set
* `bootstrap.ipaddr` - The IP address that you want set (this modifies the [dhcp config file](../templates/dhcpd.conf.j2#L17), the [dns zonefile](../templates/zonefile.j2#L26), and the [reverse dns zonefile](../templates/reverse.j2#L15))
* `bootstrap.macaddr` - The mac address for [dhcp reservation](../templates/dhcpd.conf.j2#L17). This option is not needed if you're doing static ips.


## Master Node section

Similar to the bootstrap section; this sets up master node configuration. Please note that this is an array.

```
masters:
  - name: "master0"
    ipaddr: "192.168.7.21"
    macaddr: "52:54:00:e7:9d:67"
  - name: "master1"
    ipaddr: "192.168.7.22"
    macaddr: "52:54:00:80:16:23"
  - name: "master2"
    ipaddr: "192.168.7.23"
    macaddr: "52:54:00:d5:1c:39"
```

* `masters.name` - The hostname (**__WITHOUT__** the fqdn) of the master node you want to set (x of 3).
* `masters.ipaddr` - The IP address (x of 3) that you want set (this modifies the [dhcp config file](../templates/dhcpd.conf.j2#L19), the [dns zonefile](../templates/zonefile.j2#L29), and the [reverse dns zonefile](../templates/reverse.j2#L11))
* `masters.macaddr` - The mac address for [dhcp reservation](../templates/dhcpd.conf.j2#L19). This option is not needed if you're doing static ips.


**NOTE**: 3 Masters are MANDATORY for installation of OpenShift 4

## Worker Node section

Similar to the master section; this sets up worker node configuration. Please note that this is an array.

> :rotating_light: This section is optional if you're installing a compact cluster

```
workers:
  - name: "worker0"
    ipaddr: "192.168.7.11"
    macaddr: "52:54:00:f4:26:a1"
  - name: "worker1"
    ipaddr: "192.168.7.12"
    macaddr: "52:54:00:82:90:00"
  - name: "worker2"
    ipaddr: "192.168.7.13"
    macaddr: "52:54:00:8e:10:34"
```

* `workers.name` - The hostname (**__WITHOUT__** the fqdn) of the worker node you want to set
* `workers.ipaddr` - The IP address that you want set (this modifies the [dhcp config file](../templates/dhcpd.conf.j2#L22), the [dns zonefile](../templates/zonefile.j2#L34), and the [reverse dns zonefile](../templates/reverse.j2#L20))
* `workers.macaddr` - The mac address for [dhcp reservation](../templates/dhcpd.conf.j2#L22). This option is not needed if you're doing static ips.


**NOTE**: At LEAST 2 workers is needed if you're installing a standard version of OpenShift 4

## Extra sections

Below are example of "extra" features beyond the default built-in vars that you can manipulate.

### Static IPs

In order to use static IPs, you'll need to pass `-e staticips=true` to your `ansible-playbook` command or add the following in your `vars.yaml` file

```
staticips: true
```

This effectively disables DHCP, TFTP, and PXE on the helper. This implicitly means that you will be doing an ISO/CD-ROM/USB install of RHCOS.

**NOTE**: The default setting is `staticips: false` which installs DHCP, TFTP, and PXE.

### Specifying Artifacts

You can have the helper deploy specific artifacts for a paticular version of OCP. Or, the nightly builds of OpenShift 4 or even OKD. Adding the following to your `vars.yaml` file will pull in the coresponding artifacts. Below is an example of pulling the `4.2.0-0.nightly-2019-09-16-114316` nightly build:

> :warning: note, you need to use the `ocp_bios` var for the rootfs image for 4.6+

```
ocp_bios: "https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/pre-release/latest/rhcos-42.80.20190828.2-metal-bios.raw.gz"
ocp_initramfs: "https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/pre-release/latest/rhcos-42.80.20190828.2-installer-initramfs.img"
ocp_install_kernel: "https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/pre-release/latest/rhcos-42.80.20190828.2-installer-kernel"
ocp_client: "https://mirror.openshift.com/pub/openshift-v4/clients/ocp-dev-preview/latest/openshift-client-linux-4.2.0-0.nightly-2019-09-16-114316.tar.gz"
ocp_installer: "https://mirror.openshift.com/pub/openshift-v4/clients/ocp-dev-preview/latest/openshift-install-linux-4.2.0-0.nightly-2019-09-16-114316.tar.gz"
```

To find the latest nighly build links:

* [Install Artifacts](https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/pre-release/latest/)
* [Client and Installer](https://mirror.openshift.com/pub/openshift-v4/clients/ocp-dev-preview/latest/)

You can also point these vars to files local to the helper. This is useful when doing a disconnected insall and you have "sneaker-netted" the artifacts over. For example:

```
ocp_bios: "file:///tmp/rhcos-42.80.20190828.2-metal-bios.raw.gz"
ocp_initramfs: "file:///tmp/rhcos-42.80.20190828.2-installer-initramfs.img"
ocp_install_kernel: "file:///tmp/rhcos-42.80.20190828.2-installer-kernel"
ocp_client: "file:///tmp/openshift-client-linux-4.2.0-0.nightly-2019-09-16-114316.tar.gz"
ocp_installer: "file:///tmp/openshift-install-linux-4.2.0-0.nightly-2019-09-16-114316.tar.gz"
```

The [default](../vars/main.yml#L4-L8) is to use the latest **stable** OpenShift 4 release.

> Also, you can point this to ANY apache server...not just the OpenShift 4 mirrors (Useful for disconnected installs)

### Filetranspiler 

Originally, [filetranspiler](https://github.com/ashcrow/filetranspiler) was used to write out static ip configurations. You no longer need to do this as you can just pass `ip=...` into the kernel parameters to get static IPs setup.

This tool can still be useful to write out other files. Therefore, you can install it by setting the following option:

```
install_filetranspiler: true
```

Default is set to `false` as to NOT install it.

### SSH Key

This playbook [creates an SSH key](../tasks/generate_ssh_keys.yaml) as `~/.ssh/helper_rsa` that can be used for the `install-config.yaml` file. It also creates an `~/.ssh/config` file to use this as your default key when sshing into the nodes.

```
ssh_gen_key: true
```

Default is set to `true`, set it to `false` if you don't want it to create the SSH KEY or config for you

### PXE default config

**OPTIONAL**  

This section influences the creation of pxe artifacts

```
pxe:
  generate_default: true
```

Default is false to prevent unexpected issues booting hosts in the "other" section.    
* `pxe.generate_default` - Setting to true Generates a generic default pxe config file with options for hosts not defined in the (bootstrap/master/worker) sections.  It is recommended to modify the template with appropriate boot options 
`templates/default.j2` -> `/var/lib/tftpboot/pxelinux.cfg/default`

### UEFI default config

**OPTIONAL**  

This section influences the creation of uefi artifacts for tftp boot.

```
uefi: false
```

Default is false to prevent unexpected issues.

### Other Nodes

**OPTIONAL**

If you want to have other DNS/DHCP entires managed by the helper, you can use `other` and specify the ip/mac address

```
other:
  - name: "non-cluster-vm"
    ipaddr: "192.168.7.31"
    macaddr: "52:54:00:f4:2e:2e"
```

You can omit `macaddr` if using `staticips=true`

### High Availability section

**OPTIONAL**

This section sets up the configuration for installing a high availability cluster. Please note that `high_availability.helpernodes` is an array.

```
high_availability:
  helpernodes:
    - name: "helper-0"
      ipaddr: "192.168.67.2"
      state: MASTER
      priority: 100
    - name: "helper-1"
      ipaddr: "192.168.67.3"
      state: BACKUP
      priority: 90
```

* `high_availability.helpernodes.name` - The hostname (**__WITHOUT__** the fqdn) of the helpernode you want to set
* `high_availability.helpernodes.ipaddr` - The IP address that you want to set (this modifies the [dns zonefile](../templates/zonefile.j2#L20))
* `high_availability.helpernodes.state` - The initial state of the helpernode that you want to set (MASTER|BACKUP). There must be exactly 1 helpernode with initial state as MASTER and atleast 1 helpernode with initial state as BACKUP.
* `high_availability.helpernodes.priority` - The priority of the helpernode that you want to set must be unique within range (0-100), and the helpernodes configured as state MASTER should have a priority value greater than the helpernodes configured as state BACKUP.

**NOTE**: Ensure you update `inventory` file appropriately to run the playbook on all the helpernodes. For more information refer [inventory doc](inventory-ha-doc.md).

### Local Registry

**OPTIONAL**

In order to install a local registry on the helper node:
* A pullsecret obtained at [try.openshift.com](https://cloud.redhat.com/openshift/install/pre-release) - Download and save it as `~/.openshift/pull-secret` on the helper node.
* you'll need to add the following in your `vars.yaml` file

```
setup_registry:
  deploy: false
  autosync_registry: false
  registry_image: docker.io/library/registry:2
  local_repo: "ocp4/openshift4"
  product_repo: "openshift-release-dev"
  release_name: "ocp-release"
  release_tag: "4.4.9-x86_64"
```

* `setup_registry.deploy` - Set this to true to enable registry installation.
* `setup_registry.autosync_registry` - Set this to true to enable mirroring of installation images.
* `setup_registry.registry_image` - This is the name of the image used for creating registry container.
* `setup_registry.local_repo` - This is the name of the repo in your registry.
* `setup_registry.product_repo` - Where the images are hosted in the product repo.
* `setup_registry.release_name` - This is the name of the image release.
* `setup_registry.release_tag` - The version of OpenShift you want to sync.

### Running on Power

In order to run the helper node on Power for deploying OCP on Power you'll need to pass `-e ppc64le=true` to your `ansible-playbook` command or add the following in your `vars.yaml` file

```
ppc64le: true
```

### NFS Configuration

This playbook sets up a script called [helpernodecheck](../templates/checker.sh.j2). That script can be used to set up the helpernode as an NFS server (this is the default). It uses the [k8s nfs setup in the incubator repo](https://github.com/kubernetes-incubator/external-storage/blob/master/nfs-client/README.md).

> :warning: Run `helpernodecheck nfs-info` for information after the playbook runs.

You can have that script  connect you to another NFS server, that's not the HelperNode. You do this by adding the following in your `vars.yaml` file.

```
nfs:
  server: "192.168.1.100"
  path: "/exports/helper"
```

* `nfs.server` - this is the ip address or the hostname of your nfs server
* `nfs.path` - this is the path the nfs controller is going to mount. **This path MUST exist with the proper permissions!**

### NTP Configuration

If you would like to use your own NTP servers, you can specify them in using the follwing config.

```
chronyconfig:
  enabled: true
  content:
    - server: 0.centos.pool.ntp.org
      options: iburst
    - server: 1.centos.pool.ntp.org
      options: iburst
```

* `chronyconfig.enabled` - This will flag the playbook that you want to setup chrony to use a specific config.
* `content` - This is an array of servers and their options. This eventually makes it's way to [the `chrony.conf` file](../templates/chrony.conf.j2). If you require further `options`, just put them in quotes: `options: "iburst foo bar"`

This playbook does NOT set up Chrony for you, instead it provides you with the `machineConfig` that you can load either pre-install or post-install. After the playbook has ran, you can do one of two things.

__Pre-Install__

When installing OpenShift, there is a step to create the manifests. Here is an example of creating the manifests under the `~/ocp4` directory.

```shell
openshift-install create manifests --dir=~/ocp4
```

This will create the `manifests` directory and the `openshift` directory under `~/ocp4`

```shell
# ll ~/ocp4
total 8
drwxr-x---. 2 root root 4096 Jul 16 08:08 manifests
drwxr-x---. 2 root root 4096 Jul 16 08:06 openshift
```

The playbook created the `machineConfig` files where you cloned the repo. For example, if I cloned the repo in my homedir; it'll be under `~/ocp4-helpernode/machineconfig/`.

```shell
# ll ~/ocp4-helpernode/machineconfig/
total 8
-rw-r--r--. 1 root root 748 Jul 16 07:59 99-master-chrony-configuration.yaml
-rw-r--r--. 1 root root 748 Jul 16 07:59 99-worker-chrony-configuration.yaml
```

Copy over these to the `openshift` directory in the installation directory.

```shell
cp ~/ocp4-helpernode/machineconfig/* ~/ocp4/openshift/
```

Continue on with the installation. Once done you should have chrony setup pointing to your NTP servers.

```shell
# oc get machineconfig 99-{master,worker}-chrony-configuration
NAME                             GENERATEDBYCONTROLLER   IGNITIONVERSION   AGE
99-master-chrony-configuration                           2.2.0             42s
99-worker-chrony-configuration                           2.2.0             43s
```

You can see the config if you login to one of your nodes and take a look at the file.

```shell
# oc debug node/worker1.ocp4.example.com
Starting pod/worker1ocp4examplecom-debug ...
To use host binaries, run `chroot /host`
Pod IP: 192.168.7.12
If you don't see a command prompt, try pressing enter.
sh-4.2# chroot /host
sh-4.4# cat /etc/chrony.conf 
server 0.centos.pool.ntp.org iburst
server 1.centos.pool.ntp.org iburst
driftfile /var/lib/chrony/drift
makestep 1.0 3
rtcsync
```

> :bulb: This config should be on all the masters and workers.

__Post-Install__

To set this up post-installation, just apply the `machineConfig` using `oc apply -f`. For example:

```shell
oc apply  -f ~/ocp4-helpernode/machineconfig/
```

:warning: This will reboot ALL your nodes (masters/workers) in a "rolling" fashion. You can check this with `oc get nodes`

```shell
# oc get nodes
NAME                       STATUS                     ROLES    AGE   VERSION
master0.ocp4.example.com   Ready                      master   41m   v1.17.1+912792b
master1.ocp4.example.com   Ready                      master   41m   v1.17.1+912792b
master2.ocp4.example.com   Ready,SchedulingDisabled   master   41m   v1.17.1+912792b
worker0.ocp4.example.com   Ready,SchedulingDisabled   worker   26m   v1.17.1+912792b
worker1.ocp4.example.com   Ready                      worker   26m   v1.17.1+912792b
```


This should create the `machineConfig`

```shell
# oc get machineconfig 99-{master,worker}-chrony-configuration
NAME                             GENERATEDBYCONTROLLER   IGNITIONVERSION   AGE
99-master-chrony-configuration                           2.2.0             42s
99-worker-chrony-configuration                           2.2.0             43s
```

# Example Vars file

Below are example `vars.yaml` files.

* [Default vars.yaml using DHCP](examples/vars.yaml)
* [Default vars.yaml using DHCP with Nightlies](examples/vars-nightlies.yaml)
* [Example of vars.yaml using Static IPs](examples/vars-static.yaml)
* [Example of vars.yaml using Static IPs with Nightlies](examples/vars-static-nightlies.yaml)
* [Example of vars.yaml for Power](examples/vars-ppc64le.yaml)
* [Example of vars.yaml DHCP and External NFS](examples/vars-nfs.yaml)
* [Example of vars.yaml with Chrony configuration](examples/vars-chrony.yaml)
* [Example of vars.yaml setting up a Compact Cluster](examples/vars-compact.yaml)
* [Example of vars.yaml setting up a Compact Cluster with Static IPs](examples/vars-compact-static.yaml)
