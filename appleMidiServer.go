package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

var controlSrv *net.UDPConn
var appleSessions AppleMidiSessions

func StartAppleMidi() {
	netAddr, err := net.ResolveUDPAddr("udp", ":5004")
	if err != nil {
		panic(err)
	}
	controlSrv, err = net.ListenUDP("udp", netAddr)
	if err != nil {
		panic(err)
	}

	appleSessions = make(AppleMidiSessions)

	for {
		pBuf := make([]byte, 64)
		totalBytes, remoteAddr, err := controlSrv.ReadFromUDP(pBuf)
		if err != nil {
			panic(err)
		}

		log.Printf("Got %d bytes from %s\n", totalBytes, remoteAddr.String())

		dataBuf, command, err := CheckApplePacketValid(pBuf, totalBytes)
		if err != nil {
			log.Println(err)
			continue
		}

		//command := binary.BigEndian.Uint16(dataBuf.Next(2))
		if command == 0x494e { // 'IN'
			log.Println("IN packet")
			session := HandleInvitation(dataBuf)
			appleSessions[string(session.RemoteSSRC)] = session

			var payload bytes.Buffer

			payload.Write([]byte{0xff, 0xff, 0x4f, 0x4b})
			binary.Write(&payload, binary.BigEndian, uint32(2))
			payload.Write(session.InitToken)
			payload.Write(session.LocalSSRC)
			payload.WriteString(session.LocalName)
			payload.WriteByte(0x00)

			_, err = controlSrv.WriteToUDP(payload.Bytes(), remoteAddr)

			// TODO: Start session handler

		} else if command == 0x5253 { // 'RS'
			log.Println("RS packet")
			session, ok := appleSessions[string(dataBuf.Next(4))]
			if ok {
				session.MidiAck <- binary.BigEndian.Uint16(dataBuf.Next(2))

			}
		} else if command == 0x4259 { // 'BY'
			log.Println("BY packet")
			initToken := dataBuf.Next(4)[:]
			senderSSRC := string(dataBuf.Next(4))
			if bytes.Equal(appleSessions[senderSSRC].InitToken, initToken) {
				close(appleSessions[senderSSRC].MidiAck)
				delete(appleSessions, senderSSRC)
			} else {
				log.Println("OI! Someone else tried to disconnect us! Assholes...")
			}
		} else {
			log.Println("?? packet")
		}

	}
}
