package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os/exec"
)

var usedPins = map[string]int{"in": 17, "out": 18}

type fansGroup struct {
	Id, FanSpeed, Status, Manual int
	FansGroup                    string
}

type fansData struct {
	In, Out fansGroup
}

func getPinStatus() map[string]int {
	pinStates := make(map[string]int)
	for k, v := range usedPins {
		pin := rpio.Pin(v)
		res := pin.Read()
		var pinStatus = int(res)
		pinStates[k] = pinStatus
	}

	return pinStates
}

func fansDataToDb() {
	//db := dbConn()
	fansStatus := getPinStatus()
	fmt.Println(fansStatus)

	//_, err := db.Query("UPDATE fans SET status =?, fans_speed=?, manual=? WHERE fans_group=?", hygValue1)
	//insForm, err := db.Prepare("UPDATE fans SET name=?, city=? WHERE id=?")
	//if err != nil {
	//	panic(err.Error())
	//}
	//insForm.Exec(name, city, id)
	//
	//defer db.Close()
}

func fansDataFromDb() fansData {
	db := dbConn()
	selFans, err := db.Query("SELECT * FROM fans")
	if err != nil {
		fmt.Println(err.Error())
	}
	fansData := fansData{}
	for selFans.Next() {
		fg := fansGroup{}

		var fansGroup string
		var id, fanSpeed, status, manual int
		err = selFans.Scan(&id, &fansGroup, &fanSpeed, &status, &manual)
		if err != nil {
			fmt.Println("########################### fansDataFromDb")
			fmt.Println(err.Error())
		}

		fg.Id = id
		fg.FansGroup = fansGroup
		fg.FanSpeed = fanSpeed
		fg.Status = status
		fg.Manual = manual

		switch fansGroup {
		case "out":
			fansData.Out = fg
		case "in":
			fansData.In = fg
		}

		if manual == 1 {
			autoFanControl = false
		} else {
			autoFanControl = true
		}
	}
	selFans.Close()
	db.Close()

	return fansData
}

func fansControlModeToggle() {
	fanControlStatus := make(map[string]int)
	db := dbConn()

	selFans, err := db.Query("SELECT fans_group,manual FROM fans")
	if err != nil {
		fmt.Println("########################### fansControlModeToggle")
		fmt.Println(err.Error())
	}

	for selFans.Next() {
		var fansGroup string
		var manual int
		err = selFans.Scan(&fansGroup, &manual)
		if err != nil {
			fmt.Println(err.Error())
		}

		fanControlStatus[fansGroup] = manual
	}

	for key, value := range fanControlStatus {
		var newValue int
		if value == 1 {
			newValue = 0
			autoFanControl = true
		} else {
			newValue = 1
			autoFanControl = false
		}
		_, err := db.Query("UPDATE fans SET manual=? WHERE fans_group=?", newValue, key)
		if err != nil {
			println("########################### fansControlModeToggle")
			println(err.Error())
		}
	}
	selFans.Close()
	db.Close()
}

func fansToggle(pinNum int64, command string) {

	pin := rpio.Pin(pinNum)
	pin.Output()

	switch command {
	case "toggle":
		pin.Toggle()
	case "on":
		pin.Low()
	case "off":
		pin.High()
	}

	pinStates := getPinStatus()

	db := dbConn()

	for key, value := range pinStates {
		upd, err := db.Query("UPDATE fans SET status=? WHERE fans_group=?", value, key)
		if err != nil {
			fmt.Println(err.Error())
		}
		upd.Close()
	}

	db.Close()
}

func fansSetSpeed(speed string) {
	cmd := exec.Command("gpio", "pwm", "23", speed)
	cmdErr := cmd.Run()
	if cmdErr == nil {
		db := dbConn()
		upd, err := db.Query("UPDATE fans SET fans_speed=? ", speed)
		if err != nil {
			fmt.Println(err.Error())
		}
		upd.Close()
		db.Close()
	} else {
		fmt.Println(cmdErr)
	}
}
