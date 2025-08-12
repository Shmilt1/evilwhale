package main

import (
	"evilwhale/core"
	"flag"
	"log"
)

func main() {
	host := flag.String("h", "127.0.0.1:2375", "Set the host")
	image := flag.String("i", "", "Set the docker image to be used")
	cmd := flag.String("c", "", "Set the command to be executed")
	ping := flag.Bool("p", false, "Ping the host")

	flag.Parse()

	if *ping {
		if core.Ping(*host) {
			log.Println("Host is alive!")
			return
		}
	}

	if *image != "" && *cmd != "" {
		cid, err := core.CreateContainer(*host, *image)
		if err != nil {
			log.Fatalln(err)
		}

		id, err := core.CreateExecInstance(*host, cid, *cmd)
		if err != nil {
			log.Fatalln(err)
		}

		err = core.StartExec(*host, id)
		if err != nil {
			log.Fatalln(err)
		}

		core.GetContainerLogs(*host, cid)
	}
}
