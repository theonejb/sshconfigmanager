/*
	Provides an interface to read and write ssh config files
*/
package configreader

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path"
)

type SshConfig struct {
	hostConfigs []*hostConfig
}

/*
Split the input into sections starting with the host header. To do so, we look for two
Host headers in the lower cased input, and then anything between them is the token.

If we're at the EOF and we don't have a second host header, return all remaining data. This
is because we can't find any other host header (this being the EOF and all). Otherwise, we
just wait till we find a second header and return all data between the first and second header.
*/
func hostHeaderSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var hostHeader []byte = []byte("host ")
	var firstHeaderIndex, secondHeaderIndex int

	lowerCasedData := bytes.ToLower(data)
	firstHeaderIndex = bytes.Index(lowerCasedData, hostHeader)
	// Only check for a second header if the input data is large enough to contain it
	if len(data) > firstHeaderIndex+len(hostHeader)+1 {
		secondHeaderIndex = bytes.Index(lowerCasedData[firstHeaderIndex+1:], hostHeader)
	} else {
		secondHeaderIndex = -1
	}

	if secondHeaderIndex == -1 && atEOF {
		return len(data), data[firstHeaderIndex:], nil
	} else if secondHeaderIndex != -1 {
		return secondHeaderIndex + 1, data[firstHeaderIndex:secondHeaderIndex], nil
	}

	return 0, nil, nil
}

func NewSshConfig() *SshConfig {
	return &SshConfig{make([]*hostConfig, 0)}
}

func (sc *SshConfig) ReadConfig() error {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Printf("Error in getting current user. Error: %s\n", err.Error())
		return err
	}

	homeDir := currentUser.HomeDir
	sshDir := path.Join(homeDir, ".ssh/")
	configFilePath := path.Join(sshDir, "config")

	configFile, err := os.Open(configFilePath)
	if err != nil {
		fmt.Printf("Error in opening ssh config file. Path: [%s]. Error: %s\n", configFilePath, err.Error())
		return err
	}

	fileScanner := bufio.NewScanner(configFile)
	fileScanner.Split(hostHeaderSplitFunc)
	for fileScanner.Scan() {
		hostConfigSection := fileScanner.Bytes()
		hostConfig := newHostConfig(hostConfigSection)

		sc.hostConfigs = append(sc.hostConfigs, hostConfig)
	}
	err = fileScanner.Err()
	if err != nil {
		fmt.Printf("Error while reading ssh config file lines. Error: %s\n", err.Error())
		return err
	}

	return nil
}

func (sc *SshConfig) PrintCurrentConfig() {
	for _, hostConfig := range sc.hostConfigs {
		hostConfig.PrintConfig()
	}
}
