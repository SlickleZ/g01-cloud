package mongoAdmin

import (
	"context"
	"fmt"
	"g01API/models"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// function that connect to mongoDB server
func Connect(url string) (*mongo.Client) {

	clientOptions := options.Client().ApplyURI(url)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
        log.Fatal(err)
    }
	fmt.Println("MongoDB Connected!")
	return client
}

// function that insert one doc to mongoDB
func InsertOne(client *mongo.Client, ctx context.Context, collection string, doc interface{}) (*mongo.InsertOneResult, error) {

	// select db and collection
	coll := client.Database(os.Getenv("DB_NAME")).Collection(collection)

	result, err := coll.InsertOne(ctx, doc)
	return result, err
}

func DeleteOne(client *mongo.Client, ctx context.Context, collection string, filter interface{}) (*mongo.DeleteResult, error) {

	coll := client.Database(os.Getenv("DB_NAME")).Collection(collection)

	var result *mongo.DeleteResult
	err := coll.FindOneAndDelete(ctx, filter).Decode(&result)

	return result, err
}

func FindAdminOne(client *mongo.Client, ctx context.Context, find, field interface{}) (models.AdminUser, error) {
	
	coll := client.Database(os.Getenv("DB_NAME")).Collection("admin")
     
    // collection has an method Find,
    // that returns a mongo.cursor
    // based on query and field.
	var results models.AdminUser
    err := coll.FindOne(ctx, find, options.FindOne().SetProjection(field)).Decode(&results)
	// fmt.Printf("%#v\n", results)

	return results, err
}

func FindOneProduct(client *mongo.Client, ctx context.Context, find, field interface{}) (models.Product, error) {
	
	coll := client.Database(os.Getenv("DB_NAME")).Collection("products")
     
    // collection has an method Find,
    // that returns a mongo.cursor
    // based on query and field.
	var results models.Product
    err := coll.FindOne(ctx, find, options.FindOne().SetProjection(field)).Decode(&results)
	// fmt.Printf("%#v\n", results)

	return results, err
}

func FindAllProduct(client *mongo.Client, ctx context.Context) ([]models.Product, error) {

	coll := client.Database(os.Getenv("DB_NAME")).Collection("products")

	var results []models.Product
	cur, findErr := coll.Find(ctx, bson.M{})
	if findErr != nil {
		return nil, findErr
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var result models.Product
		decodeErr := cur.Decode(&result)
		if decodeErr != nil {
			return nil, decodeErr
		}
		results = append(results, result)
	}
	return results, nil
}

func UpdateOneProduct(client *mongo.Client, ctx context.Context, filter, update interface{}) (models.Product, error) {

	coll := client.Database(os.Getenv("DB_NAME")).Collection("products")

	var result models.Product
	err := coll.FindOneAndUpdate(ctx, filter, update).Decode(&result)

	return result, err
}

func FindAllReviews(client *mongo.Client, ctx context.Context) ([]models.Review, error) {

	coll := client.Database(os.Getenv("DB_NAME")).Collection("reviews")

	var reviews []models.Review
	cur, findErr := coll.Find(ctx, bson.M{})
	if findErr != nil {
		return nil, findErr
	}

	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var review models.Review
		if decodeErr := cur.Decode(&review); decodeErr != nil {
			return nil, decodeErr
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}