package time_wheel


import (
	"time"
)

const (
	OneShot = iota
	Continual
)

var (
	slotNum = 10
)

type TmTask interface {
	Run()
}

type TmObj struct {
	mode int
	task TmTask
	leftRound int
	cntRound int
	slotNo int
}

type TmObjList struct {
	prev *TmObjList
	next *TmObjList
	tmObj *TmObj
}

type TmDevice struct {
	slots []*TmObjList
	interval int
	pivot int
	abortC chan struct{}
}

func (l *TmObjList) add(obj *TmObj) {
	head := l
	item := &TmObjList{tmObj: obj}
	if head.next == nil {
		head.next = item
		item.prev = head

	} else {
		var t_prev *TmObjList
		for t:=head.next; t!=nil; t=t.next {
			t_prev = t
		}
		t_prev.next = item
		item.prev = t_prev
	}
}

func (l *TmObjList) del(item *TmObjList) {
	prev := item.prev
	next := item.next
	prev.next = next
	if next != nil {
		next.prev = prev
	}
}

type TimerWheel interface {
	AddTimer(task TmTask, waitTime int, mode int)
	Stop()
}

func (td *TmDevice) AddTimer(task TmTask, waitTime int, mode int) {
	round := waitTime / (td.interval * slotNum)
	roundLeft := waitTime % (td.interval * slotNum)
	slotPos := roundLeft + td.pivot

	if slotPos >= slotNum {
		slotPos -= slotNum
	}

	obj := TmObj{
		mode: mode,
		task: task,
		leftRound: round,
		slotNo: slotPos,
	}

	td.addTask(&obj)
}

func (td *TmDevice) addTask(obj *TmObj) {
	l := td.slots[obj.slotNo]
	l.add(obj)
}

func (td *TmDevice) scanTask() {
	curTaskHead := td.slots[td.pivot]
	td.pivot++
	if td.pivot == slotNum {
		td.pivot = 0
	}

	for item:=curTaskHead.next; item!=nil; item=item.next {
		task := item.tmObj
		if task.leftRound > task.cntRound {
			task.cntRound++
		} else {
			if task.mode == OneShot {
				//only run once, delete
				curTaskHead.del(item)
			} else {
				//reset count
				task.cntRound = 0
			}
			go task.task.Run()
		}
	}
}

func (td *TmDevice) Stop() {
	td.abortC <- struct{}{}
}

func (td *TmDevice) wheel() {
	ticker := time.NewTicker(time.Duration(td.interval)*time.Millisecond)
	for {
		select {
		case <-ticker.C:
			td.scanTask()
		case <-td.abortC:
			return
		}
	}
}

func InitTimer(precise int)  TimerWheel {
	tmWheel := &TmDevice {
		slots: make([]*TmObjList, slotNum),
		interval: precise,
		pivot : 0,
		abortC: make(chan struct{}),
	}

	for i:=0; i<slotNum; i++ {
		//init head
		tmWheel.slots[i] = &TmObjList{}
	}
	//start routine
	go tmWheel.wheel()

	return tmWheel
}