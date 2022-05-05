package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"lpfr-o-matic/pkg/sys"
	"lpfr-o-matic/pkg/watchdog"
)

var (
	exePath       = flag.String("exepath", "Infod_LPFR.exe", "LPFR executable name")
	checkUrl      = flag.String("checkurl", "http://localhost:7555", "url of the LPFR to check")
	interval      = flag.Int("interval", 10, "interval (in seconds) to perform checks")
	pin           = flag.String("pin", "", "pin code to authorize")
	noAutoPin     = flag.Bool("nopin", false, "skip automatic pin setup")
	middlewareApp = flag.String("middleware", "", "middleware application to start on successful status")
)

func main() {

	flag.Parse()

	_, err := sys.CreateMutex("lpfr-o-matic")
	if err != nil {
		log.Println("It looks like another instance of lpfr-o-matic is running")
		return
	}

	wd := watchdog.NewWatchdog(*exePath, *checkUrl, *interval, *pin, *noAutoPin, *middlewareApp)
	err = wd.Start()
	if err != nil {
		log.Fatalln(err)
	}
	wait := make(chan os.Signal)
	fmt.Println("Started")
	signal.Notify(wait, os.Kill, os.Interrupt)
	<-wait
	fmt.Println("Stopped")

}
