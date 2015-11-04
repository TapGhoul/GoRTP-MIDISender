package main

type AppleMidiSession struct {
	MidiAck    chan uint16
	InitToken  []byte
	LocalSSRC  []byte
	RemoteSSRC []byte
	LocalName  string
	RemoteName string
}

type AppleMidiSessions map[string]AppleMidiSession
