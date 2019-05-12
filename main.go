package main

import (
	"encoding/json"
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"html/template"
	"net/http"
	"os/exec"
	"strconv"
)

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/head.html", "templates/footer.html")
	if err != nil {
		fmt.Fprint(writer, err.Error())
	}

	if request.Method != http.MethodPost {
		t.ExecuteTemplate(writer, "index", getPinStatus())
		return
	}

	pinNum, err := strconv.ParseInt(request.FormValue("pin"), 10, 64)

	if err != nil {
		fmt.Fprint(writer, err.Error())
		return
	}

	pin := rpio.Pin(pinNum)
	pin.Output()
	pin.Toggle()
	t.ExecuteTemplate(writer, "index", getPinStatus())
	return
}

func ajaxHandler(writer http.ResponseWriter, request *http.Request) {
	type Profile struct {
		Speed string
	}
	profile := Profile{Speed: request.FormValue("speed")}
	cmd := exec.Command("gpio", "pwm", "23", profile.Speed)
	cmdErr := cmd.Run()
	if cmdErr != nil {
		fmt.Println(cmdErr)
	}
	currentSpeed = profile.Speed
	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

func initRpi() {
	err := rpio.Open()

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Raspberry initialized!")
	}

	cmd := exec.Command("gpio", "mode", "23", "pwm")
	cmdErr := cmd.Run()
	if cmdErr != nil {
		fmt.Println(cmdErr)
	} else {
		fmt.Println("PWM initialized!")
	}
}

func getPinStatus() map[string]string {
	pinStates := make(map[string]string)
	for k, v := range usedPins {
		pin := rpio.Pin(v)
		res := pin.Read()
		var pinStatus int64 = int64(res)
		if pinStatus == 1 {
			pinStates[k] = "on"
		} else {
			pinStates[k] = "off"
		}
	}

	pinStates["speed"] = currentSpeed

	return pinStates
}

var usedPins = map[string]int{"in": 17, "out": 18};

var currentSpeed string

func main() {

	fmt.Println("Server UP on port 8081")
	initRpi()
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/changeSpeed", ajaxHandler)
	http.ListenAndServe(":8081", nil)
}
