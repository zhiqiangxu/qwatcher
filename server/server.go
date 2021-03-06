package server

import (
	"errors"

	"github.com/zhiqiangxu/qrpc"
)

const (
	// RegCmd for register
	RegCmd qrpc.Cmd = iota
	// RegRespCmd is resp for RegCmd
	RegRespCmd
	// LWCmd for list and watch
	LWCmd
	// LWRespCmd is resp for LWCmd
	LWRespCmd
	// LWPushRespCmd is resp for LWCmd in push mode
	LWPushRespCmd
	// JoinCmd tries to add a new node to raft cluster(it will forward to leader if needed)
	JoinCmd
	// JoinRespCmd is resp for join
	JoinRespCmd
	// SetAPIAddrCmd should be called to leader, it won't forward request
	SetAPIAddrCmd
)

var (
	// ErrNotLeader when not leader
	ErrNotLeader = errors.New("not leader")
	// ErrLeaderAPINA when leader not available
	ErrLeaderAPINA = errors.New("leader not available")
)
