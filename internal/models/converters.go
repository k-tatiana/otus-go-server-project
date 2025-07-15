package models

func ConvertUserDTOToModel(u UserDTO) User {
	return User{
		ID:        u.ID,
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
	return UserDTO{
		ID:        u.ID,
		Name:      u.Name,
		Surname:   u.Surname,
		Birthday:  u.Birthday,
		Gender:    u.Gender,
		Interests: u.Interests,
		City:      u.City,
		Login:     u.Login,
	}
}
