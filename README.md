# winquit

winquit is a golang module that supports graceful shutdown of Windows
applications through the sending and receiving of Windows quit events on Win32
message queues. This allows golang applications to implement behavior comparable
to SIGTERM signal handling on UNIX derived systems. Additionally, it supports
the graceful shutdown mechanism employed by Windows system tools, such as
`taskkill.exe`. See the [How it works](#how-it-works) section for more details.

## Overview

To aid application portability, and provide familiarity, the API follows a
similar convention and approach as the os.signal package. Additionally, the
SimulateSigTermOnQuit function supports reuse of the same underlying channel,
supporting the blending of os.signal and winquit together (a subset of
signals provided by os.signal are still relevant and desirable on Windows,
for example, break handling in console applications).

### Simple server example

The following example demonstrates usage of NotifyOnQuit() to wait for a
windows quit event before shutting down:

```golang
func server() {
    fmt.Println("Starting server")

    // Create a channel, and register it
    done := make(chan bool, 1)
    winquit.NotifyOnQuit(done)

    // Wait until we receive a quit event
    <-done

    fmt.Println("Shutting down")
    // Perform cleanup tasks
}
```

### Blended signal example

The following example demonstrates usage of SimulateSigTermOnQuit() in
concert with signal.Notify():

```golang
func server() {
    fmt.Println("Starting server")

    // Create a channel, and register it
    done := make(chan os.Signal, 1)

    // Wait on console interrupt events
    signal.Notify(done, syscall.SIGINT)

    // Simulate SIGTERM when a quit occurs
    winquit.SimulateSigTermOnQuit(done)

    // Wait until we receive a signal or quit event
    <-done

    fmt.Println("Shutting down")
    // Perform cleanup tasks
}
```

### Client example

The following example demonstrates how an application can ask or
force other windows programs to quit:

```golang
func client() {
    // Ask nicely for program "one" to quit. This request may not
    // be honored if its a console application, or if the program
    // is hung
    if err := winquit.RequestQuit(pidOne); err != nil {
        fmt.Printf("error sending quit request, %s", err.Error())
    }

    // Force program "two" to quit, but give it 20 seconds to
    // perform any cleanup tasks and quit on it's own
    timeout := time.Second * 20
    if err := winquit.QuitProcess(pidTwo, timeout); err != nil {
        fmt.Printf("error killing process, %s", err.Error())
    }
```

## Command-line tool

For demonstration and testing purposes, a command-line tool `winquit.exe` is
provided as part of the project build.

### Usage

```
Usage: winquit.exe [COMMAND] [ARG...]

  simple-server               start a server which waits on a boolean channel
  signal-server               start a server which waits on a simulated SIGTERM
  hang-server                 start a server which ignores quit messages
  multi-server                start a server with multiple channels subscribed
  request-quit  (pid)         ask another process to quit
  demand-quit   (pid) (secs)  first ask, then kill at timeout
```

### Example

#### In terminal 1:

```
PS> .\winquit.exe simple server

time="2023-09-13T23:09:14-05:00" level=info msg="Server waiting using simple boolean approach"
time="2023-09-13T23:09:14-05:00" level=info msg="Entering loop for quit"
```

#### In terminal 2:
```
PS> .\bin\winquit.exe request-quit 13332

time="2023-09-13T23:09:46-05:00" level=debug msg="Closing windows on thread 10792"
time="2023-09-13T23:09:46-05:00" level=debug msg="Closing windows on thread 1592"
time="2023-09-13T23:09:46-05:00" level=debug msg="Closing windows on thread 3500"
time="2023-09-13T23:09:46-05:00" level=debug msg="Closing windows on thread 5368"
time="2023-09-13T23:09:46-05:00" level=debug msg="Closing windows on thread 8324"
time="2023-09-13T23:09:46-05:00" level=debug msg="Closing windows on thread 12852"
```
#### Back in terminal 1:

```
time="2023-09-13T23:09:46-05:00" level=debug msg="Received QUIT notification"
time="2023-09-13T23:09:46-05:00" level=info msg="Received: true"
```


## How it works

Windows GUI applications consist of multiple components (and windows) which
intercommunicate with events over per-thread message queues and/or direct
event handoff to window procedures for cross-thread communication.
Additionally, GUI applications can use the same mechanism to communicate with
windows and threads owned by other applications, including common desktop
components.

winquit utilizes this mechanism by creating a standard win32 message loop
thread and registering a non-visible window to relay a quit message (WM_QUIT)
in the event of a window close event. WM_CLOSE is sent by Windows in response
to certain system events, or by other requesting applications. For example,
the system provided taskkill.exe (similar to the kill command on Unix), works
by iterating all windows on the system, and sending a WM_CLOSE when the
process owner matches the specified pid. Note that, unlike UNIX/X11 style
systems, on Windows the graphical APIs are built-in and accessible to all
win32 applications, including console based applications. Therefore, the APIs
provided by winquit *do not* require compilation as a windowsgui app to
effectively use them.

winquit also provides APIs to trigger a quit of another process using a WM_CLOSE
event, although in a more efficient manner than taskkill.exe. It instead
captures a thread snapshot of the target process (effectively a memory read on
Windows), and enumerates each thread's associated Windows, and sending the event
to each. In addition to supporting a graceful close of any Windows application,
which may have multiple message loops, this approach also obviates the need for
cumbersome approaches to lock code to the main thread of the application. The
message loop used by winquit does not care which thread the golang runtime
internally designates. Note that winquit purposefully relays through a thread's
windows, as opposed to posting directly to each thread's message queue, since
the former is more likely to be expected by an application, and it ensures all
window procedures have an opportunity to perform cleanup work not associated
with the thread's message loop.

## Limitations

This API is only implemented on Windows platforms. Non-operational stubs
are provided for compilation purposes.

In addition to requiring appropriate security permissions (typically a user
can only send events to other applications ran by the same user), Windows
also restricts inter-app messaging operations to programs running in the same
user logon session. While logons migrate between RDP and console sessions,
non-graphical logins (e.g sshd) typically create a logon per connection. For
this reason, tools like taskkill and winquit are normally disallowed from
crossing this boundary. Therefore, a user will not be able to gracefully stop
applications between ssh/winrm sessions, and in between ssh and graphical
logons. However, the typical user use-case of logging into Windows and
running multiple applications and terminals will work fine. Additionally,
multiple back-grounded processes in the same ssh session will be able to
communicate. Finally, it's possible to bypass this limitation by executing
code under the system user using the SeTcbPrivilege. The psexec tool does
exactly this, and can additionally be used as a workaround to this
limitation.

## Building and Testing

winquit includes a `build.ps1` build script for Windows, as well as a Makefile
for Unix/Linux. Alternatively, the Makefile supports execution on Windows if the
msys2 environment is installed. NMake can not be used to run the Makefile, since
it utilizes GNU specific grammar.  Tests and the command-line tool can only be
ran on Windows, since the implementation of winquit is Windows specific.

### Building on Windows

```
PS> .\build.ps1
removing bin
go clean -testcache
go build -v -o bin/winquit.exe ./cmd/winquit
```

### Testing on Windows

```
PS> .\build.ps1 test
go test -v ./test
=== RUN   TestTest
Running Suite: Test Suite - C:\build\winquit\test
=======================================================
Random Seed: 1694665098

Will run 7 of 7 specs
+++++++

Ran 7 of 7 Specs in 3.082 seconds
SUCCESS! -- 7 Passed | 0 Failed | 0 Pending | 0 Skipped
--- PASS: TestTest (3.09s)
PASS
ok      github.com/containers/winquit/test        3.234s
```

### Building on Linux/Unix 

```
$ make
mkdir -p bin
go build -o bin/winquit.exe ./cmd/winquit

$ file bin/winquit.exe
bin/winquit.exe: PE32+ executable (console) x86-64 (stripped to external PDB), for MS Windows
```

