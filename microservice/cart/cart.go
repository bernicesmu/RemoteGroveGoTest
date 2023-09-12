package main

import (
	"github.com/labstack/echo"
	"fmt"
	"net/http"
	"strconv"
)

type CartEntry struct{ 
	Username string 
	ItemId string 
	Quantity int
}

var allCarts = make(map[string]map[string]int)

//main function
func main() {
	// "database"
	allCarts["bernice"] = make(map[string]int) 
	allCarts["bernice"]["01"] = 3
	allCarts["bernice"]["03"] = 1

	allCarts["regine"] = make(map[string]int)
	allCarts["regine"]["01"] = 2
	allCarts["regine"]["02"] = 8
	allCarts["regine"]["03"] = 7

	// create a new echo instance
	e := echo.New()
	e.GET("/", getAllCarts)
	e.GET("/cart/:username", getCart)
	e.POST("/addCart", addToCart)
	e.Logger.Fatal(e.Start(":8000"))
}

func getAllCarts(c echo.Context) error { 
	return c.JSON(http.StatusOK, allCarts)
}

func getCart(c echo.Context) error { 
	usernameToFind := c.Param("username")
	for username, items := range allCarts { 
		if username == usernameToFind { 
			return c.JSON(http.StatusOK, items)
		}
	}
	return c.String(http.StatusOK, "Cannot find username")
}

func addToCart(c echo.Context) error {
	newCart := CartEntry{} 
	
	if err := c.Bind(&newCart); err != nil { 
		return err 
	}

	username := newCart.Username
	itemId := newCart.ItemId
	quantity := newCart.Quantity

	if _, ok := allCarts[username]; !ok { 
		allCarts[username] = make(map[string]int)
	}

	allCarts[username][itemId] = quantity 
	return c.String(http.StatusOK, fmt.Sprintf("User " + username + " added " + strconv.Itoa(quantity) + " of Item " + itemId + " to cart!"))
}