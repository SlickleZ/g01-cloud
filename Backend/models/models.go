package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AdminUser struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type Product struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Prod_Id string `json:"prod_id" bson:"prod_id"`
	Prod_Name string `json:"prod_name" bson:"prod_name"`
	Prod_Detail string `json:"prod_detail" bson:"prod_detail"`
	Prod_Price int `json:"prod_price" bson:"prod_price"`
	Prod_Quantity int `json:"prod_quantity" bson:"prod_quantity"`
}

type Review struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Review_Id string `json:"review_id" bson:"review_id"`
	Reviewer string `json:"reviewer" bson:"reviewer"`
	Rating int `json:"rating" bson:"rating"`
	Comment string `json:"comment" bson:"comment"`
}