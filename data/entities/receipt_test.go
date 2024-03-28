package entities

import "testing"

func TestUniquePoints(t *testing.T) {
	uniqueReceipt := Receipt{
		Items: []Item{
			{
				ShortDescription: "a",
			},
			{
				ShortDescription: "b",
			},
			{
				ShortDescription: "c",
			},
		},
	}

	nonUniqueReceipt := Receipt{
		Items: []Item{
			{
				ShortDescription: "a",
			},
			{
				ShortDescription: "a",
			},
			{
				ShortDescription: "a",
			},
		},
	}

	if points := uniqueNamePoints(&uniqueReceipt); points != 20 {
		t.Fatal("Unexpected number of points")
	}

	if points := uniqueNamePoints(&nonUniqueReceipt); points != 0 {
		t.Fatal("Unexpected number of points")
	}
}
