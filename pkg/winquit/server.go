package winquit

import (
	"os"
)

// NotifyOnQuit relays a Windows quit notification to the boolean done channel.
// This is a one-shot operation (will only be delivered once), however multiple
// channels may be registered. Each registered channel is sent one copy of the
// same one-shot value.
//
// This function is a no-op on non-Windows platforms. While the call will
// succeed, no notifications will be delivered to the passed channel. Each
// channel will only ever receive a "true" value.
//
// It is recommended that registered channels establish a buffer of 1, since
// values are sent non-blocking. Blocking redelivery may be attempted to reduce
// the chance of bugs; however, it should not be relied upon.
//
// If this function is called after a Windows quit notification has occurred, it
// will immediately deliver a "true" value.
func NotifyOnQuit(done chan bool) {
	notifyOnQuit(done)
}

// SimulateSigTermOnQuit relays a Windows quit notification following the same
// semantics as NotifyOnQuit; however, instead of a boolean message value, this
// function will send a SIGTERM signal to the passed channel.
//
// This function allows for the reuse of the same underlying channel used with
// in a separate os.signal.Notify method call.
func SimulateSigTermOnQuit(handler chan os.Signal) {
	simulateSigTermOnQuit(handler)
}
