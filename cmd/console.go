package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/atcdot/gov/internal/gov"
	"github.com/atcdot/gov/internal/version"
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

	minorVersions := make([]string, 0)
	for minorVersion := range regrouped {
		minorVersions = append(minorVersions, minorVersion)
	}

	sort.SliceStable(minorVersions, func(i, j int) bool {
		return minorVersions[i] > minorVersions[j]
	})

	for _, minor := range minorVersions {
		fmt.Printf(" - %s: %s\n", minor, strings.Join(regrouped[minor], ", "))
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

func cleanup() {
	err := gov.Cleanup()
	if err != nil {
		fmt.Println(err)
	}
}
