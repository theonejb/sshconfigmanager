/*
Holds the structures and functions used to read, create, and save Host section configuration
in a .ssh/config file
*/
package sshconfigmanager

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"
)

/*
Holds the config for a Host section. Fields are:
- hash: A hash of the contents of this host config used as an Id
- name: The name given with the Host header
- hostName: If given, the actual hostname which this Host config connects to
- port: Port
- user: The user name to use to connect to this Host
- identityFile: The ssh ident file to use when connecting
- otherLines: Any other lines that we don't parse. We need these because we don't parse the whole
	Host section from a config file, only a select few fields.
- indentString: The string used to indent the config lines when printed
*/
type hostConfig struct {
	hash         []byte
	name         []byte
	hostName     []byte
	port         []byte
	user         []byte
	identityFile []byte
	otherLines   [][]byte

	indentString string
}

func newHostConfig(data []byte) *hostConfig {
	hc := &hostConfig{
		otherLines:   make([][]byte, 0),
		indentString: "  ",
	}

	byteReader := bytes.NewReader(data)
	lineSplitter := bufio.NewScanner(byteReader)
	lineSplitter.Split(bufio.ScanLines)

	for lineSplitter.Scan() {
		line := lineSplitter.Bytes()
		hc.addLineToConfig(line)
	}

	if lineSplitter.Err() != nil {
		return nil
	} else {
		hc.generateHash()
		return hc
	}
}

/*
Given a line of text, sees if it is one our parsable parameters. If it is, parses it and stores it in the
appropriate field. If not, adds it to the otherLines field.

Also checks that the line is not empty or a comment
*/
func (hc *hostConfig) addLineToConfig(line []byte) {
	trimmedLine := bytes.TrimSpace(line)
	if len(trimmedLine) == 0 {
		return
	}

	hostHeader := []byte("host ")
	hostnameHeader := []byte("hostname ")
	portHeader := []byte("port ")
	userHeader := []byte("user ")
	identHeader := []byte("identityfile ")

	lowerCasedLine := bytes.ToLower(trimmedLine)
	if bytes.Index(lowerCasedLine, hostHeader) == 0 {
		name := bytes.TrimSpace(trimmedLine[len(hostHeader):])
		hc.name = name
	} else if bytes.Index(lowerCasedLine, hostnameHeader) == 0 {
		hostname := bytes.TrimSpace(trimmedLine[len(hostnameHeader):])
		hc.hostName = hostname
	} else if bytes.Index(lowerCasedLine, portHeader) == 0 {
		port := bytes.TrimSpace(trimmedLine[len(portHeader):])
		hc.port = port
	} else if bytes.Index(lowerCasedLine, userHeader) == 0 {
		user := bytes.TrimSpace(trimmedLine[len(userHeader):])
		hc.user = user
	} else if bytes.Index(lowerCasedLine, identHeader) == 0 {
		identFile := bytes.TrimSpace(trimmedLine[len(identHeader):])
		hc.identityFile = identFile
	} else if lowerCasedLine[0] == '#' {
		// Do nothing. The passed in line is a comment
	} else {
		hc.otherLines = append(hc.otherLines, trimmedLine)
	}
}

func (hc *hostConfig) generateHash() {
	// Concat the contents of the host config together so we can hash them. This hash is used later
	// to verify that the users of the API can only update a current version of the HostConfig. We
	// require a hash to verify this because ssh/config files don't have implicit ids
	var contents bytes.Buffer
	contents.Write(hc.name)
	contents.Write(hc.hostName)
	contents.Write(hc.port)
	contents.Write(hc.user)
	contents.Write(hc.identityFile)
	for _, line := range hc.otherLines {
		contents.Write(line)
	}

	// Create a hash of the contents
	var hasher hash.Hash = sha256.New()
	_, err := contents.WriteTo(hasher)
	if err != nil {
		panic("Unable to generate hash of HostConfig. Error: " + err.Error())
	}
	var hash_sum []byte = hasher.Sum(nil)
	hc.hash = hash_sum
}

func (hc *hostConfig) sprintConfig() string {
	stringParts := make([]string, 0, 5)

	hostHeaderString := fmt.Sprintf("Host %s", string(hc.name))
	stringParts = append(stringParts, hostHeaderString)

	if hc.hostName != nil {
		hostnameString := fmt.Sprintf("%sHostName %s", hc.indentString, string(hc.hostName))
		stringParts = append(stringParts, hostnameString)
	}

	if hc.port != nil {
		portString := fmt.Sprintf("%sPort %s", hc.indentString, string(hc.port))
		stringParts = append(stringParts, portString)
	}

	if hc.user != nil {
		userString := fmt.Sprintf("%sUser %s", hc.indentString, string(hc.user))
		stringParts = append(stringParts, userString)
	}

	if hc.identityFile != nil {
		identString := fmt.Sprintf("%sIdentityFile %s", hc.indentString, string(hc.identityFile))
		stringParts = append(stringParts, identString)
	}

	if hc.otherLines != nil && len(hc.otherLines) > 0 {
		for _, l := range hc.otherLines {
			lString := fmt.Sprintf("%s%s", hc.indentString, string(l))
			stringParts = append(stringParts, lString)
		}
	}

	return strings.Join(stringParts, "\n")
}

func (hc *hostConfig) printConfig() {
	fmt.Printf("$$$$$$$$$$ %s $$$$$$$$$$\n", hc.name)
	fmt.Println(hc.sprintConfig())
	fmt.Println("$$$$$$$$$$$$$$$$$$$$\n")
}

/*
A struct to hold exported configuration for a Host Config. For now, it's an almost exact
copy of hostConfig. But we use it so that if the hostConfig struct changes, the external
API remains the same
*/
type exportedHostConfig struct {
	Id           string
	Name         string
	HostName     string
	Port         string
	User         string
	IdentityFile string
	OtherLines   []string
}

func (hc *hostConfig) getExportedConfig() *exportedHostConfig {
	ret := &exportedHostConfig{
		Id:           hex.EncodeToString(hc.hash),
		Name:         string(hc.name),
		HostName:     string(hc.hostName),
		Port:         string(hc.port),
		User:         string(hc.user),
		IdentityFile: string(hc.identityFile),
		OtherLines:   make([]string, len(hc.otherLines)),
	}

	for i, line := range hc.otherLines {
		ret.OtherLines[i] = string(line)
	}

	return ret
}
