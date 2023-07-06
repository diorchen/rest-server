// package groceryItemStore testing
package groceryItemStore

import (
	"testing"
	"time"
)

func TestCreateAndGet(t *testing.T) {
	// Create a store and a single food.
	gis := New()
	id := gis.CreateFood("Strawberries", "From Costco", []string{}, time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC), Nutrition{})

	// We should be able to retrieve this food by ID, but nothing with other IDs.
	food, err := gis.GetFood(id)
	if err != nil {
		t.Fatal(err)
	}

	if food.Id != id {
		t.Errorf("got food.Id=%d, id=%d", food.Id, id)
	}
	if food.Name != "Strawberries" {
		t.Errorf("got Name=%v, want %v", food.Name, "Strawberries")
	}

	// Asking for all food, we only get the one we put in.
	allFood := gis.GetAllFood()
	if len(allFood) != 1 || allFood[0].Id != id {
		t.Errorf("got len(allFood)=%d, allFood[0].Id=%d; want 1, %d", len(allFood), allFood[0].Id, id)
	}

	_, err = gis.GetFood(id + 1)
	if err == nil {
		t.Fatal("got nil, want error")
	}

	// Add another food. Expect to find two tasks in the store.
	gis.CreateFood("Bananas", "From Costco", []string{}, time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC), Nutrition{})
	allFood2 := gis.GetAllFood()
	if len(allFood2) != 2 {
		t.Errorf("got len(allFood2)=%d; want 2", len(allFood2))
	}
}

func TestDelete(t *testing.T) {
	gis := New()
	id1 := gis.CreateFood("Apples", "From Costco", []string{}, time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC), Nutrition{})
	id2 := gis.CreateFood("Kiwis", "From Costco", []string{}, time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC), Nutrition{})

	if err := gis.DeleteFood(id1 + 1001); err == nil {
		t.Fatalf("delete food id=%d, got no error; want error", id1+1001)
	}

	if err := gis.DeleteFood(id1); err != nil {
		t.Fatal(err)
	}
	if err := gis.DeleteFood(id1); err == nil {
		t.Fatalf("delete food id=%d, got no error; want error", id1)
	}

	if err := gis.DeleteFood(id2); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteAll(t *testing.T) {
	gis := New()
	gis.CreateFood("Apples", "From Costco", []string{}, time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC), Nutrition{})
	gis.CreateFood("Kiwis", "From Costco", []string{}, time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC), Nutrition{})

	if err := gis.DeleteAllFood(); err != nil {
		t.Fatal(err)
	}

	food := gis.GetAllFood()
	if len(food) > 0 {
		t.Fatalf("want no food remaining; got %v", food)
	}
}

func TestGetFoodByIng(t *testing.T) {
	gis := New()
	gis.CreateFood("Apples", "From Costco", []string{"Apples"}, time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC), Nutrition{})
	gis.CreateFood("Kiwis", "From Costco", []string{"Kiwis"}, time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC), Nutrition{})
	gis.CreateFood("Strawberries", "From Costco", []string{"Strawberries"}, time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC), Nutrition{})
	gis.CreateFood("Guava", "From Costco", []string{}, time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC), Nutrition{})
	gis.CreateFood("Pineapple", "From Costco", []string{}, time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC), Nutrition{})
	gis.CreateFood("Oranges", "From Costco", []string{}, time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC), Nutrition{})

	var tests = []struct {
		Ingredients     string
		wantNum int
	}{
		{"Apples", 1},
		{"Strawberries", 1},
		{"Kiwis", 1},
	}

	for _, tt := range tests {
		t.Run(tt.Ingredients, func(t *testing.T) {
			numByIng := len(gis.GetFoodByIng(tt.Ingredients))

			if numByIng != tt.wantNum {
				t.Errorf("got %v, want %v", numByIng, tt.wantNum)
			}
		})
	}
}

func TestGetFoodsByExpDate(t *testing.T) {
	timeFormat := "2006-Jan-02"
	mustParseDate := func(tstr string) time.Time {
		tt, err := time.Parse(timeFormat, tstr)
		if err != nil {
			t.Fatal(err)
		}
		return tt
	}

	gis := New()
	gis.CreateFood("Apples", "From Costco", []string{"Apples"}, mustParseDate("2020-Dec-01"), Nutrition{})
	gis.CreateFood("Kiwis", "From Costco", []string{"Kiwis"}, mustParseDate("2000-Dec-21"), Nutrition{})
	gis.CreateFood("Strawberries", "From Costco", []string{"Strawberries"}, mustParseDate("2020-Dec-01"), Nutrition{})
	gis.CreateFood("Guava", "From Costco", []string{}, mustParseDate("2000-Dec-21"), Nutrition{})
	gis.CreateFood("Pineapple", "From Costco", []string{}, mustParseDate("2000-Dec-21"), Nutrition{})
	gis.CreateFood("Oranges", "From Costco", []string{}, mustParseDate("1991-Jan-01"), Nutrition{})

	// Check a single task can be fetched.
	y, m, d := mustParseDate("1991-Jan-01").Date()
	food1 := gis.GetFoodsByExpDate(y, m, d)
	if len(food1) != 1 {
		t.Errorf("got len=%d, want 1", len(food1))
	}
	if food1[0].Description != "XY5" {
		t.Errorf("got Text=%s, want XY5", food1[0].Description)
	}

	var tests = []struct {
		date    string
		wantNum int
	}{
		{"2020-Jan-01", 0},
		{"2020-Dec-01", 2},
		{"2000-Dec-21", 2},
		{"1991-Jan-01", 1},
		{"2020-Dec-21", 0},
	}

	for _, tt := range tests {
		t.Run(tt.date, func(t *testing.T) {
			y, m, d := mustParseDate(tt.date).Date()
			numByDate := len(gis.GetFoodsByExpDate(y, m, d))

			if numByDate != tt.wantNum {
				t.Errorf("got %v, want %v", numByDate, tt.wantNum)
			}
		})
	}
}