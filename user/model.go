package user

import "github.com/google/uuid"

type Model struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	DisplayName  string    `gorm:"size:32;not null"`
	PhoneNumber  string    `gorm:"size:12;unique;not null"`
	Password     string    `gorm:"size:72;not null"`
	RefreshToken string    `gorm:"size:512;not null"`
}

func NewModel(displayName, phoneNumber, password string) *Model {
	return &Model{
		ID:          uuid.New(),
		DisplayName: displayName,
		PhoneNumber: phoneNumber,
		Password:    password,
	}
}

func (m *Model) SetRefreshToken(refreshToken string) {
	m.RefreshToken = refreshToken
}

func (*Model) TableName() string {
	return "users"
}
