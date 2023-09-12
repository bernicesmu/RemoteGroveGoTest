package main

import ( 
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"encoding/json"
)

type Item struct{ 
	Name string 
	Price float32 
	Description string 
	Rating float32 
}

func main() {
	e := echo.New()
	
	e.GET("/:username", calcTotalOfUsername)

	e.Logger.Fatal(e.Start(":8002"))
}

func calcTotalOfUsername(c echo.Context) error {
	username := c.Param("username")

	cartResp, _ := http.Get("http://cart:8000/cart/" + username)
	decoder := json.NewDecoder(cartResp.Body)
	var cart map[string]int 
	if err := decoder.Decode(&cart); err != nil {
		return c.String(http.StatusInternalServerError, "Error decoding response")
	}

	itemResp, _ := http.Get("http://item:8001/")
	decoder = json.NewDecoder(itemResp.Body)
	var allItems map[string]Item
	if err := decoder.Decode(&allItems); err != nil { 
		return c.String(http.StatusInternalServerError, "Error decoding response")
	}
	
	total := float64(0.0)
	for itemId, quantity := range cart { 
		total += float64(allItems[itemId].Price) * float64(quantity)
	}
	return c.String(http.StatusOK, strconv.FormatFloat(total, 'f', 2, 64))
}