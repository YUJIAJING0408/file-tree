package fileTree

import (
	"fmt"
	"testing"
)

/*
@Date:
@Auth: YUJIAJING
@Desp:
*/

func TestByte(t *testing.T) {
	var size = 1234567890
	res := ByteString(int64(size))
	fmt.Println(res)
	toByte, err := StringToByte(res)
	fmt.Println(toByte, err)
}
