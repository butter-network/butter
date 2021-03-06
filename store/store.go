package store

import (
	"github.com/butter-network/butter/persist"
)

// NaiveStore stores information on the network naively by simply placing it on the local node. It generate a UUIS for
// the information and creates an information block and return information uuid
func NaiveStore(overlay *persist.Overlay, keywords []string, data string) string {
	uuid := overlay.AddInformation(keywords, []byte(data))
	return uuid
}
