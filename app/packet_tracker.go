package app

import (
	"context"
	"fmt"
	"time"

	"github.com/dlvlabs/ibcmon/alert"
	"github.com/dlvlabs/ibcmon/client/grpc"
	"github.com/dlvlabs/ibcmon/client/rpc"
	"github.com/dlvlabs/ibcmon/logger"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type PacketTypes int

const (
	PACKET_STATUS_SEND PacketTypes = iota
	PACKET_STATUS_RECV
	PACKET_STATUS_ACK
)

func (ps PacketTypes) String() string {
	switch ps {
	case PACKET_STATUS_SEND:
		return "send_packet"
	case PACKET_STATUS_RECV:
		return "recv_packet"
	case PACKET_STATUS_ACK:
		return "acknowledge_packet"
	default:
		return "unknown"
	}
}

type (
	IBCPacketTracker struct {
		Updated time.Time

		Health bool

		PacketType PacketTypes
		Sequence   uint64
		timeout    Timeout

		LatestSucceedPackets SucceedPackets
		MissedCnt            uint64

		Source      Chain
		Destination Chain
	}
	Timeout struct {
		timeoutHeight    uint64
		timeoutTimestamp int64
	}
	Chain struct {
		rpc       *rpc.Client
		grpc      *grpc.Client
		ChainId   string
		ChannelId string
		PortId    string
	}
	// PacketType => SucceedPacket
	SucceedPackets map[string]SucceedPacket
	SucceedPacket  struct {
		Hash     string
		Sequence uint64
		Data     string
	}
)

func NewIBCPacketTracker(
	sequence uint64,

	srcRPC *rpc.Client, srcGRPC *grpc.Client,
	srcChainId, srcChannelId, srcPortId string,

	dstRPC *rpc.Client, dstGRPC *grpc.Client,
	dstChainId, dstChannelId, dstPortId string,
) *IBCPacketTracker {
	return &IBCPacketTracker{
		Updated: time.Now().UTC(),

		Health: true,

		PacketType: PACKET_STATUS_SEND,
		Sequence:   sequence,
		timeout: Timeout{
			timeoutHeight:    0,
			timeoutTimestamp: 0,
		},

		LatestSucceedPackets: make(SucceedPackets),
		MissedCnt:            0,

		Source: Chain{
			rpc:       srcRPC,
			grpc:      srcGRPC,
			ChainId:   srcChainId,
			ChannelId: srcChannelId,
			PortId:    srcPortId,
		},

		Destination: Chain{
			rpc:       dstRPC,
			grpc:      dstGRPC,
			ChainId:   dstChainId,
			ChannelId: dstChannelId,
			PortId:    dstPortId,
		},
	}
}

func (ibcPacketTracker *IBCPacketTracker) GetPacketStatus() string {
	return ibcPacketTracker.PacketType.String()
}
func (ibcPacketTracker *IBCPacketTracker) GetSequence() uint64 {
	return ibcPacketTracker.Sequence
}
func (ibcPacketTracker *IBCPacketTracker) GetSrcInfo() (string, string, string) {
	return ibcPacketTracker.Source.ChainId, ibcPacketTracker.Source.ChannelId, ibcPacketTracker.Source.PortId
}
func (ibcPacketTracker *IBCPacketTracker) GetDstInfo() (string, string, string) {
	return ibcPacketTracker.Destination.ChainId, ibcPacketTracker.Destination.ChannelId, ibcPacketTracker.Destination.PortId
}

func (ibcPacketTracker *IBCPacketTracker) String() string {
	return fmt.Sprintf(
		"%s(%d) for %s(%s/%s) => %s(%s/%s)",
		ibcPacketTracker.PacketType.String(), ibcPacketTracker.Sequence,
		ibcPacketTracker.Source.ChainId, ibcPacketTracker.Source.ChannelId, ibcPacketTracker.Source.PortId,
		ibcPacketTracker.Destination.ChainId, ibcPacketTracker.Destination.ChannelId, ibcPacketTracker.Destination.PortId,
	)
}

func (ibcPacketTracker *IBCPacketTracker) transitStatus(timeoutHeight uint64, timeoutTimestamp int64) {
	ibcPacketTracker.timeout.timeoutHeight = timeoutHeight
	ibcPacketTracker.timeout.timeoutTimestamp = timeoutTimestamp

	if ibcPacketTracker.PacketType == PACKET_STATUS_ACK {
		ibcPacketTracker.Health = true
		ibcPacketTracker.Sequence++
		ibcPacketTracker.MissedCnt = 0
	}
	ibcPacketTracker.PacketType = (ibcPacketTracker.PacketType + 1) % 3

	ibcPacketTracker.Updated = time.Now().UTC()
}

func (ibcPacketTracker *IBCPacketTracker) isTimeout(ctx context.Context) (bool, error) {
	if ibcPacketTracker.PacketType != PACKET_STATUS_RECV {
		return false, nil
	}

	if ibcPacketTracker.timeout.timeoutHeight != 0 {
		height, err := ibcPacketTracker.Destination.rpc.GetLatestBlockHeight(ctx)
		if err != nil {
			return false, err
		}

		msg := fmt.Sprintf("timeoutHeight: %d <= current height: %d", ibcPacketTracker.timeout.timeoutHeight, height)
		logger.Debug(msg)
		return ibcPacketTracker.timeout.timeoutHeight <= uint64(height), nil
	}

	msg := fmt.Sprintf("timeout timestamp: %d <= now: %d", ibcPacketTracker.timeout.timeoutTimestamp, time.Now().UnixNano())
	logger.Debug(msg)
	return ibcPacketTracker.timeout.timeoutTimestamp <= time.Now().UnixNano(), nil
}

func (app *App) trackIBCPacket(ctx context.Context) error {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(errors.New("terminate ibc packet tracker"))

	err := app.connectGRPCs()
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)

	for chainId, clients := range app.Store.IBCInfo {
		for _, client := range clients {
			for _, channels := range client.Connections {
				for channelId, channel := range channels {
					grpcClient := app.grpcs[chainId]
					nextSequence, err := grpcClient.GetNextSequenceSend(ctx, channelId, channel.PortId)
					if err != nil {
						if errors.Is(errors.Cause(err), grpc.UNIMPLMENTED) {
							msg := fmt.Sprintf(
								"not support `NextSequenceSend` query, skip tracking IBC packet for %s(%s/%s) => %s(%s/%s)",
								chainId, channelId, channel.PortId,
								client.ChainId, channel.Counterparty.ChannelId, channel.Counterparty.PortId,
							)
							logger.Info(msg)

							continue
						}

						return err
					}

					ibcPacketTracker := NewIBCPacketTracker(
						nextSequence,

						app.rpcs[chainId], app.grpcs[chainId],
						chainId, channelId, channel.PortId,

						app.rpcs[client.ChainId], app.grpcs[client.ChainId],
						client.ChainId, channel.Counterparty.ChannelId, channel.Counterparty.PortId,
					)

					channel.IBCPacketTracker = ibcPacketTracker

					g.Go(func() error {
						ticker := time.NewTicker(app.cfg.General.PacketTrackingInterval)
						defer ticker.Stop()

						for {
							select {
							case <-ticker.C:
								missed, err := ibcPacketTracker.track(ctx)
								if err != nil {
									err := errors.Wrapf(err, "track ibc packet stopped: %s", ibcPacketTracker.String())
									logger.Error(err)

									// All ibcPacketTracker should be stopped
									cancel(errors.New("terminate ibc packet tracker"))

									return err
								}

								if missed {
									ibcPacketTracker.PacketType = PACKET_STATUS_SEND
									ibcPacketTracker.Sequence++

									ibcPacketTracker.MissedCnt++
									if ibcPacketTracker.MissedCnt >= app.cfg.Rule.ConsecutiveMissedPackets {
										ibcPacketTracker.Health = false
									}

									ibcPacketTracker.Updated = time.Now().UTC()

									msg := fmt.Sprintf("missed %d ibc tx: %s", ibcPacketTracker.MissedCnt, ibcPacketTracker.String())
									logger.Warn(msg)

									// TODO: consider remove this line
									alert.SendTg(msg)
								}

							case <-ctx.Done():
								return context.Cause(ctx)
							}
						}
					})

					msg := fmt.Sprintf("start tracking ibc packet: %s", ibcPacketTracker.String())
					logger.Info(msg)
				}
			}
		}
	}

	err = app.terminateGRPCs()
	if err != nil {
		logger.Error(err)
		return err
	}

	return g.Wait()
}

