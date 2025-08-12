package main

import (
	"evilwhale/core"
	"flag"
	"log"
)

func main() {
	host := flag.String("host", "http://127.0.0.1:2375", "Set the host")
	image := flag.String("image", "", "Set the docker image to be used")
	cmd := flag.String("cmd", "", "Set the command to be executed")
	ping := flag.Bool("ping", false, "Ping the host")
	pull := flag.Bool("pull", false, "Pull down an image")

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

	if *pull {
		core.PullImage(*host, *image)
		return
	}
}
