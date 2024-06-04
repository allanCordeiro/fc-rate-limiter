package checker

import (
	"encoding/json"
	"io"
	"os"
)

type Token struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"`
}

func GetTokenExpirationParam(tokenKey string) (int, error) {
	var tokenList []Token
	data, err := getFileData()
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(data, &tokenList)
	if err != nil {
		return 0, nil
	}

	token := getByKey(tokenKey, tokenList)
	if token != nil {
		return token.ExpiresIn, nil
	}
	return 0, nil
}

func getByKey(key string, token []Token) *Token {
	for _, v := range token {
		if v.Token == key {
			return &Token{
				Token:     v.Token,
				ExpiresIn: v.ExpiresIn,
			}
		}
	}
	return &Token{}
}

func getFileData() ([]byte, error) {
	file, err := os.Open("tokens.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}
