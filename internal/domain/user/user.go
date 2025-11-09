package user

import "time"

func New(code string, email string, profileImageURL *string, username *string) *User {
	return &User{
		Code:            code,
		Email:           email,
		ProfileImageURL: profileImageURL,
		Username:        username,
		SignedUpAt:      time.Now(),
	}
}

// TODO split into some value objects
func From(id int64, code string, email string, profileImageURL *string, username *string, signedUpAt time.Time) *User {
	return &User{
		Id:              id,
		Code:            code,
		Email:           email,
		ProfileImageURL: profileImageURL,
		Username:        username,
		SignedUpAt:      signedUpAt,
	}
}

type User struct {
	Id int64

	Code string

	Email string

	ProfileImageURL *string
	Username        *string

	SignedUpAt time.Time
}

func (u *User) IsPersisted() bool {
	var intZeroValue int64
	isIdIntZeroValue := u.Id == intZeroValue
	return !isIdIntZeroValue
}
