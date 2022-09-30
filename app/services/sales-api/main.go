package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/automaxprocs/maxprocs"
)

var build = "develop"

func main() {
	//maxprocs.New()
	// Set the correct number of threads for the service
	// based on what is available either by the machine or quotas in k8s for container.
	if _, err := maxprocs.Set(); err != nil {
		fmt.Printf("maxprocs: %w", err)
		os.Exit(1)
	}

	g := runtime.GOMAXPROCS(0) //Set cpus allowed to run here

	log.Printf("starting service build[%s] CPU[%d]", build, g)
	defer log.Println("service ended")

	shutdown := make(chan os.Signal, 1)
	// SIGINT allows for ctrl^c to close as signal
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("stopping service")
}
