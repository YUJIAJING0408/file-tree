package fileTree

import (
	"strconv"
)

/*
@Date:
@Auth: YUJIAJING
@Desp:
*/

type Node interface {
	GetName() string
	GetFullPath() string
	GetSize() int64
	IsDir() bool
	String() string
	Print(int)
}

const (
	TypeDir = iota
	TypeFile
	TypeLink
)

const (
	Byte     = 1
	KiloByte = Byte << 10
	MegaByte = KiloByte << 10
	GigaByte = MegaByte << 10
	TeraByte = GigaByte << 10
)

func ByteString(b int64) (res string) {
	if b <= 0 {
		return "0B"
	}
	T := b / TeraByte
	b = b % TeraByte
	G := b / GigaByte
	b = b % GigaByte
	M := b / MegaByte
	b = b % MegaByte
	K := b / KiloByte
	B := b % KiloByte
	if T != 0 {
		res += strconv.Itoa(int(T)) + "T"
	}
	if G != 0 {
		res += strconv.Itoa(int(G)) + "G"
	}
	if M != 0 {
		res += strconv.Itoa(int(M)) + "M"
	}
	if K != 0 {
		res += strconv.Itoa(int(K)) + "K"
	}
	if B != 0 {
		res += strconv.Itoa(int(B)) + "B"
	}
	return res
}
