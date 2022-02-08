package main

import "encoding/binary"

type Frame struct {
	fin           int
	rsv1          int
	rsv2          int
	rsv3          int
	opcode        int
	mask          int
	payloadLength int
	maskingKey    []byte
	payloadData   []byte
}

func buildFrame(msg string) *Frame {
	message := []byte(msg)
	return &Frame{
		fin:           1,
		rsv1:          0,
		rsv2:          0,
		rsv3:          0,
		opcode:        1,
		mask:          0,
		payloadLength: len(message),
		payloadData:   message,
	}
}

func (f *Frame) parse(buffer []byte) {
	//最初の byte を読み込む
	index := 0
	firstByte := int(buffer[index])

	f.fin = (firstByte & 0x80) >> 7
	f.rsv1 = (firstByte & 0x40) >> 6
	f.rsv2 = (firstByte & 0x20) >> 5
	f.rsv3 = (firstByte & 0x10) >> 4
	f.opcode = firstByte & 0x0F

	//次の byte を読み込む
	index += 1
	secondByte := int(buffer[index])

	f.mask = (secondByte & 0x80) >> 7
	f.payloadLength = secondByte & 0x7F

	//次の byte を読み込む
	index += 1

	if f.payloadLength == 126 {
		// 長さが126の場合は、次の 2byte が UInt16 として 本当の Payload length となる
		length := binary.BigEndian.Uint16(buffer[index:(index + 2)])
		f.payloadLength = int(length)
		index += 2
	} else if f.payloadLength == 127 {
		// 長さが 127 の場合は、次の 8byte が UInt64 として 本当の Payload length となる
		length := binary.BigEndian.Uint64(buffer[index:(index + 8)])
		f.payloadLength = int(length)
		index += 8
	}

	if f.mask > 0 {
		f.maskingKey = buffer[index:(index + 4)]
		index += 4
	}

	payload := buffer[index:(index + f.payloadLength)]

	if f.mask > 0 {
		for i := 0; i < f.payloadLength; i++ {
			payload[i] ^= f.maskingKey[i%4]
		}
	}

	f.payloadData = payload
}

func (f *Frame) toBytes() (data []byte) {
	bits := 0
	bits |= (f.fin << 7)
	bits |= (f.rsv1 << 6)
	bits |= (f.rsv2 << 5)
	bits |= (f.rsv3 << 4)
	bits |= f.opcode

	// first byte を追加
	data = append(data, byte(bits))

	bits = 0
	bits |= (f.mask << 7)
	bits |= f.payloadLength // 長さは 126 未満と仮定

	// second byte を追加
	data = append(data, byte(bits))

	// 実際のデータを追加
	data = append(data, f.payloadData...)

	return data
}
