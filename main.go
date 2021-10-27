package main

import (
	"context"
	"fmt"
	"log"
	"time"

	//"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Recipe struct {
	RecipeID    uint64
	Name        string
	Ingredients string
	Directions  string
	Rating      int
	CreatedAt   time.Time
}

const (
	DB_NAME           = "RecipeDB"
	RECIPE_COLLECTION = "Recipes"
	URI               = "mongodb+srv://niangmodou100:<password>@cluster0.rlimv.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
)

// // Function to populate database with dummy data
// func populateRecipieDB() {
// 	RECIPE_DB = append(RECIPE_DB, Recipe{RecipeID: 0,
// 		Name:        "Chocolate Chip Cookie",
// 		Ingredients: `1 cup Butter, 1 cup Sugar, 3 cups Flour, 1 tsp Sea Salt, 2 cups Chocolate Chips, 1 tsp Salt, 2 large eggs, 1 tsp baking soda`,
// 		Directions: `0. Preheat oven to 375 degrees F,
// 					 1. In a separate bowl mix flour,
// 					 2. baking soda, salt, baking powder,
// 					 3. Cream together butter and sugars until combined,
// 					 4. Beat in eggs and vanilla until fluffy,
// 					 5. Mix in the dry ingredients until combined,
// 					 6. Add 12 oz package of chocolate chips and mix well,
// 					 7. Roll 2-3 TBS of dough at a time into balls and place them on your prepared cookie sheets,
// 					 8. Bake in preheated oven for approximately 8-10 minutes, Let them sit on the baking pan for 2 minutes`,
// 		Rating:    2,
// 		CreatedAt: time.Now(),
// 	})

// 	RECIPE_DB = append(RECIPE_DB, Recipe{RecipeID: 1,
// 		Name:        "Spahgetti and Meatballs",
// 		Ingredients: `1lb spagetti, 1lb ground beef, 1 can crushed tomatoes, 2tbsp olive oil, 1 egg, 1/2 onion, red pepper flakes, chopped parsley, garlic cloves`,
// 		Directions: `1. In a pot of boiling salted water, cook spaghetti. Drain.
// 					 2. In a large bowl, combine beef with bread crumbs, parsley, Parmesan, egg, garlic, 1 teaspoon salt, and red pepper flakes. Mix until just combined then form into 16 balls.
// 					 3. In a large pot over medium heat, heat oil. Add meatballs and cook, turning occasionally, until browned on all sides, about 10 minutes. Transfer meatballs to a plate.
// 					 4. Add onion to pot and cook until soft, 5 minutes. Add crushed tomatoes and bay leaf. Season with salt and pepper and bring to a simmer. Return meatballs to pot and cover. Simmer until sauce has thickened, 8 to 10 minutes
// 					 5. Serve pasta with a healthy scoop of meatballs and sauce. Top with Parmesan before serving.`,
// 		Rating:    3,
// 		CreatedAt: time.Now(),
// 	})

// 	RECIPE_DB = append(RECIPE_DB, Recipe{RecipeID: 2,
// 		Name:        "Chicken Noodle Soup",
// 		Ingredients: `2lb chicken breast, chicken stock, 2tbsp butter, 1 diced onion, 2 carrots, 2 celery stalks, 2 cups egg noodles, black pepper`,
// 		Directions: `1. Melt butter in a large stockpot or Dutch oven over medium heat. Add onion, carrots and celery. Cook, stirring occasionally, until tender, about 3-4 minutes. Stir in garlic until fragrant, about 1 minute.
// 					 2. Whisk in chicken stock and bay leaves; season with salt and pepper, to taste. Add chicken and bring to boil; reduce heat and simmer, covered, until the chicken is cooked through, about 30-40 minutes. Remove chicken and let cool before dicing into bite-size pieces, discarding bones.
// 					 3. Stir in chicken and pasta and cook until tender, about 6-7 minutes.
// 					 4. Remove from heat; stir in parsley, dill and lemon juice; season with salt and pepper, to taste.
// 					 5. Serve immediately.`,
// 		Rating:    2,
// 		CreatedAt: time.Now(),
// 	})

// }

// // Function to retriev the index of a recipe given the ID
// func findRecipeIdx(recipeID uint64) int {
// 	for idx, recipe := range RECIPE_DB {
// 		if recipe.RecipeID == recipeID {
// 			return idx
// 		}
// 	}

// 	return -1
// }

// // Function to retrieve the next higest max ID
// func retrieveID() uint64 {
// 	var maxID uint64 = 0

