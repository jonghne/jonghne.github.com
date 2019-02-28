package main

import (
	"fmt"
	"math/rand"
	"time"
)

type stateHandler interface {
	doAction(ctx *context, event interface{})
	info() string
}

type stateA struct {
	tm string
}

type stateB struct {
	tm string
}

type context struct {
	st stateHandler
}

func (ctx *context) setState(st stateHandler) {
	ctx.st = st
}

func (ctx *context) getState() string {
	return ctx.st.info()
}

func (s *stateA) doAction(ctx *context, event interface{}) {
	fmt.Println("AAAAAAAAA", event)
	if v, ok := event.(uint32); ok {
		if v>5 {
			fmt.Println("change B")
			newState := &stateB{"stateB coming"}
			ctx.setState(newState)
		} else {
			fmt.Println("still in stateA")
		}
	}
}

func (s *stateA) info() string {
	return fmt.Sprint(s.tm, " working")
}

func (s *stateB) doAction(ctx *context, event interface{}) {
	fmt.Println("BBBBBBBBB", event)
	if v, ok := event.(uint32); ok {
		if v<6 {
			fmt.Println("change A")
			newState := &stateA{"stateA coming"}
			ctx.setState(newState)
		} else {
			fmt.Println("still in stateB")
		}
	}
}

func (s *stateB) info() string {
	return fmt.Sprint(s.tm, " working")
}


func main() {
	//init stateA
	s := &stateA{"stateA init"}
	ctx := &context{s}

	for {
		rand.Seed(time.Now().UnixNano())
		ev := rand.Uint32()%10
		ctx.st.doAction(ctx, ev)
		fmt.Println(ctx.getState())
		time.Sleep(time.Second)
	}
}