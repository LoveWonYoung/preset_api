//go:build windows && amd64

package preset_api

type PresetConfig struct {
	PhysicalID     uint32
	ResponseID     uint32
	FunctionalID   uint32
	NBsMS          uint32
	NCrMS          uint32
	STMinUS        uint32
	MaxPDULen      uint32 // 1-1048576 bytes
	RxCapacity     uint32 // 1-65536 frames
	IsFD           uint8
	PaddingEnabled uint8
	PaddingByte    uint8
	BlockSize      uint8
	MaxWaitFrames  uint8
	RawRxEnabled   uint8 // 1 enables the raw capture ring
	Reserved       [2]uint8
}

type PresetCanFdTiming struct {
	NominalBRP   uint32
	NominalTSEG1 uint32
	NominalTSEG2 uint32
	NominalSJW   uint32
	DataBRP      uint32
	DataTSEG1    uint32
	DataTSEG2    uint32
	DataSJW      uint32
}

type PresetAutoConfig struct {
	Channel             uint8
	IsFD                uint8
	BRS                 uint8
	Reserved            uint8
	NominalBitrate      uint32
	DataBitrate         uint32
	RxBufferSize        uint32 // 1-65536 frames
	PollBatchSize       uint32 // 1-rx_buffer_size
	MinTxIntervalUS     uint32
	MinRxPollIntervalUS uint32
	CandidateOrder      [4]uint8 // PRESET_CAN_BACKEND_* priority order; zero entries are skipped
	ReservedTail        [4]uint8
}

type PresetToomossConfig struct {
	Channel             uint8
	BRS                 uint8
	Reserved            [2]uint8
	NominalBitrate      uint32
	DataBitrate         uint32
	RxBufferSize        uint32 // 1-65536 frames
	PollBatchSize       uint32 // 1-rx_buffer_size
	MinTxIntervalUS     uint32
	MinRxPollIntervalUS uint32
}

type PresetToomossLinConfig struct {
	ChannelsMask uint8 // PRESET_TOOMOSS_CHANNEL_* bitmask
	MasterMode   uint8 // 1 master, 0 slave
	Reserved     [2]uint8
	Baudrate     uint32 // 2000-5000000, default 19200
}

type PresetToomossElinsConfig struct {
	ChannelsMask     uint8 // PRESET_TOOMOSS_CHANNEL_* bitmask
	ResistorEnabled  uint8 // 0 or 1
	Version          uint8 // PRESET_TOOMOSS_ELINS_VER_*
	Reserved         uint8
	Baudrate         uint32 // 2000-5000000, default 19200
	ReceiveTimeoutUS uint32 // 0 selects the 1000us default
}

type PresetToomossLinFrame struct {
	Timestamp       uint32
	FrameID         uint8
	DataLen         uint8
	ClassicChecksum uint8
	Reserved        uint8
	Data            [8]uint8
}

type PresetPCANLinConfig struct {
	Channel        uint8 // logical channel and default enumeration index
	Reserved       uint8
	HardwareHandle uint16 // 0 selects by channel index
	Baudrate       uint32 // 1000-20000, default 19200
}

type PresetTSMasterLinConfig struct {
	ApplicationChannel uint8 // zero-based logical channel, 0-31
	HardwareChannel    uint8 // zero-based physical LIN channel
	Protocol           uint8 // PRESET_LIN_PROTOCOL_*
	Reserved           uint8
	HardwareIndex      int32
	DeviceType         int32
	Baudrate           uint32 // 1000-20000, default 19200
}

type PresetVectorLinConfig struct {
	ApplicationChannel uint8 // logical channel exposed by preset_rs
	Version            uint8 // 1=LIN 1.3, 2=LIN 2.0, 3=LIN 2.1
	Reserved           [2]uint8
	HardwareType       int32 // PRESET_VECTOR_HWTYPE_*
	HardwareIndex      uint32
	HardwareChannel    uint32
	Baudrate           uint32 // 1000-20000, default 19200
	RxQueueSize        uint32 // power of two, 8192-524288
}

