package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func pullImage(host string, image string) {
	response, err := http.Post(host+"/images/create?fromImage="+image, "application/json", nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(data))
}

func CreateContainer(host string, image string) (string, error) {
	data := map[string]string{
		"Image": image,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}

	response, err := http.Post(host+"/containers/create", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case 404:
		pullImage(host, image)
	case 400, 409, 500:
		return "", errors.New("something went wrong")
	}

	respData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]string
	err = json.Unmarshal(respData, &result)
	if err != nil {
		log.Fatalln(err)
	}

	return result["Id"], nil
}
