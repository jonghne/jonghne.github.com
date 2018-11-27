package main

import (
	"container/list"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"strconv"
	"strings"
	"time"
)

func consumer(pack <-chan interface{}) {

	tick := time.NewTicker(1 * time.Second)
	for {
		select {
		case msg := <-pack:
			fmt.Println(msg)
		case <-tick.C:
			fmt.Println("timeout")
		}
	}
}

func producer(pack chan<- interface{}) {
	i := 0
	for {
		time.Sleep(1100 * time.Millisecond)
		if j := i % 2; j == 0 {
			pack <- i
		} else {
			pack <- "odd"
		}
		i++

	}
}

/*mainVersionNo:subVersionNo:devSerialNo*/
func GetVersion(m int, s int, d int) int {
	return (m << 16) | (s << 8) | d
}

func CurVersion() int {
	//TODO:return current version
	//return GetVersion(config.mainVersionNo, config.subVersionNo, config.devSerialNo)
	return 0 //temp
}

func VersionString(m int, s int, d int) string {
	return fmt.Sprintf("%d.%d.%d", m, s, d)
}

func ParseVersion(v string) []int {
	a := strings.Split(v, ".")
	if len(a) != 3 {
		return nil
	}

	r := make([]int, 3)
	for i := 0; i < 3; i++ {
		j, e := strconv.Atoi(a[i])
		if e != nil {
			return nil
		}
		r[i] = j
	}
	return r
}

type point struct {
	x int
	y int
}

type record struct {
	prev     point
	self     point
	distance int
}

func (r record) String() string {
	return fmt.Sprintf("%d %d %d", r.self.x, r.self.y, r.distance)
}

func makeMaze(m int, n int, mode string) [][]string {
	room := make([][]string, m)
	ps := strings.Split(mode, " ")
	for i := 0; i < m; i++ {
		slack := make([]string, n)
		for j := 0; j < n; j++ {
			slack[j] = ps[i*n+j]
		}
		room[i] = slack
	}
	fmt.Println(room)
	return room
}

const (
	INF = 10000
)

func printPath(x int, y int, rd [][]record) {
	x1 := x
	y1 := y

	//fmt.Println(rd)
	for {
		if x1 >= 0 && y1 >= 0 {
			fmt.Println(x1, y1, rd[x1][y1].distance)
			xt := rd[x1][y1].prev.x
			yt := rd[x1][y1].prev.y
			x1, y1 = xt, yt
		} else {
			break
		}
	}
}

func bfs(m int, n int, room [][]string) {
	sx := []int{1, 0, -1, 0}
	sy := []int{0, 1, 0, -1}
	var start point
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if room[i][j] == "s" {
				start.x = i
				start.y = j
				break
			}
		}
	}
	rd := make([][]record, m)
	for i := 0; i < m; i++ {
		slack := make([]record, n)
		for j := 0; j < n; j++ {
			slack[j].distance = INF
		}
		rd[i] = slack
	}
	rd[start.x][start.y].distance = 0
	rd[start.x][start.y].prev.x = -1
	rd[start.x][start.y].prev.y = -1
	rd[start.x][start.y].self.x = start.x
	rd[start.x][start.y].self.y = start.y
	fmt.Println(start.x, start.y)
	l := list.New()
	l.PushBack(rd[start.x][start.y])
	fmt.Println(rd[start.x][start.y])
	for l.Len() > 0 {
		p := l.Front()
		l.Remove(p)

		fmt.Println("from", p.Value.(record))
		for shift := 0; shift < 4; shift++ {
			//shift and search
			x := p.Value.(record).self.x + sx[shift]
			y := p.Value.(record).self.y + sy[shift]
			fmt.Println("shift to ", x, y)
			if x >= 0 && x < m && y >= 0 && y < n && room[x][y] != "#" && rd[x][y].distance == INF { //cannot turn back
				rd[x][y].distance = p.Value.(record).distance + 1
				rd[x][y].prev.x = p.Value.(record).self.x
				rd[x][y].prev.y = p.Value.(record).self.y
				rd[x][y].self.x = x
				rd[x][y].self.y = y
				if room[x][y] == "z" {
					fmt.Println("find", rd[x][y])
					printPath(x, y, rd)
					return
				}
				l.PushBack(rd[x][y])
			}
		}
	}

}

type path struct {
	line []int
	dist int
}

func updatePath(pt []path, st int, end int) {
	pt[end].line = []int{}
	pt[end].line = append(pt[end].line, pt[st].line...)
	pt[end].line = append(pt[end].line, st)
}

func Dijkstra(graph [][]int, pt []path, mid int, size int) {
	//init, begin with start point
	for i:=0; i<size; i++ {
		pt[i].dist = graph[mid][i]
		if pt[i].dist < INF {
			pt[i].line = make([]int, 1)
			pt[i].line[0] = mid
		}
		//fmt.Println(p.dist)
	}

	for i := 0; i < size; i++ {
		//fmt.Println(pt[i].dist)
		if pt[i].dist == INF {
			continue
		}
		fmt.Println("---------", i)
		for j := 0; j < size; j++ {
			fmt.Println(i, "->", j, "=", graph[i][j])
			if graph[i][j] != INF {
				//j node distance through i node
				d := graph[i][j] + pt[i].dist
				if pt[j].dist > d {
					//update distance
					fmt.Println(j, "by ", i, "dist=", d)
					pt[j].dist = d
					//update path
					updatePath(pt, i, j)
				}
			}
		}
	}
}

