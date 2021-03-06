package main

import (
	"net/http"
	"log"
	"time"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/mem"

)

func main (){
	//初始化日志服务
	logger := log.New(os.Stdout, "[Memory]", log.Lshortfile | log.Ldate | log.Ltime)

	//初始一个http handler
	http.Handle("/metrics", promhttp.Handler())

	//初始化一个容器
	diskPercent := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "memeory_percent",
		Help: "memeory use percent",
	},
		[]string {"percent"},
	)
	prometheus.MustRegister(diskPercent)

	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "test_counter",
		Help: "tracting counter",
	})
	prometheus.MustRegister(counter)

	// 启动web服务，监听1010端口
	go func() {
		logger.Println("ListenAndServe at:localhost:1010")
		err := http.ListenAndServe("localhost:1010", nil)
		if err != nil {
			logger.Fatal("ListenAndServe: ", err)
		}
	}()

	//收集内存使用的百分比
	for {
		logger.Println("start collect memory used percent!")
		v, err := mem.VirtualMemory()
		if err != nil {
			logger.Println("get memeory use percent error:%s", err)
		}
		usedPercent := v.UsedPercent
		logger.Println("get memeory use percent:", usedPercent)
		diskPercent.WithLabelValues("usedMemory").Set(usedPercent)

		counter.Add(1)
		time.Sleep(time.Second*2)
	}
}