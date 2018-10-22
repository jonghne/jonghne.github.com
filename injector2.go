package main

import (
	"fmt"
	"reflect"
)

type Injector struct {
	mappers map[reflect.Type]reflect.Value // 根据类型map实际的值
}

func (inj *Injector) SetMap(value interface{}) {
	inj.mappers[reflect.TypeOf(value)] = reflect.ValueOf(value)
}

func (inj *Injector) Get(t reflect.Type) reflect.Value {
	return inj.mappers[t]
}

func (inj *Injector) Invoke(i interface{}) []reflect.Value {
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Func {
		panic("Should invoke a function!")
	}
	inValues := make([]reflect.Value, t.NumIn())
	for k := 0; k < t.NumIn(); k++ {
		inValues[k] = inj.Get(t.In(k))
	}
	ret := reflect.ValueOf(i).Call(inValues)
	return ret
}

func New() *Injector {
	return &Injector{make(map[reflect.Type]reflect.Value)}
}

func Host(name string, i interface{}) { // 让注入的方法不受限制
	inj.Invoke(i) // 利用注入器中的环境调用f
}

func Dependency(a int, b string) {
	fmt.Println("Dependency: ", a, b)
}

var inj *Injector

func main() {
	inj = New()
	inj.SetMap(3030)
	inj.SetMap("zdd")

	d := Dependency
	Host("zddhub", d)

	inj.SetMap(8080)
	inj.SetMap("zddhub@gmail.com")
	inj.SetMap(10.24)

	Host("email", func(email string, money float64, number int) {
		fmt.Println(email, money, number)
	})

	Host("Get", func() {
		fmt.Println("Hello world!")
	})
}