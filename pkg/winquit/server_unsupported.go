//go:build !windows
// +build !windows

package winquit

import "github.com/sirupsen/logrus"

func NotifyOnQuit(done chan bool) {
	logrus.Warn("Called NotifyOnQuit(): not implemented on Non-Windows")
}
