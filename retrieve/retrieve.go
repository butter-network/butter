package retrieve

//
//// NaiveRetrieve High level entrypoint for searching for a specific piece of information on the network
//// look if I have the information else look at the most likely known host to get to that information
//// one query per piece of information (one-to-one) hence the query has to be unique i.e i.d.
//func NaiveRetrieve(node *Node, query string) string {
//	// do I have this information, if so return it
//	// else BFS (pass the query on to all known hosts (partial view)
//	node.lock.Lock()
//	defer node.lock.Unlock()
//	uuid, _ := uuid.Parse([]byte(query))
//	if val, ok := node.storage[*uuid]; ok {
//		return string(val.data)
//	} else {
//		return bfs(node, query)
//	}
//}
//
//func bfs(node *Node, query string) string {
//	// Initialise an empty queue
//	queue := make([]string, 0)
//	// Add all my known hosts to the queue
//	for _, host := range node.knownHosts {
//		queue = append(queue, host)
//	}
//	for len(queue) > 0 {
//		// Pop the first element from the queue
//		host := queue[0]
//		queue = queue[1:]
//		// Start a connection to the host
//		c, err := net.Dial("tcp", host)
//		if err != nil {
//			fmt.Println(err)
//			return "Error connecting to host"
//		}
//		c.Close()
//		// Ask host if he has data
//		fmt.Fprint(c, "/remote-retrieve "+query)
//		// Receive response
//		reply := make([]byte, 1024)
//		c.Read(reply)
//		uri, payload := parsePacket(string(reply))
//		// If the returned packet is success + the data then return it
//		// else add the known hosts of the remote node to the end of the queue
//		if uri == "/success" {
//			return payload
//		} else {
//			fmt.Fprint(c, "/get-remote-known-hosts"+query)
//			c.Read(reply)
//			// convert json list of known hosts into a slice of strings
//			remoteHosts := make([]string, 0)
//			err = json.Unmarshal(reply, &remoteHosts)
//			if err != nil {
//				fmt.Println(err)
//				return "Error decoding json"
//			}
//			// add the remote hosts to the end of the queue
//			queue = append(queue, remoteHosts...)
//		}
//		return "Information is not on the network"
//	}
//	return "This should not happen"
//}
