package data

type User struct {
	Id        int32  `json:"id"`
	Equipment string `json:"equipment"`
	Currency  int32  `json:"currency"`
}

// User list is kept server side NOT in blockchain
type Users struct {
	Users map[int32]User `json:"users"` // Maps user id to user
}

// If user has enough currency, can lower their currency.
func (users *Users) HasEnoughCurrency(userId, price int32) bool {
	user := users.Users[userId]
	return user.Currency-price >= 0
}
