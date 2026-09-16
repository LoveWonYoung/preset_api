#ifndef PRESET_RS_H
#define PRESET_RS_H

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

#if defined(_WIN32) && defined(PRESET_RS_USE_DLLIMPORT)
#define PRESET_RS_API __declspec(dllimport)
#else
#define PRESET_RS_API
#endif

typedef struct PresetDevice PresetDevice;
typedef struct PresetCanUdsClient PresetCanUdsClient;
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
};

enum {
    PRESET_CAP_CLASSIC_CAN = UINT64_C(1) << 0,
    PRESET_CAP_CAN_FD = UINT64_C(1) << 1,
    PRESET_CAP_FUNCTIONAL = UINT64_C(1) << 2,
    PRESET_CAP_RAW_CAN = UINT64_C(1) << 3,
    PRESET_CAP_STANDARD_ID_ONLY = UINT64_C(1) << 4,
    PRESET_CAP_RUNTIME_FLOW_CONTROL = UINT64_C(1) << 5,
    PRESET_CAP_TOOMOSS_SHARED_DEVICE = UINT64_C(1) << 6,
    PRESET_CAP_TOOMOSS_LIN = UINT64_C(1) << 7,
    PRESET_CAP_TOOMOSS = UINT64_C(1) << 8,
    PRESET_CAP_TSMASTER = UINT64_C(1) << 9,
    PRESET_CAP_PCAN = UINT64_C(1) << 10,
    PRESET_CAP_VECTOR = UINT64_C(1) << 11,
    PRESET_CAP_TOOMOSS_ELINS = UINT64_C(1) << 12,
    PRESET_CAP_LIN_UDS = UINT64_C(1) << 13,
    PRESET_CAP_FILE_LOGGING = UINT64_C(1) << 14,
    PRESET_CAP_CAN_CHANNEL_SELECT = UINT64_C(1) << 15,
    PRESET_CAP_BUS_LOAD = UINT64_C(1) << 16,
    PRESET_CAP_RAW_CAN_METADATA = UINT64_C(1) << 17,
    PRESET_CAP_RUNTIME_BRS = UINT64_C(1) << 18,
    PRESET_CAP_EXPLICIT_FD_TIMING = UINT64_C(1) << 19,
    PRESET_CAP_AUTO_DRIVER = UINT64_C(1) << 20,
    PRESET_CAP_PCAN_LIN = UINT64_C(1) << 21,
    PRESET_CAP_TSMASTER_LIN = UINT64_C(1) << 22,
    PRESET_CAP_VECTOR_LIN = UINT64_C(1) << 23,
    PRESET_CAP_RAW_CAN_TIMESTAMP = UINT64_C(1) << 24,
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

typedef struct PresetConfig {
    uint32_t physical_id;
    uint32_t response_id;
    uint32_t functional_id;
    uint32_t n_bs_ms;
    uint32_t n_cr_ms;
    uint32_t st_min_us;
    uint32_t max_pdu_len; /* 1-1048576 bytes */
    uint32_t rx_capacity; /* 1-65536 frames */
    uint8_t is_fd;
    uint8_t padding_enabled;
    uint8_t padding_byte;
    uint8_t block_size;
    uint8_t max_wait_frames;
    uint8_t raw_rx_enabled; /* 1 enables the raw capture ring */
    uint8_t reserved[2];
} PresetConfig;

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

typedef struct PresetAutoConfig {
    uint8_t channel;
    uint8_t is_fd;
    uint8_t brs;
    uint8_t reserved;
    uint32_t nominal_bitrate;
    uint32_t data_bitrate;
    uint32_t rx_buffer_size; /* 1-65536 frames */
    uint32_t poll_batch_size; /* 1-rx_buffer_size */
    uint32_t min_tx_interval_us;
    uint32_t min_rx_poll_interval_us;
    /* PRESET_CAN_BACKEND_* priority order; zero entries are skipped. */
    uint8_t candidate_order[4];
    uint8_t reserved_tail[4];
} PresetAutoConfig;

typedef struct PresetToomossConfig {
    uint8_t channel; /* zero-based: panel CAN1/2/3/4 = 0/1/2/3 */
    uint8_t brs;
    uint8_t reserved[2];
    uint32_t nominal_bitrate;
    uint32_t data_bitrate;
    uint32_t rx_buffer_size; /* 1-65536 frames */
    uint32_t poll_batch_size; /* 1-rx_buffer_size */
    uint32_t min_tx_interval_us;
    uint32_t min_rx_poll_interval_us;
} PresetToomossConfig;

typedef struct PresetToomossLinConfig {
    uint8_t channels_mask; /* PRESET_TOOMOSS_CHANNEL_* bitmask */
    uint8_t master_mode;   /* 1 master, 0 slave */
    uint8_t reserved[2];
    uint32_t baudrate;     /* 2000-5000000, default 19200 */
} PresetToomossLinConfig;

typedef struct PresetToomossElinsConfig {
    uint8_t channels_mask;    /* PRESET_TOOMOSS_CHANNEL_* bitmask */
    uint8_t resistor_enabled; /* 0 or 1 */
    uint8_t version;          /* PRESET_TOOMOSS_ELINS_VER_* */
    uint8_t reserved;
    uint32_t baudrate;        /* 2000-5000000, default 19200 */
    uint32_t receive_timeout_us; /* 0 selects the 1000us default */
} PresetToomossElinsConfig;

typedef struct PresetToomossLinFrame {
    uint32_t timestamp;
    uint8_t frame_id;
    uint8_t data_len;
    uint8_t classic_checksum;
    uint8_t reserved;
    uint8_t data[8];
} PresetToomossLinFrame;

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
    uint8_t version;             /* 1=LIN 1.3, 2=LIN 2.0, 3=LIN 2.1 */
    uint8_t reserved[2];
    int32_t hardware_type;       /* PRESET_VECTOR_HWTYPE_* */
    uint32_t hardware_index;
    uint32_t hardware_channel;
    uint32_t baudrate;           /* 1000-20000, default 19200 */
    uint32_t rx_queue_size;      /* power of two, 8192-524288 */
} PresetVectorLinConfig;

