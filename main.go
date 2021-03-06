package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mdlayher/wavepipe/common"
	"github.com/mdlayher/wavepipe/core"
	"github.com/mdlayher/wavepipe/env"
)

// testFlag invokes wavepipe in "test" mode, where it will start and exit shortly after.  Used for testing.
var testFlag = flag.Bool("test", false, "Starts "+core.App+" in test mode, causing it to exit shortly after starting.")

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Parse command line flags
	flag.Parse()

	// Check if wavepipe was invoked as root (which is a really bad idea)
	user := common.System.User
	if user.Uid == "0" || user.Gid == "0" || user.Username == "root" {
		log.Println(core.App, ": WARNING, it is NOT advisable to run wavepipe as root!")
	}

	// Application entry point
	log.Println(core.App, ": starting...")

	// Check if running in debug mode, which will allow bypass of certain features such as
	// API authentication.  USE THIS FOR DEVELOPMENT ONLY!
	if env.IsDebug() {
		log.Println(core.App, ": WARNING, running in debug mode; authentication disabled!")
	}

	// Gracefully handle termination via UNIX signal
	sigChan := make(chan os.Signal, 1)

	// In test mode, wait for a short time, then invoke a signal shutdown
	if *testFlag {
		// Set an environment variable to enable mocking in other areas of the program
		if err := env.SetTest(true); err != nil {
			log.Println(err)
		}

		go func() {
			// Wait a few seconds, to allow reasonable startup time
			seconds := 10
			log.Println(core.App, ": started in test mode, stopping in", seconds, "seconds.")
			<-time.After(time.Duration(seconds) * time.Second)

			// Send interrupt
			sigChan <- os.Interrupt
		}()
	}

	// Invoke the manager, with graceful termination and core.Application exit code channels
	killChan := make(chan struct{})
	exitChan := make(chan int)
	go core.Manager(killChan, exitChan)

	// Trigger a shutdown if SIGINT or SIGTERM received
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	for sig := range sigChan {
		log.Println(core.App, ": caught signal:", sig)
		killChan <- struct{}{}
		break
	}

	// Force terminate if signaled twice
	go func() {
		for sig := range sigChan {
			log.Println(core.App, ": caught signal:", sig, ", force halting now!")
			os.Exit(1)
		}
	}()

	// Graceful exit
	code := <-exitChan
	log.Println(core.App, ": graceful shutdown complete")
	os.Exit(code)
}
