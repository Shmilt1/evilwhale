package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func pullImage(host, image string) {
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

func CreateContainer(host, image string) (string, error) {
	jsonData, err := json.Marshal(map[string]string{
		"Image": image,
	})
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
		return "", errors.New("error code " + response.Status)
	}

	respData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var jsonResp map[string]string
	err = json.Unmarshal(respData, &jsonResp)
	if err != nil {
		log.Fatalln(err)
	}

	return jsonResp["Id"], nil
}

func Exec(host, id string) {

}
