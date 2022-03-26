package middleware

import "github.com/pkg/errors"

type UserAuthData struct {
	UserId     *int64   `json:"user_id"`
	Username   *string  `json:"username"`
	Privileges []string `json:"privileges"`
}

func (data UserAuthData) Valid() error {
	if data.UserId == nil || data.Username == nil {
		return errors.New("не найдены данные пользователя в токене")
	} else {
		return nil
	}
}