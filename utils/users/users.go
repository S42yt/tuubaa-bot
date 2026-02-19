package users

type DiscordUser struct {
	UserID      string
	UserName 	string
}

var Users = []DiscordUser{
	{UserID: "787306646417571860", UserName: "s42."},
	{UserID: "286765460140326913", UserName: "tsundosika1"},
	{UserID: "624623721587408896", UserName: "tuubaa"},
	{UserID: "795306274467348510", UserName: "time_root"},
}

func GetUserByID(userID string) *DiscordUser {
	for i := range Users {
		if Users[i].UserID == userID {
			return &Users[i]
		}
	}
	return nil
}

func GetUserByUserName(userName string) *DiscordUser {
	for i := range Users {
		if Users[i].UserName == userName {
			return &Users[i]
		}
	}
	return nil
}

func AddUser(userID, userName	 string) {
	Users = append(Users, DiscordUser{
		UserID:      userID,
		UserName:    userName,
	})
}
