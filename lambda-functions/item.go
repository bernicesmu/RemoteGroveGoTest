package main 

import ( 
	"fmt"
	"net/http"
)

type Item struct{ 
	Name string 
	Price float32 
	Description string 
	Rating float32 
}

type ItemEntry struct{ 
	ItemId string 
	Name string 
	Price float32 
	Description string 
	Rating float32 
}

var allItems = make(map[string]Item)

func itemMain() { 
	// "database"
	item1 := Item{ 
		Name: "Textbook",
		Price: 18.20,
		Description: "Textbook used in the IS111 module", 
		Rating: 4.8,
	}

	item2 := Item{ 
		Name: "Pen",
		Price: 1.80,
		Description: "Blue ink ball point pen in 0.38 diameter", 
		Rating: 3.4,
	}

	item3 := Item{ 
		Name: "Mug",
		Price: 9.15,
		Description: "Short classic mug", 
		Rating: 3.7,
	}

	allItems["01"] = item1
	allItems["02"] = item2
	allItems["03"] = item3

	e := echo.New() 
	e.GET("/", getAllItems)
	e.GET("/item/:id", getItemById)
	e.POST("/addItem", addItem)
	e.Logger.Fatal(e.Start(":8001"))
}

func getAllItems(c echo.Context) error { 
	return c.JSON(http.StatusOK, allItems)
}

func getItemById(c echo.Context) error { 
	itemIdToFind := c.Param("id")
	for itemId, itemDetails := range allItems { 
		if itemId == itemIdToFind { 
			return c.JSON(http.StatusOK, itemDetails)
		}
	}
	return c.String(http.StatusOK, "Item ID cannot be found")
}

func addItem(c echo.Context) error { 
	newItemEntry := ItemEntry{}

	if err := c.Bind(&newItemEntry); err != nil { 
		return err 
	}

	itemId := newItemEntry.ItemId
	name := newItemEntry.Name
	price := newItemEntry.Price
	description := newItemEntry.Description
	rating := newItemEntry.Rating

	if _, ok := allItems[itemId]; ok { 
		return c.String(http.StatusOK, "Item ID already exists")
	}

	allItems[itemId] = Item{
		Name: name, 
		Price: price, 
		Description: description, 
		Rating: rating,
	}
	return c.String(http.StatusOK, fmt.Sprintf("Item " + itemId + " has been added successfully"))
}