package time_wheel

import (
	"testing"
	"fmt"
	"time"
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

		if i.tmObj.mode == 9 {
			head.del(i)
		}
	}

	for i:=head.next;i!=nil;i=i.next {
		fmt.Println(i.tmObj)
	}
}

type T1 struct {
	cnt int
}

func (t *T1) Run() {
	t.cnt++
	fmt.Println("this is T1", t.cnt, time.Now().Format("2006-01-02 15:04:05"))
}

type T2 struct {

}

func (t *T2) Run() {
	fmt.Println("this is T2")
}

func TestInitTimer(t *testing.T) {
	tmDev := InitTimer(10)

	t1 := &T1{}
	t2 := &T2{}

	tmDev.AddTimer(t1, 1000, Continual)
	tmDev.AddTimer(t2, 500, OneShot)

	time.Sleep(20*time.Second)
}