package models

import "strconv"

func ConvertUserDTOToModel(u UserDTO) User {
	return User{
		ID:        strconv.Itoa(u.ID),
		Name:      u.Name,
		Surname:   u.Surname,
		Birthday:  u.Birthday,
		Gender:    u.Gender,
		Interests: u.Interests,
		City:      u.City,
		Login:     u.Login,
	}
}

func MustConvertUserModelToDTO(u User) UserDTO {
	id, _ := strconv.Atoi(u.ID)
	return UserDTO{
		ID:        id,
		Name:      u.Name,
		Surname:   u.Surname,
		Birthday:  u.Birthday,
		Gender:    u.Gender,
		Interests: u.Interests,
		City:      u.City,
		Login:     u.Login,
	}
}
