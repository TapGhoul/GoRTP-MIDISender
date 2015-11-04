package main

type AppleMidiSession struct {
	InitToken  []byte
	LocalSSRC  []byte
	RemoteSSRC []byte
	LocalName  string
	RemoteName string
}

type AppleMidiSessions map[string]AppleMidiSession
