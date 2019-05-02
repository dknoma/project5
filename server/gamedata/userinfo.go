package gamedata

type User struct {
	Id        int32       `json:"id"`
	Equipment []Equipment `json:"inventory"`
	Currency  int32       `json:"currency"`
}

// User list is kept server side NOT in blockchain
type Users struct {
	Users map[int32]User `json:"users"` // Maps user id to user
}

var nextWeaponId = int32(0)

func (users *Users) InitUserList() {
	users.Users = make(map[int32]User)
}

func (user *User) GenerateEquipment() {
	totalWeapons := 10
	atk := int32(5)
	def := int32(5)
	for i := 0; i < totalWeapons; i++ {
		user.Equipment = append(user.Equipment, Equipment{"Sword", nextWeaponId,
			user.Id, "This is a basic sword.", EquipmentStats{1, atk, def}})
		atk++
		def++
		nextWeaponId++
	}
}

// If user has enough currency, can lower their currency.
func (users *Users) HasEnoughCurrency(userId, price int32) bool {
	user := users.Users[userId]
	return user.Currency-price >= 0
}

// Validate that the user of the given id actually exists
func (users *Users) UserExists(userId int32) bool {
	_, exists := users.Users[userId]
	return exists
}

// Validate if the user actually owns the item
func (user *User) UserHasItem(equipment Equipment) bool {
	for _, item := range user.Equipment {
		if item == equipment {
			return true
		}
	}
	return false
}
