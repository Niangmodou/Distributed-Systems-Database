/*
Parallel and Distributed Systems
Project Part 4 Recipe Application
frontend.go
MADE WITH LOVE BY MODOU NAING
*/
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
)

var proposalNumber int
var connection net.Conn
var err error

var backends []string

type Recipe struct {
	RecipeID    uint64    `json:"recipeID`
	Name        string    `json:"name"`
	Ingredients string    `json:"ingredients`
	Directions  string    `json:"directions"`
	Rating      int       `json:"rating"`
	CreatedAt   time.Time `json:createdAt`
}

// Function to retrieve the port from the command line arguments
func retrievePort(arguments []string) string {
	return arguments[2]
}

// Function to retrieve the backend address
func retrieveBackendAddress(arguments []string) {
	backends = strings.Split(arguments[4], ",")
}

func main() {
	app := iris.Default()

	fmt.Println("Frontend Started.....")

	arguments := os.Args

	PORT := retrievePort(arguments)
	retrieveBackendAddress(arguments)

	if err != nil {
		log.Fatal("Dial Error:", err)
	}

	app.Get("/", home)
	app.Post("/addrecipe", addRecipe) // Run Paxos
	app.Get("/recipe/{id:uint64}", getRecipe)
	app.Post("/editrecipe/{id:uint64}", editRecipe)     // Run Paxos
	app.Post("/deleterecipe/{id:uint64}", deleteRecipe) // Run Paxos

	// Load all templates from the "./views" folder
	app.RegisterView(iris.HTML("./views", ".html"))

	// Listens and serves incoming http requests
	app.Listen(":" + PORT)

	// Error codes
	app.OnErrorCode(iris.StatusNotFound, notFound)
	app.OnErrorCode(iris.StatusInternalServerError, internalServerError)

}

// Function to fufill home endpoint
func home(ctx iris.Context) {
	var jsonArr []Recipe
	// Requesting available backends for the data
	for idx := 0; idx < len(backends); idx++ {
		connection, err = net.Dial("tcp", backends[idx])
		// Error connecting to the current backend
		if err != nil {
			continue
		}

		command := "ALL~"
		fmt.Fprintf(connection, command+"\n")
		connection.Write([]byte(command))

		// Retrieving the response from the backend
		message, _ := bufio.NewReader(connection).ReadString('\n')
		json.Unmarshal([]byte(message), &jsonArr)

		if jsonArr != nil {
			break
		}
		connection.Close()

	}
	ctx.ViewData("data", jsonArr)
	ctx.View("index.html")
}

// Endpoint to add a new recipe
func addRecipe(ctx iris.Context) {
	var aliveBackends []string
	var recipeJson []byte
	// Phase 1a: Prepare - Promise
	for idx := 0; idx < len(backends); idx++ {
		currConnection, err := net.Dial("tcp", backends[idx])
		// Error connecting to the current backend
		if err != nil {
			continue
		}
		command := "PREPARE~" + string(proposalNumber) + "~" + string("\n")
		fmt.Println("Making a call to #", idx)
		fmt.Fprintf(currConnection, command+"\n")

		// Retrieving response from backend
		response, _ := bufio.NewReader(currConnection).ReadString('\n')
		if response != "FAIL" {
			// Check if any responses contain accepted values
			respList := strings.Split(response, "~")
			if len(respList) == 4 {
				recipeJson = []byte(respList[3])

			} else {
				recipeJson = nil
			}
			aliveBackends = append(aliveBackends, backends[idx])
		}
		proposalNumber++
		currConnection.Close()
	}

	// Checking if we received PROMISE responses from a majority of acceptors?
	if !(len(aliveBackends) > len(backends)/2) {
		internalServerError(ctx)
		return
	}

	// Phase 2a: Propose
	var aliveAccepted []string
	for idx := 0; idx < len(aliveBackends); idx++ {
		fmt.Println("Checking for accepted backends")
		fmt.Println("Dialing #", idx)
		currConnection, err := net.Dial("tcp", aliveBackends[idx])
		// Error connecting to the current backend
		if err != nil {
			continue
		}
		// We have a value from the highest accepted ID

		// Check if any responses contain accepted values
		if recipeJson == nil {
			// Extract data from the form
			name := ctx.FormValue("name")
			ingredients := ctx.FormValue("ingredients")
			directions := ctx.FormValue("directions")
			rating, _ := strconv.Atoi(ctx.FormValue("rating"))
			currentTime := time.Now()

			// Create a Json Representation of the Recipe and marshal
			newRecipe := Recipe{Name: name, Ingredients: ingredients, Directions: directions, Rating: rating, CreatedAt: currentTime}
			recipeJson, _ = json.Marshal(newRecipe)
		}

		// Making call to the backend
		command := "CREATE+" + string(recipeJson)
		message := "PROPOSE~" + string(proposalNumber) + "~" + command
		fmt.Fprintf(currConnection, message+"\n")

		// Retrieving response from backend
		response, _ := bufio.NewReader(currConnection).ReadString('\n')

		if response != "FAIL" {
			aliveAccepted = append(aliveAccepted, aliveBackends[idx])
			fmt.Println("Succesfully added the recipe to backend #", idx)
		}
	}
	recipeJson = nil

	// Redirect to the / route
	ctx.Redirect("/")
}

