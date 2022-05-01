// Package tracker enables nodes to be tracked by the butter-tracker. This is just a means of designing a system that
// has an overview of the network. Mostly for testing purposes.
package tracker

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/butter-network/butter/node"
	"github.com/butter-network/butter/utils"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "docker"
	dbname   = "world"
)

// Track every 10s the node sends
// - it's known hosts
// - the information it contains
func Track(overlay node.Overlay) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for {
		addrs := make([]utils.SocketAddr, 0)
		time.Sleep(time.Second * 10)
		for host := range overlay.Node().KnownHosts() {
			addrs = append(addrs, host)
		}
		addrsJson, _ := json.Marshal(addrs)
		_, err1 := db.Exec("DELETE FROM nodes WHERE addr = $1;", overlay.Node().Address())
		if err1 != nil {
			//log.Print(err1)
		}
		_, err2 := db.Exec("INSERT INTO nodes(addr, peers, groups) VALUES($1,$2,NULL);", overlay.Node().Address(), string(addrsJson))
		if err2 != nil {
			//log.Print(err2)
		}
	}
}
