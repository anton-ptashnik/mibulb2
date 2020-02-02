package mibulb2

import (
	"encoding/json"
	"fmt"
	"net"
)

// Bulb represents a bulb to communicate with
type Bulb struct {
	BulbSummary
}

type controlRequest struct {
	Id     int           `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

func (bulb *Bulb) execCommand(method string, params []interface{}) {
	rawReq, _ := json.Marshal(controlRequest{Id: bulb.Id, Method: method, Params: params})
	rawReq = append(rawReq, []byte("\r\n")...)

	conn, _ := net.Dial("tcp", bulb.Ip)
	defer conn.Close()
	conn.Write(rawReq)

	var buf [1024]byte
	n, _ := conn.Read(buf[:])
	fmt.Println("response:", string(buf[:n]))
}

func (bulb *Bulb) Toggle() {
	bulb.execCommand("toggle", []interface{}{})
}
func (bulb *Bulb) SetPower(state bool) {
	var strStateRepr string
	if state {
		strStateRepr = "on"
	} else {
		strStateRepr = "off"
	}
	bulb.execCommand("set_power", []interface{}{strStateRepr, "smooth", 1000})
}
func (bulb *Bulb) GetPower() {
	bulb.execCommand("get_prop", []interface{}{"power"})
}
func (bulb *Bulb) SetColor(rgb int) {
	bulb.execCommand("set_rgb", []interface{}{rgb, "smooth", 1000})
}
func (bulb *Bulb) GetColor() {
	bulb.execCommand("get_prop", []interface{}{"rgb"})
}
func (bulb *Bulb) SetBrightness(level int) {
	bulb.execCommand("set_bright", []interface{}{level, "smooth", 1000})
}
func (bulb *Bulb) GetBrightness() {
	bulb.execCommand("get_prop", []interface{}{"bright"})
}
func (bulb *Bulb) DiscardTimer() {
	bulb.execCommand("cron_del", []interface{}{0})
}
func (bulb *Bulb) SetTimer(valInMinutes int) {
	bulb.execCommand("cron_add", []interface{}{0, valInMinutes})
}

func (bulb *Bulb) GetTimer() {
	bulb.execCommand("cron_get", []interface{}{0})
}

func (bulb *Bulb) SaveState() {
	bulb.execCommand("set_default", []interface{}{})
}
