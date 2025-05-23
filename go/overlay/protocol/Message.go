package protocol

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/nets"
	"time"
)

type MessageHeader struct {
	source      [36]byte
	vnet        [36]byte
	destination string
	serviceArea uint16
	serviceName string
}

type Transaction struct {
	id        [36]byte
	state     ifs.TransactionState
	errMsg    string
	startTime int64
}

func NewTransaction() ifs.ITransaction {
	tr := &Transaction{}
	copy(tr.id[0:36], ifs.NewUuid())
	tr.state = ifs.Create
	tr.startTime = time.Now().Unix()
	return tr
}

type Message struct {
	MessageHeader
	sequence    uint32
	priority    ifs.Priority
	action      ifs.Action
	timeout     uint16
	request     bool
	reply       bool
	failMessage string
	data        string
	tr          *Transaction
}

const (
	POS_Source       = 0
	POS_Vnet         = 36
	POS_Destination  = POS_Vnet + 36
	POS_Service_Area = POS_Destination + 37
	POS_Service_Name = POS_Service_Area + 2
)

func (this *Message) Clone() *Message {
	clone := &Message{}
	clone.source = this.source
	clone.vnet = this.vnet
	clone.destination = this.destination
	clone.serviceArea = this.serviceArea
	clone.serviceName = this.serviceName
	clone.sequence = this.sequence
	clone.priority = this.priority
	clone.action = this.action
	clone.reply = this.reply
	clone.request = this.request
	clone.data = this.data
	clone.failMessage = this.failMessage
	clone.timeout = this.timeout
	if !ifs.IsNil(this.tr) {
		clone.tr = &Transaction{
			id:        this.tr.id,
			state:     this.tr.state,
			errMsg:    this.tr.errMsg,
			startTime: this.tr.startTime,
		}
	}
	return clone
}

func (this *Message) ReplyClone(resources ifs.IResources) ifs.IMessage {
	reply := this.Clone()
	reply.action = ifs.Reply
	reply.destination = string(this.source[0:36])
	copy(reply.source[0:36], resources.SysConfig().LocalUuid)
	copy(reply.vnet[0:36], resources.SysConfig().RemoteUuid)
	reply.request = false
	reply.reply = true
	return reply
}

func (this *Message) FailClone(failMessage string) ifs.IMessage {
	fail := this.Clone()
	fail.failMessage = failMessage
	copy(fail.source[0:36], this.destination)
	fail.destination = string(this.source[0:36])
	return fail
}

func HeaderOf(data []byte) (string, string, string, string, uint16, ifs.Priority) {

	size := nets.Bytes2UInt16(data[POS_Service_Name : POS_Service_Name+2])
	POS_Sequence := POS_Service_Name + 2 + int(size)
	POS_Priority := POS_Sequence + 4

	destSize := data[POS_Destination]
	destination := ""
	if destSize != 0 {
		destination = string(data[POS_Destination+1 : POS_Service_Area])
	}

	return string(data[POS_Source:POS_Vnet]),
		string(data[POS_Vnet:POS_Destination]),
		destination,
		string(data[POS_Service_Name+2 : POS_Sequence]),
		nets.Bytes2UInt16(data[POS_Service_Area:POS_Service_Name]),
		ifs.Priority(data[POS_Priority])
}
