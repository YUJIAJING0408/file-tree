package fileTree

import (
	"strconv"
	"strings"
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

func StringToByte(s string) (res int, err error) {
	s = strings.ToUpper(s)
	split := strings.Split(s, "T")
	if len(split) == 2 {
		T, err := strconv.Atoi(split[0])
		if err != nil {
			return 0, err
		}
		res += T * TeraByte
		s = split[1]
	}
	split = strings.Split(s, "G")
	if len(split) == 2 {
		G, err := strconv.Atoi(split[0])
		if err != nil {
			return 0, err
		}
		res += G * GigaByte
		s = split[1]
	}
	split = strings.Split(s, "M")
	if len(split) == 2 {
		M, err := strconv.Atoi(split[0])
		if err != nil {
			return 0, err
		}
		res += M * MegaByte
		s = split[1]
	}
	split = strings.Split(s, "K")
	if len(split) == 2 {
		K, err := strconv.Atoi(split[0])
		if err != nil {
			return 0, err
		}
		res += K * KiloByte
		s = split[1]
	}
	split = strings.Split(s, "B")
	if len(split) == 2 && split[1] == "" {
		B, err := strconv.Atoi(split[0])
		if err != nil {
			return 0, err
		}
		res += B
	}
	return res, nil
}
