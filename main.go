package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/kidoman/embd/host/all"
	"github.com/stianeikeland/go-rpio"
	"html/template"
	"net/http"
	"os/exec"
	"strconv"
)

type outputData struct {
	FansStatus  fansData
	SensorsData SensorData
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/head.html", "templates/footer.html")
	if err != nil {
		fmt.Fprint(writer, err.Error())
	}

	output := outputData{
		FansStatus:  fansDataFromDb(),
		SensorsData: sensorsDatafromDb(),
	}

	if request.Method != http.MethodPost {
		t.ExecuteTemplate(writer, "index", output)
		return
	}

	request.ParseForm()
	formsSent := request.Form
	fmt.Println(formsSent)
	if _, ok := formsSent["controlToggle"]; ok {
		fansControlModeToggle()
	}

	if _, ok := formsSent["pin"]; ok {
		pinNum, err := strconv.ParseInt(request.FormValue("pin"), 10, 64)
		if err != nil {
			fmt.Fprint(writer, err.Error())
			return
		}
		fansToggle(pinNum, "toggle")
	}

	if _, ok := formsSent["pumpToggle"]; ok {
		pumpToggle("toggle")
	}

	output.FansStatus = fansDataFromDb()
	t.ExecuteTemplate(writer, "index", output)
	return
}

func ajaxHandler(writer http.ResponseWriter, request *http.Request) {
	speed := request.FormValue("speed")

	fansSetSpeed(speed)

	js, err := json.Marshal(speed)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

func ajaxUpdateData(writer http.ResponseWriter, request *http.Request) {
	output := outputData{
		FansStatus:  fansDataFromDb(),
		SensorsData: sensorsDatafromDb(),
	}

	js, err := json.Marshal(output)
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

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "*"
	dbPass := "*"
	dbName := "*"
	dbHost := "*"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@"+dbHost+"/"+dbName)
	if err != nil {
		panic(err.Error())
	}

	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)

	return db
}

func main() {
	initRpi()
	hyg = initHygrometer()
	upperSensor, lowerSensor = initMeteostation()
	go sensorsRun(1)
	go controlRun(1)
	fmt.Println("Server UP on port 8081")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/changeSpeed", ajaxHandler)
	http.HandleFunc("/updateData", ajaxUpdateData)
	http.ListenAndServe(":8081", nil)
}
