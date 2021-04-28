# Pluggable Services

You are able to launch other containers that are not part of the "core"
container group. This can be done by adding `pluggableService` to your
config. Below is an example:

```yaml
pluggableServices:
  nfs:
    image: docker.io/itsthenetwork/nfs-server-alpine
    ports:
      - 2049/tcp
      - 2049/udp
    startupOptions: "--privileged -v /myshare:/nfsshare:Z -e SHARED_DIRECTORY=/nfsshare"
  myweb:
    image: quay.io/christianh814/webserver-test:latest
    ports:
      - 8888/tcp
    startupOptions: "--label=foo=bar"
```

You will need:

* `pluggableServices.<name>` - Name your service something. You must name it something not in the "reserved names". For example you **CANNOT** name your service `loadbalancer` as that's reserved.
* `pluggableServices.image` - This is the registry where the image is stored. You must pass the tag as well.
* `pluggableServices.ports` - This is an array of ports with their protocol.
* `pluggableServices.startupOptions` - These are options you want to pass to the startup process (this is optional)

Note in the example I'm using an nfs image that binds to
local storage. This is only an EXAMPLE, and you should [read the documentaion](https://hub.docker.com/r/itsthenetwork/nfs-server-alpine/)
about the image to formulate your own. If you're binding to local storage,
you have to make sure that path exists first. Other nuances will exist.

> **NOTE** This is a "garbage in/garbage out" situation. We do not validate what is passed via `startupOptions`
