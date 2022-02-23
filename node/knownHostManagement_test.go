package node

import (
	"fmt"
	"testing"
)

func TestBiggestClass(t *testing.T) {
	knownHosts := KnownHosts{highUptimeHighStorageHighKH: 1, highUptimeHighStorageLowKH: 0, highUptimeLowStorageHighKH: 0, highUptimeLowStorageLowKH: 0,
		lowUptimeHighStorageHighKH: 0, lowUptimeHighStorageLowKH: 0, lowUptimeLowStorageHighKH: 0, lowUptimeLowStorageLowKH: 0}
	if knownHosts.biggestClass() != "highUptimeHighStorageHighKH" {
		t.Error("Expected biggest class to be highUptimeHighStorageHighKH")
	}
	knownHosts = KnownHosts{highUptimeHighStorageHighKH: 0, highUptimeHighStorageLowKH: 1, highUptimeLowStorageHighKH: 0, highUptimeLowStorageLowKH: 0,
		lowUptimeHighStorageHighKH: 0, lowUptimeHighStorageLowKH: 0, lowUptimeLowStorageHighKH: 0, lowUptimeLowStorageLowKH: 0}
	if knownHosts.biggestClass() != "highUptimeHighStorageLowKH" {
		t.Error("Expected biggest class to be highUptimeHighStorageLowKH")
	}
	knownHosts = KnownHosts{highUptimeHighStorageHighKH: 0, highUptimeHighStorageLowKH: 0, highUptimeLowStorageHighKH: 1, highUptimeLowStorageLowKH: 0,
		lowUptimeHighStorageHighKH: 0, lowUptimeHighStorageLowKH: 0, lowUptimeLowStorageHighKH: 0, lowUptimeLowStorageLowKH: 0}
	if knownHosts.biggestClass() != "highUptimeLowStorageHighKH" {
		t.Error("Expected biggest class to be highUptimeLowStorageHighKH")
	}
	knownHosts = KnownHosts{highUptimeHighStorageHighKH: 0, highUptimeHighStorageLowKH: 0, highUptimeLowStorageHighKH: 0, highUptimeLowStorageLowKH: 1,
		lowUptimeHighStorageHighKH: 0, lowUptimeHighStorageLowKH: 0, lowUptimeLowStorageHighKH: 0, lowUptimeLowStorageLowKH: 0}
	if knownHosts.biggestClass() != "highUptimeLowStorageLowKH" {
		t.Error("Expected biggest class to be highUptimeLowStorageLowKH")
	}
	knownHosts = KnownHosts{highUptimeHighStorageHighKH: 0, highUptimeHighStorageLowKH: 0, highUptimeLowStorageHighKH: 0, highUptimeLowStorageLowKH: 0,
		lowUptimeHighStorageHighKH: 1, lowUptimeHighStorageLowKH: 0, lowUptimeLowStorageHighKH: 0, lowUptimeLowStorageLowKH: 0}
	if knownHosts.biggestClass() != "lowUptimeHighStorageHighKH" {
		t.Error("Expected biggest class to be lowUptimeHighStorageHighKH")
	}
	knownHosts = KnownHosts{highUptimeHighStorageHighKH: 0, highUptimeHighStorageLowKH: 0, highUptimeLowStorageHighKH: 0, highUptimeLowStorageLowKH: 0,
		lowUptimeHighStorageHighKH: 0, lowUptimeHighStorageLowKH: 1, lowUptimeLowStorageHighKH: 0, lowUptimeLowStorageLowKH: 0}
	if knownHosts.biggestClass() != "lowUptimeHighStorageLowKH" {
		t.Error("Expected biggest class to be lowUptimeHighStorageLowKH")
	}
	knownHosts = KnownHosts{highUptimeHighStorageHighKH: 0, highUptimeHighStorageLowKH: 0, highUptimeLowStorageHighKH: 0, highUptimeLowStorageLowKH: 0,
		lowUptimeHighStorageHighKH: 0, lowUptimeHighStorageLowKH: 0, lowUptimeLowStorageHighKH: 1, lowUptimeLowStorageLowKH: 0}
	if knownHosts.biggestClass() != "lowUptimeLowStorageHighKH" {
		t.Error("Expected biggest class to be lowUptimeLowStorageHighKH")
	}
	knownHosts = KnownHosts{highUptimeHighStorageHighKH: 0, highUptimeHighStorageLowKH: 0, highUptimeLowStorageHighKH: 0, highUptimeLowStorageLowKH: 0,
		lowUptimeHighStorageHighKH: 0, lowUptimeHighStorageLowKH: 0, lowUptimeLowStorageHighKH: 0, lowUptimeLowStorageLowKH: 1}
	if knownHosts.biggestClass() != "lowUptimeLowStorageLowKH" {
		t.Error("Expected biggest class to be lowUptimeLowStorageLowKH")
	}

	// Two biggest classes - doest matter which one is returned as long as it is not empty
	knownHosts = KnownHosts{highUptimeHighStorageHighKH: 1, highUptimeHighStorageLowKH: 1, highUptimeLowStorageHighKH: 0, highUptimeLowStorageLowKH: 0,
		lowUptimeHighStorageHighKH: 0, lowUptimeHighStorageLowKH: 0, lowUptimeLowStorageHighKH: 0, lowUptimeLowStorageLowKH: 0}
	fmt.Println(knownHosts.biggestClass())
	if knownHosts.biggestClass() != "highUptimeHighStorageHighKH" {
		t.Error("Expected biggest class to be highUptimeHighStorageHighKH")
	}
}

