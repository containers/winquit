package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/n1hility/winquit/pkg/winquit"
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
	fmt.Printf("%s [COMMAND] [ARG...] \n", os.Args[0])
	fmt.Printf("%s simple-server\n", os.Args[0])
	fmt.Printf("%s signal-server\n", os.Args[0])
	fmt.Printf("%s hang-server\n", os.Args[0])
	fmt.Printf("%s multi-server\n", os.Args[0])
	fmt.Printf("%s request-quit (pid)\n", os.Args[0])
	fmt.Printf("%s demand-quit (pid) (timeout in secs)\n", os.Args[0])
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
