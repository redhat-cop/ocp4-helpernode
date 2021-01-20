# Disconnected Install

This quickstart will get the HelperNode running in a disconnected environment. All the HelperNode services are ran in prebuilt image containers, so getting it running disconnected can be broken down into the following steps.

* Pull the images
* Save the images into a tarball
* Sneakernet this tarball over to disconnected network.
* Load the tarball into a host.
* Push these images into your hosted registry.
* Start the HelperNode Services using the internal registry

> Note: Setting up a regsitry is beyond the scope of this document. Please
> see [Quay](https://www.projectquay.io/) or [Harbor](https://goharbor.io/)
> if you need a registry.

You can use the images hosted on your registry for the HelperNode

# Pull The Images

On your laptop, or a system with access to the internet, pull the
images. Make sure you're pulling the right tag for the version
you want.

> For the alpha build, we are using `latest`, but this will change to
> specific versioning in the future.

```shell
export TAG=latest
for hni in dns dhcp http loadbalancer pxe
do
  podman pull quay.io/helpernode/${hni}:${TAG}
done
```

# Save The Images

After you've pulled the images locally, save them into tarballs,
compressing them if you'd like.

```shell
export TAG=latest
for hni in dns dhcp http loadbalancer pxe
do
  podman pull quay.io/helpernode/${hni}:${TAG}
  podman save --compress -o helpernode-${hni}.tar.gz quay.io/helpernode/${hni}:${TAG} 
done
```

You can now "sneakernet" these `helpernode-${hni}.tar.gz` files to you
Disconnected environment and load it up to a host with access to your
image registry.

# Loading Tarball

Load the tarballs into a host on your network with the
ability to push into your image registry.

```shell
for hni in dns dhcp http loadbalancer pxe
do
  podman load -i helpernode-${hni}.tar.gz
done
```

You should now have the HelperNode service images on your host.

```shell
podman images |grep helper
```

The output should look something like this.

```
quay.io/helpernode/pxe            latest   26e2169c7c62   5 hours ago    696 MB
quay.io/helpernode/loadbalancer   latest   600c36a4cdc3   5 hours ago    588 MB
quay.io/helpernode/http           latest   a04dcc414ca6   5 hours ago    1.53 GB
quay.io/helpernode/dns            latest   1727361b9dfa   5 hours ago    666 MB
quay.io/helpernode/dhcp           latest   b651c044b62e   5 hours ago    592 MB
```

# Pushing Images To Registry

Now that you have your images locally, you can load them into your images
registry. You may need to login first, please consult with your image
registry admin.

In my example, my registry is `registry.example.com` on port `5000`

```shell
# podman login registry.example.com:5000
Username: reguser
Password: *********
Login Succeeded!
```

Now I can tag and push my local images to the registry.

> Note that I'm using a namespace called `alpha`. You may or maynot have
> a workspace, if you do; it will probably differ in name.

```shell
export TAG=latest
for hni in dns dhcp http loadbalancer pxe
do
  podman tag quay.io/helpernode/${hni}:${TAG}  registry.example.com:5000/alpha/helpernode/${hni}:${TAG}
  podman push registry.example.com:5000/alpha/helpernode/${hni}:${TAG}
done
```

# Starting HelperNode Disconnected

Once your images are uploaded, you can now start the HelperNode services
by first exporting the `HELPERNODE_IMAGE_PREFIX` environment variable
on the HelperNode.

```shell
export HELPERNODE_IMAGE_PREFIX=registry.example.com:5000/alpha
```

Next, you may need to login to your registry.

```shell
# podman login registry.example.com:5000
Username: reguser
Password: *********
Login Succeeded!
```

Now you can start the service, using a [valid YAML config file](yaml-file-forward.md) for
your environment.

First save the file

```shell
helpernodectl save -f helpernode.yaml
```

Then start your service

```shell
helpernodectl start
```

You should now be running the images from your local registry.

```shell
helpernodectl status
```

The output should look something like this.

```shell
Names                     Status                  Image
helpernode-http           Up About a minute ago   registry.example.com:5000/alpha/helpernode/http:latest
helpernode-dhcp           Up About a minute ago   registry.example.com:5000/alpha/helpernode/dhcp:latest
helpernode-dns            Up About a minute ago   registry.example.com:5000/alpha/helpernode/dns:latest
helpernode-pxe            Up 2 minutes ago        registry.example.com:5000/alpha/helpernode/pxe:latest
helpernode-loadbalancer   Up 2 minutes ago        registry.example.com:5000/alpha/helpernode/loadbalancer:latest
```

# OpenShift Installation

You can now use the HelperNode to install [OpenShift disconnected](https://docs.openshift.com/container-platform/4.6/installing/installing_bare_metal/installing-restricted-networks-bare-metal.html)
