package main

import (
	"fmt"
	"os"
)

func main() {
	var buf []byte

	tenMb, kbInMb := 10485760, 1024

	for i := 0; i < 10; i++ {
		buf = append(buf, make([]byte, tenMb)...) // 10 Mb
		msg := fmt.Sprintf("ate %d MB \n", len(buf)/kbInMb)
		_, _ = os.Stdout.WriteString(msg)
	}
}
