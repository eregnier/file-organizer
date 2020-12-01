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
	Folder     string     `json:"folder"`
	typeMap    map[string]FileType
	FoldersSet map[string]bool
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
	targetFolder := fmt.Sprintf("%s/folders", config.Folder)
	if _, err := os.Stat(targetFolder); os.IsNotExist(err) {
		os.Mkdir(targetFolder, 0755)
	}

	for _, file := range files {
		name := file.Name()
		if !file.IsDir() {
			fmt.Println("processing", name)

			var ext string
			if strings.Contains(name, ".") {
				tokens := strings.Split(name, ".")
				ext = tokens[len(tokens)-1]
			} else {
				ext = "noext"
			}
			moveFile(config, name, ext)
		} else {
			moveFolder(config, name)
		}
	}
}

func moveFolder(config Config, name string) {
	if name != "folders" && config.FoldersSet[name] != true {
		oldPath := fmt.Sprintf("%s/%s", config.Folder, name)
		newPath := fmt.Sprintf("%s/folders/%s", config.Folder, name)
		err := os.Rename(oldPath, newPath)
		if err != nil {
			log.Fatal(err)
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
	config.FoldersSet = make(map[string]bool)
	for _, ft := range config.FilesTypes {
		config.typeMap[ft.Extension] = ft
		config.FoldersSet[ft.Folder] = true
	}
	return config
}
