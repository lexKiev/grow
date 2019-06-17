package main

import (
	"fmt"
	"strconv"
	"time"
)

var autoFanControl = true
var currentFanStep = 1
var steps = map[int]int{
	1: 1,
	2: 250,
	3: 500,
	4: 750,
	5: 1023,
}

func controlRun(interval int) {
	pollIntervalMin := interval
	timerChl := time.Tick(time.Duration(pollIntervalMin) * time.Minute)

	for range timerChl {
		if autoFanControl {
			autoFans()
		}
	}
}

func autoFans() {
	targetTemp := 27.1

	minStep := 1
	maxStep := 5

	if lowerSensorTemperature != 0 && lowerSensorTemperature > targetTemp {
		if currentFanStep < maxStep {
			currentFanStep++
			t := strconv.Itoa(steps[currentFanStep])
			fansSetSpeed(t)
			fmt.Printf("Temperature is %v increasing fan speed to %v\n", lowerSensorTemperature, currentFanStep)
		}
		if lowerSensorTemperature > targetTemp && currentFanStep == minStep {
			fansToggle(17, "on")
		}
	}

	if lowerSensorTemperature != 0 && lowerSensorTemperature < targetTemp {
		if currentFanStep > minStep {
			currentFanStep--
			t := strconv.Itoa(steps[currentFanStep])
			fansSetSpeed(t)
			fmt.Printf("Temperature is %v decreasing fan speed to %v\n", lowerSensorTemperature, currentFanStep)

		}

		if lowerSensorTemperature < targetTemp && currentFanStep == minStep {
			fansToggle(17, "off")
		}
	}
}
