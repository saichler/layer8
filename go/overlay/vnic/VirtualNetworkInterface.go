package vnic

import (
	"errors"
	"github.com/saichler/layer8/go/overlay/health"
	"github.com/saichler/layer8/go/overlay/protocol"
	"github.com/saichler/shared/go/share/strings"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
	"net"
	"sync"
	"time"
)

type VirtualNetworkInterface struct {
	// Resources for this VNic such as registry, security & config
	resources common.IResources
	// The socket connection
	conn net.Conn
	// The socket connection mutex
	connMtx *sync.Mutex
	// is running
	running bool
	// Sub components/go routines
	components *SubComponents
	// The Protocol
	protocol *protocol.Protocol
	// Name for this VNic expressing the connection path in aside -->> zside
	name string
	// Indicates if this vnic in on the switch internal, hence need no keep alive
	IsVNet bool
	// Last reconnect attempt
	last_reconnect_attempt int64

	requests *Requests

	stats *types.HealthPointStats
}

func NewVirtualNetworkInterface(resources common.IResources, conn net.Conn) *VirtualNetworkInterface {
	vnic := &VirtualNetworkInterface{}
	vnic.conn = conn
	vnic.resources = resources
	vnic.connMtx = &sync.Mutex{}
	vnic.protocol = protocol.New(resources)
	vnic.components = newSubomponents()
	vnic.components.addComponent(newRX(vnic))
	vnic.components.addComponent(newTX(vnic))
	vnic.components.addComponent(newKeepAlive(vnic))
	vnic.requests = newRequests()
	vnic.stats = &types.HealthPointStats{}
	if vnic.resources.SysConfig().LocalUuid == "" {
		vnic.resources.SysConfig().LocalUuid = common.NewUuid()
	}

	if conn == nil {
		// Register the health service
		health.RegisterHealth(vnic.resources, nil)
	}

	return vnic
}

func (this *VirtualNetworkInterface) Start() {
	this.running = true
	if this.conn == nil {
		this.connectToSwitch()
	} else {
		this.receiveConnection()
	}
	this.name = this.resources.SysConfig().LocalAlias + " -->> " + this.resources.SysConfig().RemoteAlias
}

func (this *VirtualNetworkInterface) connectToSwitch() {
	this.connect()
	this.components.start()
}

func (this *VirtualNetworkInterface) connect() error {
	// Dial the destination and validate the secret and key
	destination := protocol.MachineIP
	if protocol.UsingContainers {
		// inside a containet the switch ip will be the external subnet + ".1"
		// for example if the address of the container is 172.1.1.112, the switch will be accessible
		// via 172.1.1.1
		subnet := protocol.IpSegment.ExternalSubnet()
		destination = subnet + ".1"
	}
	// Try to dial to the switch
	conn, err := this.resources.Security().CanDial(destination, this.resources.SysConfig().VnetPort)
	if err != nil {
		return errors.New("Error connecting to the vnet: " + err.Error())
	}
	// Verify that the switch accepts this connection
	if this.resources.SysConfig().LocalUuid == "" {
		panic("")
	}
	err = this.resources.Security().ValidateConnection(conn, this.resources.SysConfig())
	if err != nil {
		return errors.New("Error validating connection: " + err.Error())
	}
	this.conn = conn
	this.resources.SysConfig().Address = conn.LocalAddr().String()
	return nil
}

func (this *VirtualNetworkInterface) receiveConnection() {
	this.IsVNet = true
	this.resources.SysConfig().Address = this.conn.RemoteAddr().String()
	this.components.start()
}

func (this *VirtualNetworkInterface) Shutdown() {
	this.running = false
	if this.conn != nil {
		this.conn.Close()
	}
	this.components.shutdown()
	if this.resources.DataListener() != nil {
		go this.resources.DataListener().ShutdownVNic(this)
	}
}

func (this *VirtualNetworkInterface) Name() string {
	if this.name == "" {
		this.name = strings.New(this.resources.SysConfig().LocalUuid,
			" -->> ",
			this.resources.SysConfig().RemoteUuid).String()
	}
	return this.name
}

func (this *VirtualNetworkInterface) SendMessage(data []byte) error {
	return this.components.TX().SendMessage(data)
}

func (this *VirtualNetworkInterface) API(serviceName string, serviceArea uint16) common.API {
	return newAPI(serviceName, serviceArea, false, false)
}

func (this *VirtualNetworkInterface) Resources() common.IResources {
	return this.resources
}

func (this *VirtualNetworkInterface) reconnect() {
	if !this.running {
		return
	}
	this.connMtx.Lock()
	defer this.connMtx.Unlock()
	if time.Now().Unix()-this.last_reconnect_attempt < 5 {
		return
	}
	this.last_reconnect_attempt = time.Now().Unix()

	this.resources.Logger().Info("***** Trying to reconnect to ", this.resources.SysConfig().RemoteAlias, " *****")

	if this.conn != nil {
		this.conn.Close()
	}

	err := this.connect()
	if err != nil {
		this.resources.Logger().Error("***** Failed to reconnect to ", this.resources.SysConfig().RemoteAlias, " *****")
	} else {
		this.resources.Logger().Info("***** Reconnected to ", this.resources.SysConfig().RemoteAlias, " *****")
	}
}
