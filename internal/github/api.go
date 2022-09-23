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
	ProtoBuilderTempDir = "ProtoPort"
)

func GithubGetContent(url string, subpath string, pat string, publicRepo bool) ([]types.GithubContentApi, error) {
	var contents []types.GithubContentApi

	if subpath != "" {
		url = strings.Replace(url, "https://github.com/", "https://api.github.com/repos/", 1) + "/contents/" + subpath
	} else {
		url = strings.Replace(url, "https://github.com/", "https://api.github.com/repos/", 1) + "/contents"
	}

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
				filePath = filepath.Join(os.TempDir(), ProtoBuilderTempDir, subPath, content.Name)
			} else {
				filePath = filepath.Join(os.TempDir(), ProtoBuilderTempDir, content.Name)
			}

			file, err := os.Create(filePath)

			if err != nil {
				log.Println(err)
				return "", err
			}

			defer os.Remove(file.Name())

			_, err = file.Write(body)

			if err != nil {
				log.Println(err)
				return "", err
			}
		}
	}

	return filepath.Join(os.TempDir(), ProtoBuilderTempDir), nil
}

func generateFolderStructure(root string) {
	paths := strings.Split(root, "/")

	for i := 0; i < len(paths); i++ {
		folder := filepath.Join(os.TempDir(), ProtoBuilderTempDir, filepath.Join(paths[:i+1]...))

		if _, err := os.Stat(folder); os.IsNotExist(err) {
			os.Mkdir(folder, 0755)
		}
	}
}