typedef struct PresetLinFrame {
    uint8_t frame_id;
    uint8_t data_len;
    uint8_t reserved[2];
    uint8_t data[8];
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
    uint8_t is_fd;   /* initialize with CAN_InitializeFD */
    uint8_t brs;     /* set BRS on transmitted CAN-FD frames */
    uint8_t reserved;
    uint32_t nominal_bitrate;
    uint32_t data_bitrate;
    uint32_t rx_buffer_size; /* 1-65536 frames */
    uint32_t poll_batch_size; /* 1-rx_buffer_size */
    uint32_t min_tx_interval_us;
    uint32_t min_rx_poll_interval_us;
} PresetPCANConfig;

typedef struct PresetTSMasterConfig {
    uint8_t application_channel; /* zero-based logical channel, 0-31 */
    uint8_t hardware_channel;    /* zero-based physical channel */
    uint8_t is_fd;
    uint8_t brs;
    int32_t hardware_index;      /* zero-based enumerated device index */
    int32_t device_type;
    uint32_t nominal_bitrate;
    uint32_t data_bitrate;
    uint32_t rx_buffer_size; /* 1-65536 frames */
    uint32_t poll_batch_size; /* 1-rx_buffer_size */
    uint32_t min_tx_interval_us;
    uint32_t min_rx_poll_interval_us;
} PresetTSMasterConfig;

typedef struct PresetVectorConfig {
    uint8_t application_channel; /* zero-based logical channel, 0-255 */
    uint8_t is_fd;
    uint8_t brs;
    uint8_t reserved;
    int32_t hardware_type;       /* PRESET_VECTOR_HWTYPE_*; -1 selects any */
    uint32_t hardware_index;     /* zero-based device index */
    uint32_t hardware_channel;   /* zero-based physical channel */
    uint32_t nominal_bitrate;
    uint32_t data_bitrate;
    uint32_t nominal_sjw;
    uint32_t nominal_tseg1;
    uint32_t nominal_tseg2;
    uint32_t data_sjw;
    uint32_t data_tseg1;
    uint32_t data_tseg2;
    uint32_t driver_rx_queue_size; /* power of two, 8192-524288 */
    uint32_t rx_buffer_size; /* 1-65536 frames */
    uint32_t poll_batch_size; /* 1-rx_buffer_size */
    uint32_t min_tx_interval_us;
    uint32_t min_rx_poll_interval_us;
} PresetVectorConfig;

