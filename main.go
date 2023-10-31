package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var TlsDic = map[uint16]string{
	0x0301: "1.0",
	0x0302: "1.1",
	0x0303: "1.2",
	0x0304: "1.3",
}

const (
	defaultAddress     = "0.0.0.0"
	defaultPorts       = "443,8443,8080"
	defaultThreadCount = 20
	defaultTimeout     = 2
	outPutDef          = true
	outPutFileName     = "results.txt"
	showFailDef        = false
	ipIncrement        = 10000
)

func main() {
	addrPtr := flag.String("addr", defaultAddress, "Start scanning destination	")
	portsPtr := flag.String("ports", defaultPorts, "Ports to scan")
	threadPtr := flag.Int("thread", defaultThreadCount, "Number of parallel threads to scan")
	outPutFile := flag.Bool("o", outPutDef, "Is output to outPutFileName")
	timeOutPtr := flag.Int("timeOut", defaultTimeout, "Time out of a scan")
	showFailPtr := flag.Bool("showFail", showFailDef, "Is Show fail logs")

	flag.Parse()

	ports := strings.Split(*portsPtr, ",")
	for _, port := range ports {
		s := Scanner{
			addr:           *addrPtr,
			port:           port,
			showFail:       *showFailPtr,
			output:         *outPutFile,
			timeout:        time.Duration(*timeOutPtr) * time.Second,
			wg:             &sync.WaitGroup{},
			numberOfThread: *threadPtr,
			mu:             sync.Mutex{},
		}

		if *outPutFile {
			var err error
			s.logFile, err = os.OpenFile(outPutFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)

			if err != nil {
				fmt.Println("Failed to open log file:", err)
				return
			}

			defer s.logFile.Close()
		}

		s.Run()
	}
}

type Scanner struct {
	addr           string
	port           string
	output         bool
	showFail       bool
	logFile        *os.File
	numberOfThread int
	timeout        time.Duration
	wg             *sync.WaitGroup
	mu             sync.Mutex
	ip             net.IP
}

func (s *Scanner) Run() {
	s.mu.Lock()
	s.ip = net.ParseIP(s.addr)

	if s.ip.To4() == nil {
		s.addr = "[" + s.ip.String() + "]"
	}

	s.mu.Unlock()

	numIPsToCheck := ipIncrement

	for i := 0; i < s.numberOfThread; i++ {
		for j := 0; j < numIPsToCheck; j++ {
			nextIP := s.nextIP(i%2 == 0)
			if nextIP != nil {
				s.wg.Add(1)
				go func(ip net.IP) {
					defer s.wg.Done()
					s.Scan(ip)
				}(nextIP)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	s.wg.Wait()
}

func (s *Scanner) nextIP(increment bool) net.IP {
	s.mu.Lock()
	defer s.mu.Unlock()

	ipb := big.NewInt(0).SetBytes(s.ip.To4())
	if increment {
		ipb.Add(ipb, big.NewInt(1))
	} else {
		ipb.Sub(ipb, big.NewInt(1))
	}

	b := ipb.Bytes()
	b = append(make([]byte, 4-len(b)), b...)
	nextIP := net.IP(b)

	if nextIP.Equal(net.ParseIP("0.0.0.0")) || nextIP.Equal(net.ParseIP("255.255.255.255")) {
		return nil
	}

	s.ip = nextIP
	return s.ip
}

func (s *Scanner) Scan(ip net.IP) {
	str := ip.String()

	if ip.To4() == nil {
		str = "[" + str + "]"
	}

	dialer := &net.Dialer{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := dialer.DialContext(ctx, "tcp", str+":"+s.port)

	if err != nil {
		if s.showFail {
			s.Print(fmt.Sprint("Dial failed:", err))
		}
		return
	}

	defer conn.Close() // Ensure the connection is closed

	line := "" + conn.RemoteAddr().String() + "\t"
	conn.SetDeadline(time.Now().Add(s.timeout))
	c := tls.Client(conn, &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"h2", "http/1.1"},
	})
	err = c.Handshake()

	if err != nil {
		if s.showFail {
			s.Print(fmt.Sprint("", line, "TLS handshake failed:", err))
		}
		return
	}

	defer c.Close() // Ensure the TLS client is also properly closed

	state := c.ConnectionState()
	alpn := state.NegotiatedProtocol

	if s.showFail || (state.Version == 0x0304 && alpn == "h2") {
		certSubject := ""
		if len(state.PeerCertificates) > 0 {
			certSubject = state.PeerCertificates[0].Subject.CommonName
		}

		numPeriods := strings.Count(certSubject, ".")

		// Skip if certSubject is a wildcard domain, localhost, or not a top-level domain
		if strings.HasPrefix(certSubject, "*") || certSubject == "localhost" || numPeriods != 1 || certSubject == "invalid2.invalid" {
			return
		}

		// Print information about valid TLD with TLS v1.3 and ALPN h2
		s.Print(fmt.Sprint("  ", line, "---- TLS v", TlsDic[state.Version], "    ALPN: ", alpn, " ----    ", certSubject, ":", s.port, "\n"))
	}
}

func (s *Scanner) Print(outStr string) {
	fmt.Println(outStr)

	if s.output {
		_, err := s.logFile.WriteString(outStr + "\n")
		if err != nil {
			fmt.Println("Error writing into file:", err)
		}
	}
}
