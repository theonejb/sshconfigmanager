/*
	This file contains code to write ssh config files
*/
package sshconfigmanager

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

/*
	Updates the users ssh config file with contents from the sc SshConfig struct.
	Creates a backup file before overwriting the current file
*/
func (sc *SshConfig) UpdateSshConfigFile() error {
	backupFilePath, err := sc.createBackupFile()
	if err != nil {
		return err
	}

	if err := sc.wirteToFile(nil); err != nil {
		// TODO: Copy the backup back to the main file on error
		return err
	}

	sshFilePath, err := getSshConfigFilePath()
	if err != nil {
		return err
	}

	return nil
}

func (sc *SshConfig) wirteToFile(file *os.File) error {
	return nil
}

/*
createBackupFile copies the users existing .ssh/config file to a backup file
*/
func (sc *SshConfig) createBackupFile() (string, error) {
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
				return "", err
			}
		} else {
			return "", err
		}
	}

	sshConfigFile, err := getSshConfigFile()
	if err != nil {
		return "", err
	}

	currentTime := time.Now()
	backupFileName := fmt.Sprintf("config.%d-%d-%d-%d-%d.backup", currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), currentTime.Minute())

	backupFilePath := path.Join(backupDirPath, backupFileName)
	backupFile, err := os.OpenFile(backupFilePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(backupFile, sshConfigFile)
	if err != nil {
		return "", err
	}
	if err = backupFile.Close(); err != nil {
		return "", err
	}

	return backupFilePath, nil
}
