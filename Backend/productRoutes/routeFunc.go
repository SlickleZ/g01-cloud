package productRoutes

import (
	"context"
	// "fmt"
	"g01API/models"
	"g01API/mongoAdmin"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	
)

func GetAll(client *mongo.Client) gin.HandlerFunc {
	return func (c *gin.Context) {
		results, findAllErr := mongoAdmin.FindAllProduct(client, context.TODO())
		if findAllErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "findAllErr",
				"error": findAllErr.Error(),
			})
			return
		}
		if results == nil {
			c.JSON(http.StatusNoContent, gin.H{
				"result": results,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"result": results,
			})
		}
	}
}

func GetOne(client *mongo.Client) gin.HandlerFunc {
	return func (c *gin.Context) {
		prod_id := c.Param("prod_id")

		filter := bson.D{
			{Key: "prod_id", Value: prod_id},
		}
		opts := bson.D{
			{Key: "_id", Value: 0},
		}

		product, findOneErr := mongoAdmin.FindOneProduct(client, context.TODO(), filter, opts)
		// fmt.Printf("%#v\n", product)
		if findOneErr != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "findOneError",
				"error": findOneErr.Error(),
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"result": product,
		})
	}
}

func AddOne(client *mongo.Client) gin.HandlerFunc {
	return func (c *gin.Context)  {
		var req models.Product
	
		// convert req (JSON) to struct
		if jsonBindErr := c.ShouldBindJSON(&req); jsonBindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "jsonBindErr",
				"error": jsonBindErr.Error(),
			})
			return
		}
		
		filter := bson.D{
			{Key: "prod_id", Value: req.Prod_Id},
		}

		opts := bson.D{
			{Key: "_id", Value: 0},
		}
		
		_, findErr := mongoAdmin.FindOneProduct(client, context.TODO(), filter, opts)
		if findErr != nil { // in case of not found, insert it!
			document := bson.D{
				{Key: "prod_id", Value: req.Prod_Id},
				{Key: "prod_name", Value: req.Prod_Name},
				{Key: "prod_detail", Value: req.Prod_Detail},
				{Key: "prod_price", Value: req.Prod_Price},
				{Key: "prod_quantity", Value: req.Prod_Quantity},
			}
			_, insertErr := mongoAdmin.InsertOne(client, context.TODO(), "products",document)
			if insertErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "insertErr",
					"error": insertErr.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "success",
			})
			return
		}
		// in case of founded product exists, error response
		c.JSON(http.StatusConflict, gin.H{
			"message": "product already exists",
		})
	}
}

func UpdateOne(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.Product

		if jsonBinderErr := c.ShouldBindJSON(&req); jsonBinderErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "jsonBindErr",
				"error": jsonBinderErr.Error(),
			})
			return
		}

		filter := bson.D{
			{Key: "prod_id", Value: req.Prod_Id},
		}
		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "prod_name", Value: req.Prod_Name},
					{Key: "prod_detail", Value: req.Prod_Detail},
					{Key: "prod_price", Value: req.Prod_Price},
					{Key: "prod_quantity", Value: req.Prod_Quantity},
				},
			},
		}

		_, updateErr := mongoAdmin.UpdateOneProduct(client, context.TODO(), filter, update)
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "updateErr",
				"error": updateErr.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "product updated",
		})
	}
}

func DeleteOne(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.Product

		if jsonBinderErr := c.ShouldBindJSON(&req); jsonBinderErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "jsonBinderErr",
				"error": jsonBinderErr.Error(),
			})
			return
		}
		
		filter := bson.D{
			{Key: "prod_id", Value: req.Prod_Id},
		}

		_, deleteErr := mongoAdmin.DeleteOne(client, context.TODO(), "products", filter)
		if deleteErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "deleteErr",
				"error": deleteErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "product deleted",
		})
	}
}