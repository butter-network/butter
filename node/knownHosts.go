package node

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/a-shine/butter/utils"
)

// IDEA: A different protocol for known host management could be to just accept a known host until memory capacity is full
// and then have a certain probability of a new known host being added and random old one being dropped.

type HostQuality struct {
	Uptime           uint64
	AvailableStorage uint64
	NbHostsKnown     uint64
}

type KnownHosts struct {
	mu  sync.Mutex
	cap uint

	uptimeTally     uint64
	storageTally    uint64
	knownHostsTally uint64

	Hosts map[utils.SocketAddr]HostQuality

	highUptimeHighStorageHighKH uint
	highUptimeHighStorageLowKH  uint
	highUptimeLowStorageHighKH  uint
	highUptimeLowStorageLowKH   uint
	lowUptimeHighStorageHighKH  uint
	lowUptimeHighStorageLowKH   uint
	lowUptimeLowStorageHighKH   uint
	lowUptimeLowStorageLowKH    uint

	lastUpdate time.Time
}

func (knownHosts *KnownHosts) count() uint {
	return uint(len(knownHosts.Hosts))
}

func hostQuality(overlay Overlay, _ []byte) []byte {
	var hostQuality HostQuality

	node := overlay.Node()

	hostQuality.Uptime = uint64(node.uptime())
	hostQuality.AvailableStorage = overlay.AvailableStorage()
	hostQuality.NbHostsKnown = uint64(len(node.KnownHosts()))

	json, _ := json.Marshal(hostQuality)

	return json

}

func (knownHosts *KnownHosts) Addrs() map[utils.SocketAddr]HostQuality {
	return knownHosts.Hosts
}

func AppendHostQualityServerBehaviour(node *Node) {
	node.RegisterServerBehaviour("host-quality/", hostQuality)
}

func (knownHosts *KnownHosts) update() {
	knownHosts.mu.Lock()
	defer knownHosts.mu.Unlock()

	knownHosts.uptimeTally = 0
	knownHosts.storageTally = 0
	knownHosts.knownHostsTally = 0

	for host, _ := range knownHosts.Hosts {
		// get the host quality and add to a map of known hosts to their quality
		var hostQuality HostQuality
		response, err := utils.Request(host, []byte("host-quality/"), []byte{})
		if err != nil {
			knownHosts.Remove(host)
			continue
		}
		_ = json.Unmarshal(response, &hostQuality)
		knownHosts.Hosts[host] = hostQuality

		// tally the Uptime to obtain mean avg
		// tally the storage to obtain mean avg
		knownHosts.uptimeTally += uint64(hostQuality.Uptime)
		knownHosts.storageTally += hostQuality.AvailableStorage
		knownHosts.knownHostsTally += hostQuality.NbHostsKnown
	}
	knownHosts.classifyHosts()
	knownHosts.lastUpdate = time.Now()
}

func (knownHosts *KnownHosts) Remove(host utils.SocketAddr) {
	knownHosts.mu.Lock()
	defer knownHosts.mu.Unlock()

	hostQuality := knownHosts.Hosts[host]

	avgUptime := knownHosts.avgUptime()
	avgStorage := knownHosts.avgStorage()
	avgKH := knownHosts.avgKnownHosts()

	knownHosts.decrementHostClass(hostQuality, avgUptime, avgStorage, avgKH)

	delete(knownHosts.Hosts, host)
}

// Also prioritise adding a host that is not known by anyone else - don't want the scenario where 3 nodes are at capacity and a new node hosts and can be added to the network
func (knownHosts *KnownHosts) Add(host utils.SocketAddr) {
	knownHosts.mu.Lock()
	defer knownHosts.mu.Unlock()

	if _, ok := knownHosts.Hosts[host]; ok {
		return // already in the list
	}

	var hostQuality HostQuality
	response, _ := utils.Request(host, []byte("host-quality/"), []byte{})
	_ = json.Unmarshal(response, &hostQuality)

	// TODO: increment the correct category
	avgUptime := knownHosts.avgUptime()
	avgStorage := knownHosts.avgStorage()
	avgKnownHosts := knownHosts.avgKnownHosts()

	if knownHosts.count() < knownHosts.cap {
		knownHosts.Hosts[host] = hostQuality // if we have the memory just add the known host
		knownHosts.incrementHostClass(hostQuality, avgUptime, avgStorage, avgKnownHosts)
	} else {
		// only need to do this when known host list is at capacity
		knownHosts.intelligentAddKnownHost(host) // else figure out if its worth adding and who to remove
	}
}

