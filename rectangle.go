package main

import (
	"fmt"
	"os"
	"strconv"
	//"container/list"
	//"reflect"
)

type rectangle struct {
	x1 int
	y1 int
	x2 int
	y2 int
}

func main() {
	//var r1 []int
	//var r2 []int
	//for _, v := range os.Args[1:5] {
	//	r1=append(r1, v)
	//}
	//
	//for _, v := range os.Args[5:9] {
	//	r2=append(r2, v)
	//}
	var r1 rectangle
	var r2 rectangle
	r, _ := strconv.Atoi(os.Args[1])
	r1.x1 = r
	r, _ = strconv.Atoi(os.Args[2])
	r1.y1 = r
	r, _ = strconv.Atoi(os.Args[3])
	r1.x2 = r
	r, _ = strconv.Atoi(os.Args[4])
	r1.y2 = r
	r, _ = strconv.Atoi(os.Args[5])
	r2.x1 = r
	r, _ = strconv.Atoi(os.Args[6])
	r2.y1 = r
	r, _ = strconv.Atoi(os.Args[7])
	r2.x2 = r
	r, _ = strconv.Atoi(os.Args[8])
	r2.y2 = r

	if (r1.x2<=r2.x1) || (r1.x1 >= r2.x2) || (r1.y1 >= r2.y2) || (r1.y2<=r2.y1) {
		fmt.Println("no cover")
	} else {
		fmt.Println("cover")
	}
}