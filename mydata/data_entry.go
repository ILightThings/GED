package mydata

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"github.com/ilightthings/GED/typelib"
)

type jsonFile struct {
	filename     string
	filecontents []byte
}

//TODO Optomize this function somehow
func DatabaseToZip(entries typelib.PageEntries) ([]byte, error) {
	buf := new(bytes.Buffer)

	zipWrite := zip.NewWriter(buf)
	var fileArray []jsonFile

	credfile, err := json.Marshal(entries.CredEntries)
	if err != nil {
		return nil, err
	}
	credjson := jsonFile{
		filename:     "creds.json",
		filecontents: credfile,
	}

	hostfile, err := json.Marshal(entries.HostEntries)
	if err != nil {
		return nil, err
	}
	hostjson := jsonFile{
		filename:     "hosts.json",
		filecontents: hostfile,
	}

	cmdfile, err := json.Marshal(entries.CommandList)
	if err != nil {
		return nil, err
	}
	cmdjson := jsonFile{
		filename:     "cmds.json",
		filecontents: cmdfile,
	}
	fileArray = append(fileArray, credjson)
	fileArray = append(fileArray, hostjson)
	fileArray = append(fileArray, cmdjson)

	for _, file := range fileArray {
		zipFile, err := zipWrite.Create(file.filename)
		if err != nil {
			return nil, err
		}
		_, err = zipFile.Write(file.filecontents)
		if err != nil {
			return nil, err
		}

	}
	err = zipWrite.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil

}
