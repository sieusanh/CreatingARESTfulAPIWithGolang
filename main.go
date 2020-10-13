package main

import (
	"fmt"
	"log"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"sync"
)

type Article struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Desc string `json:"desc"`
	Content string `json:"content"`
}

// let's declare a global Articles array
// that we can then populate in our main function 
// to simulate a database
var (
	Articles []Article
	mutex sync.Mutex
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	mutex.Lock()
	json.NewEncoder(w).Encode(Articles)
	mutex.Unlock()
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	//fmt.Fprintf(w, "Key: " + key)

	// Loop over all of our Articles
	// if the article.Id equals the key we pass in
	// return the article encoded as JSON
	mutex.Lock()
	for _, article := range Articles {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
			break
		}
	}
	mutex.Unlock()
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)

	// update our global Articles array to include
	// our new Article
	mutex.Lock()
	Articles = append(Articles, article)
	mutex.Unlock()
	json.NewEncoder(w).Encode(article)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	// one again, we will need to parse the path parameters
	vars := mux.Vars(r)
	// we will need to extract the `id` of the article we
	// wish to delete
	id := vars["id"]

	mutex.Lock()
	// we then need to loop through all our articles 
	for index, article := range Articles {
		// if our id path parameter matches one of our
		// articles
		if article.Id == id {
			// updates our Articles array to remove the 
			// article
			Articles = append(Articles[:index], Articles[index+1:]...)
			break
		}
	}
	mutex.Unlock()
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)

	mutex.Lock()
	for index := range Articles {
		if Articles[index].Id == article.Id {
			Articles[index].Title = article.Title
			Articles[index].Desc = article.Desc
			Articles[index].Content = article.Content
			break
		}
	}
	mutex.Unlock()
}

// Traditional net/http package
/*
func handleRequests() {
	http.HandleFunc("/", homePage)
	// add our articles route and map it to our
	// returnAllArticles function like so
	http.HandleFunc("/articles", returnAllArticles)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
*/

// gorilla/mux package
func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAllArticles)
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	// Mặc dù khác phương thức http, nhưng trùng tên link thì vẫn ko dùng được phương thức DELETE
	//myRouter.HandleFunc("article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/delete_article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/update_article", updateArticle).Methods("PUT")

	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	Articles = []Article{
        Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
        Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
    }
	handleRequests()
}