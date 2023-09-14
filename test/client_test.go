//go:build windows
// +build windows

package e2e

import (
    "os"
    "os/exec"
    "path/filepath"
    "time"

    "github.com/n1hility/winquit/pkg/winquit"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var WINQUIT_PATH = filepath.Join("..", "bin", "winquit.exe")

const SHOULD_TIME = 10
const WONT_TIME = 5

var _ = Describe("perquisites", func() {
    It("winquit binary is built", func() {
        _, err := os.Stat(WINQUIT_PATH)
        Expect(err).ShouldNot(HaveOccurred())
    })
})

var _ = Describe("client", func() {
    It("request quit should kill thidparty(winver) process", func() {
        cmd := exec.Command("winver")
        verifyRequestQuit(cmd, SHOULD_TIME, true)
    })
})

var _ = Describe("client", func() {
    It("request quit kills winquit simple server", func() {
        cmd := exec.Command(WINQUIT_PATH, "simple-server")
        verifyRequestQuit(cmd, SHOULD_TIME, true)
    })
})

var _ = Describe("client", func() {
    It("request quit kills winquit multi-server", func() {
        cmd := exec.Command(WINQUIT_PATH, "multi-server")
        verifyRequestQuit(cmd, SHOULD_TIME, true)
    })
})

var _ = Describe("client", func() {
    It("request quit kills winquit signal server", func() {
        cmd := exec.Command(WINQUIT_PATH, "signal-server")
        verifyRequestQuit(cmd, SHOULD_TIME, true)
    })
})

var _ = Describe("client", func() {
    It("request quit does not kill winquit hang server", func() {
        cmd := exec.Command(WINQUIT_PATH, "hang-server")
        verifyRequestQuit(cmd, WONT_TIME, false)
    })
})

var _ = Describe("client", func() {
    It("demand quit does kill winquit hang server", func() {
        cmd := exec.Command(WINQUIT_PATH, "hang-server")
        verifyForceQuit(cmd, WONT_TIME, SHOULD_TIME, true)
    })
})

func verifyRequestQuit(cmd *exec.Cmd, timeout int, outcome bool) {
    verifyStart(cmd)
    winquit.RequestQuit(cmd.Process.Pid)
    verifyExit(cmd, timeout, outcome)
}

func verifyForceQuit(cmd *exec.Cmd, forceTimeout int, timeout int, outcome bool) {
    verifyStart(cmd)
    winquit.QuitProcess(cmd.Process.Pid, time.Duration(forceTimeout)*time.Second)
    verifyExit(cmd, timeout, outcome)
}

func verifyStart(cmd *exec.Cmd) {
    err := cmd.Start()
    Expect(err).ShouldNot(HaveOccurred())
    time.Sleep(100 * time.Millisecond)
    Expect(cmd.ProcessState).To(BeNil())
    _, err = os.FindProcess(cmd.Process.Pid)
    Expect(err).ShouldNot(HaveOccurred())
}

func verifyExit(cmd *exec.Cmd, timeout int, outcome bool) {
    completed := make(chan bool)
    go func() {
        cmd.Wait()
        completed <- true
    }()

    result := false
    select {
    case <-completed:
        result = true
    case <-time.After(time.Duration(timeout) * time.Second):
    }

    Expect(result).To(Equal(outcome))
    if !outcome {
        cmd.Process.Kill()
    }
}