func (knownHosts *KnownHosts) avgUptime() uint64 {
	if knownHosts.count() == 0 {
		return 0
	}
	return knownHosts.uptimeTally / uint64(len(knownHosts.Hosts))
}

func (knownHosts *KnownHosts) avgStorage() uint64 {
	if knownHosts.count() == 0 {
		return 0
	}
	return knownHosts.storageTally / uint64(len(knownHosts.Hosts))
}

func (knownHosts *KnownHosts) avgKnownHosts() uint64 {
	if knownHosts.count() == 0 {
		return 0
	}
	return knownHosts.knownHostsTally / uint64(len(knownHosts.Hosts))
}

func (knownHosts *KnownHosts) intelligentAddKnownHost(potentialHost utils.SocketAddr) {
	// OPTIMISATION - these can be cached as part of a KnownHosts struct, i.e. if we recently ran the above functions, we can assume that the distribution has not changed much and just skip to making the decision as to if we should add the known host or not

	//if time.Since(knownHosts.lastUpdate) > time.Minute*2 {
	//	knownHosts.update()
	//}

	avgUptime := knownHosts.avgUptime()
	avgStorage := knownHosts.avgStorage()
	avgKnownHosts := knownHosts.avgKnownHosts()

	// see where the new known host lies in these 4 categories (by repeating the above)
	// figure out the host Uptime and storage distribution
	var hostQuality HostQuality
	response, _ := utils.Request(potentialHost, []byte("host-quality/"), []byte{})
	_ = json.Unmarshal(response, &hostQuality)
	newHostType := hostQuality.hostType(avgUptime, avgStorage, avgKnownHosts)

	//if the new host does not improve the distribution don't add - i.e. do nothing
	if newHostType == knownHosts.biggestClass() {
		return
	} else {
		//remove host from the biggest class
		knownHosts.removeFromBiggestClass()
		// and add the new host to his class (smallest class)
		knownHosts.Add(potentialHost)
	}
	//else remove a known host in the category with teh highesr frequency and add in the known host to the category with the least frequency
}

func (knownHosts *KnownHosts) removeFromBiggestClass() {
	biggestClass := knownHosts.biggestClass()
	avgUptime := knownHosts.avgUptime()
	avgStorage := knownHosts.avgStorage()
	avgKnownHosts := knownHosts.avgKnownHosts()
	for host, hostQuality := range knownHosts.Hosts {
		if hostQuality.hostType(avgUptime, avgStorage, avgKnownHosts) == biggestClass {
			knownHosts.Remove(host)
			return
		}
	}
}

func (knownHosts *KnownHosts) biggestClass() string {
	// return the biggest class name
	classes := make(map[string]uint)
	classes["highUptimeHighStorageHighKH"] = knownHosts.highUptimeHighStorageHighKH
	classes["highUptimeHighStorageLowKH"] = knownHosts.highUptimeHighStorageLowKH
	classes["highUptimeLowStorageHighKH"] = knownHosts.highUptimeLowStorageHighKH
	classes["highUptimeLowStorageLowKH"] = knownHosts.highUptimeLowStorageLowKH
	classes["lowUptimeHighStorageHighKH"] = knownHosts.lowUptimeHighStorageHighKH
	classes["lowUptimeHighStorageLowKH"] = knownHosts.lowUptimeHighStorageLowKH
	classes["lowUptimeLowStorageHighKH"] = knownHosts.lowUptimeLowStorageHighKH
	classes["lowUptimeLowStorageLowKH"] = knownHosts.lowUptimeLowStorageLowKH
	var biggestClass string
	var biggestClassCount uint
	for class, count := range classes {
		if count > biggestClassCount {
			biggestClass = class
			biggestClassCount = count
		}
	}
	return biggestClass
}

func (knownHosts *KnownHosts) classifyHosts() {
	knownHosts.resetDistribution()

	avgUptime := knownHosts.avgUptime()
	avgStorage := knownHosts.avgStorage()
	avgKnownHosts := knownHosts.avgKnownHosts()

	for _, quality := range knownHosts.Hosts {
		knownHosts.incrementHostClass(quality, avgUptime, avgStorage, avgKnownHosts)
	}
}

func (knownHosts *KnownHosts) resetDistribution() {
	knownHosts.highUptimeHighStorageHighKH = 0
	knownHosts.highUptimeHighStorageLowKH = 0
	knownHosts.highUptimeLowStorageHighKH = 0
	knownHosts.highUptimeLowStorageLowKH = 0
	knownHosts.lowUptimeHighStorageHighKH = 0
	knownHosts.lowUptimeHighStorageLowKH = 0
	knownHosts.lowUptimeLowStorageHighKH = 0
	knownHosts.lowUptimeLowStorageLowKH = 0
}

