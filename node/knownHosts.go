package node

import (
	"encoding/json"
	"log"
	"time"

	"github.com/a-shine/butter/utils"
)

// only need to do this when known host list is at capacity

type HostQuality struct {
	// TODO: add a sense of public vs private (known public allows you to be a
	// bridge between subnets but private has lower latency)
	Uptime           uint64
	AvailableStorage uint64
	nbKnownHosts     int
}

type KnownHosts struct {
	cap uint

	uptimeTally  uint64
	storageTally uint64

	Hosts map[utils.SocketAddr]HostQuality

	highUptimeHighStorage uint
	highUptimeLowStorage  uint
	lowUptimeHighStorage  uint
	lowUptimeLowStorage   uint

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
	hostQuality.nbKnownHosts = len(node.KnownHosts())

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
	log.Println("Updating known hosts")
	knownHosts.uptimeTally = 0
	knownHosts.storageTally = 0

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
	}
	knownHosts.classifyHosts()
	knownHosts.lastUpdate = time.Now()
	log.Println("Updated known hosts")
	log.Println("Known hosts:", knownHosts.Hosts)
}

func (knownHosts *KnownHosts) Remove(host utils.SocketAddr) {
	hostQuality := knownHosts.Hosts[host]

	avgUptime := knownHosts.avgUptime()
	avgStorage := knownHosts.avgStorage()

	knownHosts.decrementHostClass(hostQuality, avgUptime, avgStorage)

	delete(knownHosts.Hosts, host)
}

func (knownHosts *KnownHosts) Add(host utils.SocketAddr) {
	if _, ok := knownHosts.Hosts[host]; ok {
		return // already in the list
	}

	var hostQuality HostQuality
	response, _ := utils.Request(host, []byte("host-quality/"), []byte{})
	_ = json.Unmarshal(response, &hostQuality)

	// TODO: increment the correct category
	avgUptime := knownHosts.avgUptime()
	avgStorage := knownHosts.avgStorage()

	if knownHosts.count() < knownHosts.cap {
		knownHosts.Hosts[host] = hostQuality // if we have the memory just add the known host
		knownHosts.incrementHostClass(hostQuality, avgUptime, avgStorage)
	} else {
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

func (knownHosts *KnownHosts) intelligentAddKnownHost(potentialHost utils.SocketAddr) {
	// OPTIMISATION - these can be cached as part of a KnownHosts struct, i.e. if we recently ran the above functions, we can assume that the distribution has not changed much and just skip to making the decision as to if we should add the known host or not

	if time.Since(knownHosts.lastUpdate) > time.Minute*2 {
		knownHosts.update()
	}

	avgUptime := knownHosts.avgUptime()
	avgStorage := knownHosts.avgStorage()

	// see where the new known host lies in these 4 categories (by repeating the above)
	// figure out the host Uptime and storage distribution
	var hostQuality HostQuality
	response, _ := utils.Request(potentialHost, []byte("host-quality/"), []byte{})
	_ = json.Unmarshal(response, &hostQuality)
	newHostType := knownHosts.hostType(hostQuality, avgUptime, avgStorage)

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
	for host, hostQuality := range knownHosts.Hosts {
		if knownHosts.hostType(hostQuality, avgUptime, avgStorage) == biggestClass {
			knownHosts.Remove(host)
			return
		}
	}
}

func (knownHosts *KnownHosts) biggestClass() string {
	if knownHosts.highUptimeHighStorage > knownHosts.highUptimeLowStorage && knownHosts.highUptimeHighStorage > knownHosts.lowUptimeHighStorage && knownHosts.highUptimeHighStorage > knownHosts.lowUptimeLowStorage {
		return "highUptimeHighStorage"
	} else if knownHosts.highUptimeLowStorage > knownHosts.highUptimeHighStorage && knownHosts.highUptimeLowStorage > knownHosts.lowUptimeHighStorage && knownHosts.highUptimeLowStorage > knownHosts.lowUptimeLowStorage {
		return "highUptimeLowStorage"
	} else if knownHosts.lowUptimeHighStorage > knownHosts.highUptimeHighStorage && knownHosts.lowUptimeHighStorage > knownHosts.lowUptimeLowStorage {
		return "lowUptimeHighStorage"
	} else if knownHosts.lowUptimeLowStorage > knownHosts.highUptimeHighStorage && knownHosts.lowUptimeLowStorage > knownHosts.lowUptimeLowStorage {
		return "lowUptimeLowStorage"
	} else {
		return "unknown"
	}
}

func (knownHosts *KnownHosts) classifyHosts() {
	knownHosts.resetDistribution()

	avgUptime := knownHosts.avgUptime()
	avgStorage := knownHosts.avgStorage()

	for _, quality := range knownHosts.Hosts {
		knownHosts.incrementHostClass(quality, avgUptime, avgStorage)
	}
}

func (knownHosts *KnownHosts) resetDistribution() {
	knownHosts.highUptimeHighStorage = 0
	knownHosts.highUptimeLowStorage = 0
	knownHosts.lowUptimeHighStorage = 0
	knownHosts.lowUptimeLowStorage = 0
}

func (knownHosts *KnownHosts) incrementHostClass(hostQuality HostQuality, avgUptime uint64, avgStorage uint64) {
	switch knownHosts.hostType(hostQuality, avgUptime, avgStorage) {
	case "highUptimeHighStorage":
		knownHosts.highUptimeHighStorage += 1
	case "highUptimeLowStorage":
		knownHosts.highUptimeLowStorage += 1
	case "lowUptimeHighStorage":
		knownHosts.lowUptimeHighStorage += 1
	case "lowUptimeLowStorage":
		knownHosts.lowUptimeLowStorage += 1
	}
}

func (knownHosts *KnownHosts) decrementHostClass(hostQuality HostQuality, avgUptime uint64, avgStorage uint64) {
	switch knownHosts.hostType(hostQuality, avgUptime, avgStorage) {
	case "highUptimeHighStorage":
		knownHosts.highUptimeHighStorage -= 1
	case "highUptimeLowStorage":
		knownHosts.highUptimeLowStorage -= 1
	case "lowUptimeHighStorage":
		knownHosts.lowUptimeHighStorage -= 1
	case "lowUptimeLowStorage":
		knownHosts.lowUptimeLowStorage -= 1
	}
}

func (knownHosts *KnownHosts) hostType(hostQuality HostQuality, avgUptime uint64, avgStorage uint64) string {
	if hostQuality.Uptime > avgUptime && hostQuality.AvailableStorage > avgStorage {
		return "highUptimeHighStorage"
	} else if hostQuality.Uptime > avgUptime && hostQuality.AvailableStorage < avgStorage {
		return "highUptimeLowStorage"
	} else if hostQuality.Uptime < avgUptime && hostQuality.AvailableStorage > avgStorage {
		return "lowUptimeHighStorage"
	} else {
		return "lowUptimeLowStorage"
	}
}

func (knownHosts *KnownHosts) JsonDigest() []byte {
	digest, _ := json.Marshal(knownHosts)
	return digest
}
