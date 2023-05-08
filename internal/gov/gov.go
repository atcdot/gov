package gov

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/atcdot/gov/internal/dlsdk"
	"github.com/atcdot/gov/internal/version"
)

type VersionInstalled struct {
	Version  string
	Bin      string
	IsMain   bool
	IsActive bool
}

func ListInstalled() ([]VersionInstalled, error) {
	actualVersion, goBin, err := GetActualVersion()
	if err != nil {
		return nil, err
	}

	ii := []VersionInstalled{
		{
			Version: actualVersion,
			Bin:     goBin,
			IsMain:  true,
		},
	}

	goPath, err := getGoPath()
	if err != nil {
		return nil, err
	}

	goBinPath := path.Join(goPath, "Bin")

	ee, err := os.ReadDir(goBinPath)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(ee, func(i, j int) bool {
		return ee[i].Name() > ee[j].Name()
	})

	for _, entry := range ee {
		if entry.IsDir() {
			continue
		}

		if isBinGo(entry.Name()) {
			goBin := path.Join(goBinPath, entry.Name())
			goBinVersion, err := getBinVersion(goBin)
			if err != nil {
				return nil, err
			}

			i := VersionInstalled{
				Version: version.ExtractFull(goBinVersion),
				Bin:     goBin,
			}

			ii = append(ii, i)
		}
	}

	activeBinPath, err := getActiveBinPath()
	for i := range ii {
		if ii[i].Bin == activeBinPath {
			ii[i].IsActive = true
			break
		}
	}

	return ii, nil
}

func ListAvailable() ([]string, error) {
	versions, err := dlsdk.GetVersions()
	if err != nil {
		return nil, err
	}

	versionsFormatted := []string{}
	for _, v := range versions {
		versionsFormatted = append(versionsFormatted, version.ExtractFull(v))
	}

	return versionsFormatted, nil
}

func Install(version string) (string, error) {
	goPath, err := getGoPath()
	if err != nil {
		return "", err
	}
	bin := path.Join(goPath, "Bin", fmt.Sprintf("go%s", version))

	// install
	{
		cmd := exec.Command(goBinCmdSystem, "install", fmt.Sprintf("golang.org/dl/go%s@latest", version))

		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			return "", err
		}
	}

	// download
	{
		cmd := exec.Command(bin, "download")

		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			return "", err
		}
	}

	return bin, nil
}

func Use(version string) (string, error) {
	vv, err := ListInstalled()
	if err != nil {
		return "", err
	}

	var bin string

	for _, v := range vv {
		if v.Version == version || (version == "system" && v.IsMain) {
			bin = v.Bin
		}
	}

	if bin == "" {
		return "", fmt.Errorf("version %s not installed", version)
	}

	// use (set alias)
	{
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.OpenFile(dirname+"/.gov", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.FileMode(0644))
		if err != nil {
			fmt.Println(err)
		}

		fileContent := fmt.Sprintf("alias go=%s\n", bin)

		_, err = f.Write([]byte(fileContent))
		if err != nil {
			fmt.Println(err)
		}
	}

	return bin, nil
}

func IsInitialised() bool {
	return fileExists(govSystemConfFilename) &&
		fileExists(govAliasFilename)
}

func SaveActualVersion() error {
	_, _, err := GetActualVersion()
	if err != nil {
		return err
	}

	return nil
}

func Remove(version string) error {
	goPath, err := getGoPath()
	if err != nil {
		return err
	}

	bin := path.Join(goPath, "Bin", fmt.Sprintf("go%s", version))

	r, err := getGoRoot(bin)
	if err != nil {
		return err
	}

	// delete go root
	{
		cmd := exec.Command("rm", "-rf", r)

		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	// delete binary
	{
		cmd := exec.Command("rm", bin)

		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func Cleanup() error {
	err := os.Remove(govSystemConfFilename)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.OpenFile(govAliasFilename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.FileMode(0644))
	if err != nil {
		fmt.Println(err)
	}

	_, err = f.Write([]byte(""))
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func GetActualVersion() (v string, goBin string, err error) {
	if IsInitialised() {
		return getSystemVersionFromFile()
	}

	goBin, err = which()
	if err != nil {
		return "", "", err
	}

	goBin = path.Clean(strings.ReplaceAll(goBin, "\n", ""))

	goBinVersion, err := getBinVersion(goBin)
	if err != nil {
		return "", "", err
	}

	v = version.ExtractFull(goBinVersion)

	if err := setSystemVersionToFile(v, goBin); err != nil {
		return "", "", err
	}

	return v, goBin, nil
}

func which() (string, error) {
	cmd := exec.Command("which", goBinCmdSystem)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func isBinGo(binName string) bool {
	return regexp.MustCompile(`go\d+(\.\d+)+`).MatchString(binName)
}

func getGoPath() (string, error) {
	cmd := exec.Command(goBinCmdSystem, `env`, `GOPATH`)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return path.Clean(strings.ReplaceAll(out.String(), "\n", "")), nil
}

func getGoRoot(goBin string) (string, error) {
	cmd := exec.Command(goBin, `env`, `GOROOT`)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(out.String(), "\n", ""), nil
}

func getBinVersion(goBin string) (string, error) {
	cmd := exec.Command(goBin, "version")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
