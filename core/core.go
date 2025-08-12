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

func GetContainerLogs(host, id string) {
	response, err := http.Get(host + "/containers/" + id + "/logs?stdout=true&stderr=true")
	if err != nil {
		log.Fatalln(err)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(data))
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
		return "", errors.New(http.StatusText(response.StatusCode))
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

	log.Println("Successfully created container" + response.Status + m["Id"])

	return m["Id"], nil
}

func CreateExecInstance(host, id, cmd string) (string, error) {
	j, err := json.Marshal(map[string]any{
		"AttachStdin":  false,
		"AttachStdout": true,
		"AttachStderr": true,
		"Tty":          false,
		"Cmd":          cmd,
	})
	if err != nil {
		log.Fatalln(err)
	}

	response, err := http.Post(host+"/containers/"+id+"/exec", "application/json", bytes.NewBuffer(j))
	if err != nil {
		log.Fatalln(err)
	}

	switch response.StatusCode {
	case 404, 500, 409:
		return "", errors.New(http.StatusText(response.StatusCode))
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var m map[string]any
	err = json.Unmarshal(data, &m)
	if err != nil {
		log.Fatalln(err)
	}

	return m["Id"].(string), nil
}

func StartExec(host, id string) error {
	j, err := json.Marshal(map[string]any{
		"Detach": false,
		"Tty":    false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	response, err := http.Post(host+"/exec/"+id+"/start", "application/json", bytes.NewBuffer(j))
	if err != nil {
		log.Fatalln(err)
	}

	switch response.StatusCode {
	case 404, 409:
		return errors.New(http.StatusText(response.StatusCode))
	}

	return nil
}
