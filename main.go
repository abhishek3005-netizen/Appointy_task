package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"fmt"
	"strings"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var db *mongo.Database
var userCollection *mongo.Collection
var postCollection *mongo.Collection
var userPostCollection *mongo.Collection

func DB_connect()(*mongo.Database){
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")


	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	database := client.Database("tasks_go")

	return database
}

// Used when we post a user object
// type user_struct struct{
// 	Name string
// 	Email string
// 	Password string
// }

// Used when we post a Post object
type post_struct struct{
	User_Id string
	Caption string
	Image_URL string
	Posted_Timestamp string
}

// new
type struct_user struct{
	_Id primitive.ObjectID
	Name string
	Email string 
	Password string
}

type struct_post struct{
	_Id primitive.ObjectID
	Caption string
	Image_URL string 
	Posted_Timestamp string
}

type struct_user_posts struct{
	User_Id primitive.ObjectID
	Post_Id []primitive.ObjectID
}


func init() {
	db  = DB_connect()
	userCollection = db.Collection("user")
	postCollection = db.Collection("post")
	userPostCollection = db.Collection("user_posts")

}

func parseid(path string) (id string){
    var lastslash = strings.LastIndex(path, "/")
    id = path[(lastslash+1):]
    return id 
}

func getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET"{
		return
	}else{
		var path string = r.URL.Path
		var user_id string = parseid(path)
		var temp_user struct_user

		userCollection = db.Collection("user")
		llid,_ := primitive.ObjectIDFromHex(user_id)


		filter := bson.M{"_id":llid}

		userCollection.FindOne(context.TODO(),filter).Decode(&temp_user)


		json.NewEncoder(w).Encode(temp_user)

	}
}

func getPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET"{
		return
	}else{
		var path string = r.URL.Path
		var post_id string = parseid(path)

		var temp_post struct_post

		postCollection = db.Collection("post")
		llid,_ := primitive.ObjectIDFromHex(post_id)


		filter := bson.M{"_id":llid}

		postCollection.FindOne(context.TODO(),filter).Decode(&temp_post)


		json.NewEncoder(w).Encode(temp_post)


	}
}

// func getUserPost(w http.ResponseWriter, r *http.Request){
// 	if r.Method != "GET"{
// 		return
// 	}else{
// 		var path string = r.URL.Path
// 		var user_id string = parseID(path)

// 		// Assuming 0th post object is of 0th user object 
// 		post,ok := posts[user_id]
// 		if ok{
// 			w.Header().Set("Content-Type","application/json")
// 			w.WriteHeader(http.StatusOK);
// 			w.Write([]byte(post))
			
// 		}
// 	}
// }

func createUser(w http.ResponseWriter, r *http.Request){
	if r.Method !="POST"{
		return 
	}else{
		// convert the json you put into the body into that of a struct
		var temp_user struct_user
		var temp_user_posts struct_user_posts

		_ = json.NewDecoder(r.Body).Decode(&temp_user)

	// insert our user model
		result, _ := userCollection.InsertOne(context.TODO(), temp_user)
		userPostCollection = db.Collection("user_posts")

		array_e := make([]primitive.ObjectID,0)

		temp_user_posts.User_Id = (result.InsertedID).(primitive.ObjectID)
		temp_user_posts.Post_Id = array_e

		result1, _ := userCollection.InsertOne(context.TODO(), temp_user_posts)

		json.NewEncoder(w).Encode(result)
		json.NewEncoder(w).Encode(result1)
	}
}

func createPost(w http.ResponseWriter, r *http.Request){
	if r.Method !="POST"{
		return 
	}else{
		// convert the json you put into the body into that of a struct
		var temp_post_1 post_struct
		_ = json.NewDecoder(r.Body).Decode(&temp_post_1)

		postCollection = db.Collection("post")
		userPostCollection = db.Collection("user_posts")

		result, err := postCollection.InsertOne(context.TODO(), struct {Caption string; Image_URL string; Posted_Timestamp string}{temp_post_1.Caption, temp_post_1.Image_URL, temp_post_1.Posted_Timestamp})

		if err !=nil{
			return
		}else{
			llid, _ := primitive.ObjectIDFromHex(temp_post_1.User_Id)
			query := bson.M {"_id":llid}
			update := bson.M {"$push":bson.M{"postids": result.InsertedID.(primitive.ObjectID)}}
			userPostCollection.FindOneAndUpdate(context.TODO(),query,update)
		}
		
	}
}

func main() {
	//s := &server{}

	http.HandleFunc("/users/", getUser)
	http.HandleFunc("/users",createUser)
	http.HandleFunc("/posts/",getPost)
	http.HandleFunc("/posts",createPost)
	// http.HandleFunc("/posts/users/",getUserPost)
	log.Fatal(http.ListenAndServe(":8080", nil))
  }
