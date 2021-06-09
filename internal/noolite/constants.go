package noolite

const respLen = 17

const (
	RespSt byte = iota + 173
	RespSp
)

const (
	CommandSt byte = iota + 171
	CommandSp
)

type Mode byte

const (
	ModeTx Mode = iota
	ModeRx
	ModeFTX
	ModeFRX
	ModeFService
	ModeFUpdate
)

type RespCtr byte

const (
	RespCtrDone RespCtr = iota
	RespCtrNoResp
	RespCtrError
	RespCtrBinded
)

type CommandCtr byte

const (
	CommandCtrSend CommandCtr = iota
	CommandCtrSendBroadcast
	CommandCtrReadResp
	CommandCtrBindOn
	CommandCtrBindOff
	CommandCtrClearChannel
	CommandCtrClearAll
	CommandCtrUnbindAddressFromChannel
	CommandCtrSendToAddress
)

type Cmd byte

const (
	CmdOff Cmd = iota
	CmdBrightDown
	CmdOn
	CmdBrightUp
	CmdSwitch
	CmdBrightBack
	CmdSetBrightness
	CmdLoadPreset
	CmdSavePreset
	CmdUnbind
	CmdStopReg
	CmdBrightStepDown
	CmdBrightStepUp
	CmdBrightReg
	CmdBind
	CmdRollColour
	CmdSwitchColour
	CmdSwitchMode
	CmdSpeedModeBack
	CmdBatteryLow
	CmdSensTempHumi
	CmdTemporaryOn Cmd = 25
	CmdModes       Cmd = 26  // only fo simple nooLite
	CmdReadState   Cmd = 128 // only for nooLite-F
	CmdWriteState  Cmd = 129 // only for nooLite-F
	CmdSendState   Cmd = 130 // only for nooLite-F
	CmdService     Cmd = 131 // only for nooLite-F
	CmdClearMemory Cmd = 132 // only for nooLite-F
)
