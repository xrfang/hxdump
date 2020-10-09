package main

import (
	"crypto/rand"
	"fmt"

	"github.com/xrfang/hxdump"
)

func main() {
	buf := make([]byte, 59)
	rand.Read(buf)
	fmt.Println(hxdump.Dump(buf))
}
