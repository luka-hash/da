// Copyright © 2023- Luka Ivanović
// This code is licensed under the terms of the MIT licence (see LICENCE for details)

package main

import (
	"flag"
	"fmt"
	"os"
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
	return time.Now().Local().Format("Date: 2006-01-02 | Time: 15:04")
}

func getBatteryPercentage() string {
	batteryOutput, _ := exec.Command("upower", "-i", "/org/freedesktop/UPower/devices/battery_BAT0").Output()
	percentageRx := regexp.MustCompile(" *percentage: *(.*)\n")
	percentage := percentageRx.FindSubmatch(batteryOutput)[1]
	return "Battery: " + string(percentage)
}

func getBatteryState() string {
	batteryOutput, _ := exec.Command("upower", "-i", "/org/freedesktop/UPower/devices/battery_BAT0").Output()
	stateRx := regexp.MustCompile(" *state: *(.*)\n")
	state := stateRx.FindSubmatch(batteryOutput)[1]
	return string(state)
}

func getVolume() string {
	volume, _ := exec.Command("wpctl", "get-volume", "@DEFAULT_SINK@").Output()
	return strings.TrimSpace(string(volume))
}

func getBrightness() string {
	brightness, _ := exec.Command("brightnessctl", "info").Output()
	brightnessRx := regexp.MustCompile(`\((.*)\)`)
	return "Brightness: " + brightnessRx.FindStringSubmatch(string(brightness))[1]
}

func getStatusLine() string {
	dateAndTime := getDateAndTime()
	batteryPercentage := getBatteryPercentage()
	batteryState := getBatteryState()
	volume := getVolume()
	brightness := getBrightness()
	var statusLine strings.Builder
	statusLine.WriteString(dateAndTime)
	statusLine.WriteString(" | ")
	statusLine.WriteString(batteryPercentage)
	if batteryState != "discharging" {
		statusLine.WriteString("(" + batteryState + ")")
	}
	statusLine.WriteString(" | ")
	statusLine.WriteString(volume)
	statusLine.WriteString(" | ")
	statusLine.WriteString(brightness)
	return statusLine.String()
}

func main() {
	noTicker := flag.Bool("notick", false, "return the status line without updates every <period> seconds")
	n := flag.Int("period", 20, "how often to refresh the status line")
	flag.Parse()
	period := time.Duration(*n) * time.Second
	if *noTicker {
		fmt.Println(getStatusLine())
		os.Exit(0)
	}
	clearScreen()
	fmt.Println(getStatusLine())
	ticker := time.NewTicker(period)
	var blocker sync.WaitGroup
	blocker.Add(1)
	go func() {
		for {
			<-ticker.C
			clearScreen()
			fmt.Println(getStatusLine())
		}
	}()
	blocker.Wait()
}
