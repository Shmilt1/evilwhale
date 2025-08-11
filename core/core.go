package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func Ping(host string) bool {
	response, err := http.Get(host + "/_ping")
	if err != nil {
		log.Fatalln(err)
	}

	if response.StatusCode == 200 {
		return true
	}

	return false
}

func pullImage(host, image string) {
	response, err := http.Post(host+"/images/create?fromImage="+image, "application/json", nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	for {
		var status map[string]any
		if err := decoder.Decode(&status); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalln(err)
		}

		j, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(string(j))
	}
}

func CreateContainer(host, image string) (string, error) {
	j, err := json.Marshal(map[string]string{
		"Image": image,
	})
	if err != nil {
		log.Fatalln(err)
	}

	response, err := http.Post(host+"/containers/create", "application/json", bytes.NewBuffer(j))
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

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var m map[string]string
	err = json.Unmarshal(data, &m)
	if err != nil {
		log.Fatalln(err)
	}

	return m["Id"], nil
}

func Exec(host, id string) {

}
