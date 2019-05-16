package gamedata

import (
	"fmt"
	"sync"
)

type User struct {
	Id        int32     `json:"id"`
	Inventory Inventory `json:"inventory"`
	Currency  float64   `json:"currency"`
}

type Inventory struct {
	Equipment []Equipment `json:"equipment"`
	mux       sync.Mutex
}

// User list is kept server side NOT in blockchain
type Users struct {
	Users map[int32]User `json:"users"` // Maps user id to user
	mux   sync.Mutex
}

var nextWeaponId = int32(0)

func (users *Users) InitUserList() {
	users.Users = make(map[int32]User)
}

func (users *Users) AddUser(id int32) {
	users.mux.Lock()
	newUser := User{Id: id, Inventory: Inventory{}, Currency: 10000}
	newUser.GenerateEquipment()
	users.Users[id] = newUser
	users.mux.Unlock()
}

// TODO: future functionality would include adding/moving equipment to a specific slot in a users inventory
func (user *User) AddEquipment(eqp Equipment) {
	user.Inventory.mux.Lock()
	user.Inventory.Equipment = append(user.Inventory.Equipment, eqp)
	user.Inventory.mux.Unlock()
}

func (user *User) RemoveEquipment(eqp Equipment) {
	user.Inventory.mux.Lock()
	defer user.Inventory.mux.Unlock()
	for slot, userEqp := range user.Inventory.Equipment {
		if userEqp == eqp {
			user.Inventory.Equipment[slot] = Equipment{}
			return
		}
	}
}

func (users *Users) AdjustCurrency(id int32, currency float64) {
	users.mux.Lock()
	user := users.Users[id]
	user.Currency += currency
	users.mux.Unlock()
}

func (users *Users) TradeItem(sellerId, buyerId int32, eqp Equipment, sellerYield, buyerYield float64) {
	users.mux.Lock()
	seller := users.Users[sellerId]
	buyer := users.Users[buyerId]
	buyer.AddEquipment(eqp)
	users.AdjustCurrency(buyerId, -buyerYield)
	seller.RemoveEquipment(eqp)
	users.AdjustCurrency(sellerId, sellerYield)
	users.mux.Unlock()
}

func (user *User) GenerateEquipment() {
	totalWeapons := 10
	atk := int32(5)
	def := int32(5)
	for i := 0; i < totalWeapons; i++ {
		user.Inventory.Equipment = append(user.Inventory.Equipment, Equipment{"Sword", nextWeaponId,
			user.Id, "This is a basic sword.", EquipmentStats{1, atk, def}})
		atk++
		def++
		nextWeaponId++
	}
	fmt.Printf("Generated equipment\n")
}

// If user has enough currency, can lower their currency.
func (users *Users) HasEnoughCurrency(userId int32, price float64) bool {
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
	if EquipmentIsEmpty(equipment) {
		return false
	}
	for _, item := range user.Inventory.Equipment {
		if item == equipment {
			return true
		}
	}
	return false
}

func (users Users) String() string {
	return fmt.Sprintf("%v", users)
}
