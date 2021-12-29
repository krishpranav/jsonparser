package jstream

var (
	smallInput  = make([]byte, 1024*12)       // 12K
	mediumInput = make([]byte, 1024*1024*12)  // 12MB
	largeInput  = make([]byte, 1024*1024*128) // 128MB
)
