package madl

import (
	"gmidarch/shared/parameters"
	"os"
	"bufio"
	"errors"
	"gmidarch/shared/shared"
)

type MADLFile struct {
	FilePath string
	FileName string
	Content  []string
}

func (m *MADLFile) Read(fileName string) {

	// Check file name
	err := m.CheckFileName(fileName)
	shared.CheckError(err, "MADL")

	// configure r
	m.FileName = fileName
	m.FilePath = parameters.DIR_CONF

	fullPathAdlFileName := m.FilePath + "/" + m.FileName

	// read file
	fileContent := []string{}
	file, err := os.Open(fullPathAdlFileName)
	shared.CheckError(err, "MADL")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}

	// configure r
	m.Content = fileContent
}

func (MADLFile) CheckFileName(fileName string) error {
	r := *new(error)

	len := len(fileName)

	if len <= 5 {
		r = errors.New("File Name Invalid")
	} else {
		if fileName[len-5:] != parameters.MADL_EXTENSION {
			r = errors.New("Invalid extension of '"+fileName+"'")
		} else {
			r = nil
		}
	}

	return r
}
