package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

//FileType : Files types
type FileType struct {
	Extension string `json:"extension"`
	Folder    string `json:"folder"`
}

//Config : General configuration for this tool
type Config struct {
	FilesTypes []FileType `json:"files_types"`
	typeMap    map[string]FileType
	Folder     string `json:"folder"`
}

func main() {
	usr, _ := user.Current()
	configFilePath := fmt.Sprintf("%s/.file-organizer.json", usr.HomeDir)

	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	var config Config

	json.Unmarshal(data, &config)

	config = parseFileType(config)
	order(config)
}

func order(config Config) {
	files, err := ioutil.ReadDir(config.Folder)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			fmt.Println("processing", name)

			var ext string
			if strings.Contains(name, ".") {
				tokens := strings.Split(name, ".")
				ext = tokens[len(tokens)-1]
			} else {
				ext = "noext"
			}
			moveFile(config, name, ext)
		}
	}
}

func moveFile(config Config, name string, ext string) {
	folder := config.typeMap[ext].Folder
	if folder == "" && config.typeMap["unknown"].Folder != "" {
		folder = config.typeMap["unknown"].Folder
	}
	if folder != "" {
		oldPath := fmt.Sprintf("%s/%s", config.Folder, name)
		newPath := fmt.Sprintf("%s/%s/%s", config.Folder, folder, name)
		targetFolder := fmt.Sprintf("%s/%s", config.Folder, folder)
		if _, err := os.Stat(targetFolder); os.IsNotExist(err) {
			os.Mkdir(targetFolder, 0755)
		}
		err := os.Rename(oldPath, newPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func parseFileType(config Config) Config {
	config.typeMap = make(map[string]FileType)
	for _, ft := range config.FilesTypes {
		config.typeMap[ft.Extension] = ft
	}
	return config
}
