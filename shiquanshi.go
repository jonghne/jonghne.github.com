package main

import (
	//"bytes"
	//"errors"
	"fmt"
	//"os"
	//"reflect"
	//"runtime"
	//"strconv"
	//"time"
)

var count int = 0

func validate(a []int, pos int, ins int) bool {
	length := pos
	if length > 0 {
		for i:=0; i<length; i++ {
			if a[i] == ins { //exist already
				return false
			}
		}
		if a[0] > 1 {
			return false
		}
		if length > 1 {
			if a[0] == 1 && a[1] > 2 {
				return false
			}
		}

		if length > 2 {
			if a[2] > 3 {
				return false
			}
		}
		if length > 3 {
			if a[2] == 3 && a[3] > 1{
				return false
			}
		}
		if length > 4 {
			if a[4] > 2 {
				return false
			}
		}
		if length > 5 {
			if a[4] == 2 && a[5] > 3{
				return false
			}
		}
		if length > 6 {
			if a[6] > 5 {
				return false
			}
		}
		if length > 8 {
			if a[8] > 5 {
				return false
			}
		}
	}
	return true

}

func scan(a []int, pos int) {
	if pos == 10 {
		fmt.Println(a)
		count++
		return
	}
	b := []int{}
	b = append(b, a...)
	//fmt.Println("save ", pos, b)
	curPos := pos
	for i:=0; i<10; i++ {
		if validate(a, curPos, i) {
			a[curPos] = i
			curPos++
			scan(a, curPos)
		}
		//restore original array
		a = append([]int{}, b...)
		curPos = pos
	}
}

func main() {
	tmp := make([]int, 10)
	fmt.Println(len(tmp), cap(tmp))
	array := make([]int, 10)
	scan(array, 0)
	fmt.Println("total: ", count)
}