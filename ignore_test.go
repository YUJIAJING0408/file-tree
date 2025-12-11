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

func TestReadIgnore(t *testing.T) {
	rules, _ := ReadIgnore(".treeignore")
	fmt.Printf("%#v\n", rules)
	//	println(rules.Ignore("a.zip", 10000))
	//	println(rules.Ignore("a.exe", 10*KiloByte))
	//	println(rules.Ignore("a.exe", 10*MegaByte))
}
