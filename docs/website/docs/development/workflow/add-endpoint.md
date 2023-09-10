# Add an Endpoint

It is essential to grasp how the server works before bringing any modifications. 
Then, we will guide you through the process of incorporating an HTTP endpoint into the server.

## The Server

Recipya's server code is located within the `internal/server` package. The main file 
is [server.go](https://github.com/reaper47/recipya/blob/main/internal/server/server.go). It exports a single with 
a receiver function and a corresponding struct.

- `Server`: This struct holds the HTTP router, the repository, the email service and the files service. You can 
  find the declaration of each service in the [internal/services/service.go](https://github.com/reaper47/recipya/blob/main/internal/services/service.go) 
  file.
- `Server.Run`: Starts the web server.
- `NewServer`: Creates a server that is ready for use. It requires the services to be passed as arguments.

The HTTP router is initialized during the server's creation. We use [chi](https://github.com/go-chi/chi) due to 
its simplicity in organizing endpoints. Please read the `mountHandlers` function to observe the router in action.

## Example

Let's walk through an example aimed at adding an endpoint that searches for recipes.

### Router

The first step involves adding the endpoint to the router. A suitable endpoint is `GET /recipes/search?q=query`. Open the 
[internal/server/server.go](https://github.com/reaper47/recipya/blob/main/internal/server/server.go) file and include the 
endpoint within the `/recipes` route block. The handler should be named `recipesSearchHandler`, following the
 `{resource}{LastWordEndpoint}{Handler}` naming convention.

```go
r.Route("/recipes", func(r chi.Router) {
  r.Use(s.mustBeLoggedInMiddleware)
  
  r.Get("/search", recipesSearchHandler)
  ...
})
```

With the route established, it is time to declare the handler. Since we are dealing with the `/recipes` resource,
add the handler to the [handlers_recipes.go](https://github.com/reaper47/recipya/blob/main/internal/server/handlers_recipes.go)
file. 

```go
func (s *Server) recipesSearchHandler(w http.ResponseWriter, r *http.Request) {
	panic("TODO: To implement")
}
```

We are now ready to create tests for our route.

### Test

Tests related to the server are written in the `handlers_{resource}_test.go` files. The tests for our handlers are 
stored in the [handlers_recipes_test.go](https://github.com/reaper47/recipya/blob/main/internal/server/handlers_recipes_test.go)
file. The naming convention for test functions is `TestHandlers_{Resource}_{Endpoint}`. Let's write the foundation
function of our tests.

```go
func TestHandlers_Recipes_AddManual(t *testing.T) {
    srv := newServerTest()

	uri := "/recipes/search"
}
```

The subsequent step involves writing the different tests that add value to the users.

```go
func TestHandlers_Recipes_AddManual(t *testing.T) {
    srv := newServerTest()

	uri := "/recipes/search"

    t.Run("must be logged in", func(t *testing.T) {
        assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
    })
	
    t.Run("search fails", func(t *testing.T) {
        t.Fail()
    })
	
    t.Run("user has no recipes", func(t *testing.T) {
        t.Fail()
    })

    t.Run("user searches empty string", func(t *testing.T) {
		t.Fail()
    })

    testcases := []struct {
        name string
		in   string
		want models.Recipes
	}{
        {name: "user searches empty string", in: "", want: ...},
        {name: "user searches for lunch", in: "lunch", want: ...},
        ...
    }
    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"?q="+tc.in, noHeader, nil)
			
            ...
        })
    }
}
```

The body of the tests is omitted for brevity. Run the tests to ensure that they fail.

### Handler

The next step entails crafting the handler's code. Return to the `handles_recipes.go` file, and implement
the `recipesSearchHandler` function that will make the tests go green. For instance, the implementation could resemble
the following. 

```go
func (s *Server) recipesSearchHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	query := chi.URLParam(r, "q")
	
	recipes, err := s.Repository.SearchRecipes(query, userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Failed to search recipes.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

    templates.RenderComponent(w, "recipes", "search-recipes", templates.Data{Recipes: recipes})
}
```

This code gets the logged-in user's ID and the search query, then passes them to the `SearchRecipes` function of the 
repository. If this function encounters an error, an HTMX toast is sent to the user, accompanied by a HTTP 500 
status code. Otherwise, the HTML containing the recipes is sent.

The `templates.RenderComponent` function displays a template from the [web/templates/components](https://github.com/reaper47/recipya/tree/main/web/templates/components)
directory. Its second parameter is the name of one of the files within that directory, excluding the extension. 
The third parameter is the name of the template within a file in that folder. Lastly, the fourth parameter is 
a [struct](https://github.com/reaper47/recipya/blob/main/internal/templates/data.go) containing data for the 
GoHTML template.

### Repository

The final piece of the puzzle involves writing the `s.Repository.SearchRecipes` function. The repository is an interface
that declares functions for interacting with a database. Currently, Recipya supports [sqlite](https://github.com/reaper47/recipya/blob/main/internal/services/sqlite_service.go)
only. To support other databases, we need define a struct satisf the [RepositoryService](https://github.com/reaper47/recipya/blob/main/internal/services/service.go)
interface.

Let's declare the function within the `RepositoryService` interface. The functions are declared alphabetically.

```go
type RepositoryService interface {
    // AddAuthToken adds an authentication token to the database.
    AddAuthToken(selector, validator string, userID int64) error
    
    // AddRecipe adds a recipe to the user's collection.
    AddRecipe(r *models.Recipe, userID int64) (int64, error)
    
    ...
	
	// SearchRecipes gets the user's recipes that include the search query.
	SearchRecipes(query string, userID int64) (models.Recipes, error)
	
	...
    
    // VerifyLogin checks whether the user provided correct login credentials.
    // If yes, their user ID will be returned. Otherwise, -1 is returned.
    VerifyLogin(email, password string) int64
    
    // Websites gets the list of supported websites from which to extract the recipe.
    Websites() models.Websites
}
```

Subsequently, let's implement the function within the `sqlite_service.go` file.

```go
func (s *SQLiteService) SearchRecipes(query string, userID int64) (models.Recipes, error) {
	// s.Mutex.Lock()
	// defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

    rows, err := s.DB.QueryContext(ctx, statements.SelectSearchRecipes, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
	var recipes models.Recipes
    for rows.Next() {
        // code to scan a recipe
		...
        recipes = append(recipes, c)
    }
    return recipes, nil
}
```

Remember, invoking `s.Mutex.Lock()` and `defer s.Mutex.Unlock()` is necessary when inserting, updating, or deleting 
database entries. However, in this scenario, we're merely fetching data, rendering the mutex unnecessary.

SQL statements are organized by action within the [internal/services/statements](https://github.com/reaper47/recipya/tree/main/internal/services/statements)
directory. The naming convention is `{Action}{Resource}`. In our case, a `SELECT` statement for fetching recipes is termed
`SelectSearchRecipes` and would reside in the 
[select.go](https://github.com/reaper47/recipya/tree/main/internal/services/statements) file.

Every statement is a `const` whose value is the SQLite statement itself. If Recipya ever supports other databases, we
shall find a way to organize the statements per database type. 

You can call it a day and open a PR once you wrote your SQL and the tests pass.