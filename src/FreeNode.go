package src

import (
	"net"
)

// FreeNode ... Structure of each Freenet node
type FreeNode struct {
	id           uint32        // My node ID
	address      net.IP        // IP address
	subscribedTo uint32        // ID of node I sent last join request to
	randomSeed   uint32        // Random seed for key
	key          uint32        // My key
	joinRequests []JoinRequest // Join requests being processed
}

// JoinRequest ... State maintained for every join request
type JoinRequest struct {
	joiningNodeID    uint32
	prevMemberNodeID uint32
	nextMemberNodeID uint32
	hopsToLive       int
}

// joinFreeNet ... Send join request to an existing FreeNet Node
func (f *FreeNode) joinFreeNet(id uint32, address net.IP) bool {
	return true
}

// processJoinAnnouncement ... Process a join request anywhere in the chain
func (f *FreeNode) processJoinAnnouncement(prevMemberNodeID uint32, joiningNodeID uint32, receivedSeed uint32, hopsToLive int) bool {
	return true
}

// addJoiningNode ... Add the newly joining node anywhere in the chain
func (f *FreeNode) addJoiningNode(joiningNodeID uint32, senderNodeID uint32, key uint32) bool {
	return true
}
