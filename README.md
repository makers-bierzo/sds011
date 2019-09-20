# go-sds011
Go library to control a SDS011 laser dust sensor.

## Example

First, open serial port
```go
c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
s, _ := serial.OpenPort(c)
```

Create a sds011 sensor instance over this connection
```go
sensor := sds011.NewSensor(s)
```

Disable sleep mode to take samples
```go
_ = sensor.Sleep(false)
```

Select sensor mode. On QueryMode, sensor only read when you request a reading 
```go
_ = sensor.SetMode(sds011.QueryMode)
```

You should wait a little while for the fan to collect an air sample
```go
time.Sleep(15 * time.Second)
```

Read and show the measure
```go
measure, _ := sensor.Query()
fmt.Printf("PM 2.5 => %f μg/m³\nPM 10 => %f μg/m³\n", measure.PM2_5, measure.PM10)
```

You can change sensor to ActiveMode. In this mode the sensor takes a sample periodically, in this case, every 1 minute.
```go
_ = sensor.SetMode(sds011.ActiveMode)
_ = sensor.SetWorkingPeriod(1)
```

To receive the samples, when they are read, you must listen new values from [channel](https://gobyexample.com/channels) of query responses.
```go
measureChannel := make(chan sds011.Measurement)

sensor.OnQuery(measureChannel)
sensor.Listen()

for true {
    measure := <- measureChannel
    fmt.Printf("PM 2.5 => %f μg/m³\nPM 10 => %f μg/m³\n", measure.PM2_5, measure.PM10)
}
```