// Copyright (c) 2023 Luka Ivanović
// This code is licensed under MIT licence (see LICENCE for details)

package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

func clearScreen() {
	print("\033[H\033[2J")
}

func getDateAndTime() string {
	return time.Now().Local().Format("2006-01-02 | 15:04")
}

func getBatteryPercentage() string {
	batteryOutput, _ := exec.Command("upower", "-i", "/org/freedesktop/UPower/devices/battery_BAT0").Output()
	percentageRx := regexp.MustCompile(" *percentage: *(.*)\n")
	percentage := percentageRx.FindSubmatch(batteryOutput)[1]
	return string(percentage)
}

func getBatteryState() string {
	batteryOutput, _ := exec.Command("upower", "-i", "/org/freedesktop/UPower/devices/battery_BAT0").Output()
	stateRx := regexp.MustCompile(" *state: *(.*)\n")
	state := stateRx.FindSubmatch(batteryOutput)[1]
	return string(state)
}

func getVolume() string {
	volume, _ := exec.Command("wpctl", "get-volume", "@DEFAULT_SINK@").Output()
	return string(volume)
}

func getStatusLine() string {
	dateAndTime := getDateAndTime()
	batteryPercentage := getBatteryPercentage()
	batteryState := getBatteryState()
	volume := getVolume()
	var statusLine strings.Builder
	statusLine.WriteString(dateAndTime)
	statusLine.WriteString(" | ")
	statusLine.WriteString(batteryPercentage)
	if batteryState != "discharging" {
		statusLine.WriteString("(" + batteryState + ")")
	}
	statusLine.WriteString(" | ")
	statusLine.WriteString(volume)
	return statusLine.String()
}

func main() {
	clearScreen()
	fmt.Println(getStatusLine())
	ticker := time.NewTicker(10 * time.Second)
	var blocker sync.WaitGroup
	blocker.Add(1)
	go func() {
		for {
			select {
			case <-ticker.C:
				clearScreen()
				fmt.Println(getStatusLine())
			}
		}
	}()
	blocker.Wait()
}