// // Endpoint to read a recipe data - READ
func getRecipe(ctx iris.Context) {
	recipeID, _ := ctx.Params().GetUint64("id")
	// Sending the recipeID to the backend
	recipeIDString := strconv.FormatUint(recipeID, 10)

	var jsonArr []Recipe
	for idx := 0; idx < len(backends); idx++ {
		connection, err = net.Dial("tcp", backends[idx])
		// Error connecting to the current backend
		if err != nil {
			continue
		}

		message := "READ+" + recipeIDString
		fmt.Fprintf(connection, message+"\n")

		// Retrieving the response from the backend
		response, _ := bufio.NewReader(connection).ReadString('\n')
		json.Unmarshal([]byte(response), &jsonArr)

		// We have retrieved a valid response
		if jsonArr != nil {
			break
		}
		connection.Close()
	}

	ctx.ViewData("data", jsonArr)

	ctx.View("recipe.html")
}

// Endpoint to delete a recipe - DELETE
func deleteRecipe(ctx iris.Context) {
	var aliveBackends []string
	// Phase 1a: Prepare - Promise
	for idx := 0; idx < len(backends); idx++ {
		currConnection, err := net.Dial("tcp", backends[idx])
		// Error connecting to the current backend
		if err != nil {
			continue
		}
		command := "PREPARE~" + string(proposalNumber) + "~" + string("\n")
		fmt.Fprintf(currConnection, command+"\n")

		// Retrieving response from backend
		response, _ := bufio.NewReader(currConnection).ReadString('\n')
		if response != "FAIL" {
			// Check if any responses contain accepted values
			aliveBackends = append(aliveBackends, backends[idx])
		}
		currConnection.Close()
		proposalNumber++
	}

	// Checking if we received PROMISE responses from a majority of acceptors?
	if !(len(aliveBackends) > len(backends)/2) {
		internalServerError(ctx)
		return
	}

	// Phase 2a: Propose
	var aliveAccepted []string
	for idx := 0; idx < len(aliveBackends); idx++ {
		currConnection, err := net.Dial("tcp", aliveBackends[idx])
		// Error connecting to the current backend
		if err != nil {
			continue
		}

		// Extract request data
		recipeID, _ := ctx.Params().GetUint64("id")
		recipeIDString := strconv.FormatUint(recipeID, 10)

		// Making call to the backend
		command := "DELETE+" + recipeIDString
		message := "PROPOSE~" + string(proposalNumber) + "~" + command
		fmt.Fprintf(currConnection, message+"\n")

		// Retrieving response from backend
		response, _ := bufio.NewReader(currConnection).ReadString('\n')

		if response != "FAIL" {
			aliveAccepted = append(aliveAccepted, aliveBackends[idx])
		}
	}

	// Redirect to the / route
	ctx.Redirect("/")
}

// Endpoint to edit a given receipe - UPDATE
func editRecipe(ctx iris.Context) {
	var aliveBackends []string
	var recipeJson []byte
	// Phase 1a: Prepare - Promise
	for idx := 0; idx < len(backends); idx++ {
		fmt.Println("Checking for alive backends")
		fmt.Println("Dialing #", idx)
		currConnection, err := net.Dial("tcp", backends[idx])
		// Error connecting to the current backend
		if err != nil {
			continue
		}
		command := "PREPARE~" + string(proposalNumber) + "~" + string("\n")
		fmt.Fprintf(currConnection, command+"\n")

		// Retrieving response from backend
		response, _ := bufio.NewReader(currConnection).ReadString('\n')
		fmt.Println("Retrieved response from #", idx)
		if response != "FAIL" {
			// Check if any responses contain accepted values
			respList := strings.Split(response, "~")
			if len(respList) == 4 {
				recipeJson = []byte(respList[3])

			} else {
				recipeJson = nil
			}
			aliveBackends = append(aliveBackends, backends[idx])
		} else {
			internalServerError(ctx)
			return
		}
		proposalNumber++
		currConnection.Close()
	}

	// Checking if we received PROMISE responses from a majority of acceptors?
	if !(len(aliveBackends) > len(backends)/2) {
		internalServerError(ctx)
		return
	}

	// Phase 2a: Propose
	var aliveAccepted []string
	for idx := 0; idx < len(aliveBackends); idx++ {
		currConnection, _ := net.Dial("tcp", aliveBackends[idx])
		// We have a value from the highest accepted ID

		// Check if any responses contain accepted values
		var recipeIDString string
		if recipeJson == nil {
			// Extract request data
			recipeID, _ := ctx.Params().GetUint64("id")
			recipeIDString = strconv.FormatUint(recipeID, 10)

			name := ctx.FormValue("name")
			ingredients := ctx.FormValue("ingredients")
			directions := ctx.FormValue("directions")
			rating, _ := strconv.Atoi(ctx.FormValue("rating"))

			// Current time retrieval
			currentTime := time.Now()

			editedRecipe := Recipe{Name: name, Ingredients: ingredients, Directions: directions, Rating: rating, CreatedAt: currentTime}

			recipeJson, _ = json.Marshal(editedRecipe)
		}

		// Making call to the backend
		command := "UPDATE+" + recipeIDString + "+" + string(recipeJson)
		message := "PROPOSE~" + string(proposalNumber) + "~" + command
		fmt.Fprintf(currConnection, message+"\n")

		// Retrieving response from backend
		response, _ := bufio.NewReader(currConnection).ReadString('\n')

		if response != "FAIL" {
			aliveAccepted = append(aliveAccepted, aliveBackends[idx])
			fmt.Println("Succesfully edited the recipe to backend #", idx)
		}
	}
	recipeJson = nil

	// Redirect to the / route
	ctx.Redirect("/")
}

// NOT FOUND ERRORS
func notFound(ctx iris.Context) {
	ctx.WriteString("URL not found, 404")
}

// SERVER ERRORS
func internalServerError(ctx iris.Context) {
	ctx.WriteString("Something went wrong, try again")
}
