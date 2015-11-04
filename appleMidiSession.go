package main

type AppleMidiSession struct {
	initToken  []byte
	localSSRC  []byte
	remoteSSRC []byte
	localName  string
	remoteName string
}

type AppleMidiSessions map[string]AppleMidiSession
