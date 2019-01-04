package time_wheel

import (
	"testing"
	"fmt"
)

func TestAdd(t *testing.T) {
	head := &TmObjList{}

	for i:=0; i<10; i++ {
		it := &TmObj{mode:i}
		head.add(it)
	}
	for i:=head.next;i!=nil;i=i.next {
		fmt.Println(i.tmObj)
	}

	for i:=head.next;i!=nil;i=i.next {
		if i.tmObj.mode %2 == 0 {
			head.del(i)
		}
	}

	for i:=head.next;i!=nil;i=i.next {
		fmt.Println(i.tmObj)
	}
}