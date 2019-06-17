package main

import "github.com/stianeikeland/go-rpio"

var pumpPin = 27

func getPumpStatus() int {
	pin := rpio.Pin(pumpPin)
	res := pin.Read()
	var pumpStatus = int(res)

	return pumpStatus
}

func pumpToggle(command string) {

	pin := rpio.Pin(pumpPin)
	pin.Output()

	switch command {
	case "toggle":
		pin.Toggle()
	case "on":
		pin.Low()
	case "off":
		pin.High()
	}

	//pumpStates := getPumpStatus()
	//
	//db := dbConn()
	//upd, err := db.Query("UPDATE fans SET status=? WHERE fans_group=?", value, key)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//upd.Close()
	//db.Close()
}
