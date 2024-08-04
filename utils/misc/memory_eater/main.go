package main

import (
	"fmt"
	"os"
)

func main() {
	var buf []byte

	tenMb, kbInMb := 10485760, 1024

	for range 10 {
		buf = append(buf, make([]byte, tenMb)...) // 10 Mb
		msg := fmt.Sprintf("ate %d MB \n", len(buf)/kbInMb)
		_, _ = os.Stdout.WriteString(msg)
	}
}