func makeGraph() ([][]int, int) {
	graph := [][]int {
		{0, 10, INF, 30, 100},
		{INF, 0, 50, INF, INF},
		{INF, INF, 0, INF, 10},
		{INF, INF, 20, 0, 60},
		{INF, INF, INF, INF, 0},
	}
	return graph, 5
}

type persist struct {
	db map[interface{}]interface{}
}

func newPersist() *persist {
	return &persist{
		db: make(map[interface{}]interface{}),
	}
}

func (disk *persist) writePersist(k, v interface{}) bool {
	disk.db[k] = v
	return true
}

func (disk *persist) readPersist(k interface{}) (interface{}, bool) {
	r, ok := disk.db[k]
	fmt.Println("in read disk", r, ok)
	if ok {
		return r, true
	}
	return nil, false
}

type cache struct {
	cdb map[interface{}]interface{}
}

func newCache() *cache {
	return &cache{
		cdb: make(map[interface{}]interface{}),
	}
}

func (ca *cache) writeCache(k, v interface{}) bool {
	ca.cdb[k] = v
	return true
}

func (ca *cache) readCache(k interface{}) (interface{}, bool) {
	r, ok := ca.cdb[k]
	fmt.Println("in read cache", r, ok)
	if ok {
		return r, true
	}
	return nil, false
}

type block struct {
	key   interface{}
	value interface{}
}

type stateDb struct {
	cache     *cache
	disk      *persist
	pipe      *list.List
	msg       chan struct{}
	threshold int
}

func newStateDb(lt int) *stateDb {
	return &stateDb{
		cache:     newCache(),
		disk:      newPersist(),
		pipe:      list.New(),
		msg:       make(chan struct{}),
		threshold: lt,
	}
}

func (statedb *stateDb) start() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				fmt.Println("timeout")
			case <-statedb.msg:
				pack := statedb.pipe.Front()
				b := pack.Value.(block)
				fmt.Println("wd", b)
				statedb.disk.writePersist(b.key, b.value)
				statedb.pipe.Remove(pack)
			}
		}
	}()
}

func (statedb *stateDb) stop() {

}

func (statedb *stateDb) write(k, v interface{}) bool {
	statedb.cache.writeCache(k, v)
	b := block{
		key:   k,
		value: v,
	}
	statedb.pipe.PushBack(b)
	go func() {
		//notify persist write goroutine
		statedb.msg <- struct{}{}
	}()
	if len(statedb.cache.cdb) > statedb.threshold {
		//clear
		for k := range statedb.cache.cdb {
			//for k, v := range statedb.cache.cdb {
			//	fmt.Println(v)
			delete(statedb.cache.cdb, k)
		}
	}
	return true
}

func (statedb *stateDb) read(k interface{}) (interface{}, bool) {
	if v, ok := statedb.cache.readCache(k); ok {
		fmt.Println("rd cache", v)
		return v, true
	} else if v, ok := statedb.disk.readPersist(k); ok {
		fmt.Println("rd disk", v)
		return v, true
	}
	return nil, false
}

func loop() {
	for {
		fmt.Println("loop...")
		time.Sleep(1 * time.Second)
	}
}

func startloop() {
	go loop()
}

func testLevelDb() {
	db, _ := leveldb.OpenFile("/home/qydev/softwares", nil)
	defer db.Close()

	db.Put([]byte("key"), []byte("value"), nil)
	if data, err := db.Get([]byte("key"), nil); err == nil {
		fmt.Println("saved:", string(data))
	}

}

func main() {

	g, sz := makeGraph()
	point := make([]path, sz)

	Dijkstra(g, point, 0, sz)
	for _, p := range point {
		fmt.Println("go ", p.line, "distance=", p.dist)
	}

	var flag bool
	fmt.Println(flag)
	testLevelDb()

	statedb := newStateDb(4)
	statedb.write(1, 2)
	statedb.write("hello", "world")
	statedb.start()
	statedb.write(5, "qiu")
	statedb.write("hi", "world")
	statedb.write(18, 9)
	statedb.write("ji", "earth")
	time.Sleep(1 * time.Second)
	if v, e := statedb.read(1); e {
		fmt.Println(v)
	}

	if v, e := statedb.read("ji"); e {
		fmt.Println(v)
	}

	//	buf := make(chan interface{})
	//	fmt.Println(GetVersion(3, 10, 4), VersionString(3, 10, 4), ParseVersion("4.5.312"), ParseVersion("4,5,312"))
	//	go producer(buf)
	//	go consumer(buf)
	//
	//	fmt.Println(`
	//		sljfla
	//		sljflsa
	//31231
	//sdfs
	//	`, string(0x3214321))
	//	mode := "# s . # . # # . . . . # # # . # # # . z"
	//	room := makeMaze(4, 5, mode)
	//	bfs(4, 5, room)
	//	l := list.New()
	//	// 入队, 压栈
	//	l.PushBack(1)
	//	l.PushBack(2)
	//	l.PushBack(3)
	//	l.PushBack(4)
	//	fmt.Println(l.Len())
	//	// 出队
	//	for i := 0; i < 4; i++ {
	//		i1 := l.Front()
	//		l.Remove(i1)
	//		fmt.Println(i1.Value.(int), l.Len())
	//	}
	select {}
}
