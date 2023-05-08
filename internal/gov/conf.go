package gov

import (
	"errors"
	"log"
	"os"
	"regexp"

	"github.com/atcdot/gov/internal/version"
	"github.com/joho/godotenv"
)

const (
	binPathKey = "GOV_SYSTEM_GO_BIN_PATH"
)

var (
	govSystemConfFilename = func() string {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		return dirname + "/.gov_system"
	}()

	govAliasFilename = func() string {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		return dirname + "/.gov"
	}()

	goBinCmdSystem = func() string {
		cmd := "go"
		if IsInitialised() {
			var err error
			_, cmd, err = getSystemVersionFromFile()
			if err != nil {
				log.Fatal(err)
			}
		}

		return cmd
	}()
)

func readSystemConf() (map[string]string, error) {
	if !fileExists(govSystemConfFilename) {
		return nil, errors.New("undefined")
	}

	return godotenv.Read(govSystemConfFilename)
}

func getSystemVersionFromFile() (v string, goBin string, err error) {
	systemConf, err := readSystemConf()
	if err != nil {
		return "", "", err
	}

	binPath, ok := systemConf[binPathKey]
	if !ok {
		return "", "", errors.New("system config invalid")
	}

	binVersion, err := getBinVersion(binPath)
	if err != nil {
		return "", "", err
	}

	return version.ExtractFull(binVersion), binPath, nil
}

func setSystemVersionToFile(v string, goBin string) error {
	return godotenv.Write(map[string]string{
		binPathKey: goBin,
	}, govSystemConfFilename)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getActiveBinPath() (string, error) {
	aliasBytes, err := os.ReadFile(govAliasFilename)
	switch {
	case err == nil:
		if len(aliasBytes) == 0 {
			_, bin, err := GetActualVersion()
			if err != nil {
				return "", err
			}

			return bin, nil
		}
	case os.IsNotExist(err):
		_, bin, err := GetActualVersion()
		if err != nil {
			return "", err
		}

		return bin, nil
	case err != nil:
		return "", err
	}

	ss := regexp.
		MustCompile("alias\\sgo=(.*)").
		FindStringSubmatch(string(aliasBytes))

	if len(ss) == 0 {
		return "", errors.New("alias invalid")
	}

	return ss[1], nil
}
