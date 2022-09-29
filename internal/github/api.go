package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/TechMDW/ProtoPort/internal/github/types"
	"github.com/TechMDW/ProtoPort/internal/utilities"
)

const (
	ProtoPortDir = "ProtoPort"
)

func GithubUrlParser(url string, subpath string) string {

	url = strings.Replace(url, "/tree/main", "/contents", 1)

	if !strings.Contains(url, "/contents") {
		url = url + "/contents"
	}

	if subpath != "" {
		url = strings.Replace(url, "https://github.com/", "https://api.github.com/repos/", 1) + subpath
	} else {
		url = strings.Replace(url, "https://github.com/", "https://api.github.com/repos/", 1)
	}

	return url
}

func GithubGetContent(url string, subpath string, pat string, publicRepo bool) ([]types.GithubContentApi, error) {
	var contents []types.GithubContentApi

	url = GithubUrlParser(url, subpath)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return contents, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")

	if !publicRepo {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pat))
	}

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		return contents, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return contents, err
	}

	err = json.Unmarshal(body, &contents)
	if err != nil {
		log.Println(err)
		return contents, err
	}

	return contents, nil
}

func GithubReadAndGenrateProtos(githubUrl string, subPath string, pat string, publicRepo bool) (string, error) {
	var contents []types.GithubContentApi
	if subPath != "" {
		subContents, err := GithubGetContent(githubUrl, subPath, pat, publicRepo)
		if err != nil {
			log.Println(err)
			return "", err
		}
		contents = subContents
	} else {
		firstCall, err := GithubGetContent(githubUrl, "", pat, publicRepo)
		if err != nil {
			log.Println(err)
			return "", err
		}
		contents = firstCall
	}

	configPath, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	techMDWPath := filepath.Join(configPath, "TechMDW")

	if _, err := os.Stat(techMDWPath); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(techMDWPath, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	ProtoPortPath := filepath.Join(techMDWPath, ProtoPortDir)

	if _, err := os.Stat(ProtoPortPath); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(ProtoPortPath, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, content := range contents {
		if content.Type == "dir" {
			if subPath != "" {
				GithubReadAndGenrateProtos(githubUrl, subPath+"/"+content.Name, pat, publicRepo)
			} else {
				GithubReadAndGenrateProtos(githubUrl, content.Name, pat, publicRepo)
			}
		} else if utilities.CheckForFileExtension(content.Name, ".proto") && content.Type == "file" {
			var filePath string
			req, err := http.NewRequest("GET", content.URL, nil)

			if err != nil {
				log.Println(err)
				return "", err
			}

			req.Header.Set("Accept", "application/vnd.github.raw+json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pat))

			response, err := http.DefaultClient.Do(req)

			if err != nil {
				log.Println(err)
				return "", err
			}

			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)

			if err != nil {
				log.Println(err)
				return "", err
			}

			if subPath != "" {
				generateFolderStructure(subPath)
				filePath = filepath.Join(ProtoPortPath, subPath, content.Name)
			} else {
				filePath = filepath.Join(ProtoPortPath, content.Name)
			}

			file, err := os.Create(filePath)
			
			if err != nil {
				log.Println(err)
				return "", err
			}

			_, err = file.Write(body)

			if err != nil {
				log.Println(err)
				return "", err
			}
		}
	}

	return filepath.Join(ProtoPortPath), nil
}

func generateFolderStructure(root string) {
	paths := strings.Split(root, "/")

	configPath, err := os.UserConfigDir()

	if err != nil {
		log.Println(err)
	}

	techMDWPath := filepath.Join(configPath, "TechMDW")

	ProtoPortPath := filepath.Join(techMDWPath, ProtoPortDir)

	for i := 0; i < len(paths); i++ {
		folder := filepath.Join(ProtoPortPath, filepath.Join(paths[:i+1]...))

		if _, err := os.Stat(folder); os.IsNotExist(err) {
			os.Mkdir(folder, 0755)
		}
	}
}
