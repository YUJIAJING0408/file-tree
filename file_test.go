package fileTree

import (
	"fmt"
	"math/rand"
	"testing"
)

/*
@Date:
@Auth: YUJIAJING
@Desp:
*/

func TestFileHeap(t *testing.T) {
	var tk = NewTopK(10)
	for i := 0; i < 1000; i++ {
		tk.Push(File{
			Size: rand.Int63n(1000),
		})
	}
	fmt.Println(tk.TopK())       // [{  992 0 0 } {  994 0 0 } {  995 0 0 } {  997 0 0 } {  995 0 0 } {  998 0 0 } {  995 0 0 } {  999 0 0 } {  999 0 0 } {  999 0 0 }]
	fmt.Println(tk.TopKSorted()) // [{  999 0 0 } {  999 0 0 } {  999 0 0 } {  998 0 0 } {  997 0 0 } {  995 0 0 } {  995 0 0 } {  995 0 0 } {  994 0 0 } {  992 0 0 }]
}
