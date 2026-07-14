package db

type User struct {
	ID       int
	Nickname string
	Email    string
	Password string
}

func GetUserByIdentifier(identifier string) (*User, error) {
	var user User
	err := DB.QueryRow("SELECT id, Nickname, Password FROM users WHERE Nickname = ? or Email = ?", identifier, identifier).Scan(&user.ID, &user.Nickname, &user.Password)
	if err != nil{
		return nil, err
	}
	return &user, nil
}
