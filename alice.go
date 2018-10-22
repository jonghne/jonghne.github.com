package main

import (
	"container/list"
	"fmt"
	"os"
	"strconv"
	//"reflect"
)

func partition(buck []int, length int) int {
	pivot := buck[0]
	left := 0
	right := length - 1
	index := left
	for left < right {
		for ; left < right; right-- {
			if pivot > buck[right] {
				buck[index] = buck[right]
				left++
				index = right

				break
			}
		}

		for ; left < right; left++ {
			if pivot < buck[left] {
				buck[index] = buck[left]
				right--
				index = left
				break
			}
		}

	}
	buck[index] = pivot
	return index
}

func qsort(buck []int, length int) {
	//fmt.Println(len(buck), buck)
	if length < 2 {
		//fmt.Println("end")
		return
	}
	index := partition(buck, length)
	//fmt.Println("index=",index)
	//divide 2 segments
	qsort(buck[:index], index)
	//fmt.Println("===========")
	qsort(buck[index+1:], length-index-1)
}

func search(buck []int, w int) int {
	var cnt int = 0
	var length int = len(buck)
	var i int
	var invalid int = 0
	for invalid != length {
		var ary []int
		j := 0
		//fmt.Println("*********")
		for i = 0; i < length; {
			//fmt.Println("i=",i)
			if buck[i] != 0x7fffffff {
				for k := i + 1; k < length; k++ {
					//fmt.Println("k=",k)
					if buck[k] != 0x7fffffff {
						if buck[k] == (buck[i] + 1) {
							//fmt.Println("find ", i, k)
							ary = append(ary, buck[i])
							j++
							buck[i] = 0x7fffffff
							invalid++
							i = k
							break
						} else if buck[k] == buck[i] {
							i++
							continue
						} else {
							//fmt.Println("nothing error ", ary, buck, i, buck[i], k, buck[k])
							return 0
						}
					}
				}
			} else {
				i++
			}
			if j == (w - 1) {
				cnt++
				ary = append(ary, buck[i])
				buck[i] = 0x7fffffff
				invalid++
				fmt.Println(ary)
				break
			}
		}
		if i >= length {
			//fmt.Println("scan end error ", ary)
			return 0
		}
	}
	return cnt
}

func scanSubSquence(l *list.List, w int) bool {
	j := 0
	var ary []int
	var e *list.Element
	var e1 *list.Element
	for e = l.Front(); e != nil; {//from first node, fix head every time
		for e1 = e.Next(); e1 != nil; e1 = e1.Next() {//scan nodes after head
			if e1.Value == (e.Value.(int) + 1) {//valid successor
				ary = append(ary, e.Value.(int)) //save head
				//delete head
				l.Remove(e)
				e = e1
				j++
				break
			} else if e1.Value == e.Value {//successor equal to head, ignore, goto next
				e = e.Next()
			} else {//invalid successor, exit cycle
				return false
			}
		}
		if e1 == nil {
			return false
		}
		if j == (w - 1) {
			ary = append(ary, e.Value.(int))
			l.Remove(e)
			fmt.Println(ary)
			break
		}
	}
	return true
}

func scanSubSequenceOpt(l *list.List, w int) bool {
	var ary []int
	var e *list.Element
	var e1 *list.Element
	e = l.Front()//from first node, fix head every time

	for j := 0; j < (w-1); j++ {//find w members
		for e1 = e.Next(); e1 != nil; e1 = e1.Next() {//scan nodes after head
			if e1.Value == (e.Value.(int) + 1) {//valid successor
				ary = append(ary, e.Value.(int)) //save head
				//delete head
				l.Remove(e)
				e = e1
				break
			} else if e1.Value == e.Value {//successor equal to head, ignore, goto next
				e = e.Next()
			} else {//invalid successor, exit cycle
				return false
			}
		}
		if e1 == nil {
			return false
		}
	}
	ary = append(ary, e.Value.(int))
	l.Remove(e)
	fmt.Println(ary)
	return true
}

func search1(buck []int, w int) int {

	var cnt int = 0

	if w == 1 {
		return len(buck)
	}
	l := list.New() //创建一个新的list
	for i := 0; i < len(buck); i++ {
		l.PushBack(buck[i])
	}

	for l.Len() > 0 {//scan whole list
		//j := 0
		//var ary []int
		//var e *list.Element
		//var e1 *list.Element
		//for e = l.Front(); e != nil; {//from first node, fix head every time
		//	for e1 = e.Next(); e1 != nil; e1 = e1.Next() {//scan nodes after head
		//		if e1.Value == (e.Value.(int) + 1) {//valid successor
		//			ary = append(ary, e.Value.(int)) //save head
		//			//delete head
		//			l.Remove(e)
		//			e = e1
		//			j++
		//			break
		//		} else if e1.Value == e.Value {//successor equal to head, ignore
		//			e = e.Next()
		//		} else {//invalid successor, exit cycle
		//			return 0
		//		}
		//	}
		//	if e1 == nil {
		//		return 0
		//	}
		//	if j == (w - 1) {
		//		cnt++
		//		ary = append(ary, e.Value.(int))
		//		l.Remove(e)
		//		fmt.Println(ary)
		//		break
		//	}
		//}
		if err := scanSubSequenceOpt(l, w); err {
			cnt++
		} else {
			return 0
		}
	}
	return cnt
}

func main() {

	var sc []int

	w, _ := strconv.Atoi(os.Args[1])
	for _, v := range os.Args[2:] {
		r, _ := strconv.Atoi(v)
		//fmt.Println(r)
		sc = append(sc, r)
	}
	//fmt.Println(sc, reflect.TypeOf(sc))
	qsort(sc, len(sc))
	//fmt.Println(sc)
	if (len(sc) % w) > 0 {
		fmt.Println("not match, w error")
		return
	}
	//count := search(sc, w)
	count := search1(sc, w)
	if count > 0 {
		fmt.Println("find split = ", count)
	} else {
		fmt.Println("error hand")
	}
}
