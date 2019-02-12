package sds011

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"
)

type sensor struct {
	conn          Connection
	id            uint16
	onMeasurement []chan Measurement
	firmware      time.Time
}

func NewSensor(conn Connection) *sensor {

	return &sensor{
		conn: conn,
		id:   0xffff,
	}
}

func (sensor *sensor) OnQuery(ch chan Measurement) {
	for i := range sensor.onMeasurement {
		if sensor.onMeasurement[i] == ch {
			return
		}
	}
	sensor.onMeasurement = append(sensor.onMeasurement, ch)
}

func (sensor *sensor) Listen() {
	go func(handlers []chan Measurement) {
		for true {
			measure, err := sensor.queryResponse()
			if err != nil {
				for i := range handlers {
					handlers[i] <- *measure
				}
			}
		}
	}(sensor.onMeasurement)
}

func (sensor *sensor) emit(measure Measurement) {
	for i := range sensor.onMeasurement {
		go func(handler chan Measurement) {
			handler <- measure
		}(sensor.onMeasurement[i])
	}
}

func (sensor *sensor) write(cmd Command, data []byte) error {
	var commandData = make([]byte, 12)
	copy(commandData, data)
	checksum := (sum(data) + byte(cmd) - 0x02) % 0xff

	var ret bytes.Buffer
	ret.Write([]byte{0xaa, 0xb4}) // Sender head
	ret.WriteByte(byte(cmd))
	ret.Write(commandData)
	ret.Write(parseId(sensor.id)) // Device ID
	ret.WriteByte(checksum)       // Checksum
	ret.WriteByte(0xab)        // Tail

	_, err := sensor.conn.Write(ret.Bytes())
	return err
}

func (sensor *sensor) read() []byte {
	responseHead := make([]byte, 1)
	_, err := sensor.conn.Read(responseHead)
	for err == nil && responseHead[0] != 0xaa {
		_, err = sensor.conn.Read(responseHead)
	}

	response := make([]byte, 10)
	count, err := sensor.conn.Read(response)

	return response[:count]
}

func (sensor *sensor) Sleep(sleep bool) error {
	var mode byte = 0x00
	if !sleep {
		mode = 0x01
	}

	err := sensor.write(SleepCommand, []byte{0x01, mode})
	if err != nil {
		return err
	}

	if sensor.read()[1] == byte(SleepCommand) {
		return nil
	} else {
		return errors.New("sleep bad response")
	}
}

func (sensor *sensor) SetMode(mode Mode) error {
	err := sensor.write(ModeCommand, []byte{0x01, byte(mode)})
	if err != nil {
		return err
	}

	if sensor.read()[1] == byte(ModeCommand) {
		return nil
	} else {
		return errors.New("set mode bad response")
	}
}

func (sensor *sensor) queryResponse() (*Measurement, error) {
	response := sensor.read()

	if response[0] == byte(0xc0) {
		measure := Measurement{
			PM2_5: float32(binary.LittleEndian.Uint16(response[1:3])) / 10.0,
			PM10:  float32(binary.LittleEndian.Uint16(response[3:5])) / 10.0,
		}

		sensor.emit(measure)

		return &measure, nil
	} else {
		return nil, errors.New("query bad response")
	}
}

func (sensor *sensor) Query() (*Measurement, error) {
	err := sensor.write(QueryDataCommand, nil)
	if err != nil {
		return nil, err
	}

	return sensor.queryResponse()
}

func (sensor *sensor) GetFirmwareVersion() (time.Time, error) {
	if !sensor.firmware.IsZero() {
		return sensor.firmware, nil
	}

	err := sensor.write(FirmwareCommand, nil)
	if err != nil {
		return time.Time{}, err
	}

	response := sensor.read()

	if response[1] == byte(FirmwareCommand) {
		sensor.firmware = time.Date(2000+int(response[2]), time.Month(response[3]), int(response[4]), 0, 0, 0, 0, time.UTC)
		return sensor.firmware, nil
	} else {
		return time.Time{}, errors.New("get firmware bad response")
	}
}

func (sensor *sensor) SetWorkingPeriod(minutes uint) error {
	if minutes < 0 {
		minutes = 0
	} else if minutes > 30 {
		minutes = 30
	}
	err := sensor.write(WorkingPeriodCommand, []byte{0x1, byte(minutes)})
	if err != nil {
		return err
	}

	if sensor.read()[1] == byte(WorkingPeriodCommand) {
		return nil
	} else {
		return errors.New("set working period bad response")
	}
}

func (sensor *sensor) SetId(id uint16) error {
	err := sensor.write(DeviceIdCommand, append(make([]byte, 10), parseId(id)...))
	if err != nil {
		return err
	}

	if sensor.read()[1] == byte(DeviceIdCommand) {
		sensor.id = id
		return nil
	} else {
		return errors.New("set id bad response")
	}
}

func (sensor *sensor) GetId() uint16 {
	return sensor.id
}
