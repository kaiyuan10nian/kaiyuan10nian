package dto

import "kaiyuan10nian/model"

type UserDto struct {
	ID  uint `json:"id"`
	Name string `json:"name"`
	Telephone string `json:"telephone"`
}

func ToUserDto(user model.User) UserDto  {
	return UserDto{
		ID:user.ID,
		Name: user.Name,
		Telephone: user.Mobile,
	}
}
