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

type hostConfig struct {
	lines [][]byte
}

type SshConfig struct {
	hostConfigs []*hostConfig
}

/*
	Split the input into sections starting with the host header. To do so, we look for two
	Host headers in the input, and then anything between them is the token.

	If we're at the EOF, we just return stuff up to the first Host header (if it is not at position zero),
	and then at the next iteration, since we don't have any more Host headers in the input, we just return
	everything
*/
func hostHeaderSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var hostHeader []byte = []byte("Host ")
	var firstHeaderIndex, secondHeaderIndex int

	firstHeaderIndex = bytes.Index(data, hostHeader)
	if len(data) > firstHeaderIndex+len(hostHeader)+1 {
		secondHeaderIndex = bytes.Index(data[firstHeaderIndex+1:], hostHeader)
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

func newHostConfig(data []byte) *hostConfig {
	hc := &hostConfig{make([][]byte, 0)}

	byteReader := bytes.NewReader(data)
	lineSplitter := bufio.NewScanner(byteReader)
	lineSplitter.Split(bufio.ScanLines)

	for lineSplitter.Scan() {
		line := lineSplitter.Bytes()
		line = bytes.TrimSpace(line)

		if len(line) > 0 {
			hc.lines = append(hc.lines, line)
		}
	}

	if lineSplitter.Err() != nil {
		return nil
	} else {
		return hc
	}
}

func NewSshConfig() *SshConfig {
	return &SshConfig{make([]*hostConfig, 0)}
}

func (sshconfig *SshConfig) ReadConfig() error {
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

		sshconfig.hostConfigs = append(sshconfig.hostConfigs, hostConfig)
	}
	err = fileScanner.Err()
	if err != nil {
		fmt.Printf("Error while reading ssh config file lines. Error: %s\n", err.Error())
		return err
	}

	return nil
}

func (sshconfig *SshConfig) PrintCurrentConfig() {
	for _, hostConfig := range sshconfig.hostConfigs {
		fmt.Println("Starting new Host section")
		fmt.Println("********************************************************************************")
		for _, line := range hostConfig.lines {
			fmt.Println(string(line))
		}
		fmt.Println("********************************************************************************")
	}
}
