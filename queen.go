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
	"math"
)

var count int = 0

func showQueen(qPos int, rN int) {
	str := ""
	for i:=0;i<qPos;i++ {
		//fmt.Printf("*")
		str += "*"
	}
	//fmt.Println("#")
	str += "#"
	for i:=0;i<(rN-qPos-1);i++ {
		//fmt.Printf("*")
		str += "*"
	}
	fmt.Println(str)
}

func showQueenAry(qAry []int, rN int) {
	fmt.Println("==========")
	for _, v := range qAry {
		showQueen(v, rN)
	}
	fmt.Println("----------")
}

func validate(houses []int, pos int, room int) bool {
	//comparing from first one to current one
	for i:=0; i<pos; i++ {
		if houses[i] == room {
			//same column
			return false
		}

		if math.Abs(float64(pos-i)) == math.Abs(float64(houses[i] - room)) {
			//diagnonal
			return false
		}
	}
	return true
}

func arrangeQueens(houses []int, pos int, num int) {
	if pos == num {
		//fmt.Println(houses)
		showQueenAry(houses, num)
		count++
		return
	}
	temp := append([]int{}, houses...)
	curPos := pos
	//fmt.Println(pos, temp)
	for i:=0;i<num;i++ {
		if validate(houses, curPos, i) {
			//if curPos >= num {
			//	fmt.Println(curPos, num, len(houses))
			//	break
			//}
			houses[curPos] = i
			curPos++
			arrangeQueens(houses, curPos, num)
		}
		//restore original env for next search
		houses = append([]int{}, temp...)
		curPos = pos
	}

}

func main() {
	num := 8
	houses := make([]int, num)
	//showQueenAry(houses, 8)

	arrangeQueens(houses, 0, num)
	fmt.Println(count)
}