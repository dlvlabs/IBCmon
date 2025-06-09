package rpc

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dlvlabs/ibcmon/client/rpc/exported"
	"github.com/dlvlabs/ibcmon/logger"
	"github.com/pkg/errors"

	coreTypes "github.com/cometbft/cometbft/rpc/core/types"
	clientTypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
)

// return code, hash, data, timeoutHeight, timeoutTimestamp, isFound, err
func (c *Client) SearchIBCPacket(ctx context.Context, ibcPacketTracker exported.IBCPacketTracker) (uint32, string, string, uint64, int64, bool, error) {
	packet := ibcPacketTracker.GetPacketStatus()
	sequence := ibcPacketTracker.GetSequence()
	_, srcChannelId, srcPortId := ibcPacketTracker.GetSrcInfo()
	_, dstChannelId, dstPortId := ibcPacketTracker.GetDstInfo()

	query := fmt.Sprintf(
		"%s.packet_sequence='%d' AND %s.packet_src_channel='%s' AND %s.packet_src_port='%s' AND %s.packet_dst_channel='%s' AND %s.packet_dst_port='%s'",
		packet, sequence,
		packet, srcChannelId, packet, srcPortId,
		packet, dstChannelId, packet, dstPortId,
	)

	resp, err := c.rpcClient.TxSearch(ctx, query, false, nil, nil, "asc")
	if err != nil {
		return 1, "", "", 0, 0, false, errors.Wrapf(err, "failed to search tx: %s", query)
	} else if len(resp.Txs) == 0 {
		return 0, "", "", 0, 0, false, nil
	}

	var data string = ""
	var timeoutHeight uint64 = 0
	var timeoutTimestamp int64 = 0
	for _, event := range resp.Txs[0].TxResult.Events {
		for _, attr := range event.Attributes {
			switch tryBase64Decoding(attr.Key) {
			case "packet_data":
				data = tryBase64Decoding(attr.Value)
			case "packet_timeout_height":
				timeoutHeight = clientTypes.MustParseHeight(tryBase64Decoding(attr.Value)).GetRevisionHeight()
			case "packet_timeout_timestamp":
				timeoutTimestamp, err = strconv.ParseInt(tryBase64Decoding(attr.Value), 10, 64)
				if err != nil {
					err = errors.Wrapf(err, "failed to parse timeout timestamp: %+v", event.Attributes)
					logger.Error(err)

					panic(err)
				}
			}
		}
	}

	return resp.Txs[0].TxResult.Code, resp.Txs[0].Hash.String(), data, timeoutHeight, timeoutTimestamp, true, nil
}

func (c *Client) GetLatestBlockHeight(ctx context.Context) (int64, error) {
	abciInfo, err := c.rpcClient.ABCIInfo(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get ABCI info")
	}

	return abciInfo.Response.LastBlockHeight, nil
}

func (c *Client) Subscribe(ctx context.Context, query string) (<-chan coreTypes.ResultEvent, error) {
	resultEvent, err := c.rpcClient.Subscribe(ctx, "subscribe", query)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to subscribe to query: %s", query)
	}

	return resultEvent, nil
}
