package exported

type IBCPacketTracker interface {
	GetPacketStatus() string
	GetSequence() uint64

	GetSrcInfo() (string, string, string)
	GetDstInfo() (string, string, string)
}
