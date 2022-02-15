package node

import (
	"encoding/json"
	"github.com/a-shine/butter/utils"
	"time"
)

// only need to do this when known host list is at capacity

type HostQuality struct {
	uptime           time.Duration
	availableStorage uint64
	knownHosts       int
}

func hostQuality(overlay Overlay, _ []byte) []byte {
	var hostQuality HostQuality

	node := overlay.Node()

	hostQuality.uptime = node.uptime()
	hostQuality.availableStorage = overlay.AvailableStorage()
	hostQuality.knownHosts = len(node.KnownHosts())

	json, _ := json.Marshal(hostQuality)

	return json

}

func AppendHostQualityServerBehaviour(node *Node) {
	node.RegisterServerBehaviour("host-quality/", hostQuality)
}

type knownHosts struct {
	state []utils.SocketAddr
}

func (node *Node) intelligentAddKnownHost(potentialHost utils.SocketAddr) {
	var uptimeTally uint64
	var storageTally uint64

	var highUptimeHighStorage []utils.SocketAddr
	var highUptimeLowStorage []utils.SocketAddr
	var lowUptimeHighStorage []utils.SocketAddr
	var lowUptimeLowStorage []utils.SocketAddr

	var knownHostsQuality map[utils.SocketAddr]HostQuality{}

	for _, host := range node.KnownHosts() {
		// get the host quality and add to a map of known hosts to their quality
		var hostQuality HostQuality
		response, _ := utils.Request(host, []byte("host-quality/"), []byte{})
		_ = json.Unmarshal(response, &hostQuality)
		knownHostsQuality[host] = hostQuality

		// tally the uptime to obtain mean avg
		// tally the storage to obtain mean avg
		uptimeTally += uint64(hostQuality.uptime)
		storageTally += hostQuality.availableStorage
	}

	avgUptime := uptimeTally / uint64(len(node.KnownHosts()))
	avgStorage := storageTally / uint64(len(node.KnownHosts()))

	// classify hosts into 4 categories - add then to appropriate list
	for _, host := range node.KnownHosts() {
		if knownHostsQuality[host].uptime > avgUptime && knownHostsQuality[host].availableStorage > avgStorage {
			highUptimeHighStorage = append(highUptimeHighStorage, host)
		} else if knownHostsQuality[host].uptime > avgUptime && knownHostsQuality[host].availableStorage < avgStorage {
			highUptimeLowStorage = append(highUptimeLowStorage, host)
		} else if knownHostsQuality[host].uptime < avgUptime && knownHostsQuality[host].availableStorage > avgStorage {
			lowUptimeHighStorage = append(lowUptimeHighStorage, host)
		} else if knownHostsQuality[host].uptime < avgUptime && knownHostsQuality[host].availableStorage < avgStorage {
			lowUptimeLowStorage = append(lowUptimeLowStorage, host)
		}
	}


	// see where the new knwown host lies in these 4 categories (by repeating the above)

	//figure out the host uotime and storage distribution
	//if the new host does not improve the distribution don't add - i.e. do nothing
	//else remove a known host in the category with teh highesr frequency and add in the known host to the category with the least frequency
}

func (k *knownHosts) Copy() interface{} {
	copiedState := make([]utils.SocketAddr, len(k.state))
	copy(copiedState, k.state)
	return &knownHosts{state: copiedState}
}

func (k *knownHosts) Move() {
	// take the host from possible host
	// add it to state, remove another host from the state and set it as the possibleHost
}

// Value function for the list of known hosts - needs to be diverse (minimisae)
func (k *knownHosts) Energy() float64 {
	var energy float64

	// 4 possible types of host - (high uptime, high storage), (high uptime, low storage), (low uptime, high storage), (low uptime, low storage)
	// The value function switches between these 4 states as to what it values
	valueHighUptime := true
	valueHighStorage := true

	// for each host in the state, get the uptime and available storage
	for _, host := range k.state {
		var hostQuality HostQuality
		response, _ := utils.Request(host, []byte("host-quality/"), []byte{})
		_ = json.Unmarshal(response, &hostQuality)
		if valueHighUptime && valueHighStorage {
			energy -= float64(hostQuality.uptime.Milliseconds()) // we like high uptime (the bigger the better)
			energy -= float64(hostQuality.availableStorage)      // we like lots of available storage
			valueHighUptime = true
			valueHighStorage = false // move to next state
		} else if valueHighUptime && !valueHighStorage {
			energy -= float64(hostQuality.uptime.Milliseconds()) // we dislike high uptime (the bigger the better)
			energy += float64(hostQuality.availableStorage)      // we dislike lots of available storage
			valueHighUptime = false
			valueHighStorage = true // move to next state
		} else if !valueHighUptime && valueHighStorage {
			energy += float64(hostQuality.uptime.Milliseconds()) // we like high uptime (the bigger the better)
			energy -= float64(hostQuality.availableStorage)      // we like lots of available storage
			valueHighUptime = false
			valueHighStorage = false // move to next state
		} else if !valueHighUptime && !valueHighStorage {
			energy += float64(hostQuality.uptime.Milliseconds()) // we dislike high uptime (the bigger the better)
			energy += float64(hostQuality.availableStorage)      // we dislike lots of available storage
			valueHighUptime = true
			valueHighStorage = true // move to next state
		}
	}
	return energy
}
