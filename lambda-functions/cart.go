package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"log"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

	lambda.Start(cartHandler)
}

func cartHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) { 
	log.Println("Cart Lambda function starting")
	
	switch request.Path {
		case "/":
			return getAllCarts(request)
		case "/cart/{username}":
			return getCart(request)
		case "/addCart": 
			return addToCart(request)
		default: 
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound, 
				Body: "Route not found",
			}, nil 
	}
}

func getAllCarts(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) { 
	reponseBody, _ := json.Marshal(allCarts)
	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK, 
		Body: string(reponseBody),
	}
	return response, nil 
}

func getCart(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) { 
	usernameToFind := request.PathParameters["username"]
	for username, items := range allCarts { 
		if username == usernameToFind { 
			response, _ := json.Marshal(items)
			return events.APIGatewayProxyResponse{ 
				StatusCode: http.StatusOK, 
				Body: string(response),
			}, nil 
		}
	}
	return events.APIGatewayProxyResponse{ 
		StatusCode: http.StatusOK, 
		Body: "Cannot find username",
	}, nil 
}

func addToCart(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	newCart := CartEntry{} 
	
	if err := json.Unmarshal([]byte(request.Body), &newCart); err != nil { 
		return events.APIGatewayProxyResponse{ 
			StatusCode: http.StatusInternalServerError, 
			Body: "Invalid request body",
		}, nil
	}

	username := newCart.Username
	itemId := newCart.ItemId
	quantity := newCart.Quantity

	if _, ok := allCarts[username]; !ok { 
		allCarts[username] = make(map[string]int)
	}

	allCarts[username][itemId] = quantity 
	return events.APIGatewayProxyResponse{ 
		StatusCode: http.StatusOK, 
		Body: fmt.Sprintf("User " + username + " added " + strconv.Itoa(quantity) + " of Item " + itemId + " to cart!"),
	}, nil
}