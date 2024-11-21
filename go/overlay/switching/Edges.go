package switching

import (
	"github.com/saichler/shared/go/share/interfaces"
	"sync"
)

type Edges struct {
	internal map[string]interfaces.IEdge
	external map[string]interfaces.IEdge
	mtx      *sync.RWMutex
}

func newEdges() *Edges {
	edges := &Edges{}
	edges.internal = make(map[string]interfaces.IEdge)
	edges.external = make(map[string]interfaces.IEdge)
	edges.mtx = &sync.RWMutex{}
	return edges
}

func (edges *Edges) addInternal(uuid string, edge interfaces.IEdge) {
	edges.mtx.Lock()
	defer edges.mtx.Unlock()
	exist, ok := edges.internal[uuid]
	if ok {
		exist.Shutdown()
	}
	edges.internal[uuid] = edge
}

func (edges *Edges) addExternal(uuid string, edge interfaces.IEdge) {
	edges.mtx.Lock()
	defer edges.mtx.Unlock()
	exist, ok := edges.external[uuid]
	if ok {
		exist.Shutdown()
	}
	edges.external[uuid] = edge
}

func (edges *Edges) getEdge(edgeUuid, remoteUuid string, isHope0 bool) (string, interfaces.IEdge) {
	edges.mtx.RLock()
	defer edges.mtx.RUnlock()
	edge, ok := edges.internal[edgeUuid]
	if ok {
		return edgeUuid, edge
	}
	edge, ok = edges.external[edgeUuid]
	if ok {
		return edgeUuid, edge
	}
	// only if this is hope0, e.g. the source of the message is from this switch sources,
	// fetch the external port.
	if isHope0 && remoteUuid != "" {
		edge, ok = edges.internal[remoteUuid]
		if ok {
			return remoteUuid, edge
		}
		edge, ok = edges.external[remoteUuid]
		if ok {
			return remoteUuid, edge
		}
	}
	return "", nil
}

func (edges *Edges) all() map[string]interfaces.IEdge {
	all := make(map[string]interfaces.IEdge)
	edges.mtx.RLock()
	defer edges.mtx.RUnlock()
	for uuid, edge := range edges.internal {
		all[uuid] = edge
	}
	for uuid, edge := range edges.external {
		all[uuid] = edge
	}
	return all
}

func (edges *Edges) removeExternals(uuids map[string]string) {
	edges.mtx.RLock()
	defer edges.mtx.RUnlock()
	for uuid, _ := range edges.external {
		delete(uuids, uuid)
	}
}