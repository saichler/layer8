package switching

import (
	"github.com/saichler/layer8/go/overlay/edge"
	"github.com/saichler/shared/go/share/interfaces"
)

func (switchService *SwitchService) ConnectTo(host string, destPort uint32) error {
	// Dial the destination and validate the secret and key
	conn, err := interfaces.SecurityProvider().CanDial(host, destPort)
	if err != nil {
		return err
	}

	config := interfaces.SwitchConfig()
	config.SwitchPort = destPort
	config.Local_Uuid = switchService.switchConfig.Local_Uuid
	config.IsSwitchSide = true
	config.IsAdjacentASwitch = true

	err = interfaces.SecurityProvider().ValidateConnection(conn, config)
	if err != nil {
		return err
	}

	edge := edge.NewEdgeImpl(conn, switchService, switchService.registry, nil, config)

	//Below attributes are only for the port initiating the connection
	/* @TODO implement reconnect between switches
	edge.reconnectInfo = &ReconnectInfo{
		host:         host,
		port:         destPort,
		reconnectMtx: &sync.Mutex{},
	} */

	//We have only one go routing per each because we want to keep the order of incoming and outgoing messages
	edge.Start()

	switchService.notifyNewEdge(edge)
	return nil
}
