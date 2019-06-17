package main

import (
	"fmt"
	"github.com/MichaelS11/go-dht"
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/convertors/mcp3008"
	"time"
)

var hyg *mcp3008.MCP3008
var upperSensor, lowerSensor *dht.DHT
var upperSensorHumidity, upperSensorTemperature, lowerSensorHumidity, lowerSensorTemperature float64

type Sensor struct {
	Id                 int
	Value              float32
	Name, Type, Update string
}

type SensorData struct {
	Humidity     map[string][]Sensor
	Temperature  map[string][]Sensor
	SoilMoisture map[string][]Sensor
}

func initHygrometer() *mcp3008.MCP3008 {
	bus := embd.NewSPIBus(0, 0, 5000, 0, 1)

	mpc3008 := mcp3008.New(1, bus)
	fmt.Println("Hygrometer initialized!")
	return mpc3008
}

func initMeteostation() (*dht.DHT, *dht.DHT) {
	err := dht.HostInit()
	if err != nil {
		fmt.Println("DHT Init error:", err)
	}

	lowerSensor, err := dht.NewDHT("GPIO6", dht.Celsius, "")
	upperSensor, err := dht.NewDHT("GPIO5", dht.Celsius, "")
	if err != nil {
		fmt.Println("NewDHT error:", err)
	}
	fmt.Println("Meteostation initialized!")
	return lowerSensor, upperSensor
}

func sensorsDataToDb() {
	db := dbConn()
	hygValue1, _ := hyg.AnalogValueAt(0)
	hygValue2, _ := hyg.AnalogValueAt(1)
	upperSensorHumidity, upperSensorTemperature, _ = upperSensor.ReadRetry(10)
	lowerSensorHumidity, lowerSensorTemperature, _ = lowerSensor.ReadRetry(10)

	type sd struct {
		sZone  string
		sType  string
		sValue float64
	}

	sensor1 := sd{
		"Surface",
		"hygrometer",
		float64(hygValue1),
	}
	sensor2 := sd{
		"Root",
		"hygrometer",
		float64(hygValue2),
	}
	sensor3 := sd{
		"Upper",
		"humidity",
		upperSensorHumidity,
	}
	sensor4 := sd{
		"Lower",
		"humidity",
		lowerSensorHumidity,
	}
	sensor5 := sd{
		"Upper",
		"temperature",
		upperSensorTemperature,
	}
	sensor6 := sd{
		"Lower",
		"temperature",
		lowerSensorTemperature,
	}

	values := []sd{
		sensor1, sensor2, sensor3, sensor4, sensor5, sensor6,
	}

	for _, value := range values {
		ins, err := db.Query("INSERT INTO sensors (sensor_name,sensor_type,sensor_value) VALUES (?,?,?)", value.sZone, value.sType, value.sValue)
		if err != nil {
			fmt.Println("########################### sensorsDataToDb")
			fmt.Println(err.Error())
		}
		ins.Close()
	}

	db.Close()
}

func sensorsDatafromDb() SensorData {
	db := dbConn()

	rm, _ := db.Query("DELETE FROM sensors WHERE NOT YEARWEEK(`date`, 1) = YEARWEEK(CURDATE(), 1)")
	rm.Close()
	selDB, err := db.Query("SELECT * FROM sensors ORDER BY id")
	if err != nil {
		fmt.Println(err.Error())
	}

	sensor := Sensor{}
	sensorData := SensorData{}
	var humSensor = make([]Sensor, 0)
	var tempSensor = make([]Sensor, 0)
	var hygSensor = make([]Sensor, 0)

	for selDB.Next() {
		var id int
		var sensorValue float32
		var sensorName, sensorType, updated string
		err = selDB.Scan(&id, &sensorName, &sensorType, &sensorValue, &updated)
		if err != nil {
			fmt.Println("########################### sensorsDatafromDb")
			fmt.Println(err.Error())
		}
		sensor.Id = id
		sensor.Name = sensorName
		sensor.Type = sensorType
		sensor.Value = sensorValue
		sensor.Update = updated

		switch sensorType {
		case "humidity":
			humSensor = append(humSensor, sensor)
		case "temperature":
			tempSensor = append(tempSensor, sensor)
		case "hygrometer":
			hygSensor = append(hygSensor, sensor)
		}

	}

	if len(humSensor) > 1 {
		sensorData.Humidity = map[string][]Sensor{
			"current": humSensor[len(humSensor)-2:],
			"all":     humSensor,
		}
	}

	if len(hygSensor) > 1 {
		sensorData.SoilMoisture = map[string][]Sensor{
			"current": hygSensor[len(hygSensor)-2:],
			"all":     hygSensor,
		}
	}

	if len(tempSensor) > 1 {
		sensorData.Temperature = map[string][]Sensor{
			"current": tempSensor[len(tempSensor)-2:],
			"all":     tempSensor,
		}
	}

	selDB.Close()
	db.Close()

	return sensorData
}

func sensorsRun(interval int) {
	pollInterval := interval
	timerCh := time.Tick(time.Duration(pollInterval) * time.Minute)

	for range timerCh {
		sensorsDataToDb()
	}
}
