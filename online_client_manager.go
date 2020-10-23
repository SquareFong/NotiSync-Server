package main

type client struct {
	detail         phoneDetail
	allMessages    []message
	command        int
	lastActiveTime int
	UUID           string
}
