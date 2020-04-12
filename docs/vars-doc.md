# How to use vars.yaml

This page gives you an explanation of the variables found in the [vars.yaml](examples/vars.yaml) example given in this repo to help you formulate your own/edit the provided example.

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

* `helper.name` - This needs to be set to the hostname you want your helper to be (some people leave it as "helper" others change it to "bastion")
* `helper.ipaddr` - Set this to the current IP address of the helper. This is used to set up the [reverse dns definition](../templates/named.conf.j2#L65)
* `helper.networkifacename` - This is set to the network interface of the helper (what you see when you do `ip addr`)

**NOTE**: The `helper.networkifacename` is the ACTUAL name of the interface, NOT the NetworkManager name (you should _NEVER_ need to set it to something like `System eth0`. Set it to what you see in `ip addr`)


## DNS Section

This section sets up your DNS server.

```
dns:
  domain: "example.com"
  clusterid: "ocp4"
  forwarder1: "8.8.8.8"
  forwarder2: "8.8.4.4"
```

Explanation of the DNS variables:

* `dns.domain` - This is what domain the installed DNS server will have. This needs to match what you will put for the [baseDomain](examples/install-config-example.yaml#L2) inside the `install-config.yaml` file.
* `dns.clsuterid` - This is what your clusterid will be named and needs to match what you will for [metadata.name](examples/install-config-example.yaml#L12) inside the `install-config.yaml` file.
* `dns.forwarder1` - Tis will be set up as the DNS forwarder. This is usually one of the corprate (or "upstream") DNS servers.
* `dns.forwarder2` - Tis will be set up as the second DNS forwarder. This is usually one of the corprate (or "upstream") DNS servers.

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
  poolstart: "192.168.7.10"
  poolend: "192.168.7.30"
  ipid: "192.168.7.0"
  netmaskid: "255.255.255.0"
```

Explanation of the options you can set:

* `dhcp.router` - This is the default gateway of your network you're going to assign to the masters/workers
* `dhcp.bcast` - This is the broadcast address for your network
* `dhcp.netmask` - This is the netmask that gets assigned to your masters/workers
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


**NOTE**: At LEAST 1 worker is needed for installation of OpenShift 4

## Extra sections

Below are example of "extra" features beyond the default built-in vars that you can manipulate.

### Static IPs

In order to use static IPs, you'll need to pass `-e staticips=true` to your `ansible-playbook` command or add the following in your `vars.yaml` file

```
staticips: true
```

This effectively disables DHCP, TFTP, and PXE on the helper. This implicitly means that you will be doing an ISO/CD-ROM/USB install of RHCOS.

**NOTE**: The default setting is `staticips: false` which installs DHCP, TFTP, and PXE.

### Nightly Builds

You can have the helper deploy the nightly builds of OpenShift 4. Adding the following to your `vars.yaml` files will pull in the coresponding artifacts. Below is an example of pulling the `4.2.0-0.nightly-2019-09-16-114316` nightly

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

The [default](../vars/main.yml#L4-L8) is to use the latest stable OpenShift 4 release

> Also, you can point this to ANY apache server...not just the OpenShift 4 mirrors (*cough* *cough* disconnected hint here *cough* *cough*)

### Filetranspiler 

Originally, [filetranspiler](https://github.com/ashcrow/filetranspiler) was used to write out static ip configurations. You no longer need to do this as you can just pass `ip=...` into the kernel parameters to get static IPs setup.

This tool can still be useful to write out other files. Therefore, you can install it by setting the following option:

```
install_filetranspiler: true
```

Default is set to `false` as to NOT install it.

### SSH Key

This playbook [creates an SSH key](../tasks/generate_ssh_keys.yaml) as `~/.ssk/helper_rsa` that can be used for the `install-config.yaml` file. It also creates an `~/.ssh/config` file to use this as your default key when sshing into the nodes.

```
ssh_gen_key: true
```

Default is set to `true`, set it to `false` if you don't want it to create the SSH KEY or config for you

# Example Vars file

Below are example `vars.yaml` files.

* [Default vars.yaml using DHCP](examples/vars.yaml)
* [Default vars.yaml using DHCP with Nightlies](examples/vars-nightlies.yaml)
* [Example of vars.yaml using Static IPs](examples/vars-static.yaml)
* [Example of vars.yaml using Static IPs with Nightlies](examples/vars-static-nightlies.yaml)
