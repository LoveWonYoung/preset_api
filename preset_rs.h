#ifndef PRESET_RS_H
#define PRESET_RS_H

#include <stddef.h>
#include <stdint.h>

#define PRESET_RS_ABI_VERSION 12

#ifdef __cplusplus
extern "C" {
#endif

#if defined(_WIN32) && defined(PRESET_RS_USE_DLLIMPORT)
#define PRESET_RS_API __declspec(dllimport)
#else
#define PRESET_RS_API
#endif

typedef struct PresetDevice PresetDevice;
typedef struct PresetLinUdsClient PresetLinUdsClient;

enum {
    PRESET_OK = 0,
    PRESET_ERR_NULL_PTR = -1,
    PRESET_ERR_BUFFER_TOO_SMALL = -2,
    PRESET_ERR_NOT_INIT = -3,
    PRESET_ERR_TIMEOUT = -4,
    PRESET_ERR_UDS_NEGATIVE = -5,
    PRESET_ERR_TRANSPORT = -6,
    PRESET_ERR_UNSUPPORTED = -7,
    PRESET_ERR_INVALID_ARG = -8,
    PRESET_ERR_BUSY = -9,
    PRESET_ERR_ISOTP = -10,
    PRESET_ERR_PANIC = -11,
    PRESET_ERR_CANCELLED = -12,
    PRESET_ERR_WAIT_TIMEOUT = -13,
    PRESET_ERR_NOT_READY = -14,
    PRESET_ERR_DIAGNOSTIC_OCCUPIED = -15,
};

/* Platform backends. Zero on macOS. Bit positions stay stable across ABI versions. */
enum {
    PRESET_CAP_TSMASTER = UINT64_C(1) << 9,
    PRESET_CAP_PCAN = UINT64_C(1) << 10,
    PRESET_CAP_VECTOR = UINT64_C(1) << 11,
    PRESET_CAP_AUTO_DRIVER = UINT64_C(1) << 20,
    PRESET_CAP_PCAN_LIN = UINT64_C(1) << 21,
    PRESET_CAP_TSMASTER_LIN = UINT64_C(1) << 22,
    PRESET_CAP_VECTOR_LIN = UINT64_C(1) << 23,
};

enum {
    PRESET_CAN_DIRECTION_TX = 0,
    PRESET_CAN_DIRECTION_RX = 1,
};

enum {
    PRESET_CAN_BACKEND_NONE = 0,
    PRESET_CAN_BACKEND_TOOMOSS = 1,
    PRESET_CAN_BACKEND_TSMASTER = 2,
    PRESET_CAN_BACKEND_PCAN = 3,
    PRESET_CAN_BACKEND_VECTOR = 4,
};

enum {
    PRESET_LOG_FILTER_OFF = 0,
    PRESET_LOG_FILTER_LIST = 1,
};

enum {
    /* Panel LIN1/2/3/4 = bits 0/1/2/3 (same zero-based index as CAN). */
    PRESET_TOOMOSS_CHANNEL_1 = UINT8_C(1) << 0,
    PRESET_TOOMOSS_CHANNEL_2 = UINT8_C(1) << 1,
    PRESET_TOOMOSS_CHANNEL_3 = UINT8_C(1) << 2,
    PRESET_TOOMOSS_CHANNEL_4 = UINT8_C(1) << 3,
};

enum {
    PRESET_TOOMOSS_LIN_POWER_0V = 0,
    PRESET_TOOMOSS_LIN_POWER_12V = 1,
    PRESET_TOOMOSS_LIN_POWER_5V = 2,
};

enum {
    PRESET_TOOMOSS_ELINS_VER_IND83080 = 0,
    PRESET_TOOMOSS_ELINS_VER_IND83220 = 1,
    PRESET_TOOMOSS_ELINS_VER_IND83010 = 2,
};

enum {
    PRESET_LIN_PROTOCOL_13 = 0,
    PRESET_LIN_PROTOCOL_20 = 1,
    PRESET_LIN_PROTOCOL_21 = 2,
};

enum {
    PRESET_LIN_FUNCTIONAL_NAD = 0x7E,
    PRESET_LIN_BROADCAST_NAD = 0x7F,
};

/* Vector XL hardware type values from vxlapi.h. */
enum {
    PRESET_VECTOR_HWTYPE_ANY = -1,
    PRESET_VECTOR_HWTYPE_NONE = 0,
    PRESET_VECTOR_HWTYPE_VIRTUAL = 1,
    PRESET_VECTOR_HWTYPE_CANCARDX = 2,
    PRESET_VECTOR_HWTYPE_CANAC2PCI = 6,
    PRESET_VECTOR_HWTYPE_CANCARDY = 12,
    PRESET_VECTOR_HWTYPE_CANCARDXL = 15,
    PRESET_VECTOR_HWTYPE_CANCASEXL = 21,
    PRESET_VECTOR_HWTYPE_CANCASEXL_LOG_OBSOLETE = 23,
    PRESET_VECTOR_HWTYPE_CANBOARDXL = 25,
    PRESET_VECTOR_HWTYPE_CANBOARDXL_PXI = 27,
    PRESET_VECTOR_HWTYPE_VN2600 = 29,
    PRESET_VECTOR_HWTYPE_VN2610 = 29,
    PRESET_VECTOR_HWTYPE_VN3300 = 37,
    PRESET_VECTOR_HWTYPE_VN3600 = 39,
    PRESET_VECTOR_HWTYPE_VN7600 = 41,
    PRESET_VECTOR_HWTYPE_CANCARDXLE = 43,
    PRESET_VECTOR_HWTYPE_VN8900 = 45,
    PRESET_VECTOR_HWTYPE_VN8950 = 47,
    PRESET_VECTOR_HWTYPE_VN2640 = 53,
    PRESET_VECTOR_HWTYPE_VN1610 = 55,
    PRESET_VECTOR_HWTYPE_VN1630 = 57,
    PRESET_VECTOR_HWTYPE_VN1640 = 59,
    PRESET_VECTOR_HWTYPE_VN8970 = 61,
    PRESET_VECTOR_HWTYPE_VN1611 = 63,
    PRESET_VECTOR_HWTYPE_VN5240 = 64,
    PRESET_VECTOR_HWTYPE_VN5610 = 65,
    PRESET_VECTOR_HWTYPE_VN5620 = 66,
    PRESET_VECTOR_HWTYPE_VN7570 = 67,
    PRESET_VECTOR_HWTYPE_VN5650 = 68,
    PRESET_VECTOR_HWTYPE_IPCLIENT = 69,
    PRESET_VECTOR_HWTYPE_IPSERVER = 71,
    PRESET_VECTOR_HWTYPE_VX1121 = 73,
    PRESET_VECTOR_HWTYPE_VX1131 = 75,
    PRESET_VECTOR_HWTYPE_VT6204 = 77,
    PRESET_VECTOR_HWTYPE_VN1630_LOG = 79,
    PRESET_VECTOR_HWTYPE_VN7610 = 81,
    PRESET_VECTOR_HWTYPE_VN7572 = 83,
    PRESET_VECTOR_HWTYPE_VN8972 = 85,
    PRESET_VECTOR_HWTYPE_VN0601 = 87,
    PRESET_VECTOR_HWTYPE_VN5640 = 89,
    PRESET_VECTOR_HWTYPE_VX0312 = 91,
    PRESET_VECTOR_HWTYPE_VH6501 = 94,
    PRESET_VECTOR_HWTYPE_VN8800 = 95,
    PRESET_VECTOR_HWTYPE_IPCL8800 = 96,
    PRESET_VECTOR_HWTYPE_IPSRV8800 = 97,
    PRESET_VECTOR_HWTYPE_CSMCAN = 98,
    PRESET_VECTOR_HWTYPE_VN5610A = 101,
    PRESET_VECTOR_HWTYPE_VN7640 = 102,
    PRESET_VECTOR_HWTYPE_VX1135 = 104,
    PRESET_VECTOR_HWTYPE_VN4610 = 105,
    PRESET_VECTOR_HWTYPE_VT6306 = 107,
    PRESET_VECTOR_HWTYPE_VT6104A = 108,
    PRESET_VECTOR_HWTYPE_VN5430 = 109,
    PRESET_VECTOR_HWTYPE_VTSSERVICE = 110,
    PRESET_VECTOR_HWTYPE_VN1530 = 112,
    PRESET_VECTOR_HWTYPE_VN1531 = 113,
    PRESET_VECTOR_HWTYPE_VX1161A = 114,
    PRESET_VECTOR_HWTYPE_VX1161B = 115,
};

/* TSMaster hardware subtype values, matching TSMaster.h/canbuskit. */
enum {
    PRESET_TSMASTER_TSCAN_PRO = 1,
    PRESET_TSMASTER_TSCAN_LITE1 = 2,
    PRESET_TSMASTER_TC1001 = 3,
    PRESET_TSMASTER_TL1001 = 4,
    PRESET_TSMASTER_TC1011 = 5,
    PRESET_TSMASTER_TM5011 = 6,
    PRESET_TSMASTER_TC1002 = 7,
    PRESET_TSMASTER_TC1014 = 8,
    PRESET_TSMASTER_TSCANFD2517 = 9,
    PRESET_TSMASTER_TC1026 = 10,
    PRESET_TSMASTER_TC1016 = 11,
    PRESET_TSMASTER_TC1012 = 12,
    PRESET_TSMASTER_TC1013 = 13,
    PRESET_TSMASTER_TLOG1002 = 14,
    PRESET_TSMASTER_TC1034 = 15,
    PRESET_TSMASTER_TC1018 = 16,
    PRESET_TSMASTER_GW2116 = 17,
    PRESET_TSMASTER_TC2115 = 18,
    PRESET_TSMASTER_MP1013 = 19,
    PRESET_TSMASTER_TC1113 = 20,
    PRESET_TSMASTER_TC1114 = 21,
    PRESET_TSMASTER_TP1013 = 22,
    PRESET_TSMASTER_TC1017 = 23,
    PRESET_TSMASTER_TP1018 = 24,
    PRESET_TSMASTER_TF10XX = 25,
    PRESET_TSMASTER_TL1004_FD_4_LIN_2 = 26,
    PRESET_TSMASTER_TE1051 = 27,
    PRESET_TSMASTER_TP1051 = 28,
    PRESET_TSMASTER_TP1034 = 29,
    PRESET_TSMASTER_TTS9015 = 30,
    PRESET_TSMASTER_TP1026 = 31,
    PRESET_TSMASTER_TTS1026 = 32,
    PRESET_TSMASTER_TTS1034 = 33,
    PRESET_TSMASTER_TTS1018 = 34,
    PRESET_TSMASTER_TL1011 = 35,
    PRESET_TSMASTER_TTS1015_LIAUTO = 36,
    PRESET_TSMASTER_TTS1013_LIAUTO = 37,
    PRESET_TSMASTER_TTS1016PRO = 38,
    PRESET_TSMASTER_TC1054PRO = 39,
    PRESET_TSMASTER_TC1054 = 40,
    PRESET_TSMASTER_TLOG1038 = 41,
    PRESET_TSMASTER_TO1013 = 42,
    PRESET_TSMASTER_TC1034PRO = 43,
    PRESET_TSMASTER_TC1018PRO = 44,
    PRESET_TSMASTER_TC1038PRO = 45,
    PRESET_TSMASTER_TC1014PRO = 46,
    PRESET_TSMASTER_TC1034PROPLUS = 47,
    PRESET_TSMASTER_TA1038 = 48,
    PRESET_TSMASTER_TC1055PRO = 49,
    PRESET_TSMASTER_TC1056PRO = 50,
    PRESET_TSMASTER_TC1057PRO = 51,
    PRESET_TSMASTER_TC4016 = 52,
    PRESET_TSMASTER_GW2208 = 53,
    PRESET_TSMASTER_TLOG1039 = 54,
    PRESET_TSMASTER_GW1040 = 55,
    PRESET_TSMASTER_TC3014 = 56,
    PRESET_TSMASTER_TP1014 = 57,
    PRESET_TSMASTER_TA825_4 = 58,
    PRESET_TSMASTER_TC1013HV = 59,
    PRESET_TSMASTER_TC1052 = 60,
    PRESET_TSMASTER_TTS1017PRO = 61,
    PRESET_TSMASTER_TLOG1057 = 62,
    PRESET_TSMASTER_TC1017PRO = 63,
    PRESET_TSMASTER_GW2202 = 64,
    PRESET_TSMASTER_GW2204 = 65,
    PRESET_TSMASTER_GW2212 = 66,
    PRESET_TSMASTER_TA821 = 67,
    PRESET_TSMASTER_TX1000 = 68,
    PRESET_TSMASTER_TC1055PROPLUS = 69,
    PRESET_TSMASTER_TC1043 = 70,
};

typedef struct PresetCanFdTiming {
    uint32_t nominal_brp;
    uint32_t nominal_tseg1;
    uint32_t nominal_tseg2;
    uint32_t nominal_sjw;
    uint32_t data_brp;
    uint32_t data_tseg1;
    uint32_t data_tseg2;
    uint32_t data_sjw;
} PresetCanFdTiming;

/* Shared bus settings. is_fd and brs describe the channel, not one request. */
typedef struct PresetCanChannelConfig {
    uint8_t is_fd;
    uint8_t brs;
    uint32_t nominal_bitrate;
    uint32_t data_bitrate;
    uint32_t rx_buffer_size; /* 1-65536 frames */
    uint32_t poll_batch_size; /* 1-rx_buffer_size */
    uint32_t min_tx_interval_us;
    uint32_t min_rx_poll_interval_us; /* zero selects the 1000-us default */
} PresetCanChannelConfig;

typedef struct PresetAutoConfig {
    uint8_t channel;
    PresetCanChannelConfig can;
    /* PRESET_CAN_BACKEND_* priority order; zero entries are skipped. */
    uint8_t candidate_order[4];
} PresetAutoConfig;

typedef struct PresetToomossConfig {
    uint8_t channel; /* zero-based: panel CAN1/2/3/4 = 0/1/2/3 */
    /* Default helper uses 8192 for can.rx_buffer_size and can.poll_batch_size. */
    PresetCanChannelConfig can;
} PresetToomossConfig;

typedef struct PresetToomossLinConfig {
    uint8_t channel;     /* zero-based: panel LIN1/2/3/4 = 0/1/2/3 */
    uint8_t master_mode; /* 1 master, 0 slave */
    uint32_t baudrate;   /* 2000-100000, default 19200 */
} PresetToomossLinConfig;

typedef struct PresetToomossElinsConfig {
    uint8_t channels_mask;    /* PRESET_TOOMOSS_CHANNEL_* bitmask */
    uint8_t resistor_enabled; /* 0 or 1 */
    uint8_t version;          /* PRESET_TOOMOSS_ELINS_VER_* */
    uint8_t reserved;
    uint32_t baudrate;        /* 2000-5000000, default 19200 */
    uint32_t receive_timeout_us; /* 0 selects the 1000us default */
} PresetToomossElinsConfig;

typedef struct PresetPCANLinConfig {
    uint8_t channel;          /* logical channel and default enumeration index */
    uint8_t reserved;
    uint16_t hardware_handle; /* 0 selects by channel index */
    uint32_t baudrate;        /* 1000-20000, default 19200 */
} PresetPCANLinConfig;

typedef struct PresetTSMasterLinConfig {
    uint8_t application_channel; /* zero-based logical channel, 0-31 */
    uint8_t hardware_channel;    /* zero-based physical LIN channel */
    uint8_t protocol;            /* PRESET_LIN_PROTOCOL_* */
    uint8_t reserved;
    int32_t hardware_index;
    int32_t device_type;
    uint32_t baudrate;           /* 1000-20000, default 19200 */
} PresetTSMasterLinConfig;

typedef struct PresetVectorLinConfig {
    uint8_t application_channel; /* logical channel exposed by preset_rs */
    uint8_t version;             /* PRESET_LIN_PROTOCOL_* */
    uint8_t reserved[2];
    int32_t hardware_type;       /* PRESET_VECTOR_HWTYPE_* */
    uint32_t hardware_index;
    uint32_t hardware_channel;
    uint32_t baudrate;           /* 1000-20000, default 19200 */
    uint32_t rx_queue_size;      /* power of two, 8192-524288 */
} PresetVectorLinConfig;

/* classic_checksum is meaningful only when checksum_known is 1.
 * checksum_known is 1 when the driver reports the type, or when the frame
 * follows the checksum table configured for the channel (LIN 1.3 is classic;
 * later protocols are enhanced except IDs 0x3C and 0x3D). It is not recovered
 * from the bytes on the wire. timestamp is the vendor counter (Toomoss copies
 * LIN_EX_MSG.Timestamp) and is zero when that counter is unavailable.
 * timestamp_us is microseconds and is zero when the backend does not provide
 * that unit. */
typedef struct PresetLinFrame {
    uint8_t frame_id;
    uint8_t data_len;
    uint8_t classic_checksum; /* 1 classic, 0 enhanced */
    uint8_t checksum_known;   /* 0 unknown, 1 driver or channel configuration */
    uint8_t data[8];
    uint32_t timestamp;
    uint64_t timestamp_us;
} PresetLinFrame;

/* Matches the vendor ELINS_MSG layout (88 bytes). */
typedef struct PresetToomossElinsMessage {
    uint8_t data_len;
    uint8_t break_bits;
    uint8_t status;
    uint8_t flags;
    uint8_t sync;
    uint8_t timestamp_high;
    uint16_t msg_send_times;
    uint32_t timestamp;
    uint8_t cmd_code;
    uint8_t device_id;
    uint16_t register_address;
    uint16_t crc16;
    uint8_t data[64];
    uint8_t ack_value[4];
} PresetToomossElinsMessage;

typedef struct PresetPCANConfig {
    uint8_t channel; /* zero-based PCAN USB channel, 0-15 */
    PresetCanChannelConfig can;
} PresetPCANConfig;

typedef struct PresetTSMasterConfig {
    uint8_t application_channel; /* zero-based logical channel, 0-31 */
    uint8_t hardware_channel;    /* zero-based physical channel */
    int32_t hardware_index;      /* zero-based enumerated device index */
    int32_t device_type;
    PresetCanChannelConfig can;
} PresetTSMasterConfig;

typedef struct PresetVectorConfig {
    uint8_t application_channel; /* zero-based logical channel, 0-255 */
    int32_t hardware_type;       /* PRESET_VECTOR_HWTYPE_*; -1 selects any */
    uint32_t hardware_index;     /* zero-based device index */
    uint32_t hardware_channel;   /* zero-based physical channel */
    uint32_t nominal_sjw;
    uint32_t nominal_tseg1;
    uint32_t nominal_tseg2;
    uint32_t data_sjw;
    uint32_t data_tseg1;
    uint32_t data_tseg2;
    uint32_t driver_rx_queue_size; /* power of two, 8192-524288 */
    PresetCanChannelConfig can;
} PresetVectorConfig;

/* timestamp_us is the vendor hardware or driver monotonic timestamp in
 * microseconds; zero means unavailable. Direct reads drain the channel queue
 * and report RX. Observer reads drain only that observer. */
typedef struct PresetCanFrameEx {
    uint32_t id;
    uint8_t data_len;
    uint8_t direction; /* PRESET_CAN_DIRECTION_* */
    uint8_t is_fd;
    uint8_t brs;
    uint8_t is_extended; /* 29-bit identifier format */
    uint8_t is_remote;   /* classic CAN RTR; data_len retains the requested DLC */
    uint8_t esi;         /* receive-only CAN-FD error-state indicator */
    uint8_t reserved;    /* must be zero when transmitting */
    uint64_t timestamp_us;
    uint8_t data[64];
} PresetCanFrameEx;

typedef struct PresetRxStats {
    uint64_t received;
    uint64_t dropped;
    uint32_t queued;
    uint32_t capacity;
} PresetRxStats;

typedef struct PresetBusLoad {
    double load;             /* occupied-time ratio, 0.0-1.0 */
    uint64_t window_ns;      /* effective measurement interval */
    uint32_t nominal_bitrate;
    uint32_t data_bitrate;
    uint64_t frame_count;
} PresetBusLoad;

PRESET_RS_API const char *preset_version(void);
PRESET_RS_API uint32_t preset_abi_version(void);
PRESET_RS_API uint64_t preset_get_capabilities(void);
PRESET_RS_API PresetAutoConfig preset_auto_default_config(void);
/* Explicit 500 kbit/s / 2 Mbit/s examples for each vendor clock/layout.
 * Keep the device config bitrates consistent for bus-load reporting. */
PRESET_RS_API PresetCanFdTiming preset_pcan_default_fd_timing(void);
PRESET_RS_API PresetCanFdTiming preset_toomoss_default_fd_timing(void);
PRESET_RS_API PresetToomossConfig preset_toomoss_default_config(void);
PRESET_RS_API PresetToomossLinConfig preset_toomoss_lin_default_config(void);
PRESET_RS_API PresetToomossElinsConfig preset_toomoss_elins_default_config(void);
PRESET_RS_API PresetPCANConfig preset_pcan_default_config(void);
PRESET_RS_API PresetTSMasterConfig preset_tsmaster_default_config(void);
PRESET_RS_API PresetVectorConfig preset_vector_default_config(void);
PRESET_RS_API PresetPCANLinConfig preset_pcan_lin_default_config(void);
PRESET_RS_API PresetTSMasterLinConfig preset_tsmaster_lin_default_config(void);
PRESET_RS_API PresetVectorLinConfig preset_vector_lin_default_config(void);

/* Process-wide CAN/CAN-FD/LIN frame logging. Init creates
 * ./YYYY_MM_DD/<name><timestamp>.log; frame output remains disabled until
 * preset_set_print_log(1). Writes use a bounded asynchronous queue and may be
 * dropped rather than blocking transport timing. Shutdown drains the queue,
 * joins the logger worker, and also disables frame logging; a later init starts
 * a fresh worker. Timestamps are captured when frames enter the queue, not when
 * the logger thread writes them. */
PRESET_RS_API int32_t preset_log_init(const char *name);
PRESET_RS_API void preset_log_shutdown(void);
PRESET_RS_API void preset_set_print_log(uint8_t enable);
/* Dropped entries since the latest successful preset_log_init call. */
PRESET_RS_API uint64_t preset_log_dropped_count(void);

/* LIST accepts at most 2048 standard 11-bit CAN IDs. An empty list suppresses
 * all CAN/CAN-FD frame logs; LIN logs are not filtered. OFF removes filtering. */
PRESET_RS_API int32_t preset_set_log_filter(
    uint8_t mode,
    const uint32_t *ids,
    size_t id_count);

/* Every backend uses explicit device ownership. Opening a device never creates
 * a runtime or LIN client; releasing either never closes the device.
 * Configuration pointers are required and NULL returns PRESET_ERR_NULL_PTR.
 * Use preset_*_default_config() for defaults. A null CAN-FD timing pointer
 * derives controller timing from the channel bitrates. */
/* Windows AutoDriver tries candidate_order until a mode-compatible backend
 * opens successfully. Defaults: Toomoss, TSMaster/TC1016, PCAN, Vector/VN1640. */
PRESET_RS_API int32_t preset_auto_open(
    const PresetAutoConfig *config,
    PresetDevice **out_device);

PRESET_RS_API int32_t preset_device_get_backend(
    const PresetDevice *device,
    uint8_t *out_backend);

PRESET_RS_API int32_t preset_toomoss_open(PresetDevice **out_device);

/* May be called repeatedly to initialize distinct CAN channels.
 * can.is_fd = 1 attempts CAN-FD and falls back to classic CAN when the adapter
 * rejects it. A channel that fell back cannot serve FD diagnostic requests.
 * A null timing pointer derives BRP/TSEG/SJW from the bitrates. Toomoss CAN
 * initialization enables internal 120-ohm termination and automatic
 * retransmission in both CAN-FD and classic/fallback modes. */
PRESET_RS_API int32_t preset_toomoss_can_init(
    PresetDevice *device,
    const PresetToomossConfig *device_config,
    const PresetCanFdTiming *timing);

PRESET_RS_API int32_t preset_toomoss_lin_init(
    PresetDevice *device,
    const PresetToomossLinConfig *config);

/* Read frames buffered by LIN_EX_SlaveGetData. On input, *inout_len is the
 * frame capacity (1-512); on success it is the number of frames written. */
PRESET_RS_API int32_t preset_toomoss_lin_try_read(
    PresetDevice *device,
    uint8_t channel,
    PresetLinFrame *frames,
    size_t *inout_len);

PRESET_RS_API int32_t preset_toomoss_lin_break(
    PresetDevice *device,
    uint8_t channel);

PRESET_RS_API int32_t preset_toomoss_lin_set_power(
    PresetDevice *device,
    uint8_t channel,
    uint8_t voltage);

PRESET_RS_API int32_t preset_toomoss_elins_init(
    PresetDevice *device,
    const PresetToomossElinsConfig *config);

PRESET_RS_API int32_t preset_toomoss_elins_read(
    PresetDevice *device,
    uint8_t channel,
    PresetToomossElinsMessage *messages,
    size_t *inout_len);

/* Windows LIN master init adds one channel to an already opened device of the
 * same vendor. CAN and LIN channels may share that device handle. Add all
 * TSMaster mappings before creating any client. */
PRESET_RS_API int32_t preset_pcan_lin_init(
    PresetDevice *device,
    const PresetPCANLinConfig *config);
PRESET_RS_API int32_t preset_tsmaster_lin_init(
    PresetDevice *device,
    const PresetTSMasterLinConfig *config);
PRESET_RS_API int32_t preset_vector_lin_init(
    PresetDevice *device,
    const PresetVectorLinConfig *config);

/* Same configs as the master init functions, but the channel is a LIN slave.
 * Toomoss slave mode stays PresetToomossLinConfig.master_mode = 0. */
PRESET_RS_API int32_t preset_pcan_lin_init_slave(
    PresetDevice *device,
    const PresetPCANLinConfig *config);
PRESET_RS_API int32_t preset_tsmaster_lin_init_slave(
    PresetDevice *device,
    const PresetTSMasterLinConfig *config);
PRESET_RS_API int32_t preset_vector_lin_init_slave(
    PresetDevice *device,
    const PresetVectorLinConfig *config);

/* Publish 1-8 response bytes for a 6-bit ID. classic_checksum and single_shot
 * are 0 or 1. Toomoss, PCAN, and Vector require a slave channel. Vector
 * rejects single_shot. TSMaster returns PRESET_ERR_UNSUPPORTED. */
PRESET_RS_API int32_t preset_lin_slave_set_response(
    PresetDevice *device,
    uint8_t channel,
    uint8_t frame_id,
    const uint8_t *data,
    size_t data_len,
    uint8_t classic_checksum,
    uint8_t single_shot);

PRESET_RS_API int32_t preset_lin_slave_clear_response(
    PresetDevice *device,
    uint8_t channel,
    uint8_t frame_id);

/* Wake a sleeping bus. Toomoss sends LIN_EX_MasterBreak and requires master
 * mode. Vector sends xlLinWakeUp. PCAN sends LIN_XmtWakeUp and requires slave
 * mode. TSMaster returns PRESET_ERR_UNSUPPORTED. */
PRESET_RS_API int32_t preset_lin_wakeup(
    PresetDevice *device,
    uint8_t channel);

PRESET_RS_API int32_t preset_lin_master_write(
    PresetDevice *device,
    uint8_t channel,
    uint8_t frame_id,
    const uint8_t *data,
    size_t data_len);

/* Sends the requested master header and waits up to timeout_ms.
 * PCAN uses subscriber auto-length. TSMaster still requests 8 data bytes.
 * Toomoss polls LIN_EX_MasterSendMsg and LIN_EX_GetMsg and releases the
 * adapter lock between polls. checksum_known follows the comment on
 * PresetLinFrame. A checksum or synchronization error returns
 * PRESET_ERR_TRANSPORT immediately. No slave response is PRESET_ERR_TIMEOUT. */
PRESET_RS_API int32_t preset_lin_master_read(
    PresetDevice *device,
    uint8_t channel,
    uint8_t frame_id,
    uint32_t timeout_ms,
    PresetLinFrame *out_frame);

/* Create a LIN UDS client that borrows an initialized LIN master channel from
 * the selected device. The caller retains ownership of the physical device;
 * at most one client may consume a channel at a time. nad accepts physical
 * addresses 0x01-0x7D, functional address 0x7E, or broadcast/configuration
 * address 0x7F. */
PRESET_RS_API int32_t preset_lin_uds_client_new(
    PresetDevice *device,
    uint8_t channel,
    uint8_t nad,
    PresetLinUdsClient **out_client);

/* Set both the LIN TP inter-frame transmit delay and response polling interval.
 * The default is 10 ms. Valid values are 1 through 1000 ms, and a new value is
 * applied to requests started after this call. */
PRESET_RS_API int32_t preset_lin_uds_set_poll_interval(
    PresetLinUdsClient *client,
    uint32_t interval_ms);

/* Releasing a LIN UDS client never closes its physical device. */
PRESET_RS_API int32_t preset_lin_uds_client_close(
    PresetLinUdsClient **client);

/* PCAN, TSMaster, and Vector open calls return an empty device and
 * PRESET_ERR_UNSUPPORTED off Windows. can_init adds one channel. A null PCAN
 * timing pointer derives BRP/TSEG/SJW from the bitrates. Add all TSMaster
 * channels before creating CAN runtimes. */
PRESET_RS_API int32_t preset_pcan_open(PresetDevice **out_device);
PRESET_RS_API int32_t preset_pcan_can_init(
    PresetDevice *device,
    const PresetPCANConfig *device_config,
    const PresetCanFdTiming *timing);

PRESET_RS_API int32_t preset_tsmaster_open(PresetDevice **out_device);
PRESET_RS_API int32_t preset_tsmaster_can_init(
    PresetDevice *device,
    const PresetTSMasterConfig *device_config);

PRESET_RS_API int32_t preset_vector_open(PresetDevice **out_device);
PRESET_RS_API int32_t preset_vector_can_init(
    PresetDevice *device,
    const PresetVectorConfig *device_config);

/* Raw CAN I/O directly on an initialized channel. These functions do not
 * create/register an ISO-TP or UDS client. They return PRESET_ERR_BUSY while a
 * CAN runtime owns the same channel. Direct reads report received frames only
 * and return zero frames when called before the backend's next RX poll time.
 * brs must be 0 or 1. Classic CAN ignores BRS. The channel's configured BRS
 * setting is unchanged. */
PRESET_RS_API int32_t preset_can_write_with_brs(
    PresetDevice *device,
    uint8_t channel,
    uint32_t id,
    uint8_t is_fd,
    uint8_t brs,
    const uint8_t *data,
    size_t data_len);

/* Structured raw write supporting standard/extended data frames and classic
 * CAN remote frames. direction must be PRESET_CAN_DIRECTION_TX; esi and
 * reserved must be zero. CAN-FD remote frames are invalid. */
PRESET_RS_API int32_t preset_can_write_ex(
    PresetDevice *device,
    uint8_t channel,
    const PresetCanFrameEx *frame);

PRESET_RS_API int32_t preset_can_try_read_ex(
    PresetDevice *device,
    uint8_t channel,
    PresetCanFrameEx *frames,
    size_t *inout_len);

/* Close any backend. Returns PRESET_ERR_BUSY and leaves *device unchanged while
 * a CAN runtime or LIN UDS client is active. Release those handles, then retry. */
PRESET_RS_API int32_t preset_device_close(PresetDevice **device);

/* Latest synchronous error on the calling thread, including open/init/write
 * and validation failures. Status-call success clears it. This getter preserves
 * it across size queries. Workers never modify another thread's error storage.
 * Async diagnostic outcomes/text belong to request_get_info/copy_error; CAN
 * faults and raw driver rejections also belong to runtime_copy_error. Error
 * text is UTF-8 without NUL termination; NULL/0 queries the required byte count. */
PRESET_RS_API int32_t preset_last_error(
    uint8_t *out_data,
    size_t out_capacity,
    size_t *out_len);

/* ABI v12. Diagnostic STmin values are milliseconds, 0..127. CAN-FD and BRS
 * for a request come from the initialized channel. Calls on live handles may
 * run concurrently, but close must not race a call using the same raw pointer.
 * Close clears *handle. Check preset_abi_version() before creating handles.
 */
typedef struct PresetCanRuntime PresetCanRuntime;
typedef struct PresetCanObserver PresetCanObserver;
typedef struct PresetCanPeriodic PresetCanPeriodic;
typedef struct PresetDiagnosticRequest PresetDiagnosticRequest;

enum {
    PRESET_RUNTIME_RUNNING = 0, PRESET_RUNTIME_RECOVERING = 1,
    PRESET_RUNTIME_STOPPING = 2, PRESET_RUNTIME_STOPPED = 3,
    PRESET_RUNTIME_FAULTED = 4
};
enum {
    PRESET_REQUEST_ACCEPTED = 0, PRESET_REQUEST_TRANSMITTING = 1,
    PRESET_REQUEST_WAIT_RESPONSE = 2, PRESET_REQUEST_RECEIVING = 3,
    PRESET_REQUEST_RESPONSE_PENDING = 4, PRESET_REQUEST_CANCEL_REQUESTED = 5,
    PRESET_REQUEST_COMPLETED = 6, PRESET_REQUEST_CANCELLED = 7,
    PRESET_REQUEST_TIMED_OUT = 8, PRESET_REQUEST_FAILED = 9
};
enum {
    PRESET_OUTCOME_NONE = 0, PRESET_OUTCOME_POSITIVE = 1,
    PRESET_OUTCOME_NEGATIVE = 2, PRESET_OUTCOME_SUPPRESSED_POSITIVE = 3,
    PRESET_OUTCOME_CANCELLED = 4, PRESET_OUTCOME_TIMEOUT = 5,
    PRESET_OUTCOME_ERROR = 6
};
enum {
    PRESET_TIMEOUT_NONE = 0, PRESET_TIMEOUT_P2 = 1, PRESET_TIMEOUT_P2_STAR = 2,
    PRESET_TIMEOUT_TOTAL = 3, PRESET_TIMEOUT_N_BS = 4, PRESET_TIMEOUT_N_CR = 5
};
enum {
    PRESET_ADDRESS_PHYSICAL = 0, PRESET_ADDRESS_FUNCTIONAL = 1
};
typedef struct PresetCanRuntimeOptions {
    uint32_t command_capacity, recovery_ms;
} PresetCanRuntimeOptions;
typedef struct PresetCanObserverOptions {
    uint32_t capacity;
    uint8_t rx, tx;
} PresetCanObserverOptions;
typedef struct PresetRequestTimeouts {
    uint32_t p2_timeout_ms, p2_star_timeout_ms, total_timeout_ms;
} PresetRequestTimeouts;
typedef struct PresetCanDiagnosticConfig {
    uint32_t physical_id, response_id, functional_id;
    uint32_t n_bs_ms, n_cr_ms, st_min_ms, max_pdu_len, tx_st_min_override_ms;
    uint8_t addressing, padding_enabled, padding_byte, block_size;
    uint8_t max_wait_frames, tx_st_min_override_enabled, manual_flow_control;
} PresetCanDiagnosticConfig;
typedef struct PresetDiagnosticRequestInfo {
    uint64_t transaction_id, runtime_epoch, elapsed_us, response_len;
    uint32_t state, outcome;
    int32_t error_code;
    uint32_t timeout_kind, pending_count, response_address;
    uint8_t last_nrc;
} PresetDiagnosticRequestInfo;
typedef struct PresetCanCaptureEvent {
    uint64_t sequence, runtime_epoch, transaction_id;
    PresetCanFrameEx frame;
} PresetCanCaptureEvent;
typedef struct PresetCanRuntimeInfo {
    uint64_t runtime_epoch, transaction_id, recovery_remaining_us;
    uint32_t state, physical_id, response_id, functional_id;
    uint8_t call_in_progress;
} PresetCanRuntimeInfo;

/* Every channel argument selects an initialized logical slot, not a vendor
 * handle or physical port. Toomoss/PCAN: config.channel; TSMaster/Vector:
 * config.application_channel. Toomoss panel CAN1 is slot 0. AutoDriver uses
 * config.channel for both its initial logical slot and physical selection.
 * options: command_capacity=1..4096 (suggested 256),
 * recovery_ms=0..60000 (suggested 100). A second runtime on the channel is BUSY.
 * Runtime creation borrows an initialized channel; it never initializes hardware.
 */
PRESET_RS_API int32_t preset_can_runtime_new(PresetDevice *, uint8_t channel,
    const PresetCanRuntimeOptions *, PresetCanRuntime **);
PRESET_RS_API int32_t preset_can_runtime_get_info(PresetCanRuntime *, PresetCanRuntimeInfo *);
PRESET_RS_API int32_t preset_can_runtime_copy_error(PresetCanRuntime *, uint8_t *, size_t, size_t *);
/* Synchronous driver write, bounded queue admission. Occupied IDs return -15
 * except valid manual FC for the active manual-flow-control transaction.
 * Driver write rejection is returned without faulting the runtime. */
PRESET_RS_API int32_t preset_can_runtime_write_with_brs(PresetCanRuntime *, uint32_t id,
    uint8_t is_fd, uint8_t brs, const uint8_t *, size_t);
PRESET_RS_API int32_t preset_can_runtime_write_ex(
    PresetCanRuntime *, const PresetCanFrameEx *frame);
/* Nonblocking cached software estimate; updated at least once each 10 ms worker opportunity. */
PRESET_RS_API int32_t preset_can_runtime_get_bus_load(PresetCanRuntime *, PresetBusLoad *);
/* At most 16 observers, 1..65536 events each; drop oldest. rx/tx are 0 or 1.
 * try_read capacity <= 4096. Frames retain extended/RTR/ESI metadata.
 * Software TX mirrors are emitted after driver success; driver TX echoes are
 * deduplicated, not misrepresented as additional confirmed transmissions.
 */
PRESET_RS_API int32_t preset_can_observer_new(PresetCanRuntime *, const PresetCanObserverOptions *, PresetCanObserver **);
PRESET_RS_API int32_t preset_can_observer_try_read(PresetCanObserver *, PresetCanCaptureEvent *, size_t, size_t *);
PRESET_RS_API int32_t preset_can_observer_get_stats(PresetCanObserver *, PresetRxStats *);
PRESET_RS_API int32_t preset_can_observer_close(PresetCanObserver **);
/* Worker-scheduled periodic transmission. At most 64 groups of 1..64 frames
 * per runtime; each group transmits all its frames once per period_ms (>= 1)
 * and groups run independently. The first tick is immediate. All fields are
 * copied before the call returns. direction must be PRESET_CAN_DIRECTION_TX;
 * flags follow preset_can_runtime_write_ex rules, and scheduling CAN-FD on a
 * classic channel returns PRESET_ERR_INVALID_ARG. Frames whose ID is reserved
 * by the active diagnostic transaction are skipped for that tick, including
 * recovery. Successful writes enter TX observers, logging and bus-load
 * accounting. A driver write rejection is retained in runtime_copy_error
 * without faulting the runtime. Ticks missed while a driver call is slow are
 * skipped rather than burst-transmitted.
 */
PRESET_RS_API int32_t preset_can_periodic_start(
    PresetCanRuntime *,
    const PresetCanFrameEx *frames,
    size_t count,
    uint32_t period_ms,
    PresetCanPeriodic **);
/* Nonblocking and idempotent; stops future ticks. The handle stays valid. */
PRESET_RS_API int32_t preset_can_periodic_cancel(PresetCanPeriodic *);
/* Cancels and releases the handle. Releasing is also required before the
 * runtime can close; a NULL handle or repeated close is PRESET_OK. */
PRESET_RS_API int32_t preset_can_periodic_close(PresetCanPeriodic **);
/* timeouts must all be positive (suggested 1000/5000/30000 ms).
 * One active request per channel, at most 64 retained request results.
 * Start copies PDU and configuration before returning, with no ECU wait.
 * st_min_ms and tx_st_min_override_ms are 0..127 milliseconds;
 * block_size is 0..255 (0 means unlimited). manual_flow_control is 0 or 1.
 * Manual FC suppresses automatic FC but retains reassembly and N_Cr timeout.
 * Send caller-encoded FC bytes through preset_can_runtime_write_with_brs on
 * physical_id with matching CAN/FD and BRS flags. The first three bytes are
 * 0x30|flow_status (0..2), block_size and the ISO-TP wire STmin code.
 * Other writes to diagnostic-reserved IDs remain rejected.
 */
PRESET_RS_API int32_t preset_can_request_start(PresetCanRuntime *, const PresetCanDiagnosticConfig *,
    const PresetRequestTimeouts *, const uint8_t *, size_t, PresetDiagnosticRequest **);
/* LIN payload is SID followed by service data. The client uses 0x3C/0x3D,
 * handles NRC 0x78, and records the actual response NAD. Functional 0x7E and
 * broadcast 0x7F require the network to have one responder. Cancel with
 * preset_request_cancel. A suppressed positive response completes empty when
 * no negative response arrives before the timeout. */
PRESET_RS_API int32_t preset_lin_request_start(PresetLinUdsClient *, const PresetRequestTimeouts *,
    const uint8_t *, size_t, PresetDiagnosticRequest **);
/* Query success is independent of transaction outcome (including UDS negatives). */
PRESET_RS_API int32_t preset_request_get_info(PresetDiagnosticRequest *, PresetDiagnosticRequestInfo *);
PRESET_RS_API int32_t preset_request_cancel(PresetDiagnosticRequest *);
/* wait_timeout_ms=0 polls; -13 means only this wait expired. */
PRESET_RS_API int32_t preset_request_wait(PresetDiagnosticRequest *, uint32_t wait_timeout_ms);
/* Non-consuming, repeatable reads. Response before terminal state returns -14.
 * NULL/0 queries required length. -2 reports insufficient capacity. UTF-8 error
 * text has no terminating NUL; out_len is its exact byte length.
 */
PRESET_RS_API int32_t preset_request_copy_response(PresetDiagnosticRequest *, uint8_t *, size_t, size_t *);
PRESET_RS_API int32_t preset_request_copy_error(PresetDiagnosticRequest *, uint8_t *, size_t, size_t *);
PRESET_RS_API int32_t preset_request_close(PresetDiagnosticRequest **);
/* stop only registers the request. wait_stopped is bounded; it never frees live
 * resources on timeout. close requires STOPPED/FAULTED and no child handles
 * (requests, observers and periodic groups).
 */
PRESET_RS_API int32_t preset_can_runtime_stop(PresetCanRuntime *);
PRESET_RS_API int32_t preset_can_runtime_wait_stopped(PresetCanRuntime *, uint32_t wait_timeout_ms);
PRESET_RS_API int32_t preset_can_runtime_close(PresetCanRuntime **);
/* Recovery defaults to 100 ms, applies to subsequent failed/cancelled requests,
 * and is restricted to 0..60000 ms. Close also waits for recovery via BUSY.
 */
PRESET_RS_API int32_t preset_lin_uds_set_recovery_interval(PresetLinUdsClient *, uint32_t recovery_ms);
PRESET_RS_API int32_t preset_lin_uds_wait_idle(PresetLinUdsClient *, uint32_t wait_timeout_ms);

#undef PRESET_RS_API

#ifdef __cplusplus
}
#endif

#endif
