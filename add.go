package main

import (
	"fmt"
	"os"
	"strconv"
	//"container/list"
	//"reflect"
)

type ListNode struct {
	Val int
	Next *ListNode
}

func add(n1 *ListNode, n2 *ListNode) *ListNode {
	var ret *ListNode = nil
	var tmp *ListNode = nil
	var v1 int
	var v2 int
	var carry int = 0

	//for n:=n1;n!=nil;n=n.Next {
	//	fmt.Println(n.Val)
	//}
	//for n:=n2;n!=nil;n=n.Next {
	//	fmt.Println(n.Val)
	//}
	for {
		if n1 == nil && n2 == nil {
			if carry > 0 { //有进位
				l := new(ListNode)
				l.Val = carry
				l.Next = nil
				tmp.Next = l
			}
			return ret
		} else {
			if n1 != nil {//get new node
				v1 = n1.Val
				n1 = n1.Next
			} else {//nothing more, as 0
				v1 = 0
			}
			if n2 != nil {//get new node
				v2 = n2.Val
				n2 = n2.Next
			} else {//nothing more, as 0
				v2 = 0
			}
			//fmt.Println(v1, " ", v2)
			sum := v1 + v2 +carry
			l := new(ListNode)
			carry = sum/10
			l.Val = sum%10
			l.Next = nil
			if ret == nil {//first time, create head node
				ret = l
				tmp = l
			} else {//add new node
				tmp.Next = l
				tmp = l
			}
		}
	}
}

func main() {
	var br int = 0
	var n1 *ListNode = nil
	var h1 *ListNode = nil
	var n2 *ListNode = nil
	var h2 *ListNode = nil
	for _, v := range os.Args[1:] {
		br++
		//fmt.Println(v)
		if v == "next" {
			break
		}
		n := new(ListNode)
		r, _ := strconv.Atoi(v)
		//fmt.Println(r)
		n.Val = r
		n.Next = nil
		if h1 == nil {
			h1 = n
			n1 = n
		} else {
			n1.Next = n
			n1 = n
		}
	}
	br++
	//fmt.Println("break at ", br)
	for _, v := range os.Args[br:] {
		//fmt.Println(v)
		n := new(ListNode)
		r, _ := strconv.Atoi(v)
		//fmt.Println(r)
		n.Val = r
		n.Next = nil
		if h2 == nil {
			h2 = n
			n2 = n
		} else {
			n2.Next = n
			n2 = n
		}
	}
	ret := add(h1, h2)
	for n:=ret;n!=nil;n=n.Next {
		fmt.Println(n.Val)
	}
}