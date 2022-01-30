package icmpSurvive

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type goIcmpJob struct {
	addrs              []string
	finishNum, gos, id int
	job, finish        chan string
}

type Icmp struct {
	Type        uint8
	Code        uint8
	CheckSum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func ping(id uint16) Icmp {
	icmp := Icmp{Type: 8, Code: 0, CheckSum: 0, Identifier: 0, SequenceNum: id}
	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.BigEndian, icmp); err != nil {
		log.Panic(err)
	}
	icmp.CheckSum = check(buffer.Bytes())
	buffer.Reset()
	return icmp
}

func check(data []byte) uint16 {
	var (
		sum    uint32
		lenght int = len(data)
		index  int
	)

	for lenght > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		lenght -= 2
		index += 2
	}
	if lenght > 0 {
		sum += uint32(data[index])
	}
	return ^uint16(sum + sum>>16)
}

func (g goIcmpJob) sendICMP(domain string) {
	var (
		m        string
		icmp     = ping(uint16(g.id))
		laddr    = net.IPAddr{IP: net.ParseIP("0.0.0.0")}
		raddr, _ = net.ResolveIPAddr("ip4", domain)
	)
	conn, err := net.DialIP("ip4:icmp", &laddr, raddr)

	if err != nil {
		log.Panic(err)
		return
	}

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, icmp); err != nil {
		log.Panic(err)
		return
	}

	if _, err := conn.Write(buf.Bytes()); err != nil {
		log.Panic(err)
		return
	}

	if err := conn.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
		log.Panic(err)
		return
	}

	recv := make([]byte, 1024)
	if _, err := conn.Read(recv); err != nil {
		m = fmt.Sprintf("[-]%v down", domain)
	} else {
		m = fmt.Sprintf("[+]%v   up", domain)
	}
	g.finish <- m
	_ = conn.Close()
}

func (g goIcmpJob) Start() {
	log.Println("开始ICMP存活扫描……")
	var wg sync.WaitGroup
	for i := 0; i < g.gos; i++ {
		go g.goroutine()
	}
	wg.Add(2)
	go g.pushJob(&wg)
	go g.getFinish(&wg)
	wg.Wait()
	log.Println("扫描结束")
}

func (g goIcmpJob) Close() {
	close(g.job)
	close(g.finish)
}

func (g goIcmpJob) goroutine() {
	for i := range g.job {
		go g.sendICMP(i)
	}
}

func (g goIcmpJob) pushJob(wg *sync.WaitGroup) {
	for i := 1; i <= len(g.addrs); i++ {
		g.job <- g.addrs[i-1]
	}
	wg.Done()
}

func (g goIcmpJob) getFinish(wg *sync.WaitGroup) {
	for i := range g.finish {
		fmt.Println(i)
		g.finishNum++
		if g.finishNum == len(g.addrs) {
			wg.Done()
			return
		}
	}
}

func NewScan(gos int, addrs []string, jobNum, finishNum int) *goIcmpJob {
	s := &goIcmpJob{
		addrs:  addrs,
		gos:    gos,
		job:    make(chan string, jobNum),
		finish: make(chan string, finishNum),
	}

	return s
}
