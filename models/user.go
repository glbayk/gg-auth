package models

type User struct {
	BaseModel
	Email    string `gorm:"unique_index;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
}

func (user *User) Create() error {
	err := DB.Create(&user).Error
	i := 0

	for err != nil && i < 5 {
		err = DB.Create(&user).Error
		i++
	}

	return err
}

func (user *User) Find() error {
	err := DB.Where("email = ?", user.Email).First(&user).Error
	return err
}

func (user *User) Update() error {
	err := DB.Save(&user).Error
	return err
}