func TestHostType(t *testing.T) {
	avgUptime := uint64(10000)
	avgAvailableStorage := uint64(10000)
	avgHostsKnown := uint64(50)

	// generate permutations of HostQuality
	hostQuality := HostQuality{Uptime: 10000, AvailableStorage: 10000, NbHostsKnown: 50}
	if hostQuality.hostType(avgUptime, avgAvailableStorage, avgHostsKnown) != "highUptimeHighStorageHighKH" {
		t.Error("Expected host type to be highUptimeHighStorageHighKH")
	}

	hostQuality = HostQuality{Uptime: 10000, AvailableStorage: 10000, NbHostsKnown: 3}
	if hostQuality.hostType(avgUptime, avgAvailableStorage, avgHostsKnown) != "highUptimeHighStorageLowKH" {
		t.Error("Expected host type to be highUptimeHighStorageLowKH")
	}

	hostQuality = HostQuality{Uptime: 10000, AvailableStorage: 0, NbHostsKnown: 50}
	if hostQuality.hostType(avgUptime, avgAvailableStorage, avgHostsKnown) != "highUptimeLowStorageHighKH" {
		t.Error("Expected host type to be highUptimeLowStorageHighKH")
	}

	hostQuality = HostQuality{Uptime: 10000, AvailableStorage: 0, NbHostsKnown: 3}
	if hostQuality.hostType(avgUptime, avgAvailableStorage, avgHostsKnown) != "highUptimeLowStorageLowKH" {
		t.Error("Expected host type to be highUptimeLowStorageLowKH")
	}

	hostQuality = HostQuality{Uptime: 10, AvailableStorage: 10000, NbHostsKnown: 50}
	if hostQuality.hostType(avgUptime, avgAvailableStorage, avgHostsKnown) != "lowUptimeHighStorageHighKH" {
		t.Error("Expected host type to be lowUptimeHighStorageHighKH")
	}

	hostQuality = HostQuality{Uptime: 10, AvailableStorage: 10000, NbHostsKnown: 3}
	if hostQuality.hostType(avgUptime, avgAvailableStorage, avgHostsKnown) != "lowUptimeHighStorageLowKH" {
		t.Error("Expected host type to be lowUptimeHighStorageLowKH")
	}

	hostQuality = HostQuality{Uptime: 10, AvailableStorage: 0, NbHostsKnown: 50}
	if hostQuality.hostType(avgUptime, avgAvailableStorage, avgHostsKnown) != "lowUptimeLowStorageHighKH" {
		t.Error("Expected host type to be lowUptimeLowStorageHighKH")
	}

	hostQuality = HostQuality{Uptime: 10, AvailableStorage: 0, NbHostsKnown: 3}
	if hostQuality.hostType(avgUptime, avgAvailableStorage, avgHostsKnown) != "lowUptimeLowStorageLowKH" {
		t.Error("Expected host type to be lowUptimeLowStorageLowKH")
	}

}

func TestKnownHostCountManagement(t *testing.T) {
	// Spawn 3 nodes with a knownhost cap of 1
	// Make sure they all known each other
}
