package main

import (
	// "os"
	// "context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"log"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	// "github.com/aws/aws-sdk-go-v2/aws"
	// "github.com/aws/aws-sdk-go-v2/config"
	// "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	// "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	// "github.com/joho/godotenv"
)

type CartEntry struct{ 
	Username string 
	ItemId string 
	Quantity int
}

var allCarts = map[string]map[string]int{
	"bernice": {"01": 3, "03": 1},
	"regine":  {"01": 2, "02": 8, "03": 7},
}

//main function
func main() {
	// // initialising dynamodb
	// if err := godotenv.Load(); err != nil {
    //     log.Fatalf("Error loading .env file: %v", err)
    // }

    // // Read AWS credentials and region from environment variables
    // accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
    // secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
    // region := os.Getenv("AWS_REGION")

    // // Configure AWS SDK using the loaded credentials and region
	// // cfg is an instance of aws.Config 
	// ctx := context.TODO() 
    // cfg, err := config.LoadDefaultConfig(
	// 	ctx,
    //     config.WithRegion(region),
    //     config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
	// 		accessKey,
	// 		secretKey,
	// 		"",
	// 	)),
    // )
    // if err != nil {
    //     log.Fatalf("Error configuring AWS SDK: %v", err)
    // }

    // // Create a DynamoDB client using the loaded configuration
    // client := dynamodb.NewFromConfig(cfg)
	
	// initialiseDb(client)

	lambda.Start(getAllCartsHandler)
	// lambda.Start(getCartHandler)
	// lambda.Start(addToCartHandler)
}

// func initialiseDb(dynamoDBClient *dynamodb.Client) { 
// 	// "database"
// 	allCarts["bernice"] = make(map[string]int) 
// 	allCarts["bernice"]["01"] = 3
// 	allCarts["bernice"]["03"] = 1

// 	allCarts["regine"] = make(map[string]int)
// 	allCarts["regine"]["01"] = 2
// 	allCarts["regine"]["02"] = 8
// 	allCarts["regine"]["03"] = 7
// }

func getAllCartsHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) { 
	log.Println("Get All Cart Lambda function starting")
	if (request.HTTPMethod == "GET") { 
		return getAllCarts(request) 
	} else { 
		return events.APIGatewayProxyResponse{
            StatusCode: http.StatusMethodNotAllowed,
            Body: "Method not allowed",
        }, nil
	}
}

// func getCartHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) { 
// 	log.Println("Get Cart Lambda function starting")
// 	return getCart(request)
// }

// func addToCartHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) { 
// 	log.Println("Add Cart Lambda function starting")
// 	return addToCart(request)
// }

func getAllCarts(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) { 
	log.Println("Received request: ", request)
	responseBody, err := json.Marshal(allCarts)
	if err != nil {
        log.Println("Error marshaling response: ", err)
		return events.APIGatewayProxyResponse{
            StatusCode: http.StatusInternalServerError, 
            Body: "Error marshaling response",
        }, nil
    }
	log.Println("Response body: ", string(responseBody))
	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK, 
		Body: string(responseBody),
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
		Body: fmt.Sprintf("User " + newCart.Username + " added " + strconv.Itoa(newCart.Quantity) + " of Item " + newCart.ItemId + " to cart!"),
	}, nil
}

// func putItem(ctx context.Context, client *dynamodb.Client, cart CartEntry) error {
//     av, err := attributevalue.MarshalMap(cart)
//     if err != nil {
//         return fmt.Errorf("failed to marshal cart: %w", err)
//     }

//     input := &dynamodb.PutItemInput{
//         TableName: aws.String("berniceTest_carts"),
//         Item:      av,
//     }

//     _, err = client.PutItem(ctx, input)
//     if err != nil {
//         return fmt.Errorf("failed to put item: %w", err)
//     }

//     return nil
// }