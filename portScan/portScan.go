package portScan

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type goScanJob struct {
	addrs          []string
	ports          []int64
	gos, finishNum int
	job, finish    chan string
}

func (g goScanJob) allScan(host string) {
	conn, err := net.DialTimeout("tcp", host, 1*time.Second)
	var m string
	if err != nil {
		m = fmt.Sprintf("[-]%s CLOSE\r", host)
	} else {
		m = fmt.Sprintf("[+]%s OPEN\n", host)
		_ = conn.Close()
	}
	g.finish <- m
}

func (g goScanJob) Start() {
	log.Println("开始扫描端口……")
	var wg = sync.WaitGroup{}
	for i := 0; i < g.gos; i++ {
		go g.goroutine()
	}
	wg.Add(2)
	go g.pushJob(&wg)
	go g.getFinish(&wg)
	wg.Wait()
	log.Println("扫描结束以上端口开启")
}

func (g goScanJob) goroutine() {
	for i := range g.job {
		go g.allScan(i)
	}
}

func (g goScanJob) Close() {
	close(g.finish)
	close(g.job)
}

func (g goScanJob) pushJob(wg *sync.WaitGroup) {
	for i := 1; i <= len(g.addrs); i++ {
		for p := 1; p <= len(g.ports); p++ {
			g.job <- fmt.Sprintf("%s:%d", g.addrs[i-1], g.ports[p-1])
		}
	}
	wg.Done()
}

func (g goScanJob) getFinish(wg *sync.WaitGroup) {
	for i := range g.finish {
		fmt.Printf("%s", i)
		g.finishNum++
		if g.finishNum == len(g.addrs)*len(g.ports) {
			wg.Done()
			return
		}
	}
}

func getPorts(ports string) []int64 {
	var portArr []int64
	if !strings.Contains(ports, "-") {
		port, _ := strconv.ParseInt(ports, 10, 0)
		if port < 1 || port > 65535 {
			panic("最小端口不要小于1，最大端口不要大于65535")
		}
		portArr = append(portArr, port)
	} else {
		portList := strings.Split(ports, "-")
		mini, _ := strconv.ParseInt(portList[0], 10, 0)
		max, _ := strconv.ParseInt(portList[1], 10, 0)
		if mini < 1 || max > 65535 {
			panic("最小端口不要小于1，最大端口不要大于65535")
		}
		for i := mini; i < max+1; i++ {
			portArr = append(portArr, i)
		}

	}
	return portArr
}

func NewScan(gos int, ports string, addrs []string, jobNum, finishNum int) *goScanJob {
	var portArr = getPorts(ports)
	s := &goScanJob{
		addrs:  addrs,
		ports:  portArr,
		gos:    gos,
		job:    make(chan string, jobNum),
		finish: make(chan string, finishNum),
	}
	return s

}
