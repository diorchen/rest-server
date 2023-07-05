// Basic REST server using standard library
// GroceryItemStore keeps inventory of groceries bought with expiration dates and other related information.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"     //  Multipurpose Internet Mail Extensions (MIME) type detection and extensions
	"net/http" // HTTP client and server implementations
	"strconv"
	"strings"
	"time"

	"github.com/diorchen/rest-server/internal/groceryItemStore"
)

type foodServer struct {
	groceryItemStore *groceryItemStore.GroceryItemStore
}
 // Creates a new instance of FoodServer
func NewFoodServer() *foodServer {
	store := groceryItemStore.New()
	return &foodServer{groceryItemStore: store}
}
// Handles incoming HTTP requests related to food items
func (fs *foodServer) foodHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/food/" { 
		// Request is plain "/food/", without trailing ID.
		if req.Method == http.MethodPost {
			fs.createFoodHandler(w, req)
		} else if req.Method == http.MethodGet {
			fs.getAllFoodHandler(w, req)
		} else if req.Method == http.MethodDelete {
			fs.deleteAllFoodHandler(w, req)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE or POST at /food/, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
		// Request has an ID, as in "/food/<id>".
		path := strings.Trim(req.URL.Path, "/") // Trims the '/'
		pathParts := strings.Split(path, "/") // splits the string into parts
		if len(pathParts) < 2 {
			http.Error(w, "expect /food/<id> in food handler", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(pathParts[1]) // converts the string into integer
		if err != nil { // checks if there is an error during this conversion
			http.Error(w, err.Error(), http.StatusBadRequest) // return error
			return
		}

		if req.Method == http.MethodDelete { // checks the HTTP method and performs the corresponding action
			fs.deleteFoodHandler(w, req, id)
		} else if req.Method == http.MethodGet {
			fs.getFoodHandler(w, req, id)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET or DELETE at /food/<id>, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func (fs *foodServer) createFoodHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling food creation at %s\n", req.URL.Path)

	// Types used internally in this handler to (de-)serialize the request and
	// response from/to JSON.
	type requestFood struct {
		Name        string                     `json:"name"`
		Description string                     `json:"description"`
		Ingredients []string                   `json:"ingredients"`
		Expiration  time.Time                  `json:"expiration"`
		Nutrition   groceryItemStore.Nutrition `json:"nutrition"`
	}

	// data structure representing the expected payload format for creating a food item
	type responseId struct {
		Id int `json:"id"`
	}

	// Enforces a JSON Content-Type.
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil { // checks for error while parsing the media type
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" { // checks if media type is not equal to JSON
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	// Decode the JSON request body into go struct 'requestFood'
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rf requestFood // holds the decoded JSON data in 'rf'
	if err := dec.Decode(&rf); err != nil { // Checks for error
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := fs.groceryItemStore.CreateFood(rf.Name, rf.Description, rf.Ingredients, rf.Expiration, rf.Nutrition)
	js, err := json.Marshal(responseId{Id: id}) // creates a new struct 'respondID' with 'id' value and marsals it into JSON format
	if err != nil { // checks for error during JSON marshaling process
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json") // sets the 'Content-Type' header of the HTTP response to indicate that the response body contains JSON data
	w.Write(js) // writes the JSON data stored in 'js' as the response body
}

func (fs *foodServer) getAllFoodHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all food items at %s\n", req.URL.Path)

	allFood := fs.groceryItemStore.GetAllFood()
	js, err := json.Marshal(allFood)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (fs *foodServer) getFoodHandler(w http.ResponseWriter, req *http.Request, id int) {
	log.Printf("handling get food item at %s\n", req.URL.Path)

	food, err := fs.groceryItemStore.GetFood(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(food)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (fs *foodServer) deleteFoodHandler(w http.ResponseWriter, req *http.Request, id int) {
	log.Printf("handling deletion of food item at %s\n", req.URL.Path)

	err := fs.groceryItemStore.DeleteFood(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (fs *foodServer) deleteAllFoodHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling deletion of all foods at %s\n", req.URL.Path)
	fs.groceryItemStore.DeleteAllFood()
}

func (fs *foodServer) ingHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling foods by ingredients at %s\n", req.URL.Path)

	if req.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("expect method GET /ing/<ing>, got %v", req.Method), http.StatusMethodNotAllowed)
		return
	}

	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /ing/<ingredient> path", http.StatusBadRequest)
		return
	}
	tag := pathParts[1]

	food := fs.groceryItemStore.GetFoodByIng(tag)
	js, err := json.Marshal(food)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (fs *foodServer) expHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling food items by expiration date at %s\n", req.URL.Path)

	if req.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("expect method GET /exp/<date>, got %v", req.Method), http.StatusMethodNotAllowed)
		return
	}

	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	badRequestError := func() {
		http.Error(w, fmt.Sprintf("expect /exp/<year>/<month>/<day>, got %v", req.URL.Path), http.StatusBadRequest)
	}
	if len(pathParts) != 4 {
		badRequestError()
		return
	}

	year, err := strconv.Atoi(pathParts[1])
	if err != nil {
		badRequestError()
		return
	}
	month, err := strconv.Atoi(pathParts[2])
	if err != nil || month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}
	day, err := strconv.Atoi(pathParts[3])
	if err != nil {
		badRequestError()
		return
	}

	food := fs.groceryItemStore.GetFoodsByExpDate(year, time.Month(month), day)
	js, err := json.Marshal(food)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	mux := http.NewServeMux() // Creates new instance of HTTP request multiplexer
	server := NewFoodServer() // Creates new instance of FoodServer
	mux.HandleFunc("/food/", server.foodHandler)
	mux.HandleFunc("/ing", server.ingHandler)
	mux.HandleFunc("/exp", server.expHandler)
	http.ListenAndServe("localhost:8080", mux)
	fmt.Println("Running Server at port 8080")
}

