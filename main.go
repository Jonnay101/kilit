package main

import (
	"flag"
	"fmt"

	"github.com/jonnay101/kilit/command"
	"github.com/jonnay101/kilit/domain"
)

func main() {
	port := flag.String("p", "8080", "the port you want to kill")
	flag.Parse()

	Do("lsof", "-i", fmt.Sprintf(":%s", *port))
}

func Do(name string, params ...string) {
	store := new(domain.DataStore)
	command.CreateData(store, name, params...)

	killer := new(command.ProcessKiller)
	domain.GetAllPIDsFromData(store.Data(), killer)

	killer.Kill()
}
