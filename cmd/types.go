package cmd

//TODO finalize version
const VERSION string = "latest"

var coreImageNames = []string{"dns", "dhcp", "http", "loadbalancer", "pxe"}
var containerRuntime string = "podman"
var repository string = "helpernode"
var registry string = "quay.io"
var logLevel string

//TODO figure out how to not have this as a global var
var imageList []string //this is used in start/stop

var portlist = map[string][]string{
	"53": {"udp", "tcp"},
	"67": {"udp"},
	"69": {"udp"},
	"80": {"tcp"},
	"443": {"tcp"},
	"546": {"udp"},
	"6443": {"tcp"},
	"8080": {"tcp"},
	"9000": {"tcp"},
	"9090": {"tcp"},
	"22623": {"tcp"},
}

// Define ports needed for preflight check
var ports = []string{
	"67",
	"546",
	"53",
	"80",
	"443",
	"69",
	"6443",
	"22623",
	"8080",
	"9000",
}

//Default images
//TODO Add disconnected
//TODO Add an image struct later...too many code changes for it now
var images = make(map[string]string)

//TODO lets build something that builds the image list. should have add/remove functions
var coreImages = map[string]string{
	"dns":          "quay.io/helpernode/dns",
	"dhcp":         "quay.io/helpernode/dhcp",
	"http":         "quay.io/helpernode/http",
	"loadbalancer": "quay.io/helpernode/loadbalancer",
	"pxe":          "quay.io/helpernode/pxe",
}

// Define systemd services we will check
var systemdsvc = map[string]string{
	"resolved": "systemd-resolved.service",
	"dnsmasq":  "dnsmasq.service",
}

var fwrule = []string{
	"6443/tcp",
	"22623/tcp",
	"8080/tcp",
	"9000/tcp",
	"9090/tcp",
	"67/udp",
	"546/udp",
	"53/tcp",
	"53/udp",
	"80/tcp",
	"443/tcp",
	"22/tcp",
	"69/udp",
}

var clients = map[string]string{
	"oc":                "openshift-client-linux.tar.gz",
	"openshift-install": "openshift-install-linux.tar.gz",
	"helm":              "helm.tar.gz",
}
