package dto

import "kaiyuan10nian/model"

type UserDto struct {
	Name string `json:"name"`
	Telephone string `json:"telephone"`
}

func ToUserDto(user model.User) UserDto  {
	return UserDto{
		Name: user.Name,
		Telephone: user.Mobile,
	}
}
