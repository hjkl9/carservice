package data

// user public info.
type UserPublicInfo struct {
	Username  string `json:"username"`
	AvatarUrl string `json:"avatarUrl"`
}

type UserID struct {
	ID uint `db:"id"`
}

type CreateUser struct {
	PhoneNumber string
	Username    string
}
