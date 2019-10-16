package models

import (
	. "bytes"
	"fmt"
)

func welcome(buf *Buffer) {
	if b, exrr := buf.ReadByte(); exrr != nil {
		fmt.Println(b)
	}
}

func main() {
	arrB := []byte{0, 2, 147, 22, 56, 127}
	b := NewBuffer(arrB)
	welcome(b)
}