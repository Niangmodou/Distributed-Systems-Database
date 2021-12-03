/*
Parallel and Distributed Systems
Project Part 4 Recipe Application
backend.go
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
)

// Recipe Data Type
type Recipe struct {
	RecipeID    uint64
	Name        string
	Ingredients string
	Directions  string
	Rating      int
	CreatedAt   time.Time
}

// Message queue data structure
type Queue struct {
	data []string
}

func (queue *Queue) isEmpty() bool {
	return len(queue.data) == 0
}

func (queue *Queue) enqueue(element string) {
	queue.data = append(queue.data, element)
}

func (queue *Queue) size() int {
	return len(queue.data)
}

func (queue *Queue) dequeue() string {
	element := queue.data[0]

	if queue.size() > 1 {
		queue.data = queue.data[1:]
	} else {
		queue.data = []string{}
	}

	return element
}

// Database of all recipies
var RECIPE_DB []Recipe

// Log of Messages
var LOG []string = []string{}

// Storage of other backend addresses
var BACKENDS []string

// Queue of all the messages
var messageQueue Queue

// Function to populate database with dummy data
func populateRecipieDB() {
	RECIPE_DB = append(RECIPE_DB, Recipe{RecipeID: 0,
		Name:        "Chocolate Chip Cookie",
		Ingredients: `1 cup Butter, 1 cup Sugar, 3 cups Flour, 1 tsp Sea Salt, 2 cups Chocolate Chips, 1 tsp Salt, 2 large eggs, 1 tsp baking soda`,
		Directions: `0. Preheat oven to 375 degrees F, 
					 1. In a separate bowl mix flour, 
					 2. baking soda, salt, baking powder, 
					 3. Cream together butter and sugars until combined, 
					 4. Beat in eggs and vanilla until fluffy, 
					 5. Mix in the dry ingredients until combined, 
					 6. Add 12 oz package of chocolate chips and mix well, 
					 7. Roll 2-3 TBS of dough at a time into balls and place them on your prepared cookie sheets, 
					 8. Bake in preheated oven for approximately 8-10 minutes, Let them sit on the baking pan for 2 minutes`,
		Rating:    2,
		CreatedAt: time.Now(),
	})

	RECIPE_DB = append(RECIPE_DB, Recipe{RecipeID: 1,
		Name:        "Spahgetti and Meatballs",
		Ingredients: `1lb spagetti, 1lb ground beef, 1 can crushed tomatoes, 2tbsp olive oil, 1 egg, 1/2 onion, red pepper flakes, chopped parsley, garlic cloves`,
		Directions: `1. In a pot of boiling salted water, cook spaghetti. Drain.
					 2. In a large bowl, combine beef with bread crumbs, parsley, Parmesan, egg, garlic, 1 teaspoon salt, and red pepper flakes. Mix until just combined then form into 16 balls.
					 3. In a large pot over medium heat, heat oil. Add meatballs and cook, turning occasionally, until browned on all sides, about 10 minutes. Transfer meatballs to a plate.
					 4. Add onion to pot and cook until soft, 5 minutes. Add crushed tomatoes and bay leaf. Season with salt and pepper and bring to a simmer. Return meatballs to pot and cover. Simmer until sauce has thickened, 8 to 10 minutes
					 5. Serve pasta with a healthy scoop of meatballs and sauce. Top with Parmesan before serving.`,
		Rating:    3,
		CreatedAt: time.Now(),
	})

	RECIPE_DB = append(RECIPE_DB, Recipe{RecipeID: 2,
		Name:        "Chicken Noodle Soup",
		Ingredients: `2lb chicken breast, chicken stock, 2tbsp butter, 1 diced onion, 2 carrots, 2 celery stalks, 2 cups egg noodles, black pepper`,
		Directions: `1. Melt butter in a large stockpot or Dutch oven over medium heat. Add onion, carrots and celery. Cook, stirring occasionally, until tender, about 3-4 minutes. Stir in garlic until fragrant, about 1 minute.
					 2. Whisk in chicken stock and bay leaves; season with salt and pepper, to taste. Add chicken and bring to boil; reduce heat and simmer, covered, until the chicken is cooked through, about 30-40 minutes. Remove chicken and let cool before dicing into bite-size pieces, discarding bones.
					 3. Stir in chicken and pasta and cook until tender, about 6-7 minutes.
					 4. Remove from heat; stir in parsley, dill and lemon juice; season with salt and pepper, to taste.
					 5. Serve immediately.`,
		Rating:    2,
		CreatedAt: time.Now(),
	})

}

// // Function to retriev the index of a recipe given the ID
func findRecipeIdx(recipeID uint64) int {
	for idx, recipe := range RECIPE_DB {
		if recipe.RecipeID == recipeID {
			return idx
		}
	}

	return -1
}

// Function to retrieve the next higest max ID
func retrieveID() uint64 {
	var maxID uint64 = 0

	for _, recipe := range RECIPE_DB {
		if recipe.RecipeID > maxID {
			maxID = recipe.RecipeID
		}
	}

	return maxID + 1
}

// Function to retrieve the port from the command line arguments
func retrievePort(arguments []string) string {
	port := arguments[2]

	return port
}

// Function to handle the home route by returning all recipes
func handleHome(conn net.Conn) {
	recipeJson, _ := json.Marshal(RECIPE_DB)
	payload := string(recipeJson) + string("\n")

	conn.Write([]byte(payload))
}

// Function to retrieve a given recipe
func handleGet(conn net.Conn, recipeIDString string) {
	recipeID, _ := strconv.Atoi(recipeIDString)

	idx := findRecipeIdx(uint64(recipeID))
	recipeArr := []Recipe{RECIPE_DB[idx]}
	recipeJson, _ := json.Marshal(recipeArr)

	payload := string(recipeJson) + string("\n")

	conn.Write([]byte(payload))
}

// Function to create a new recipe
func handleCreate(jsonString string) {
	var recipe Recipe
	json.Unmarshal([]byte(jsonString), &recipe)

	recipe.RecipeID = retrieveID()

	RECIPE_DB = append(RECIPE_DB, recipe)
}

// Function to edit recipe from DB
func handleEdit(recipeIDString, jsonString string) {
	recipeID, _ := strconv.Atoi(recipeIDString)
	recipeIdx := findRecipeIdx(uint64(recipeID))

	jsonString = jsonString[:len(jsonString)-1]

	var editedRecipe Recipe
	json.Unmarshal([]byte(jsonString), &editedRecipe)

	RECIPE_DB[recipeIdx] = editedRecipe
}

// Function to delete recipe from DB
func handleDelete(recipeIDString string) {
	recipeID, _ := strconv.Atoi(recipeIDString)
	indexToRemove := findRecipeIdx(uint64(recipeID))

	// Assigning last element to current element
	N := len(RECIPE_DB)
	RECIPE_DB[indexToRemove] = RECIPE_DB[N-1]
	RECIPE_DB = RECIPE_DB[:N-1]
}

// Function to contact other backends and sync logs
func syncLogs() {
	for idx := 0; idx < len(BACKENDS); idx++ {
		currConnection, err := net.Dial("tcp", BACKENDS[idx])
		// Checking whether the backend is running
		if err != nil {
			continue
		}
		// Marshalling our Log to send to other backends
		message := "SYNC~" //+ string(marshaledLog)
		fmt.Fprintf(currConnection, message+"\n")
		// We get the current backends log
		response, _ := bufio.NewReader(currConnection).ReadString('\n')

		var currLog []string
		json.Unmarshal([]byte(response), &currLog)

		// We then compare that size with ours to see if we have to reassign
		if len(currLog) > len(LOG) {
			LOG = currLog
			executeMessagesInLog()
		}

		currConnection.Close()
	}
}

// Function to handle the server connection given a client connection
func handleServerConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	log.Println("Client connected from", remoteAddr)
	// variables used in PAXOS
	var maxId int = 0
	var proposalAccepted bool
	var acceptedID int
	var acceptedValue []byte

	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		dataList := strings.Split(string(message), "~")

		if dataList[0] == "ALL" {
			handleHome(conn)
		} else if dataList[0] == "READ" {
			recipeIDString := dataList[1][:len(dataList[1])-1]

			handleGet(conn, recipeIDString)
		} else if dataList[0] == "PREPARE" {
			// Phase 1b: Promise
			proposalNumber, _ := strconv.Atoi(dataList[1])

			if proposalNumber > maxId {
				conn.Write([]byte("FAIL" + "\n"))
			} else {
				maxId = proposalNumber
				// Check whether the proposal was already accepted
				if proposalAccepted {
					payload := "PROMISE~" + string(maxId) + "~" + string(acceptedID) + "~" + string(acceptedValue)
					conn.Write([]byte(payload + "\n"))
				} else {
					conn.Write([]byte("PROMISE~" + string(maxId) + "\n"))
				}
			}

		} else if dataList[0] == "PROPOSE" {
			proposalNumber, _ := strconv.Atoi(dataList[1])
			if proposalNumber == maxId {

				proposalAccepted = true
				acceptedID = proposalNumber
				acceptedValue = []byte(dataList[2])
				command := dataList[2]

				// Adding the current command to the message queue
				messageQueue.enqueue(command)

				// Adding messages to the log
				LOG = append(LOG, command)

				payload := "ACCEPT~" + string(acceptedValue)
				conn.Write([]byte(payload + "\n"))

				// Spawning a thread that executes the messages in message queue
				go executeMessages()

				// Contacting other backends to sync logs
				go syncLogs()

			} else {
				conn.Write([]byte("FAIL" + "\n"))
			}
		} else if dataList[0] == "SYNC" {
			marshaledLog, err := json.Marshal(LOG)

			if err != nil {
				fmt.Println("Error marshaling")
			}

			conn.Write([]byte(string(marshaledLog) + "\n"))
		}
	}
}

// Function to execute the messages in our log incase of a discrepancy
func executeMessagesInLog() {
	// Erasing our database and executing messages in correct order
	RECIPE_DB = []Recipe{}
	populateRecipieDB()

	for idx := 0; idx < len(LOG); idx++ {
		currCommand := LOG[idx]
		executeCommand(currCommand)
	}
}

// Function to retrieve the other backend addresses within the system
func retrieveAddreses(arguments []string) []string {
	var backends []string
	backends = strings.Split(arguments[4], ",")

	return backends
}

// Function to execute a given command from the frontend
func executeCommand(command string) {
	dataList := strings.Split(string(command), "+")
	if dataList[0] == "CREATE" {
		jsonString := dataList[1]
		handleCreate(jsonString)
	} else if dataList[0] == "UPDATE" {
		recipeIDString := dataList[1]
		handleEdit(recipeIDString, dataList[2])
	} else if dataList[0] == "DELETE" {
		recipeIDString := dataList[1][:len(dataList[1])-1]
		handleDelete(recipeIDString)
	}
}

// Function to execute the messages in the message queue
func executeMessages() {
	// Execute the messages in the queue while the queue isn't empty
	for !messageQueue.isEmpty() {
		message := messageQueue.dequeue()
		executeCommand(message)
	}
}

func main() {
	populateRecipieDB()

	arguments := os.Args
	fmt.Println(arguments)

	port := retrievePort(arguments)
	BACKENDS = retrieveAddreses(arguments)

	// Sync with other ones in case we crashed
	go syncLogs()

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer listener.Close()
	fmt.Println("Listening on port #", port)
	for {
		fmt.Println("Backend Started......")
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal("Listener Error:", err)
		}

		// Handling individual server connections
		go handleServerConnection(conn)
	}
}
