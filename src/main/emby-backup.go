package main

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type configuration struct {
	Extensions  string `yaml:"video-extensions"`
	Logo        string `yaml:"logo"`
	BasePath    string `yaml:"base-path"`
	ArchivePath string `yaml:"archive-path"`
}

func (c *configuration) initialize() *configuration {
	yamlFile, err := ioutil.ReadFile("configuration.yaml")
	if err != nil {
		log.Println("A problem occurred reading configuration: #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalln("A problem occurred unmarshalling configuration: %v", err)
	}
	return c
}

func IsVideo(file string, conf configuration) bool {
	var video = false
	var extension = filepath.Ext(file)
	var videoExtensions = strings.Split(conf.Extensions, ",")
	for i := 0; i < len(videoExtensions); i++ {
		if extension == videoExtensions[i] {
			video = true
		}
	}
	return video
}

func main() {
	var err error

	// Initialize Configuration
	var conf configuration
	conf.initialize()

	var zipFile *os.File
	zipFile, err = os.Create(conf.ArchivePath + "\\emby-backup.zip")
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk the directory structure of the configured base path
	err = filepath.Walk(conf.BasePath,
		func(sourcePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() || IsVideo(sourcePath, conf) {
				return nil
			}

			if err != nil {
				return err
			}
			relativePath := strings.TrimPrefix(sourcePath, filepath.Dir(conf.BasePath))
			zipFile, err := zipWriter.Create(relativePath)
			if err != nil {
				return err
			}
			fileToZip, err := os.Open(sourcePath)
			if err != nil {
				return err
			}
			_, err = io.Copy(zipFile, fileToZip)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}
