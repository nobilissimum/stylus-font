package main

import (
	"io"
	"log"
	"os"
)

const (
	sourceDirectory = "src"
	distributable   = "dist.scss"
	global          = "global.scss"
)

var directories = []string{"extensions", "websites"}

func main() {
	inputFilename := sourceDirectory + "/" + global
	inputFile, err := os.Open(inputFilename)
	if err != nil {
		log.Panicf("Unable to open global style file %v. %v", inputFilename, err.Error())
	}
	defer func() {
		if err := inputFile.Close(); err != nil {
			log.Panicf("Unable to close global style file %v. %v", inputFilename, err.Error())
		}
	}()

	outputFile, err := os.Create(distributable)
	if err != nil {
		log.Panicf("Unable to create output file %v. %v", distributable, err.Error())
	}
	defer func() {
		if err := outputFile.Close(); err != nil {
			log.Panicf("Unable to close output file %v. %v", distributable, err.Error())
		}
	}()

	buf := make([]byte, 1024)
	for true {
		bytesNum, err := inputFile.Read(buf)
		if err != nil && err != io.EOF {
			log.Panicf("Unable to read from input file. %v", err.Error())
		}

		if bytesNum == 0 {
			break
		}

		if _, err := outputFile.Write(buf[:bytesNum]); err != nil {
			log.Panicf("Unable to write to output file. %v", err.Error())
		}

		if _, err := outputFile.WriteString("\n"); err != nil {
			log.Panicf("Unable to write a blank line to output file. %v", err.Error())
		}
	}

	for _, directory := range directories {
		directoryPath := sourceDirectory + "/" + directory
		files, err := os.ReadDir(directoryPath)
		if err != nil {
			log.Panicf("Unable to read directory %v, %v", directoryPath, err.Error())
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			styleFilename := directoryPath + "/" + file.Name()
			styleFile, err := os.Open(styleFilename)
			if err != nil {
				log.Panicf("Unable to open style file %v. %v", styleFilename, err.Error())
			}
			defer func() {
				if err := styleFile.Close(); err != nil {
					log.Panicf("Unable to close input file %v. %v", styleFilename, err.Error())
				}
			}()

			for true {
				bytesNum, err := styleFile.Read(buf)
				if err != nil && err != io.EOF {
					log.Panicf("Unable to read %v. %v", styleFilename, err.Error())
				}

				if bytesNum == 0 {
					break
				}

				if _, err := outputFile.Write(buf[:bytesNum]); err != nil {
					log.Panicf("Unable to write %v to output file. %v", styleFilename, err.Error())
				}

				if _, err := outputFile.WriteString("\n"); err != nil {
					log.Panicf("Unable to write a blank line to output file. %v", err.Error())
				}
			}
		}
	}
}
