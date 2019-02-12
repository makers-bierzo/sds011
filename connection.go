package sds011

type Connection interface {
	Write(data []byte) (int, error)
	Read(buffer []byte) (int, error)
}
