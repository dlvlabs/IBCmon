# Prometheus Metrics Documentation

## 1. IBCInfo

### Metric: `ibcmon_ibc_tao_up`

- **Type:** Gauge
- **Description:** Indicates if the IBC client, connection, and channel are normal (1 if normal).
- **Labels:**
  - `src_chain_id`: `ChainId` of the source chain
  - `src_path`: IBC path of source chain
  - `dst_chain_id`: `ChainId` of the destination chain
  - `dst_path`: IBC path of destination chain

**Example:**
```
ibcmon_ibc_tao_up{src_chain_id="milkyway", src_path="milkyway(07-tendermint-1/connection-0/channel-0/transfer)", dst_chain_id="osmosis-1", dst_path="osmosis-1(07-tendermint-3364/connection-2821/channel-89298/transfer)"} 1
```

---

## 2. ClientHealth

### Metric: `ibcmon_client_health`

- **Type:** Gauge
- **Description:** Health status of the IBC client (1 if healthy, 0 otherwise).
- **Labels:**
  - `src_chain_id`: `ChainId` of the source chain
  - `dst_chain_id`: `ChainId` of the destination chain
  - `client_id`: Identifier for the IBC client
  - `updated_at`: Last updated timestamp in UTC timezone

**Example:**
```
ibcmon_client_health{src_chain_id="milkyway", dst_chain_id="osmosis-1", client_id="07-tendermint-1", updated_at="2025-06-16 08:22:06.810185842 +0000 UTC"} 1
```

---

## 3. IBCPacketTracker

### Metrics

| Metric Name                                           | Type   | Description                                                      | Labels                                                    |
|-------------------------------------------------------|--------|------------------------------------------------------------------|-----------------------------------------------------------|
| `ibcmon_channel_sequence`                           | Gauge  | Sequence number of the next packet to be sent on the channel     | src_chain_id, src_path, dst_chain_id, dst_path            |
| `ibcmon_consecutive_missed`                         | Gauge  | Number of consecutive IBC transactions that have been missed     | src_chain_id, src_path, dst_chain_id, dst_path            |
| `ibcmon_observed_succeed_send_packet_sequence`      | Gauge  | Sequence number of the last successfully sent packet             | src_chain_id, src_path, dst_chain_id, dst_path            |
| `ibcmon_observed_succeed_recv_packet_sequence`      | Gauge  | Sequence number of the last successfully received packet         | src_chain_id, src_path, dst_chain_id, dst_path            |
| `ibcmon_observed_succeed_ack_packet_sequence`       | Gauge  | Sequence number of the last successfully acknowledged packet     | src_chain_id, src_path, dst_chain_id, dst_path            |

**Examples:**
```text
ibcmon_channel_sequence{src_chain_id="milkyway", src_path="milkyway(07-tendermint-1/connection-0/channel-0/transfer)", dst_chain_id="osmosis-1", dst_path="osmosis-1(07-tendermint-3364/connection-2821/channel-89298/transfer)"} 16731
ibcmon_consecutive_missed{src_chain_id="milkyway", src_path="milkyway(07-tendermint-1/connection-0/channel-0/transfer)", dst_chain_id="osmosis-1", dst_path="osmosis-1(07-tendermint-3364/connection-2821/channel-89298/transfer)"} 0
ibcmon_observed_succeed_send_packet_sequence{src_chain_id="osmosis-1", src_path="osmosis-1(07-tendermint-3364/connection-2821/channel-89298/transfer)", dst_chain_id="milkyway", dst_path="milkyway(07-tendermint-1/connection-0/channel-0/transfer)"} 26888
ibcmon_observed_succeed_recv_packet_sequence{src_chain_id="osmosis-1", src_path="osmosis-1(07-tendermint-3364/connection-2821/channel-89298/transfer)", dst_chain_id="milkyway", dst_path="milkyway(07-tendermint-1/connection-0/channel-0/transfer)"} 26888
ibcmon_observed_succeed_ack_packet_sequence{src_chain_id="osmosis-1", src_path="osmosis-1(07-tendermint-3364/connection-2821/channel-89298/transfer)", dst_chain_id="milkyway", dst_path="milkyway(07-tendermint-1/connection-0/channel-0/transfer)"} 26888
```

---

## Labels Description

- `src_chain_id`: `ChainId` of the source chain
- `src_path`: IBC path of the source chain, formatted as `chain_id(client_id/connection_id/channel_id/port_id)`
- `dst_chain_id`: `ChainId` of the destination chain
- `dst_path`: IBC path of the destination chain, formatted similarly to `src_path`
- `client_id`: The identifier for the IBC client
- `updated_at`: Last updated timestamp in UTC timezone
