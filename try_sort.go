package main

import (
	"fmt"
	"math"
)

func BubbleSort(a []int) []int {
	length := len(a)
	var temp int
	for i:=0; i<length; i++ {
		for j:=0; j<length-1-i; j++ {
			if a[j] > a[j+1] {
				temp = a[j+1]
				a[j+1] = a[j]
				a[j] = temp
			}
		}
	}
	return a
}

func BubbleSortImprove(a []int) []int {
	length := len(a)
	var temp int
	var pos int

	i := length-1

	for i>0 {
		pos = 0
		for j:=0;j<i;j++ {
			if a[j] > a[j+1] {
				temp = a[j+1]
				a[j+1] = a[j]
				a[j] = temp
				pos = j
			}
		}
		i = pos
	}
	return a
}

func BubbleSortImprove2Dir(a []int) []int {
	length := len(a)
	var temp int

	low := 0
	high := length-1
	for low < high {
		for j:=low; j<high; j++ {
			if a[j] > a[j+1] {
				temp = a[j+1]
				a[j+1] = a[j]
				a[j] = temp
			}
		}
		high--

		for j:=high;j>low; j-- {
			if a[j]<a[j-1] {
				temp = a[j]
				a[j] = a[j-1]
				a[j-1] = temp
			}
		}
		low++
	}
	return a
}

func SelectSort(a []int) []int {
	length := len(a)
	var temp int

	for i:=0;i<length;i++ {
		minIndex := i
		for j:=i+1;j<length;j++ {
			if a[j]<a[minIndex] {
				minIndex=j
			}
		}

		temp = a[i]
		a[i] = a[minIndex]
		a[minIndex] = temp
	}
	return a
}

func InsertSort(a []int) []int {
	length := len(a)
	var key int

	for i:=1;i<length;i++ {
		key = a[i]
		j:=i-1
		for (j>=0)&&(a[j]>key) {
			a[j+1]=a[j]
			j--
		}
		a[j+1] = key
	}
	return a
}

func merge(a1 []int, a2 []int) []int {
	l1 := len(a1)
	l2 := len(a2)
	ret := make([]int, l1+l2)
	i:=0
	j:=0
	k:=0
	for i<l1 && j<l2 {
		if a1[i]<=a2[j] {
			ret[k]=a1[i]
			i++
		} else {
			ret[k]=a2[j]
			j++
		}
		k++
	}
	if i<l1 {
		for i<l1 {
			ret[k]=a1[i]
			i++
			k++
		}
	}
	if j<l2 {
		for j<l2 {
			ret[k]=a2[j]
			j++
			k++
		}
	}
	return ret
}

func MergeSort(a []int) []int {
	length := len(a)
	if length < 2 {
		return a
	}
	mid := length/2
	aLeft := a[:mid]
	aRight := a[mid:]
	return merge(MergeSort(aLeft), MergeSort(aRight))
}

func quickSortQueue(a []int) int {
	length := len(a)
	pivot := a[0]
	index := 0
	left := 1
	right := length-1
	for left<right {
		for left<right {
			if a[right] < pivot {
				a[index] = a[right]
				index = right
				right--
				break
			} else {
				right--
			}
		}
		for left<right {
			if a[left] > pivot {
				a[index] = a[left]
				index = left
				left++
				break
			} else {
				left++
			}
		}
	}
	a[index]=pivot
	return index
}

func QuickSort(a []int) []int {
	if len(a)<2 {
		return a
	}
	index := quickSortQueue(a)
	aL := QuickSort(a[:index])
	aR := QuickSort(a[index+1:])
	aL = append(aL, a[index])
	aL =append(aL, aR...)
	return aL
}

func heapify(a []int, index int, size int) {
	l := 2*index+1
	r := 2*index+2
	largest := index
	if l<size && a[l]>a[largest] {
		largest = l
	}
	if r<size && a[r]>a[largest] {
		largest = r
	}
	if largest != index {
		temp := a[index]
		a[index] = a[largest]
		a[largest] = temp
		heapify(a, largest, size)
	}
	//fmt.Println("h:", a)
}

func heapMake(a []int, size int) []int {
	for i:= size/2-1; i>=0; i-- { //从最后一个非叶子节点开始
		heapify(a, i, size)
	}
	//fmt.Println("m:", a)
	return a
}

func HeapSort(a []int) []int {
	length := len(a)
	//create heap
	heapMake(a, length)
	for i:=length-1; i>0; i-- {
		temp := a[0]
		a[0] = a[i]
		a[i] = temp
		heapify(a, 0, i)
	}
	return a
}

