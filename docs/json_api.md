# JSON API Documentation

## 1. `/ibc-info`

### Response

```json
[
  {
    "updated": "2025-06-05T12:00:20.055331397Z",
    "source": {
      "path": "milkyway(07-tendermint-1/connection-0/channel-0/transfer)",
      "ChainId": "milkyway",
      "ClientId": "07-tendermint-1",
      "ConnectionId": "connection-0",
      "ChannelId": "channel-0",
      "PortId": "transfer"
    },
    "destination": {
      "path": "osmosis-1(07-tendermint-3364/connection-2821/channel-89298/transfer)",
      "ChainId": "osmosis-1",
      "ClientId": "07-tendermint-3364",
      "ConnectionId": "connection-2821",
      "ChannelId": "channel-89298",
      "PortId": "transfer"
    }
  },

  ...

]
```

- **updated**: Timestamp when the info was last updated (UTC timezone)
- **source/destination**: IBC information for source and destination (see [IBC Object](#ibc-object))

---

## 2. `/client-health`

### Response

```json
[
  {
    "health": true,
    "client_updated": "2025-06-05T11:46:54.218181246Z",
    "source": "milkyway",
    "destination": "osmosis-1",
    "client_id": "07-tendermint-1",
    "trusting_period": 1209600
  },

  ...

]
```

- **health**: Boolean indicating if the client is healthy
- **client_updated**: Timestamp for last client update (UTC timezone)
- **source/destination**: `ChainId` for source and destination
- **client_id**: The client identifier
- **trusting_period**: Trusting period of client in seconds

---

## 3. `/ibc-packet`

### Response

```json
[
  {
    "updated": "2025-06-05T12:09:14.305655367Z",
    "health": true,
    "source": {
      "path": "milkyway(07-tendermint-1/connection-0/channel-0/transfer)",
      "ChainId": "milkyway",
      "ClientId": "07-tendermint-1",
      "ConnectionId": "connection-0",
      "ChannelId": "channel-0",
      "PortId": "transfer"
    },
    "destination": {
      "path": "osmosis-1(07-tendermint-3364/connection-2821/channel-89298/transfer)",
      "ChainId": "osmosis-1",
      "ClientId": "07-tendermint-3364",
      "ConnectionId": "connection-2821",
      "ChannelId": "channel-89298",
      "PortId": "transfer"
    },
    "sequence": 16087,
    "consecutive_missed": 0,
    "latest_succeed_packets": {
      "acknowledge_packet": {
        "hash": "86966325D2B8D26DABB22640DA87FC7DF20A076B68640F7E9771F4F9C3923873",
        "sequence": 16086,
        "data": ""
      },
      "recv_packet": {
        "hash": "5AA7D65A234AE2B5609AB3FD920D38493968BBCB99492CDA3040BC58DA651820",
        "sequence": 16086,
        "data": "{\"amount\":\"21595556\",\"denom\":\"transfer/channel-0/transfer/channel-874/factory/neutron1ut4c6pv4u6vyu97yw48y8g7mle0cat54848v6m97k977022lzxtsaqsgmq/udtia\",\"receiver\":\"osmo1hn7f4x23xtajz3hhevy83pcm7n0m0wpj6cpyap\",\"sender\":\"milk1hn7f4x23xtajz3hhevy83pcm7n0m0wpjus52rp\"}"
      },
      "send_packet": {
        "hash": "FB4A911B549B42BBAD320A72F482111F54AB39576C24F58BA98D71E9534E2597",
        "sequence": 16086,
        "data": "{\"amount\":\"21595556\",\"denom\":\"transfer/channel-0/transfer/channel-874/factory/neutron1ut4c6pv4u6vyu97yw48y8g7mle0cat54848v6m97k977022lzxtsaqsgmq/udtia\",\"receiver\":\"osmo1hn7f4x23xtajz3hhevy83pcm7n0m0wpj6cpyap\",\"sender\":\"milk1hn7f4x23xtajz3hhevy83pcm7n0m0wpjus52rp\"}"
      }
    }
  },

  ...


]
```

- **updated**: Timestamp when packet tracking was last executed and updated (UTC timezone)
- **health**: Boolean for packet health
- **source/destination**: IBC information for source and destination (see [IBC Object](#ibc-object))
- **sequence**: Sequence number being tracked currently
- **consecutive_missed**: Number of consecutively missed packets
- **latest_succeed_packets**: Map of packet types to their latest succeed packets (see [SucceedPacket Object](#succeedpacket-object))

### SucceedPacket Object

```json
{
  "hash": "FB4A911B549B42BBAD320A72F482111F54AB39576C24F58BA98D71E9534E2597",
  "sequence": 16086,
  "data": "{\"amount\":\"21595556\",\"denom\":\"transfer/channel-0/transfer/channel-874/factory/neutron1ut4c6pv4u6vyu97yw48y8g7mle0cat54848v6m97k977022lzxtsaqsgmq/udtia\",\"receiver\":\"osmo1hn7f4x23xtajz3hhevy83pcm7n0m0wpj6cpyap\",\"sender\":\"milk1hn7f4x23xtajz3hhevy83pcm7n0m0wpjus52rp\"}"
}
```

- **hash**: Hash of the transaction that packet succeeded
- **sequence**: Sequence number of the packet
- **data**: Details of the packet in JSON format

---

## IBC Object

```json
{
  "path": "milkyway(07-tendermint-1/connection-0/channel-0/transfer)",
  "ChainId": "milkyway",
  "ClientId": "07-tendermint-1",
  "ConnectionId": "connection-0",
  "ChannelId": "channel-0",
  "PortId": "transfer"
}
```

- **path**: IBC path of the chain, formatted as `chain_id(client_id/connection_id/channel_id/port_id)`
- **ChainId**: Chain identifier
- **ClientId**: Client identifier
- **ConnectionId**: Connection identifier
- **ChannelId**: Channel identifier
- **PortId**: Port identifier
