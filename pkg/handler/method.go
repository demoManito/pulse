package handler

type Method int

const (
	GET  Method = 1 << iota
	POST Method = 1 << iota
)
