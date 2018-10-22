package main

import (
    "flag"
    "log"
    "net/http"
    _ "net/http/pprof"
    "sync"
    "time"
    "strconv"
    "runtime"
    //"runtime/pprof"
)

//func handler(w http.ResponseWriter, r *http.Request) {
//    w.Header().Set("Content-Type", "text/plain")
//
//    p := pprof.Lookup("goroutine")
//    p.WriteTo(w, 1)
//}

//for debug
func startPprof() {
    http.HandleFunc("/goroutines", func(w http.ResponseWriter, r *http.Request) {
        num := strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
        w.Write([]byte(num))
    })
    log.Println(http.ListenAndServe("localhost:8008", nil))
}

func Counter(wg *sync.WaitGroup) {
    time.Sleep(time.Second)

    var counter int
    for i := 0; i < 1000000; i++ {
        time.Sleep(time.Millisecond * 200)
        counter++
    }
    wg.Done()
}

func startPprof() {
    log.Println(http.ListenAndServe("localhost:8008", nil))
}

func main() {
    flag.Parse()

    //远程获取pprof数据
    go startPprof()

    var wg sync.WaitGroup
    wg.Add(10)
    for i := 0; i < 10; i++ {
        go Counter(&wg)
    }
    wg.Wait()

    // sleep 10mins, 在程序退出之前可以查看性能参数.
    time.Sleep(60 * time.Second)
}
