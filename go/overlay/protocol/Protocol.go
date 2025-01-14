package protocol

import (
	"github.com/saichler/serializer/go/serializers"
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/types"
	"google.golang.org/protobuf/proto"
	"reflect"
	"sync/atomic"
)

type Protocol struct {
	sequence   atomic.Int32
	serializer interfaces.Serializer
	providers  *interfaces.Providers
}

func New(providers *interfaces.Providers,
	serializer interfaces.Serializer) *Protocol {
	p := &Protocol{}
	p.providers = providers
	p.serializer = serializer
	if p.serializer == nil {
		p.serializer = &serializers.ProtoBuffBinary{}
	}
	return p
}

func (this *Protocol) MessageOf(data []byte) (*types.Message, error) {
	msg, err := this.serializer.Unmarshal(data[109:], "Message", this.providers.Registry())
	if err != nil {
		panic(err)
	}
	return msg.(*types.Message), err
}

func (this *Protocol) ProtoOf(msg *types.Message) (proto.Message, error) {
	data, err := this.providers.Security().Decrypt(msg.Data)
	if err != nil {
		return nil, err
	}

	info, err := this.providers.Registry().TypeInfo(msg.Type)
	if err != nil {
		panic(err)
		return nil, interfaces.Error(err)
	}
	pbIns, err := info.NewInstance()
	if err != nil {
		return nil, err
	}

	pb := pbIns.(proto.Message)
	err = proto.Unmarshal(data, pb)
	return pb, err
}

func (this *Protocol) CreateMessageFor(priority types.Priority, action types.Action, source, sourceSwitch, dest string, pb proto.Message) ([]byte, error) {
	//first marshal the protobuf into bytes
	data, err := this.serializer.Marshal(pb, nil)
	if err != nil {
		return nil, err
	}
	//Encode the data
	encData, err := this.providers.Security().Encrypt(data)
	if err != nil {
		return nil, err
	}
	//create the wrapping message for the destination
	msg := &types.Message{}
	msg.SourceUuid = source
	msg.SourceSwitchUuid = sourceSwitch
	msg.Destination = dest
	msg.Sequence = this.sequence.Add(1)
	msg.Priority = priority
	msg.Data = encData
	msg.Type = reflect.ValueOf(pb).Elem().Type().Name()
	msg.Action = action
	//Now serialize the message
	msgData, err := this.serializer.Marshal(msg, nil)
	if err != nil {
		return nil, err
	}
	//Create the header for the switch
	header := CreateHeader(msg)
	//Append the msgData to the header
	header = append(header, msgData...)
	return header, nil
}

func (this *Protocol) Serializer() interfaces.Serializer {
	return this.serializer
}

func (this *Protocol) Providers() *interfaces.Providers {
	return this.providers
}
