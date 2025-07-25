package vnic

import (
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types"
	"github.com/saichler/layer8/go/overlay/health"
	"google.golang.org/protobuf/proto"
	"reflect"
)

func (this *VirtualNetworkInterface) Unicast(destination, serviceName string, serviceArea byte,
	action ifs.Action, any interface{}) error {
	if destination == "" {
		destination = ifs.DESTINATION_Single
	}
	elems, err := createElements(any, this.resources)
	if err != nil {
		return err
	}
	return this.components.TX().Unicast(destination, serviceName, serviceArea, action, elems, 0,
		false, false, this.protocol.NextMessageNumber(), ifs.Empty, "", "", -1, "")
}

func (this *VirtualNetworkInterface) Request(destination, serviceName string, serviceArea byte,
	action ifs.Action, any interface{}, tokens ...string) ifs.IElements {

	if destination == "" {
		destination = ifs.DESTINATION_Single
	}

	request := this.requests.NewRequest(this.protocol.NextMessageNumber(), this.resources.SysConfig().LocalUuid, 5, this.resources.Logger())

	request.Lock()
	defer request.Unlock()

	elements, err := createElements(any, this.resources)
	if err != nil {
		return object.NewError(err.Error())
	}
	token := ""
	if tokens != nil && len(tokens) > 0 {
		token = tokens[0]
	}
	e := this.components.TX().Unicast(destination, serviceName, serviceArea, action, elements, 0,
		true, false, request.MsgNum(), ifs.Empty, "", "", -1, token)
	if e != nil {
		return object.NewError(e.Error())
	}
	request.Wait()
	return request.Response()
}

func (this *VirtualNetworkInterface) Reply(msg *ifs.Message, response ifs.IElements) error {
	reply := msg.CloneReply(this.resources.SysConfig().LocalUuid, this.resources.SysConfig().RemoteUuid)
	data, e := this.protocol.CreateMessageForm(reply, response)
	if e != nil {
		this.resources.Logger().Error(e)
		return e
	}
	hc := health.Health(this.resources)
	hp := hc.Health(msg.Source())
	alias := " No Alias Yet"
	if hp != nil {
		alias = hp.Alias
	}
	this.resources.Logger().Debug("Replying to ", msg.Source(), " ", alias)
	return this.SendMessage(data)
}

func (this *VirtualNetworkInterface) Multicast(serviceName string, serviceArea byte, action ifs.Action, any interface{}) error {
	elems, err := createElements(any, this.resources)
	if err != nil {
		return err
	}
	return this.components.TX().Multicast("", serviceName, serviceArea, action, elems, 0,
		false, false, this.protocol.NextMessageNumber(), ifs.Empty, "", "", -1, "")
}

func (this *VirtualNetworkInterface) Single(serviceName string, serviceArea byte, action ifs.Action, any interface{}) (string, error) {
	hc := health.Health(this.resources)
	destination := hc.DestinationFor(serviceName, serviceArea, this.resources.SysConfig().LocalUuid, false, false)
	if destination == "" {
		destination = ifs.DESTINATION_Single
	}

	hp := hc.Health(destination)
	alias := "Unknown Yet"
	if hp != nil {
		alias = hp.Alias
	}
	this.Resources().Logger().Info("Sending Single to ", destination, " alias ", alias)

	return destination, this.Unicast(destination, serviceName, serviceArea, action, any)
}

func (this *VirtualNetworkInterface) SingleRequest(serviceName string, serviceArea byte, action ifs.Action, any interface{}) ifs.IElements {
	hc := health.Health(this.resources)
	destination := hc.DestinationFor(serviceName, serviceArea, this.resources.SysConfig().LocalUuid, false, false)
	if destination == "" {
		destination = ifs.DESTINATION_Single
	}

	hp := hc.Health(destination)
	alias := "Unknown Yet"
	if hp != nil {
		alias = hp.Alias
	}
	this.Resources().Logger().Info("Sending Single Request to ", destination, " alias ", alias)
	return this.Request(destination, serviceName, serviceArea, action, any)
}

func (this *VirtualNetworkInterface) Leader(serviceName string, serviceArea byte, action ifs.Action, any interface{}) ifs.IElements {
	hc := health.Health(this.resources)
	destination := hc.DestinationFor(serviceName, serviceArea, this.resources.SysConfig().LocalUuid, false, true)
	if destination == "" {
		destination = ifs.DESTINATION_Leader
	}
	return this.Request(destination, serviceName, serviceArea, action, any)
}

func (this *VirtualNetworkInterface) Forward(msg *ifs.Message, destination string) ifs.IElements {
	pb, err := this.protocol.ElementsOf(msg)
	if err != nil {
		return object.NewError(err.Error())
	}

	request := this.requests.NewRequest(this.protocol.NextMessageNumber(), this.resources.SysConfig().LocalUuid, 5, this.resources.Logger())
	request.Lock()
	defer request.Unlock()

	e := this.components.TX().Unicast(destination, msg.ServiceName(), msg.ServiceArea(), msg.Action(),
		pb, 0, true, false, request.MsgNum(),
		msg.Tr_State(), msg.Tr_Id(), msg.Tr_ErrMsg(), msg.Tr_StartTime(), msg.AAAId())
	if e != nil {
		return object.NewError(e.Error())
	}
	request.Wait()
	return request.Response()
}

func createElements(any interface{}, resources ifs.IResources) (ifs.IElements, error) {
	if any == nil {
		return object.New(nil, nil), nil
	}
	pq, ok := any.(*types.Query)
	if ok {
		return object.NewQuery(pq.Text, resources)
	}

	gsql, ok := any.(string)
	if ok {
		return object.NewQuery(gsql, resources)
	}

	elems, ok := any.(ifs.IElements)
	if ok {
		return elems, nil
	}

	pb, ok := any.(proto.Message)
	if ok {
		return object.New(nil, pb), nil
	}

	v := reflect.ValueOf(any)

	if v.Kind() == reflect.Slice {
		pbs := make([]proto.Message, v.Len())
		for i := 0; i < v.Len(); i++ {
			elm := v.Index(i)
			pb, ok = elm.Interface().(proto.Message)
			if ok {
				pbs[i] = pb
			} else {
				panic("Uknown input type " + reflect.ValueOf(pb).String())
			}
		}
		return object.New(nil, pbs), nil
	}
	panic("Uknown input type " + reflect.ValueOf(any).String())
}
