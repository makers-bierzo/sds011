package sds011

type Command byte

const (
	ModeCommand          Command = 0x02
	QueryDataCommand     Command = 0x04
	DeviceIdCommand      Command = 0x05
	SleepCommand         Command = 0x06
	FirmwareCommand      Command = 0x07
	WorkingPeriodCommand Command = 0x08
)
