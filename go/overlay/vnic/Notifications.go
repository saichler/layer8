package vnic

import (
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types"
	"github.com/saichler/layer8/go/overlay/health"
	"github.com/saichler/layer8/go/overlay/protocol"
	"time"
)

func (this *VirtualNetworkInterface) NotifyServiceAdded(serviceNames []string, serviceArea byte) error {
	hc := health.Health(this.resources)
	curr := hc.Health(this.resources.SysConfig().LocalUuid)
	hp := &types.Health{}
	hp.AUuid = curr.AUuid
	hp.Services = curr.Services
	mergeServices(hp, this.resources.SysConfig().Services)
	//send notification for health service
	err := this.Unicast(this.resources.SysConfig().RemoteUuid, health.ServiceName, 0, ifs.PATCH, hp)
	for _, serviceName := range serviceNames {
		{
			go this.requestCacheSync(serviceName, serviceArea)
		}
	}
	return err
}

func (this *VirtualNetworkInterface) requestCacheSync(serviceName string, serviceArea byte) {
	time.Sleep(time.Second)
	err := this.Multicast(serviceName, serviceArea, ifs.Sync, object.New(nil, nil))
	if err != nil {
		this.resources.Logger().Error("Failed to send cache sync multicast:", err.Error())
	}
}

func (this *VirtualNetworkInterface) NotifyServiceRemoved(serviceName string, serviceArea byte) error {
	hc := health.Health(this.resources)
	curr := hc.Health(this.resources.SysConfig().LocalUuid)
	hp := &types.Health{}
	hp.AUuid = curr.AUuid
	hp.Services = curr.Services
	mergeServices(hp, this.resources.SysConfig().Services)
	ifs.RemoveService(hp.Services, serviceName, int32(serviceArea))
	return this.Unicast(this.resources.SysConfig().RemoteUuid, health.ServiceName, serviceArea, ifs.PATCH, hp)
}

func (this *VirtualNetworkInterface) PropertyChangeNotification(set *types.NotificationSet) {
	protocol.AddPropertyChangeCalled(set, this.resources.SysConfig().LocalAlias)
	this.Multicast(set.ServiceName, byte(set.ServiceArea), ifs.Notify, set)
}

func mergeServices(hp *types.Health, services *types.Services) {
	if hp.Services == nil {
		hp.Services = services
		return

	}
	for serviceName, serviceAreas := range services.ServiceToAreas {
		_, ok := hp.Services.ServiceToAreas[serviceName]
		if !ok {
			hp.Services.ServiceToAreas[serviceName] = serviceAreas
			continue
		}
		if hp.Services.ServiceToAreas[serviceName].Areas == nil {
			hp.Services.ServiceToAreas[serviceName].Areas = serviceAreas.Areas
			continue
		}
		for svArea, score := range serviceAreas.Areas {
			serviceArea := svArea
			hp.Services.ServiceToAreas[serviceName].Areas[serviceArea] = score
		}
	}
}
