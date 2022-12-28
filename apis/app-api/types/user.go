package types

import (
	"strings"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (e SignupRequest) GetInvalidFields() string {
	if !strings.Contains(e.Email, "@") {
		return "Verifique seu email"
	} else if len(strings.Split(e.Name, " ")) <= 1 {
		return "Digite seu nome completo"
	} else if len(e.Password) < 6 {
		return "Sua senha deve conter pelo menos 6 caracteres"
	}

	return ""
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPassword struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (e ResetPassword) GetInvalidFields() string {
	if e.Token == "" {
		return "Este link não está mais disponível"
	} else if len(e.Password) < 6 {
		return "Sua senha deve conter pelo menos 6 caracteres"
	}

	return ""
}
