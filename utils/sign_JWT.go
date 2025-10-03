package utils

import (
	// "time"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	//

	// return &tokenString, err
	return &jwtToken, err
}

type Payload struct {
	Handle   string `json:"handle"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
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

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%s/api/login", config.GetGlobalConfig().App.Port), body)
	if err != nil {
		// handle err

		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		return "", err

	}
	defer resp.Body.Close()

	fmt.Printf("resp data: %+v\n", resp.Body)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res loginResponse
	if err = json.Unmarshal(respBody, &res); err != nil {
		return "", fmt.Errorf("failed to umarshal: %s", err)
	}
	if res.Token == "" {
		return "", fmt.Errorf("failed to login")
	}

	return res.Token, nil
}
