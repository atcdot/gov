package cmd

import (
	"fmt"
	"strings"

	"gov/internal/gov"
	"gov/internal/version"
)

func listInstalled() {
	ii, err := gov.ListInstalled()
	if err != nil {
		fmt.Println(err)
	}

	for _, installed := range ii {
		fmt.Printf(" - %s, bin: %s", installed.Version, installed.Bin)

		if installed.IsMain {
			fmt.Print(" (system)")
		}

		fmt.Println()
	}
}

func listAvailable() {
	ai, err := gov.ListAvailable()
	if err != nil {
		fmt.Println(err)
	}

	regrouped := make(map[string][]string)
	for _, available := range ai {
		minorVersion := version.ExtractMinor(available)
		regrouped[minorVersion] = append(regrouped[minorVersion], available)
	}

	for minor, group := range regrouped {
		fmt.Printf(" - %s: %s\n", minor, strings.Join(group, ", "))
	}
}

func install(version string) {
	bin, err := gov.Install(version)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("bin: %s\n", bin)
}

func use(version string) {
	bin, err := gov.Use(version)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("bin: %s\n", bin)
}

func remove(version string) {
	err := gov.Remove(version)
	if err != nil {
		fmt.Println(err)
	}
}

func clear() {
	err := gov.Clear()
	if err != nil {
		fmt.Println(err)
	}
}
