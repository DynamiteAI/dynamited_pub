// dynamited main package
package main

import (
	// Built-ins
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	// Dynamite packages
	"dynamite_daemon_core/pkg/common"
	"dynamite_daemon_core/pkg/conf"
	"dynamite_daemon_core/pkg/logging"
	"dynamite_daemon_core/pkg/watcher"
)

var (
	confFile string
)

func main() {
	// Set up a global context and get a closer func
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Attach a channel to receive interrupts
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Kill)

	// Deferred func to ensure we stop receiving signals and cancel the context before exit
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	// Background task to receive and forward interrupts
	go func() {
		select {
		case msg := <-c:
			common.Quit <- []byte(msg.String())
		}
	}()

	// Declare and parse command line options
	flag.StringVar(&confFile, "c", "/etc/dynamite/dynamited/config.yml", "Location of the Dynamite Manage configuration file.")
	flag.Parse()

	// Load the provided conf file into the global conf.Conf struct variable
	conf.Load(confFile)

	// Initialize the configured logging directory
	if !logging.Init() {
		fmt.Println("Unable to write logs. Exiting.")
		signal.Stop(c)
		cancel()
		os.Exit(1)
	}
	// Initialize the dynamited application log
	logging.SetupAppLogger()

	logging.LogEntry.WithField("roles", strings.Join(conf.Conf.Roles, ", ")).Info("loading_roles")

	// Always run watcher routines for configured roles
	watcher.Init(ctx)

	fmt.Printf("dynamited is running. log directory: %v\n", logging.LogDir)
	logging.LogEntry.Info("dynamited_running.")

	// Main loop
	for {
		// Run until signaled on common.Quit channel
		select {
		case msg := <-common.Quit:
			fmt.Println(string(msg))
			fmt.Println("shutting down")
			cancel()
			// Close the dynamited log file
			if file, ok := logging.Log.Out.(*os.File); ok {
				file.Sync()
				file.Close()
			}
		case <-ctx.Done():
			os.Exit(0)
		}
	}
}