typedef struct PresetCanFrame {
    uint32_t id;
    uint8_t data_len;
    uint8_t is_fd;
    uint8_t reserved[2];
    uint8_t data[64];
} PresetCanFrame;

/* Extended raw-frame representation. timestamp_us is the vendor hardware or
 * driver monotonic timestamp normalized to microseconds; zero means unavailable.
 * Do not mix preset_can_uds_try_read and preset_can_uds_try_read_ex on the same
 * client because both drain the same raw capture ring. */
typedef struct PresetCanFrameEx {
    uint32_t id;
    uint8_t dlc;
    uint8_t data_len;
    uint8_t direction; /* PRESET_CAN_DIRECTION_* */
    uint8_t is_fd;
    uint8_t brs;
    uint8_t reserved[7];
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
PRESET_RS_API PresetConfig preset_default_config(void);
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
 * a UDS client, and releasing a UDS client never closes the device. */
/* Windows AutoDriver tries candidate_order until a mode-compatible backend
 * opens successfully. Defaults: Toomoss, TSMaster/TC1016, PCAN, Vector/VN1640. */
PRESET_RS_API int32_t preset_auto_open(
    const PresetAutoConfig *config,
    PresetDevice **out_device);

PRESET_RS_API int32_t preset_device_get_backend(
    const PresetDevice *device,
    uint8_t *out_backend);

PRESET_RS_API int32_t preset_toomoss_open(PresetDevice **out_device);

/* May be called repeatedly to initialize distinct CAN channels. CAN-FD is
 * attempted first; adapters that reject it automatically fall back to the
 * classic CAN API. A channel that fell back cannot serve an FD UDS client. */
PRESET_RS_API int32_t preset_toomoss_can_init(
    PresetDevice *device,
    const PresetToomossConfig *device_config);
PRESET_RS_API int32_t preset_toomoss_can_init_with_timing(
    PresetDevice *device,
    const PresetToomossConfig *device_config,
    const PresetCanFdTiming *timing);

PRESET_RS_API int32_t preset_toomoss_lin_init(
    PresetDevice *device,
    const PresetToomossLinConfig *config);

PRESET_RS_API int32_t preset_toomoss_lin_write(
    PresetDevice *device,
    uint8_t channel,
    uint8_t frame_id,
    const uint8_t *data,
    size_t data_len);

/* Returns PRESET_ERR_TIMEOUT when no slave responds to the master header. */
PRESET_RS_API int32_t preset_toomoss_lin_read(
    PresetDevice *device,
    uint8_t channel,
    uint8_t frame_id,
    PresetToomossLinFrame *out_frame);

/* Read frames buffered by LIN_EX_SlaveGetData. On input, *inout_len is the
 * frame capacity (1-512); on success it is the number of frames written. */
PRESET_RS_API int32_t preset_toomoss_lin_try_read(
    PresetDevice *device,
    uint8_t channel,
    PresetToomossLinFrame *frames,
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

/* Windows LIN master backends. Each open explicitly initializes one channel
 * and returns an owning device handle. The matching init function adds another
 * channel to a device of the same vendor. CAN and LIN channels may share that
 * vendor device handle. Add all TSMaster mappings before creating any client. */
PRESET_RS_API int32_t preset_pcan_lin_open(
    const PresetPCANLinConfig *config,
    PresetDevice **out_device);
PRESET_RS_API int32_t preset_pcan_lin_init(
    PresetDevice *device,
    const PresetPCANLinConfig *config);
PRESET_RS_API int32_t preset_tsmaster_lin_open(
    const PresetTSMasterLinConfig *config,
    PresetDevice **out_device);
PRESET_RS_API int32_t preset_tsmaster_lin_init(
    PresetDevice *device,
    const PresetTSMasterLinConfig *config);
PRESET_RS_API int32_t preset_vector_lin_open(
    const PresetVectorLinConfig *config,
    PresetDevice **out_device);
PRESET_RS_API int32_t preset_vector_lin_init(
    PresetDevice *device,
    const PresetVectorLinConfig *config);

PRESET_RS_API int32_t preset_lin_master_write(
    PresetDevice *device,
    uint8_t channel,
    uint8_t frame_id,
    const uint8_t *data,
    size_t data_len);

/* Sends the requested master header and waits up to timeout_ms. */
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

/* payload contains SID followed by service data. The timeout covers request
 * transmission and response reception. The client builds LIN TP frames on
 * 0x3C, polls 0x3D, reassembles the response, handles NRC 0x78 and returns the
 * actual response NAD. Physical NADs must be 0x01-0x7D. Functional NAD 0x7E
 * and broadcast/configuration NAD 0x7F accept a response from an actual node
 * NAD; callers must use them only when the network guarantees one responder.
 * A request that suppresses its positive response succeeds with length zero
 * when the response timeout expires without a negative response. */
PRESET_RS_API int32_t preset_lin_uds_request(
    PresetLinUdsClient *client,
    const uint8_t *payload,
    size_t payload_len,
    uint32_t timeout_ms,
    uint8_t *out_nad,
    uint8_t *out_data,
    size_t out_capacity,
    size_t *out_len);

/* Set both the LIN TP inter-frame transmit delay and response polling interval.
 * The default is 10 ms. Valid values are 1 through 1000 ms, and a new value is
 * applied to requests started after this call. */
PRESET_RS_API int32_t preset_lin_uds_set_poll_interval(
    PresetLinUdsClient *client,
    uint32_t interval_ms);

/* Releasing a LIN UDS client never closes its physical device. */
PRESET_RS_API int32_t preset_lin_uds_client_close(
    PresetLinUdsClient **client);

/* PCAN, TSMaster, and Vector open/init calls return PRESET_ERR_UNSUPPORTED off
 * Windows. Each open initializes the first channel. The matching can_init may
 * then add distinct channels to the same PresetDevice. Add all TSMaster
 * channels before creating CAN clients. */
PRESET_RS_API int32_t preset_pcan_open(
    const PresetPCANConfig *device_config,
    PresetDevice **out_device);
PRESET_RS_API int32_t preset_pcan_open_with_timing(
    const PresetPCANConfig *device_config,
    const PresetCanFdTiming *timing,
    PresetDevice **out_device);
PRESET_RS_API int32_t preset_pcan_can_init(
    PresetDevice *device,
    const PresetPCANConfig *device_config);
PRESET_RS_API int32_t preset_pcan_can_init_with_timing(
    PresetDevice *device,
    const PresetPCANConfig *device_config,
    const PresetCanFdTiming *timing);

PRESET_RS_API int32_t preset_tsmaster_open(
    const PresetTSMasterConfig *device_config,
    PresetDevice **out_device);
PRESET_RS_API int32_t preset_tsmaster_can_init(
    PresetDevice *device,
    const PresetTSMasterConfig *device_config);

PRESET_RS_API int32_t preset_vector_open(
    const PresetVectorConfig *device_config,
    PresetDevice **out_device);
PRESET_RS_API int32_t preset_vector_can_init(
    PresetDevice *device,
    const PresetVectorConfig *device_config);

/* The CAN UDS client borrows RX/TX from one initialized channel. Different
 * channels may have active clients concurrently; each individual channel may
 * have at most one active client. For Toomoss and PCAN, channel is the
 * config channel (Toomoss: panel CAN1 = 0). TSMaster uses the corresponding
 * config channel; Vector uses application_channel. */
PRESET_RS_API int32_t preset_can_uds_client_new(
    PresetDevice *device,
    uint8_t channel,
    const PresetConfig *config,
    PresetCanUdsClient **out_client);

/* Close any backend. Returns PRESET_ERR_BUSY and leaves *device unchanged while
 * a CAN or LIN UDS client is active. Release clients, then retry the close. */
PRESET_RS_API int32_t preset_device_close(PresetDevice **device);

/* For recognized sub-function services, a suppressed positive response
 * completes successfully with an empty output buffer if no negative response
 * arrives before timeout. SID 0x2F is not a sub-function service. */
PRESET_RS_API int32_t preset_can_uds_request(
    PresetCanUdsClient *client,
    const uint8_t *payload,
    size_t payload_len,
    uint32_t timeout_ms,
    uint8_t *out_data,
    size_t out_capacity,
    size_t *out_len);

PRESET_RS_API int32_t preset_can_uds_functional_request(
    PresetCanUdsClient *client,
    const uint8_t *payload,
    size_t payload_len,
    uint32_t timeout_ms,
    uint8_t *out_data,
    size_t out_capacity,
    size_t *out_len);

/* Update values advertised by subsequently generated ISO-TP flow-control
 * frames. STmin is expressed in whole milliseconds (0-127); BlockSize is
 * 0-255, where 0 allows all remaining consecutive frames. */
PRESET_RS_API int32_t preset_can_uds_set_default_st_min(
    PresetCanUdsClient *client,
    uint32_t st_min_ms);

PRESET_RS_API int32_t preset_can_uds_set_default_block_size(
    PresetCanUdsClient *client,
    uint32_t block_size);

/* When enabled is 1, automatic flow-control transmission is disabled while
 * ISO-TP receive state and timers remain active. The caller must send flow-
 * control frames through preset_can_uds_write. Pass 0 to restore automatic
 * mode. */
PRESET_RS_API int32_t preset_can_uds_set_manual_flow_control(
    PresetCanUdsClient *client,
    uint8_t enabled);

/* Changes/reads BRS for subsequently transmitted CAN-FD frames. Classic CAN
 * frames are unaffected. */
PRESET_RS_API int32_t preset_can_uds_set_brs(
    PresetCanUdsClient *client,
    uint8_t enabled);
PRESET_RS_API int32_t preset_can_uds_get_brs(
    PresetCanUdsClient *client,
    uint8_t *out_enabled);

/* When enabled, successfully transmitted frames also enter the raw ring with
 * direction PRESET_CAN_DIRECTION_TX. Raw capture must have been enabled in
 * PresetConfig before creating the client. */
PRESET_RS_API int32_t preset_can_uds_set_raw_tx_echo(
    PresetCanUdsClient *client,
    uint8_t enabled);

PRESET_RS_API int32_t preset_can_uds_write(
    PresetCanUdsClient *client,
    uint32_t id,
    uint8_t is_fd,
    const uint8_t *data,
    size_t data_len);

/* Raw capture must be enabled with PresetConfig.raw_rx_enabled before the
 * client is created. ISO-TP response processing is independent of this ring. */
PRESET_RS_API int32_t preset_can_uds_try_read(
    PresetCanUdsClient *client,
    PresetCanFrame *frames,
    size_t *inout_len);

PRESET_RS_API int32_t preset_can_uds_try_read_ex(
    PresetCanUdsClient *client,
    PresetCanFrameEx *frames,
    size_t *inout_len);

PRESET_RS_API int32_t preset_can_uds_rx_get_stats(
    PresetCanUdsClient *client,
    PresetRxStats *out_stats);

/* Enable or disable continuous bus-load measurement. It is enabled by default
 * for compatibility. Changing the state clears the current window. Pass only
 * 0 or 1. */
PRESET_RS_API int32_t preset_can_uds_set_bus_load_enabled(
    PresetCanUdsClient *client,
    uint8_t enabled);

/* Estimated total bus occupancy over a one-second sliding window. Both
 * received frames and successfully transmitted frames are included. TX echo
 * frames matching a recent transmission are counted only once. When
 * measurement is disabled, load and frame_count are zero. */
PRESET_RS_API int32_t preset_can_uds_get_bus_load(
    PresetCanUdsClient *client,
    PresetBusLoad *out_load);

/* Includes asynchronous errors reported by the CAN worker. */
PRESET_RS_API int32_t preset_can_uds_last_error(
    PresetCanUdsClient *client,
    uint8_t *out_data,
    size_t out_capacity,
    size_t *out_len);

/* Latest synchronous C API error on the calling thread, including device
 * open/init failures. */
PRESET_RS_API int32_t preset_last_error(
    uint8_t *out_data,
    size_t out_capacity,
    size_t *out_len);

/* Releasing a CAN UDS client stops its worker but never closes its device. */
PRESET_RS_API int32_t preset_can_uds_client_close(
    PresetCanUdsClient **client);

#undef PRESET_RS_API

#ifdef __cplusplus
}
#endif

#endif
