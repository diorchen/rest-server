// Simple data store for food items identifiable by numeric ID.
package groceryItemStore

import (
	"fmt"
	"sync" // synchronization primitives for managing concurrent access to shared resources
	"time"
)

type FoodItem struct {
	Id   			int       `json:"id"`
	Name 			string    `json:"name"`
	Description 	string	  `json:"description"`
	Ingredients 	[]string  `json:"ingredients"` // slice of stirngs
	Expiration   	time.Time `json:"expiration"`
	Nutrition		Nutrition `json:"nutrition"`
}

type Nutrition struct {
	Calories		int			`json:"calories"`
	Protein			float64		`json:"protein"`
	Carbohydrates	float64		`json:"carbohydrates"`
	Fat				float64		`json:"fat"`
	Fiber			float64		`json:"fiber"`
}
// GroceryItemStore is a simple in-memory database of food items; GroceryItemStore methods are
// safe to call concurrently.
type GroceryItemStore struct { // collection of food items
	sync.Mutex // mutual exclusion lock to protect shared resources from concurrent access by multiple goroutines

	food  map[int]FoodItem // each groceryItem associated with ID by mapping keys of type 'int' to values of type 'FoodItem'
	nextId int // ensures ID uniqueness, keeps track of next available ID to be assigned
}

func New() *GroceryItemStore { // Func 'New' returns pointer to (*) struct GroceryItemStore
	gis := &GroceryItemStore{} // new var 'gis' assigned to newly allocated 'GroceryItemStore' object (empty), initialized with {}.
	gis.food = make(map[int]FoodItem) // initializes 'food' of the 'GroceryItemStore' as an empty map, providing a storage container for grocery items.
	gis.nextId = 0
	return gis
}

// CreateFood creates a new food in the store.
	// method receiver, indicates CreateFood is associated with GroceryItemStore object, gis = name of receiver variable
func (gis *GroceryItemStore) CreateFood(name string, description string, ingredients []string, expiration time.Time, nutrition Nutrition) int {
	gis.Lock() // lock synchronizes access to resource 'item' variable
	defer gis.Unlock() // ensure lock is released when function returns

	food := FoodItem{ // Creates new FoodItem and initializes fields
		Id:   gis.nextId,	
		Name: name,
		Description: description,
		Ingredients: make([]string, len(ingredients)),
		Expiration: expiration,
		Nutrition: nutrition}

	copy(food.Ingredients, ingredients)

	gis.food[gis.nextId] = food // associates new created food with new ID
	gis.nextId++ // increments new ID
	return food.Id
}

// GetFood retrieves a food from the store, by id. If no such id exists, an
// error is returned.
func (gis *GroceryItemStore) GetFood(id int) (FoodItem, error) {
	gis.Lock()
	defer gis.Unlock()

	food, ok := gis.food[id] // food = key, ok = boolean flag
	if ok {
		return food, nil
	} else {
		return FoodItem{}, fmt.Errorf("food with id=%d not found", id)
	}
}

// DeleteFood deletes the food with the given id. If no such id exists, an error
// is returned.
func (gis *GroceryItemStore) DeleteFood(id int) error {
	gis.Lock()
	defer gis.Unlock()

	if _, ok := gis.food[id]; !ok { // check if food item with given id exists in store.food map, if not, return error
		return fmt.Errorf("food with id=%d not found", id)
	}

	delete(gis.food, id)
	return nil
}

// DeleteAllFood deletes all food in the store.
func (gis *GroceryItemStore) DeleteAllFood() error {
	gis.Lock()
	defer gis.Unlock()

	gis.food = make(map[int]FoodItem) // reset the store.food map to an empty map
	return nil // return nil to indicate successful deletion
}

// GetAllFood returns all the food in the store, in arbitrary order.
func (gis *GroceryItemStore) GetAllFood() []FoodItem {
	gis.Lock()
	defer gis.Unlock()

	allFood := make([]FoodItem, 0, len(gis.food))
	for _, food := range gis.food {
		allFood = append(allFood, food)
	}
	return allFood
}

// GetFoodByIng returns all the food that have the given ingredients, in arbitrary
// order.
func (gis *GroceryItemStore) GetFoodByIng(ingredients string) []FoodItem {
	gis.Lock()
	defer gis.Unlock()

	var foods []FoodItem

foodloop:
	for _, food := range gis.food {
		for _, foodIng := range food.Ingredients {
			if foodIng == ingredients {
				foods = append(foods, food)
				continue foodloop
			}
		}
	}
	return foods
}

// GetFoodByExpDate returns all the food that have the given exp date, in
// arbitrary order.
func (gis *GroceryItemStore) GetFoodsByExpDate(year int, month time.Month, day int) []FoodItem {
	gis.Lock()
	defer gis.Unlock()

	var foods []FoodItem

	for _, food := range gis.food {
		y, m, d := food.Expiration.Date()
		if y == year && m == month && d == day {
			foods = append(foods, food)
		}
	}

	return foods
}