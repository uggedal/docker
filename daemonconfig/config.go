package daemonconfig

import (
	"github.com/dotcloud/docker/engine"
	"github.com/dotcloud/docker/runtime/networkdriver"
	"net"
)

const (
	defaultNetworkMtu    = 1500
	DisableNetworkBridge = "none"
)

// FIXME: separate runtime configuration from http api configuration
type Config struct {
	Pidfile                     string
	Root                        string
	AutoRestart                 bool
	Dns                         []string
	DnsSearch                   []string
	EnableIptables              bool
	EnableIpForward             bool
	DefaultIp                   net.IP
	BridgeIface                 string
	BridgeIP                    string
	InterContainerCommunication bool
	GraphDriver                 string
	ExecDriver                  string
	Mtu                         int
	DisableNetwork              bool
	EnableSelinuxSupport        bool
}

// ConfigFromJob creates and returns a new DaemonConfig object
// by parsing the contents of a job's environment.
func ConfigFromJob(job *engine.Job) *Config {
	config := &Config{
		Pidfile:                     job.Getenv("Pidfile"),
		Root:                        job.Getenv("Root"),
		AutoRestart:                 job.GetenvBool("AutoRestart"),
		EnableIptables:              job.GetenvBool("EnableIptables"),
		EnableIpForward:             job.GetenvBool("EnableIpForward"),
		BridgeIP:                    job.Getenv("BridgeIP"),
		BridgeIface:                 job.Getenv("BridgeIface"),
		DefaultIp:                   net.ParseIP(job.Getenv("DefaultIp")),
		InterContainerCommunication: job.GetenvBool("InterContainerCommunication"),
		GraphDriver:                 job.Getenv("GraphDriver"),
		ExecDriver:                  job.Getenv("ExecDriver"),
		EnableSelinuxSupport:        false, // FIXME: hardcoded default to disable selinux for .10 release
	}
	if dns := job.GetenvList("Dns"); dns != nil {
		config.Dns = dns
	}
	if dnsSearch := job.GetenvList("DnsSearch"); dnsSearch != nil {
		config.DnsSearch = dnsSearch
	}
	if mtu := job.GetenvInt("Mtu"); mtu != 0 {
		config.Mtu = mtu
	} else {
		config.Mtu = GetDefaultNetworkMtu()
	}
	config.DisableNetwork = config.BridgeIface == DisableNetworkBridge

	return config
}

func GetDefaultNetworkMtu() int {
	if iface, err := networkdriver.GetDefaultRouteIface(); err == nil {
		return iface.MTU
	}
	return defaultNetworkMtu
}
