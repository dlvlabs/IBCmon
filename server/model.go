package server

func (server *Server) QueryIBCInfo() IBCInfos {
	ibcInfos := make(IBCInfos, 0, len(server.Store.IBCInfo))

	for chainId, clients := range server.Store.IBCInfo {
		for clientId, client := range clients {
			for connectionId, channels := range client.Connections {
				for channelId, channel := range channels {
					source := newIBC(
						chainId, clientId, connectionId,
						channelId, channel.PortId,
					)
					destination := newIBC(
						client.ChainId, channel.Counterparty.ClientId, channel.Counterparty.ConnectionId,
						channel.Counterparty.ChannelId, channel.Counterparty.PortId,
					)
					ibcInfos = append(ibcInfos, IBCInfo{
						Updated: server.Store.Updated,

						Source:      source,
						Destination: destination,
					})
				}
			}
		}
	}

	return ibcInfos
}

func (server *Server) QueryClientHealth() ClientHealths {
	clientHealths := make(ClientHealths, 0, len(server.Store.IBCInfo))

	for chainId, clients := range server.Store.IBCInfo {
		for clientId, client := range clients {
			clientHealths = append(clientHealths, ClientHealth{
				Health:        client.Health,
				ClientUpdated: client.ClientUpdated,

				Source:      chainId,
				Destination: client.ChainId,

				ClientId:       clientId,
				TrustingPeriod: client.TrustingPeriod,
			})
		}
	}

	return clientHealths
}

func (server *Server) QueryIBCPacket() IBCPackets {
	ibcPackets := make(IBCPackets, 0, len(server.Store.IBCInfo))

	for chainId, clients := range server.Store.IBCInfo {
		for clientId, client := range clients {
			for connectionId, channels := range client.Connections {
				for channelId, channel := range channels {
					// exception handling for chains not support `NextSequenceSend` query
					if channel.IBCPacketTracker == nil {
						continue
					}

					latestSucceedPackets := make(map[string]SucceedPacket)
					for packetType, succeedPacket := range channel.IBCPacketTracker.LatestSucceedPackets {
						latestSucceedPackets[packetType] = SucceedPacket{
							Hash:     succeedPacket.Hash,
							Sequence: succeedPacket.Sequence,
							Data:     succeedPacket.Data,
						}
					}

					source := newIBC(
						chainId, clientId, connectionId,
						channelId, channel.PortId,
					)
					destination := newIBC(
						client.ChainId, channel.Counterparty.ClientId, channel.Counterparty.ConnectionId,
						channel.Counterparty.ChannelId, channel.Counterparty.PortId,
					)
					ibcPackets = append(ibcPackets, IBCPacket{
						Updated: channel.IBCPacketTracker.Updated,

						Health: channel.IBCPacketTracker.Health,

						Source:      source,
						Destination: destination,

						Sequence:          channel.IBCPacketTracker.Sequence,
						ConsecutiveMissed: channel.IBCPacketTracker.MissedCnt,

						LatestSucceedPackets: latestSucceedPackets,
					})
				}
			}
		}
	}

	return ibcPackets
}