// if packet is missed return true
func (ibcPacketTracker *IBCPacketTracker) track(ctx context.Context) (bool, error) {
	rpc := ibcPacketTracker.Source.rpc
	if ibcPacketTracker.PacketType == PACKET_STATUS_RECV {
		rpc = ibcPacketTracker.Destination.rpc
	}
	code, hash, data, timeoutHeight, timeoutTimestamp, isFound, err := rpc.SearchIBCPacket(ctx, ibcPacketTracker, 0)
	if err != nil {
		return false, err
	} else if !isFound {
		// timeout
		timeout, err := ibcPacketTracker.isTimeout(ctx)
		if err != nil {
			return false, err
		}
		if timeout {
			msg := fmt.Sprintf("timeout ibc tx: %s", ibcPacketTracker.String())
			logger.Debug(msg)

			return true, nil
		}

		// not found
		msg := fmt.Sprintf("no ibc tx: %s", ibcPacketTracker.String())
		logger.Debug(msg)

		return false, nil
	}
	if code != 0 {
		msg := fmt.Sprintf("ibc tx not successed: %s", ibcPacketTracker.String())
		logger.Debug(msg)

		return true, nil
	}

	msg := fmt.Sprintf("ibc packet succeed: %s", ibcPacketTracker.String())
	logger.Info(msg)

	ibcPacketTracker.LatestSucceedPackets[ibcPacketTracker.PacketType.String()] = SucceedPacket{
		Hash:     hash,
		Sequence: ibcPacketTracker.Sequence,
		Data:     data,
	}

	ibcPacketTracker.transitStatus(timeoutHeight, timeoutTimestamp)

	return false, nil
}
