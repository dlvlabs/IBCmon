package app

import (
	"context"
	"sync"
	"time"

	"github.com/dlvlabs/ibcmon/client/grpc"
	"github.com/dlvlabs/ibcmon/client/rpc"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	tendermint "github.com/cosmos/ibc-go/v10/modules/light-clients/07-tendermint"
)

type (
	App struct {
		cfg Config
		cdc codecTypes.InterfaceRegistry

		rpcs RPCs

		grpcsMutex sync.Mutex
		grpcs      GRPCs

		storeMutex sync.Mutex
		Store      Store
	}

	// chainId => rpc | grpc client
	RPCs  map[string]*rpc.Client
	GRPCs map[string]*grpc.Client
)

type (
	Config struct {
		General General `toml:"general"`
		TG      TG      `toml:"tg"`
		Rule    Rule    `toml:"rule"`

		BaseChain Endpoints `toml:"base_chain"`

		// chainId => endpoint info
		Counterparties map[string]Endpoints `toml:"counterparties"`
	}

	General struct {
		baseChainId string

		LogLevel   string `toml:"log_level"`
		ListenPort int    `toml:"listen_port"`

		IbcInfoUpdateInterval  time.Duration `toml:"ibc_info_update_interval"`
		ClientCheckInterval    time.Duration `toml:"client_check_interval"`
		PacketTrackingInterval time.Duration `toml:"packet_tracking_interval"`
	}
	TG struct {
		Enable bool   `toml:"enable"`
		Token  string `toml:"token"`
		ChatID string `toml:"chat_id"`
	}
	Rule struct {
		ClientExpiredWarningTime time.Duration `toml:"client_expired_warning_time"`
		ConsecutiveMissedPackets uint64        `toml:"consecutive_missed_packets"`
	}
	Endpoints struct {
		GRPC    GRPC   `toml:"grpc"`
		RPCAddr string `toml:"rpc_addr"`
	}
	GRPC struct {
		Addr    string `toml:"addr"`
		TLSConn bool   `toml:"tls_conn"`
	}
)

type (
	Store struct {
		Updated time.Time

		IBCInfo IBCInfo
	}
)

func NewApp(ctx context.Context, cfg Config) (*App, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	bcGRPC := grpc.New(cfg.BaseChain.GRPC.Addr, cfg.BaseChain.GRPC.TLSConn)
	bcGRPC.Connect()
	defer bcGRPC.Terminate()

	bcChainId, err := bcGRPC.GetChainId(ctx)
	if err != nil {
		return nil, err
	}

	cfg.General.baseChainId = bcChainId

	rpcs := make(RPCs)
	grpcs := make(GRPCs)

	rpcs[bcChainId], err = rpc.New(cfg.BaseChain.RPCAddr)
	if err != nil {
		return nil, err
	}
	grpcs[bcChainId] = grpc.New(cfg.BaseChain.GRPC.Addr, cfg.BaseChain.GRPC.TLSConn)

	for chainId, endpoints := range cfg.Counterparties {
		rpcs[chainId], err = rpc.New(endpoints.RPCAddr)
		if err != nil {
			return nil, err
		}
		grpcs[chainId] = grpc.New(endpoints.GRPC.Addr, endpoints.GRPC.TLSConn)
	}

	cdc := codecTypes.NewInterfaceRegistry()
	tendermint.RegisterInterfaces(cdc)

	app := &App{
		cfg: cfg,
		cdc: cdc,

		rpcs:  rpcs,
		grpcs: grpcs,

		Store: Store{
			IBCInfo: make(IBCInfo),
		},
	}
	return app, nil
}
