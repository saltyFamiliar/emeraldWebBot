package commands

import (
	"fmt"
	"net"
	"slices"
	"time"
)

type PortData struct {
	Number int
	IsOpen bool
}

func ScanPorts(address string, startPort, endPort, numWorkers int) ([]int, []int) {
	jobCh := make(chan int, numWorkers)
	resCh := make(chan *PortData, numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				p := <-jobCh
				fmt.Println("Checking ", p)
				addressStr := fmt.Sprintf("%s:%d", address, p)
				conn, err := net.DialTimeout("tcp", addressStr, time.Second*2)
				if err != nil {
					resCh <- &PortData{Number: p, IsOpen: false}
				} else {
					resCh <- &PortData{Number: p, IsOpen: true}
					if err = conn.Close(); err != nil {
						panic(err)
					}
				}
			}
		}()
	}

	go func() {
		for p := startPort; p <= endPort; p++ {
			jobCh <- p
		}
	}()

	open, closed := make([]int, 0, endPort-startPort), make([]int, 0, endPort-startPort)

	for i := 0; i <= endPort-startPort; i++ {
		p := <-resCh
		if p.IsOpen {
			open = append(open, p.Number)
		} else {
			closed = append(closed, p.Number)
		}
	}

	slices.Sort(open)
	slices.Sort(closed)

	return open, closed
}