// 	for _, recipe := range RECIPE_DB {
// 		if recipe.RecipeID > maxID {
// 			maxID = recipe.RecipeID
// 		}
// 	}

// 	return maxID + 1
// }

// Function to retrieve the port from the command line arguments
func retrievePort(arguments []string) string {
	port := "8080"

	if len(arguments) == 3 {
		port = arguments[2]
	}

	return port
}

func main() {
	ctx := context.Background()

	// Connecting to the cluster
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://niangmodou100:3826719mn@cluster0.rlimv.mongodb.net/RecipeDB?retryWrites=true&w=majority")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	databases, _ := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Println(databases)

	// Connect to the collection
	// collection := client.Database(DB_NAME).Collection(RECIPE_COLLECTION)

	// app := iris.Default()
	// port := retrievePort(os.Args)

	// // Load all templates from the "./views" folder
	// app.RegisterView(iris.HTML("./views", ".html"))

	// // Error codes
	// app.OnErrorCode(iris.StatusNotFound, notFound)
	// app.OnErrorCode(iris.StatusInternalServerError, internalServerError)

	// // API Endpoints
	// app.Get("/", home)
	// app.Get("/recipe/{id:uint64}", getRecipe)
	// app.Post("/deleterecipe/{id:uint64}", deleteRecipe)
	// app.Post("/editrecipe/{id:uint64}", editRecipe)
	// app.Post("/addrecipe", addRecipe)

	// // Listens and serves incoming http requests
	// app.Listen(":" + port)
}

// // Endpoint to read all recipies within our database
// func home(ctx iris.Context) {
// 	ctx.ViewData("data", RECIPE_DB)

// 	ctx.View("index.html")
// }

// // Endpoint to read a recipe data - READ
// func getRecipe(ctx iris.Context) {
// 	recipeID, _ := ctx.Params().GetUint64("id")

// 	recipeIdx := findRecipeIdx(recipeID)

// 	recipeArr := []Recipe{RECIPE_DB[recipeIdx]}

// 	ctx.ViewData("data", recipeArr)

// 	ctx.View("recipe.html")
// }

// // Endpoint to delete a recipe - DELETE
// func deleteRecipe(ctx iris.Context) {
// 	recipeID, _ := ctx.Params().GetUint64("id")

// 	indexToRemove := findRecipeIdx(recipeID)

// 	// Assigning last element to current element
// 	N := len(RECIPE_DB)
// 	RECIPE_DB[indexToRemove] = RECIPE_DB[N-1]

// 	RECIPE_DB = RECIPE_DB[:N-1]

// 	// Redirect to the / route
// 	ctx.Redirect("/")
// }

// // Endpoint to add a new recipe - CREATE
// func addRecipe(ctx iris.Context) {
// 	// Extract data from the form
// 	name := ctx.FormValue("name")
// 	ingredients := ctx.FormValue("ingredients")
// 	directions := ctx.FormValue("directions")
// 	rating, _ := strconv.Atoi(ctx.FormValue("rating"))

// 	// Get new ID
// 	newID := retrieveID()

// 	// Current time retrieval
// 	currentTime := time.Now()

// 	// Creating a new Recipe
// 	newRecipe := Recipe{RecipeID: newID, Name: name, Ingredients: ingredients, Directions: directions, Rating: rating, CreatedAt: currentTime}

// 	RECIPE_DB = append(RECIPE_DB, newRecipe)

// 	// Redirect to the / route
// 	ctx.Redirect("/")
// }

// // Endpoint to edit a given receipe - UPDATE
// func editRecipe(ctx iris.Context) {
// 	// Extract request data
// 	recipeID, _ := ctx.Params().GetUint64("id")

// 	name := ctx.FormValue("name")
// 	ingredients := ctx.FormValue("ingredients")
// 	directions := ctx.FormValue("directions")
// 	rating, _ := strconv.Atoi(ctx.FormValue("rating"))

// 	// Current time retrieval
// 	currentTime := time.Now()

// 	recipeIdx := findRecipeIdx(recipeID)

// 	editedRecipe := Recipe{Name: name, Ingredients: ingredients, Directions: directions, Rating: rating, CreatedAt: currentTime}

// 	RECIPE_DB[recipeIdx] = editedRecipe

// 	// Redirect to the / route
// 	ctx.Redirect("/")
// }

// // NOT FOUND ERRORS
// func notFound(ctx iris.Context) {
// 	ctx.WriteString("URL not found, 404")
// }

// // SERVER ERRORS
// func internalServerError(ctx iris.Context) {
// 	ctx.WriteString("Something went wrong, try again")
// }
