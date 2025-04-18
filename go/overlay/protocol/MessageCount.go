package protocol

import (
	"bytes"
	"github.com/saichler/shared/go/share/logger"
	"github.com/saichler/types/go/types"
	"sync/atomic"
)

var CountMessages = false
var messagesCreated atomic.Uint64
var propertyChangeCalled atomic.Uint64
var ExplicitLog = logger.NewLoggerDirectImpl(logger.NewFileLogMethod("/tmp/Explicit.log"))

func AddMessageCreated() {
	if CountMessages {
		messagesCreated.Add(1)
	}
}

func AddPropertyChangeCalled(set *types.NotificationSet) {
	if CountMessages {
		propertyChangeCalled.Add(1)
		props := ""
		if set.Type == types.NotificationType_Update {
			buff := bytes.Buffer{}
			buff.WriteString(" - ")
			for _, chg := range set.NotificationList {
				buff.WriteString(chg.PropertyId)
				buff.WriteString(" ")
			}
			props = buff.String()
		}
		ExplicitLog.Trace("*** Property Change: ", set.ServiceName, " ", set.Type.String(), ":", props)

	}
}

func MessagesCreated() uint64 {
	return messagesCreated.Load()
}

func PropertyChangedCalled() uint64 {
	return propertyChangeCalled.Load()
}
