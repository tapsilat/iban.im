package utils

import (
	// "time"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tapsilat/iban.im/config"
)

// SignJWT : func to generate JWT
func SignJWT(userMail, userPass *string) (*string, error) {
	jwtToken, err := loginJwtToken(userMail, userPass)
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"userID": *userID,
	// 	"exp":    time.Now().Add(time.Second * 30 * 24 * 60 * 60),
	// })

	// tokenString, err := token.SignedString([]byte("my_secret"))
	// fmt.Println("signJWT token: ",tokenString)

	// return &tokenString, err
	return &jwtToken, err
}

type Payload struct {
	Handle   string `json:"handle"`
	Password string `json:"password"`
}

func loginJwtToken(userMail, userPass *string) (string, error) {

	data := Payload{
		Handle:   *userMail,
		Password: *userPass,
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
		return "", err
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/api/login", config.Config.App.Port), body)
	if err != nil {
		// handle err
		fmt.Println("error in req: ", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		fmt.Println("error in resp: ", err)
		return "", err

	}
	fmt.Printf("resp data: %+v\n", resp.Body)
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error in resp body read: ", err)
		return "", err
	}
	fmt.Println("resp: ", string(respbody))
	m := make(map[string]interface{})
	errUnMarshal := json.Unmarshal(respbody, &m)
	if errUnMarshal != nil {
		fmt.Println("error in resp body unmarshal: ", errUnMarshal)
		return "", errUnMarshal
	}
	fmt.Println("token: ", m["token"])

	defer resp.Body.Close()
	token := fmt.Sprintf("%s", m["token"])
	return token, err
}
