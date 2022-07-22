package pscan

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"time"
)

const (
	SourcePortFirstDefault = 32768

	SourcePortLastDefault = 61000

	PacketStreamsDefault = 1
)

type PValidate struct {
	sourcePortFirst uint32

	sourcePortLast uint32

	packetStreams uint32

	numPorts uint32

	key []byte

	aes cipher.Block
}

type Validator struct {
	output []byte
}

var pv *PValidate

func makeKey() []byte {

	key := make([]byte, 16)
	send := time.Now().UnixNano()

	binary.LittleEndian.PutUint64(key, uint64(send))

	return key
}

func InitPValidate(sourcePortFirst, sourcePortLast, packetStreams uint32) {

	key := makeKey()

	aes, _ := aes.NewCipher(key)

	if sourcePortFirst == 0 {
		sourcePortFirst = SourcePortFirstDefault
	}

	if sourcePortLast == 0 {
		sourcePortLast = SourcePortLastDefault
	}

	if packetStreams == 0 {

		packetStreams = PacketStreamsDefault
	}

	numPorts := sourcePortLast - sourcePortFirst + 1

	pv = &PValidate{
		sourcePortFirst: sourcePortFirst,
		sourcePortLast:  sourcePortLast,
		packetStreams:   packetStreams,
		numPorts:        numPorts,
		key:             key,
		aes:             aes,
	}
}

func enCrypto(src, dst uint32) []byte {

	var srcPort uint32 = 0
	var dstPort uint32 = 0

	input := make([]byte, 16)
	output := make([]byte, 16)

	bbuf := bytes.NewBuffer(input)

	binary.Write(bbuf, binary.BigEndian, &src)
	binary.Write(bbuf, binary.BigEndian, &dst)
	binary.Write(bbuf, binary.BigEndian, &srcPort)
	binary.Write(bbuf, binary.BigEndian, &dstPort)

	pv.aes.Encrypt(output, input)

	return output
}

func GenValidator(src, dst uint32) *Validator {

	enBytes := enCrypto(src, dst)

	return &Validator{output: enBytes}

}

func (v *Validator) GetSequence() uint32 {

	return binary.BigEndian.Uint32(v.output[:4])
}

func (v *Validator) GetDstPort(probeNum uint32) uint32 {

	portInt := binary.BigEndian.Uint32(v.output[4:8])
	return pv.sourcePortFirst + ((portInt + probeNum) % pv.numPorts)
}

func (v *Validator) IsValidDstPort(port uint32) bool {

	if port > pv.sourcePortLast || port < pv.sourcePortFirst {
		return false
	}

	toValidate := port - pv.sourcePortFirst
	portInt := binary.BigEndian.Uint32(v.output[4:8])

	min := portInt % pv.numPorts

	max := (portInt + pv.packetStreams) % pv.numPorts

	return ((max - min) % pv.numPorts) >= ((toValidate - min) % pv.numPorts)
}
