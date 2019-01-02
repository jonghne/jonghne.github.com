package model

import (
	"fmt"
	"errors"
)

var (
	ErrNoOpFunc = errors.New("未设置操作函数")
	ErrNoData = errors.New("未提供数据")
)

type MapWork struct {
	thdNum int
	op func(interface{}) interface{}
	//result chan interface{}
	done  chan WorkResp
	inputs chan WorkReq
	stop chan struct{}
	abort chan struct{}
}

type WorkReq struct {
	index int
	req interface{}
}

type WorkResp struct {
	index int
	result interface{}
}

func InitWorkMap(tNum int) *MapWork {
	work := MapWork{thdNum: tNum}

	work.done = make(chan WorkResp, tNum)
	work.inputs = make(chan WorkReq, tNum)
	//work.result = make(chan interface{}, tNum)
	work.stop = make(chan struct{}, tNum)
	work.abort = make(chan struct{})

	//start threads
	for i:=0; i<tNum; i++ {
		go work.workThread()
	}

	return &work
}

func (work *MapWork) SetOperation(f func(interface{}) interface{}) {
	work.op = f
}

func (work *MapWork) workThread() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}
	}()

	for {
		select {
		case in := <-work.inputs:
			ret := work.op(in.req)
			resp := WorkResp{in.index, ret}
			//fmt.Println("++++", in.index, in.req)
			work.done <- resp

		case <- work.stop:
			return
		}
	}
}

func (work *MapWork) StartWorks(items []interface{}) (chan interface{}, error) {
	var (
		in, out = 0, 0
		checked = make([]bool, len(items))
		inputs  = work.inputs
		cache = make([]interface{}, len(items))
	)

	length := len(items)

	if length < 1 {
		return nil, ErrNoData
	}

	if work.op == nil {
		return nil, ErrNoOpFunc
	}
	result := make(chan interface{}, length)

	go func() {
		item := WorkReq{0, items[0]}
		for {
			select {
			case inputs <- item:
				if in++; in == length {
					//no more work to do
					inputs = nil
				} else {
					item.index = in
					item.req = items[in]
				}
			case resp := <-work.done:
				//fmt.Println(out, resp.index)
				cache[resp.index] = resp.result
				//fmt.Println(resp.index, out)
				for checked[resp.index]=true; checked[out]; out++ {
					//fmt.Println("@@", out, resp)
					result <- cache[out]
					if out == length-1 {
						//finish, end
						return
					}
				}
				//fmt.Println("$$$", out)

			case <-work.abort:
				return
			}
		}
	}()

	return result, nil
}

func (work *MapWork) AbortWorks() {
	work.abort <- struct{}{}
}