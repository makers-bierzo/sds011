package main

import (
	"fmt"
	"github.com/makers-bierzo/sds011"
	"github.com/tarm/serial"
	"log"
	"time"
)

func main() {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	sensor := sds011.NewSensor(s)

	_ = sensor.Sleep(false)
	_ = sensor.SetWorkingPeriod(1)
	_ = sensor.SetMode(sds011.ActiveMode)

	measureChannel := make(chan sds011.Measurement)
	sensor.OnQuery(measureChannel)

	sensor.Listen()

	for true {
		measure := <- measureChannel
		fmt.Printf("[%s]\nPM 2.5 => %f μg/m³\nPM 10 => %f μg/m³\n", time.Now().Format("2006-01-02 15:04:05"), measure.PM2_5, measure.PM10)
	}
}

