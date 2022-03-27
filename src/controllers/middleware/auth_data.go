package middleware

import "github.com/pkg/errors"

const CREATE_NOTE_PRIVILEGE = "CREATE_NOTE"
// TODO - метод получения истории изменения Note
const VIEW_NOTE_VERSION_HISTORY_PRIVILEGE = "VIEW_NOTE_VERSION_HISTORY"
const CHANGE_NOTE_VERSION_PRIVILEGE = "CHANGE_NOTE_VERSION"

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

func (data UserAuthData) HasPrivilege(privilege string) bool {
	for i:=0; i < len(data.Privileges); i++ {
		if (data.Privileges[i] == privilege) {
			return true
		}
	}
	return false
}