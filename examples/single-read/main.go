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
	_ = sensor.SetMode(sds011.QueryMode)

	// Wait to spin fan
	time.Sleep(15 * time.Second)

	measure, _ := sensor.Query()

	fmt.Printf("[%s]\nPM 2.5 => %f μg/m³\nPM 10 => %f μg/m³\n", time.Now().Format("2006-01-02 15:04:05"), measure.PM2_5, measure.PM10)

	_ = sensor.Sleep(true)
}
