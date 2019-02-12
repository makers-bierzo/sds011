package go_laser_dust

type Connection interface {
	Write(data []byte) (int, error)
	Read(buffer []byte) (int, error)
}
