package mibulb2

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

//BulbSummary contains bulb info enough to communicate with a bulb
type BulbSummary struct {
	Id    int
	Ip    string
	Model string
}

func (b BulbSummary) String() string {
	return fmt.Sprintln(b.Id, b.Ip)
}

func parseResponce(r []byte) BulbSummary {
	str := string(r[:])

	m := make(map[string]string)
	for _, keyVal := range strings.Split(str, "\n") {
		keyValSeparated := strings.Split(keyVal, ": ")
		if len(keyValSeparated) < 2 {
			continue
		}
		k, v := keyValSeparated[0], keyValSeparated[1]
		m[k] = v[:len(v)-1]
	}
	result := BulbSummary{}
	rawId := m["id"]
	id64, _ := strconv.ParseInt(rawId[2:], 16, 0)
	result.Id = int(id64)
	result.Model = m["model"]
	result.Ip = strings.Split(m["Location"], "://")[1]
	return result
}

// Search does bulb discovery and returns BulbSummary for the found bulb
// it executes until any value is sent to @stopIndicator
func Search(stopIndicator <-chan bool) <-chan BulbSummary {
	cres := make(chan BulbSummary)

	go func() {
		discoverMsg := []byte("M-SEARCH * HTTP/1.1\r\nMAN: \"ssdp:discover\"\r\nST: wifi_bulb\r\n")
		lAddr, _ := net.ResolveUDPAddr("udp", ":50000")
		rAddr, _ := net.ResolveUDPAddr("udp", "239.255.255.250:1982")
		conn, _ := net.ListenUDP("udp", lAddr)
		defer conn.Close()
		var responseBuf = make([]byte, 2048)
		defer close(cres)
		for {
			println("iiiiiiiiiii")
			select {
			case <-stopIndicator:
				break
			default:
			}
			conn.WriteTo(discoverMsg, rAddr)

			conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(3)))
			n, _, _ := conn.ReadFrom(responseBuf)
			if n == 0 {
				continue
			}
			res := parseResponce(responseBuf)
			cres <- res
		}
	}()

	return cres
}