func CountSort(a []int) []int {
	length := len(a)
	ret := make([]int, length)
	max := 0
	record := make(map[int]int)
	for _, e := range a {
		if max < e {
			max = e
		}
		if cnt, ok := record[e]; ok {
			record[e] = cnt+1
		} else {
			record[e] = 1
		}
	}
	offset := make([]int, max+1)
	for i:=0; i<max+1; i++ {
		cnt, ok := record[i]
		if !ok {
			cnt = 0
		}
		if i>0 {
			offset[i] = offset[i-1]+ cnt
		} else {
			offset[i] = cnt
		}
	}
	//fmt.Println(len(offset), offset)

	for i:=length-1;i>=0;i-- {
		//fmt.Println(a[i], offset[a[i]])
		ret[offset[a[i]]-1] = a[i]
		offset[a[i]]--
	}
	return ret
}

type node struct {
	v int
	prev *node
	next *node
}

func (o *node) insertBefore(n *node) {
	n.prev = o.prev
	if n.prev != nil {
		n.prev.next = n
	}
	n.next = o
	o.prev = n
}

func (o *node) insertAfter(n *node) {
	n.next = o.next
	if n.next != nil {
		n.next.prev = n
	}
	n.prev = o
	o.next = n
}

func (head *node) insertAndSort(value int) {
	item := &node{value, nil, nil}
	if head.next == nil {
		head.insertAfter(item)
		return
	}
	iter := head.next
	cur := head
	for iter != nil {
		if iter.v > value {
			iter.insertBefore(item)
			//fmt.Println(item.prev, item, item.next)
			return
		}
		cur = iter
		iter = iter.next
	}
	cur.insertAfter(item)
}

func (head *node) through() {
	iter := head.next
	for iter != nil {
		//fmt.Printf("%d->", iter.v)
		iter = iter.next
	}
	//fmt.Printf("\n")
}

func pos(v int, min int, max int, seg int) int {
	temp := seg*(v-min)/(max-min)
	if temp == seg {
		temp--
	}
	return temp
}

func BucketSort(a []int) []int {
	length := len(a)

	max := a[0]
	min := a[length-1]
	for _, i := range a {
		if max < i {
			max = i
		}
		if min > i {
			min = i
		}
	}

	bN := int(math.Sqrt(float64(length)))
	//fmt.Println("need bucket:", bN)

	bucket := make([]*node, bN)
	for i:=0; i<bN; i++ {
		bucket[i]=&node{-1, nil, nil}
	}

	for i:=0;i<length;i++ {
		index := pos(a[i], min, max, bN)
		//fmt.Println(index, a[i])
		bucket[index].insertAndSort(a[i])
		bucket[index].through()
	}

	//merge
	j:=0
	ret := make([]int, length)
	for i:=0; i< bN; i++ {
		start := bucket[i].next
		for start != nil {
			ret[j] = start.v
			start = start.next
			j++
		}
	}
	return ret
}


func main() {
	arr := []int{3,44,38,5,47,15,36,26,15,27,2,46,28,4,19,50,48}
	fmt.Println(BubbleSort(arr))
	arr = []int{3,44,38,5,47,15,36,26,15,27,2,46,28,4,19,50,48}
	fmt.Println(BubbleSortImprove(arr))
	arr = []int{3,44,38,5,47,15,36,26,15,27,2,46,28,4,19,50,48}
	fmt.Println(BubbleSortImprove2Dir(arr))
	arr = []int{3,44,38,5,47,15,36,26,15,27,2,46,28,4,19,50,48}
	fmt.Println(SelectSort(arr))
	arr = []int{3,44,38,5,47,15,36,26,15,27,2,46,28,4,19,50,48}
	fmt.Println(InsertSort(arr))

	arr = []int{3,44,38,5,47,15,36,26,15,27,2,46,28,4,19,50,48}
	fmt.Println(MergeSort(arr))

	arr = []int{3,44,38,5,47,15,36,26,15,27,2,46,28,4,19,50,48}
	fmt.Println(QuickSort(arr))

	arr = []int{3,44,38,5,47,15,36,26,15,27,2,46,28,4,19,50,48}
	fmt.Println(HeapSort(arr))

	arr = []int{3,44,38,5,47,15,36,26,15,27,2,46,28,4,19,50,48}
	fmt.Println(CountSort(arr))

	arr = []int{3,44,38,5,47,15,36,26,15,27,2,46,28,4,19,50,48}
	fmt.Println(BucketSort(arr))
}