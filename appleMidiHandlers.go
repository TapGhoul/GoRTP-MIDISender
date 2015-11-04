package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"log"
)

func CheckApplePacketValid(pBuf []byte, totalBytes int) (dataBuf *bytes.Buffer, command uint16, err error) {
	dataBuf = bytes.NewBuffer(pBuf[:totalBytes])

	if signature := dataBuf.Next(2); !bytes.Equal(signature, []byte{0xff, 0xff}) {
		log.Printf("Signature invalid! 0x%.4x\n", signature)
		return nil, 0, errors.New("appleMidiHandler: Signature invalid")
	}

	command = binary.BigEndian.Uint16(dataBuf.Next(2))

	if protoVersion := binary.BigEndian.Uint32(dataBuf.Next(4)); protoVersion != 2 {
		log.Printf("Invalid protocol version %d\n", protoVersion)
		return nil, 0, errors.New("appleMidiHandler: Version invalid")
	}

	return dataBuf, command, nil
}

func HandleInvitation(dataBuf *bytes.Buffer) (session AppleMidiSession) {
	var err error

	session.InitToken = dataBuf.Next(4)[:]
	session.LocalSSRC = make([]byte, 4)
	session.RemoteSSRC = dataBuf.Next(4)[:]
	session.LocalName = "PicartoTVPikachu"
	session.RemoteName, err = dataBuf.ReadString(0)
	if err != nil {
		panic(err)
	}

	_, err = rand.Read(session.LocalSSRC)
	if err != nil {
		panic(err)
	}

	return session
}
