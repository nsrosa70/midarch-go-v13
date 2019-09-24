package interfaces

import "newsolution/miop"

type Marshaller interface {
	Marshall(packet miop.Packet) []byte
	Unmarshall ([]byte) miop.Packet
}
