package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/containers/winquit/pkg/winquit"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "simple-server":
		simpleServer()
	case "signal-server":
		signalServer()
	case "multi-server":
		multiServer()
	case "hang-server":
		hangServer()
	case "request-quit":
		requestQuit()
	case "demand-quit":
		demandQuit()
	}
}

func parseInt(pos int) int {
	if len(os.Args) < pos+1 {
		printUsage()
		os.Exit(1)
	}

	pid, err := strconv.ParseUint(os.Args[pos], 10, 0)
	if err != nil {
		panic(err)
	}

	return int(pid)
}

func requestQuit() {
	pid := parseInt(2)

	if err := winquit.RequestQuit(pid); err != nil {
		panic(err)
	}
}

func demandQuit() {
	pid := parseInt(2)
	timeout := parseInt(3)

	if err := winquit.QuitProcess(pid, time.Second*time.Duration(timeout)); err != nil {
		panic(err)
	}
}

func printUsage() {
	executable := filepath.Base(os.Args[0])
	fmt.Printf("Usage: %s [COMMAND] [ARG...] \n\n", executable)
	fmt.Printf("  simple-server               start a server which waits on a boolean channel\n")
	fmt.Printf("  signal-server               start a server which waits on a simulated SIGTERM\n")
	fmt.Printf("  hang-server                 start a server which ignores quit messages\n")
	fmt.Printf("  multi-server                start a server with multiple channels subscribed\n")
	fmt.Printf("  request-quit  (pid)         ask another process to quit\n")
	fmt.Printf("  demand-quit   (pid) (secs)  first ask, then kill at timeout\n")
}

func signalServer() {
	logrus.Info("Server waiting using signal approach")
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM)
	winquit.SimulateSigTermOnQuit(done)
	logrus.Infof("Received: %v", <-done)
}

func simpleServer() {
	logrus.Info("Server waiting using simple boolean approach")
	done := make(chan bool)
	winquit.NotifyOnQuit(done)
	logrus.Infof("Received: %v", <-done)
}

func multiServer() {
	logrus.Info("Server waiting using multiple boolean approach")
	var chans []chan bool

	for i := 0; i < 5; i++ {
		channel := make(chan bool, 1)
		chans = append(chans, channel)
		winquit.NotifyOnQuit(channel)
	}

	for _, channel := range chans {
		logrus.Infof("Received: %v", <-channel)
	}
}

func hangServer() {
	logrus.Info("Hanging server waiting forever")

	for {
		time.Sleep(time.Second * 100)
	}
}
