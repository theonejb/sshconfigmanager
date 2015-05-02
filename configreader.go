/*
This file contains code to read ssh config files
*/
package sshconfigmanager

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"os/user"
	"path"
)

type SshConfig struct {
	fileHash    string
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

	if secondHeaderIndex == -1 && atEOF && len(data) > 0 {
		return len(data), data[firstHeaderIndex:], nil
	} else if secondHeaderIndex != -1 {
		return secondHeaderIndex + 1, data[firstHeaderIndex:secondHeaderIndex], nil
	}

	return 0, nil, nil
}

func NewSshConfig() *SshConfig {
	return &SshConfig{"", make([]*hostConfig, 0)}
}

func (sc *SshConfig) ReadConfig() error {
	configFile, err := getSshConfigFile()
	defer configFile.Close()
	if err != nil {
		return err
	}

	// Create a TeeReader that writes all bytes read from the config file to a SHA256 Hash instance. We calculate a hash of the
	// file without having to read it twice
	hasher := sha256.New()
	configFileReader := io.TeeReader(configFile, hasher)

	// Since we're reading from a file, each section is given an incremental Id to keep track of our configs for update/delete
	fileScanner := bufio.NewScanner(configFileReader)
	fileScanner.Split(hostHeaderSplitFunc)
	for fileScanner.Scan() {
		hostConfigSection := fileScanner.Bytes()
		hostConfig := newHostConfig(hostConfigSection)

		sc.hostConfigs = append(sc.hostConfigs, hostConfig)
	}
	err = fileScanner.Err()
	if err != nil {
		return err
	}

	sc.fileHash = hex.EncodeToString(hasher.Sum(nil))
	return nil
}

func (sc *SshConfig) PrintCurrentConfig() {
	for _, hostConfig := range sc.hostConfigs {
		hostConfig.printConfig()
	}
}

func (sc *SshConfig) GetAllHostNames() []string {
	names := make([]string, len(sc.hostConfigs))
	for i, hostConfig := range sc.hostConfigs {
		names[i] = string(hostConfig.name)
	}

	return names
}

type ExportedSshConfig struct {
	FileHash    string
	HostConfigs []*exportedHostConfig
}

func (sc *SshConfig) ExportFileContents() *ExportedSshConfig {
	ret := &ExportedSshConfig{
		FileHash:    sc.fileHash,
		HostConfigs: make([]*exportedHostConfig, len(sc.hostConfigs)),
	}

	for i, hostConfig := range sc.hostConfigs {
		ret.HostConfigs[i] = hostConfig.getExportedConfig()
	}

	return ret
}

func getUserSshDir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	homeDir := currentUser.HomeDir
	sshDir := path.Join(homeDir, ".ssh/")

	return sshDir, nil
}

func getSshConfigFilePath() (string, error) {
	sshDir, err := getUserSshDir()
	if err != nil {
		return "", err
	}

	configFilePath := path.Join(sshDir, "config")
	return configFilePath, nil
}

func getSshConfigFile() (*os.File, error) {
	configFilePath, err := getSshConfigFilePath()
	if err != nil {
		return nil, err
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}

	return configFile, nil
}
