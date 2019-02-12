package main

import (
	"fmt"
	"github.com/tarm/serial"
	"github.com/tokkenno/go-laser-dust"
	"log"
	"time"
)

func main() {
	c := &serial.Config{Name: "COM6", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	sensor := go_laser_dust.NewSensor(s)

	_ = sensor.Sleep(false)
	_ = sensor.SetMode(go_laser_dust.QueryMode)

	// Wait to spin fan
	time.Sleep(15 * time.Second)

	measure, _ := sensor.Query()

	fmt.Printf("[%s]\nPM 2.5 => %f μg/m³\nPM 10 => %f μg/m³\n", time.Now().Format("2006-01-02 15:04:05"), measure.PM2_5, measure.PM10)

	_ = sensor.Sleep(true)
}
