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

type Camera struct { // 定义一个单反相机
	Name string
	// @Inject // <-－你在这里说明了L是需要注入的，通常用注释来说明注入
	L Lens
}

type Lens struct {
	LensType int // 0: 普通镜头，1: 广角镜头，2: 长焦镜头
}

func (c Camera) Capture() { // 使用不同的镜头拍照
	m := make(map[int]string, 3)
	m[0], m[1], m[2] = "normal lens", "wide-angle lens", "telephoto lens"

	//* fmt.Println("capture with", m[c.L.LensType]) // 这是你在程序中的代码，框架看到你使用了c.L
	// 默默的给你做了如下替换：
	fmt.Println("capture with", m[getLens(c.L).LensType])
}

// 由框架在编译时插入的方法
func getLens(L Lens) Lens {
	value := inj.Get(reflect.TypeOf(L))
	var l Lens
	if value.Kind() != reflect.Invalid {
		l = value.Interface().(Lens)
	} else {
		l = Lens{0}
	}
	return l
}

var inj *Injector // 持有一个全局的注入器，复用之前实现的注入器。

func main() {
	inj = New()

	c := Camera{"ZddCamera", Lens{0}}

	wLens := Lens{1}
	inj.SetMap(wLens) // 你需要用某种方式进行注入

	c.Capture()
}