// Copyright © 2023- Luka Ivanović
// This code is licensed under the terms of the MIT licence (see LICENCE for details)

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func getDate() string {
	return time.Now().Local().Format("Date: 2006-01-02")
}

func getTime() string {
	return time.Now().Local().Format("Time: 15:04:05")
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

func getStatusLine(separator string) string {
	date := getDate()
	time := getTime()
	// batteryPercentage := getBatteryPercentage()
	// batteryState := getBatteryState()
	volume := getVolume()
	brightness := getBrightness()
	var statusLine strings.Builder
	statusLine.WriteString(date)
	statusLine.WriteString(separator)
	statusLine.WriteString(time)
	statusLine.WriteString(separator)
	// statusLine.WriteString(batteryPercentage)
	// if batteryState != "discharging" {
	// 	statusLine.WriteString("(" + batteryState + ")")
	// }
	// statusLine.WriteString(separator)
	statusLine.WriteString(volume)
	statusLine.WriteString(separator)
	statusLine.WriteString(brightness)
	return statusLine.String()
}

func main() {
	noTicker := flag.Bool("notick", false, "return the status line without updates every <period> seconds")
	separator := flag.String("separator", " | ", "separator to use between status items")
	n := flag.Int("period", 20, "how often to refresh the status line")
	flag.Parse()

	if *noTicker {
		fmt.Println(getStatusLine(*separator))
		os.Exit(0)
	}

	period := time.Duration(*n) * time.Second
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			fmt.Print("\r\033[K" + getStatusLine(*separator))
		case <-sigs:
			return
		}
	}

}
