package store

// NaiveStore stores information on the network naively by simply placing it on the local node. It generate a UUIS for
// the information and creates an information block and return information uuid
//func NaiveStore(node *Node, keywords []string, information string) string {
//	node.lock.Lock()
//	// Generate UUID
//	u, _ := uuid.NewV4()
//	node.storage[u.String()] = Block{
//		keywords: keywords,
//		part:     0,
//		parts:    0,
//		data:     information,
//	}
//	node.lock.Unlock()
//	return u.String()
//}