func (knownHosts *KnownHosts) incrementHostClass(hostQuality HostQuality, avgUptime uint64, avgStorage uint64, avgKnownHosts uint64) {
	switch hostQuality.hostType(avgUptime, avgStorage, avgKnownHosts) {
	case "highUptimeHighStorageHighKH":
		knownHosts.highUptimeHighStorageHighKH++
	case "highUptimeHighStorageLowKH":
		knownHosts.highUptimeHighStorageLowKH++
	case "highUptimeLowStorageHighKH":
		knownHosts.highUptimeLowStorageHighKH++
	case "highUptimeLowStorageLowKH":
		knownHosts.highUptimeLowStorageLowKH++
	case "lowUptimeHighStorageHighKH":
		knownHosts.lowUptimeHighStorageHighKH++
	case "lowUptimeHighStorageLowKH":
		knownHosts.lowUptimeHighStorageLowKH++
	case "lowUptimeLowStorageHighKH":
		knownHosts.lowUptimeLowStorageHighKH++
	case "lowUptimeLowStorageLowKH":
		knownHosts.lowUptimeLowStorageLowKH++
	}
}

func (knownHosts *KnownHosts) decrementHostClass(hostQuality HostQuality, avgUptime uint64, avgStorage uint64, avgKnownHosts uint64) {
	switch hostQuality.hostType(avgUptime, avgStorage, avgKnownHosts) {
	case "highUptimeHighStorageHighKH":
		knownHosts.highUptimeHighStorageHighKH--
	case "highUptimeHighStorageLowKH":
		knownHosts.highUptimeHighStorageLowKH--
	case "highUptimeLowStorageHighKH":
		knownHosts.highUptimeLowStorageHighKH--
	case "highUptimeLowStorageLowKH":
		knownHosts.highUptimeLowStorageLowKH--
	case "lowUptimeHighStorageHighKH":
		knownHosts.lowUptimeHighStorageHighKH--
	case "lowUptimeHighStorageLowKH":
		knownHosts.lowUptimeHighStorageLowKH--
	case "lowUptimeLowStorageHighKH":
		knownHosts.lowUptimeLowStorageHighKH--
	case "lowUptimeLowStorageLowKH":
		knownHosts.lowUptimeLowStorageLowKH--
	}
}

func (hostQuality *HostQuality) hostType(avgUptime uint64, avgStorage uint64, avgKnownHosts uint64) string {
	if hostQuality.Uptime >= avgUptime && hostQuality.AvailableStorage >= avgStorage && hostQuality.NbHostsKnown >= avgKnownHosts {
		return "highUptimeHighStorageHighKH"
	} else if hostQuality.Uptime >= avgUptime && hostQuality.AvailableStorage >= avgStorage && hostQuality.NbHostsKnown < avgKnownHosts {
		return "highUptimeHighStorageLowKH"
	} else if hostQuality.Uptime >= avgUptime && hostQuality.AvailableStorage < avgStorage && hostQuality.NbHostsKnown >= avgKnownHosts {
		return "highUptimeLowStorageHighKH"
	} else if hostQuality.Uptime >= avgUptime && hostQuality.AvailableStorage < avgStorage && hostQuality.NbHostsKnown < avgKnownHosts {
		return "highUptimeLowStorageLowKH"
	} else if hostQuality.Uptime < avgUptime && hostQuality.AvailableStorage >= avgStorage && hostQuality.NbHostsKnown >= avgKnownHosts {
		return "lowUptimeHighStorageHighKH"
	} else if hostQuality.Uptime < avgUptime && hostQuality.AvailableStorage >= avgStorage && hostQuality.NbHostsKnown < avgKnownHosts {
		return "lowUptimeHighStorageLowKH"
	} else if hostQuality.Uptime < avgUptime && hostQuality.AvailableStorage < avgStorage && hostQuality.NbHostsKnown >= avgKnownHosts {
		return "lowUptimeLowStorageHighKH"
	} else if hostQuality.Uptime < avgUptime && hostQuality.AvailableStorage < avgStorage && hostQuality.NbHostsKnown < avgKnownHosts {
		return "lowUptimeLowStorageLowKH"
	} else {
		return "unknown"
	}
}

func (knownHosts *KnownHosts) JsonDigest() []byte {
	digest, _ := json.Marshal(knownHosts)
	return digest
}
