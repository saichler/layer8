package protocol

import (
	"github.com/saichler/types/go/nets"
	"github.com/saichler/types/go/types"
)

const (
	UNICAST_ADDRESS_SIZE = 36
	INT_SIZE             = 4
	PRIORITY_SIZE        = 1
	HEADER_SIZE          = UNICAST_ADDRESS_SIZE*3 + INT_SIZE + PRIORITY_SIZE
	SOURCE_VNET_POS      = UNICAST_ADDRESS_SIZE
	TOPIC_POS            = UNICAST_ADDRESS_SIZE * 2
	VLAN_POS             = UNICAST_ADDRESS_SIZE * 3
	PRIORITY_POS         = UNICAST_ADDRESS_SIZE*3 + 4
)

func CreateHeader(msg *types.Message) []byte {
	header := make([]byte, HEADER_SIZE)
	for i, c := range msg.SourceUuid {
		header[i] = byte(c)
	}
	for i, c := range msg.SourceVnetUuid {
		header[i+SOURCE_VNET_POS] = byte(c)
	}
	for i, c := range msg.Topic {
		header[i+TOPIC_POS] = byte(c)
	}
	vlan := nets.Int2Bytes(msg.Vlan)
	for i := 0; i < len(vlan); i++ {
		header[VLAN_POS+i] = vlan[i]
	}
	header[PRIORITY_POS] = byte(msg.Priority)
	return header
}

func HeaderOf(data []byte) (string, string, string, int32, types.Priority) {
	// Source will always be Uuid
	source := string(data[0:UNICAST_ADDRESS_SIZE])
	// Source vnet will always be Uuid
	sourceVnet := string(data[SOURCE_VNET_POS:TOPIC_POS])
	//Topic, either than being a uuid can also be a topic/multicast so it might not be full 16 bytes
	topic := make([]byte, 0)
	for i := TOPIC_POS; i < TOPIC_POS+UNICAST_ADDRESS_SIZE; i++ {
		if data[i] == 0 {
			break
		}
		topic = append(topic, data[i])
	}
	vlan := nets.Bytes2Int(data[VLAN_POS : VLAN_POS+4])
	pri := types.Priority(data[PRIORITY_POS])
	return source, sourceVnet, string(topic), vlan, pri
}
