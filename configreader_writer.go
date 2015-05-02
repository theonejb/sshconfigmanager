/*
This file contains code to write ssh config files
*/
package sshconfigmanager

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

/*
	Updates the users ssh config file with contents from the sc SshConfig struct.
	Creates a backup file before overwriting the current file
*/
func (sc *SshConfig) UpdateSshConfigFile(oldFileHash string) (string, error) {
	currentConfig := NewSshConfig()
	if err := currentConfig.ReadConfig(); err != nil {
		return "", err
	}

	backupFilePath, err := sc.createBackupFile()
	if err != nil {
		return "", err
	}

	if err := sc.wirteToFile(nil); err != nil {
		// TODO: Copy the backup back to the main file on error
		return "", err
	}

	return backupFilePath, nil
}

/*
Return a string with the contents of the SshConfig that can be used to create a new
config file
*/
func (sc *SshConfig) writeToString() string {
	hostConfigStrings := make([]string, len(sc.hostConfigs))

	for i, hostConfig := range sc.hostConfigs {
		hostConfigStrings[i] = hostConfig.sprintConfig()
	}

	return strings.Join(hostConfigStrings, "\n\n")
}

func (sc *SshConfig) wirteToFile(file *os.File) error {
	newFileContents := sc.writeToString()
	fmt.Println(newFileContents)
	return nil
}

/*
createBackupFile copies the users existing .ssh/config file to a backup file. Returns the new file path
*/
func (sc *SshConfig) createBackupFile() (string, error) {
	sshConfigFilePath, err := getSshConfigFilePath()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error in getSshConfigFilePath(): %s", err.Error()))
	}

	backupDirPath, err := ensureBackupDirExists()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error in ensureBackupDirExists(): %s", err.Error()))
	}

	currentTime := time.Now()
	backupFileName := fmt.Sprintf("config.%d-%d-%d-%d.backup", currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Unix())

	backupFilePath := path.Join(backupDirPath, backupFileName)
	err = copyFile(backupFilePath, sshConfigFilePath, false)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error while copying current ssh file to backup file: %s", err.Error()))
	}

	return backupFilePath, nil
}

/*
ensureBackupDirExists checks if the dir where we write our backup files exists. If it
doesn't it is created. Return the path to the dir
*/
func ensureBackupDirExists() (string, error) {
	sshDir, err := getUserSshDir()
	if err != nil {
		return "", err
	}

	backupDirPath := path.Join(sshDir, "sshconfigmanager_backups")
	_, err = os.Stat(backupDirPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Since the dir doesn't exist, create it
			if err = os.Mkdir(backupDirPath, os.ModeDir|0700); err != nil {
				return backupDirPath, err
			}
		}

		return "", err
	}

	return backupDirPath, nil
}

/*
copyFile copies the given src file to the dest file. It can override the dest file,
or return an err if the dest file already exists.
*/
func copyFile(dstFilePath string, srcFilePath string, overwrite bool) error {
	var fileFlags int
	if overwrite {
		fileFlags = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	} else {
		fileFlags = os.O_CREATE | os.O_EXCL | os.O_WRONLY
	}

	dstFile, err := os.OpenFile(dstFilePath, fileFlags, 0600)
	if err != nil {
		return err
	}

	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	if err = srcFile.Close(); err != nil {
		return err
	}

	if err = dstFile.Close(); err != nil {
		return err
	}

	return nil
}
