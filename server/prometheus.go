package server

import (
	"github.com/prometheus/client_golang/prometheus"
)

type IBCInfoCollector struct {
	server *Server

	Up *prometheus.Desc
}

func newIBCInfoCollector(server *Server) *IBCInfoCollector {
	labels := []string{"src_chain_id", "src_path", "dst_chain_id", "dst_path"}

	return &IBCInfoCollector{
		server: server,

		Up: prometheus.NewDesc(
			server.MetricPrefix+"_ibc_tao_up",
			"If 1 client, connection, channel is normal",
			labels, nil,
		),
	}
}

func (c *IBCInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Up
}

func (c *IBCInfoCollector) Collect(ch chan<- prometheus.Metric) {
	resp := c.server.QueryIBCInfo()

	for _, ibcInfo := range resp {
		labels := []string{
			ibcInfo.Source.ChainId,
			ibcInfo.Source.Path,
			ibcInfo.Destination.ChainId,
			ibcInfo.Destination.Path,
		}

		var up float64 = 1
		ch <- prometheus.MustNewConstMetric(
			c.Up,
			prometheus.GaugeValue,
			up,
			labels...,
		)
	}
}

type ClientHealthCollector struct {
	server *Server

	Health *prometheus.Desc
}

func newClientHealthCollector(server *Server) *ClientHealthCollector {
	labels := []string{"src_chain_id", "dst_chain_id", "client_id", "updated_at"}

	return &ClientHealthCollector{
		server: server,

		Health: prometheus.NewDesc(
			server.MetricPrefix+"_client_health",
			"Health status of the ibc client",
			labels, nil,
		),
	}
}

func (c *ClientHealthCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Health
}

func (c *ClientHealthCollector) Collect(ch chan<- prometheus.Metric) {
	resp := c.server.QueryClientHealth()

	for _, clientHealth := range resp {
		labels := []string{
			clientHealth.Source,
			clientHealth.Destination,
			clientHealth.ClientId,
			clientHealth.ClientUpdated.String(),
		}

		var health float64 = 0
		if clientHealth.Health {
			health = 1
		}

		ch <- prometheus.MustNewConstMetric(
			c.Health,
			prometheus.GaugeValue,
			float64(health),
			labels...,
		)
	}
}

type IBCPacketCollector struct {
	server *Server

	ChannelSequence                   *prometheus.Desc
	ConsecutiveMissed                 *prometheus.Desc
	ObservedSucceedSendPacketSequence *prometheus.Desc
	ObservedSucceedRecvPacketSequence *prometheus.Desc
	ObservedSucceedAckPacketSequence  *prometheus.Desc
}

func newIBCPacketCollector(server *Server) *IBCPacketCollector {
	labels := []string{"src_chain_id", "src_path", "dst_chain_id", "dst_path"}

	return &IBCPacketCollector{
		server: server,

		ChannelSequence: prometheus.NewDesc(
			server.MetricPrefix+"_channel_sequence",
			"Sequence number of the next packet to be sent on the channel",
			labels, nil,
		),
		ConsecutiveMissed: prometheus.NewDesc(
			server.MetricPrefix+"_consecutive_missed",
			"Number of consecutive ibc tx that have been missed",
			labels, nil,
		),
		ObservedSucceedSendPacketSequence: prometheus.NewDesc(
			server.MetricPrefix+"_observed_succeed_send_packet_sequence",
			"Sequence number of the last successfully sent packet",
			labels, nil,
		),
		ObservedSucceedRecvPacketSequence: prometheus.NewDesc(
			server.MetricPrefix+"_observed_succeed_recv_packet_sequence",
			"Sequence number of the last successfully received packet",
			labels, nil,
		),
		ObservedSucceedAckPacketSequence: prometheus.NewDesc(
			server.MetricPrefix+"_observed_succeed_ack_packet_sequence",
			"Sequence number of the last successfully acknowledged packet",
			labels, nil,
		),
	}
}

func (c *IBCPacketCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.ChannelSequence
	ch <- c.ConsecutiveMissed
	ch <- c.ObservedSucceedSendPacketSequence
	ch <- c.ObservedSucceedRecvPacketSequence
	ch <- c.ObservedSucceedAckPacketSequence
}

func (c *IBCPacketCollector) Collect(ch chan<- prometheus.Metric) {
	resp := c.server.QueryIBCPacket()

	for _, ibcPacket := range resp {
		labels := []string{
			ibcPacket.Source.ChainId,
			ibcPacket.Source.Path,
			ibcPacket.Destination.ChainId,
			ibcPacket.Destination.Path,
		}

		ch <- prometheus.MustNewConstMetric(
			c.ChannelSequence,
			prometheus.GaugeValue,
			float64(ibcPacket.Sequence),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.ConsecutiveMissed,
			prometheus.GaugeValue,
			float64(ibcPacket.ConsecutiveMissed),
			labels...,
		)

		for packetType, succeedPacket := range ibcPacket.LatestSucceedPackets {
			switch packetType {
			case "send_packet":
				ch <- prometheus.MustNewConstMetric(
					c.ObservedSucceedSendPacketSequence,
					prometheus.GaugeValue,
					float64(succeedPacket.Sequence),
					labels...,
				)
			case "recv_packet":
				ch <- prometheus.MustNewConstMetric(
					c.ObservedSucceedRecvPacketSequence,
					prometheus.GaugeValue,
					float64(succeedPacket.Sequence),
					labels...,
				)
			case "acknowledge_packet":
				ch <- prometheus.MustNewConstMetric(
					c.ObservedSucceedAckPacketSequence,
					prometheus.GaugeValue,
					float64(succeedPacket.Sequence),
					labels...,
				)
			}
		}
	}
}
