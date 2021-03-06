package internalcmd

import (
	"github.com/zhiqiangxu/qrpc"
	"github.com/zhiqiangxu/qwatch/pkg/bson"
	"github.com/zhiqiangxu/qwatch/pkg/entity"
	"github.com/zhiqiangxu/qwatch/pkg/logger"
	"github.com/zhiqiangxu/qwatch/server"
	"github.com/zhiqiangxu/qwatch/store"
)

// JoinCmd do join
type JoinCmd struct {
	store *store.Store
}

// NewJoinCmd returns a JoinCmd
func NewJoinCmd(store *store.Store) *JoinCmd {
	return &JoinCmd{store: store}
}

// ServeQRPC implements qrpc.Handler
func (cmd *JoinCmd) ServeQRPC(writer qrpc.FrameWriter, frame *qrpc.RequestFrame) {
	var joinCmd entity.JoinCmd
	err := bson.FromBytes(frame.Payload, &joinCmd)
	if err != nil {
		frame.Close()
		return
	}

	if cmd.store.IsLeader() {
		err = cmd.store.Join(joinCmd.NodeID, joinCmd.RaftAddr)
		if err != nil {
			cmd.writeResp(writer, frame, err)
		}

		return
	}

	leaderAPIAddr := cmd.store.LeaderAPIAddr()
	if leaderAPIAddr == "" {
		cmd.writeResp(writer, frame, server.ErrLeaderAPINA)
		return
	}

	err = store.JoinPeerByQrpc(leaderAPIAddr, joinCmd.NodeID, joinCmd.RaftAddr)
	cmd.writeResp(writer, frame, err)
}

func (cmd *JoinCmd) writeResp(writer qrpc.FrameWriter, frame *qrpc.RequestFrame, err error) {
	var resp entity.JoinResp
	if err == nil {
		resp.OK = true
	} else {
		resp.OK = false
		resp.Msg = err.Error()
	}

	bytes, err := bson.ToBytes(resp)
	if err != nil {
		logger.Error("ToBytes", err)
		frame.Close()
		return
	}

	writer.StartWrite(frame.RequestID, server.JoinRespCmd, 0)
	writer.WriteBytes(bytes)
	writer.EndWrite()
}
