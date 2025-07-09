package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrFileAlreadyExists = errors.New("file already exist")
)

func getTableFile(fileName string) (f *os.File, err error) {
	// trim `_` character from fileName
	fileName = strings.Replace(fileName, "_", "", -1)
	if fileName == "" {
		return f, errors.New("error: fileName can not be empty")
	}
	fileName += ".go"
	filePath := filepath.Join(outputDir, fileName)
	if _, err = os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return os.Create(filePath)
		}
		return
	}
	fmt.Printf("file %s already exist, do you want to overwrite it? [y/n]: ", fileName)
	var op string
	fmt.Scanf("%s", &op)
	if strings.ToLower(op) == "y" {
		return os.Create(filePath)
	}
	return f, ErrFileAlreadyExists
}

func getTableCondsFile(fileName string) (f *os.File, err error) {
	// trim `_` character from fileName
	fileName = strings.Replace(fileName, "_", "", -1)
	if fileName == "" {
		return f, errors.New("error: fileName can not be empty")
	}
	fileName += "conds.go"
	filePath := filepath.Join(outputDir, fileName)
	if _, err = os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return os.Create(filePath)
		}
		return
	}
	fmt.Printf("file %s already exist, do you want to overwrite it? [y/n]: ", fileName)
	var op string
	fmt.Scanf("%s", &op)
	if strings.ToLower(op) == "y" {
		return os.Create(filePath)
	}
	return f, ErrFileAlreadyExists
}

func getInitDaoFile() (f *os.File, err error) {
	fileName := "dao.go"
	filePath := filepath.Join(outputDir, fileName)
	if _, err = os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return os.Create(filePath)
		}
		return
	}
	fmt.Printf("file %s already exist, do you want to overwrite it? [y/n]: ", fileName)
	var op string
	fmt.Scanf("%s", &op)
	if strings.ToLower(op) == "y" {
		return os.Create(filePath)
	}
	return f, ErrFileAlreadyExists
}
