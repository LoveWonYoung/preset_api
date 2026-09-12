package preset_api

import "syscall"

var (
	lib, _                                 = syscall.LoadLibrary("preset_rs.dll")
	preset_default_config, _               = syscall.GetProcessAddress(lib, "preset_default_config")
	preset_toomoss_default_config, _       = syscall.GetProcessAddress(lib, "preset_toomoss_default_config")
	preset_toomoss_lin_default_config, _   = syscall.GetProcessAddress(lib, "preset_toomoss_lin_default_config")
	preset_toomoss_elins_default_config, _ = syscall.GetProcessAddress(lib, "preset_toomoss_elins_default_config")
	preset_pcan_default_config, _          = syscall.GetProcessAddress(lib, "preset_pcan_default_config")
	preset_tsmaster_default_config, _      = syscall.GetProcessAddress(lib, "preset_tsmaster_default_config")
	preset_vector_default_config, _        = syscall.GetProcessAddress(lib, "preset_vector_default_config")
)

func DefaultConfig() (c PresetConfig) {

}

func ToomossDefaultConfig() (c PresetToomossConfig) {

}

func ToomossLinDefaultConfig() (c PresetToomossLinConfig) {

}

func ToomossElinsDefaultConfig() (c PresetToomossElinsConfig) {

}

func PcanDefaultConfig() (c PresetPCANConfig) {

}

func TsmasterDefaultConfig() (c PresetTSMasterConfig) {

}

func VectorDefaultConfig() (c PresetVectorConfig) {

}