type PresetLinFrame struct {
	FrameID  uint8
	DataLen  uint8
	Reserved [2]uint8
	Data     [8]uint8
}

// PresetToomossElinsMessage Matches the vendor ELINS_MSG layout (88 bytes).
type PresetToomossElinsMessage struct {
	DataLen         uint8
	BreakBits       uint8
	Status          uint8
	Flags           uint8
	Sync            uint8
	TimestampHigh   uint8
	MsgSendTimes    uint16
	Timestamp       uint32
	CmdCode         uint8
	DeviceID        uint8
	RegisterAddress uint16
	CRC16           uint16
	Data            [64]uint8
	AckValue        [4]uint8
}

type PresetPCANConfig struct {
	Channel             uint8 // zero-based PCAN USB channel, 0-15
	IsFD                uint8 // initialize with CAN_InitializeFD
	BRS                 uint8 // set BRS on transmitted CAN-FD frames
	Reserved            uint8
	NominalBitrate      uint32
	DataBitrate         uint32
	RxBufferSize        uint32 // 1-65536 frames
	PollBatchSize       uint32 // 1-rx_buffer_size
	MinTxIntervalUS     uint32
	MinRxPollIntervalUS uint32
}

type PresetTSMasterConfig struct {
	ApplicationChannel  uint8 // zero-based logical channel, 0-31
	HardwareChannel     uint8 // zero-based physical channel
	IsFD                uint8
	BRS                 uint8
	HardwareIndex       int32 // zero-based enumerated device index
	DeviceType          int32
	NominalBitrate      uint32
	DataBitrate         uint32
	RxBufferSize        uint32 // 1-65536 frames
	PollBatchSize       uint32 // 1-rx_buffer_size
	MinTxIntervalUS     uint32
	MinRxPollIntervalUS uint32
}

type PresetVectorConfig struct {
	ApplicationChannel  uint8 // zero-based logical channel, 0-255
	IsFD                uint8
	BRS                 uint8
	Reserved            uint8
	HardwareType        int32  // PRESET_VECTOR_HWTYPE_*; -1 selects any
	HardwareIndex       uint32 // zero-based device index
	HardwareChannel     uint32 // zero-based physical channel
	NominalBitrate      uint32
	DataBitrate         uint32
	NominalSJW          uint32
	NominalTSEG1        uint32
	NominalTSEG2        uint32
	DataSJW             uint32
	DataTSEG1           uint32
	DataTSEG2           uint32
	DriverRxQueueSize   uint32 // power of two, 8192-524288
	RxBufferSize        uint32 // 1-65536 frames
	PollBatchSize       uint32 // 1-rx_buffer_size
	MinTxIntervalUS     uint32
	MinRxPollIntervalUS uint32
}

type PresetCanFrame struct {
	ID       uint32
	DataLen  uint8
	IsFD     uint8
	Reserved [2]uint8
	Data     [64]uint8
}

// PresetCanFrameEx includes DLC, direction, CAN-FD, and BRS metadata. Do not
// mix CanUdsTryRead and CanUdsTryReadEx on one client; both drain the same ring.
type PresetCanFrameEx struct {
	ID        uint32
	DLC       uint8
	DataLen   uint8
	Direction uint8 // PRESET_CAN_DIRECTION_*
	IsFD      uint8
	BRS       uint8
	Reserved  [3]uint8
	Data      [64]uint8
}

type PresetRxStats struct {
	Received uint64
	Dropped  uint64
	Queued   uint32
	Capacity uint32
}

type PresetBusLoad struct {
	Load           float64 // occupied-time ratio, 0.0-1.0
	WindowNS       uint64  // effective measurement interval
	NominalBitrate uint32
	DataBitrate    uint32
	FrameCount     uint64
}
