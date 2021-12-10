package adminRoutes

import (
	"context"
	"g01API/models"
	"g01API/mongoAdmin"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func AddAdmin(client *mongo.Client) gin.HandlerFunc {
	return func (c *gin.Context) {
		
		adminHeader := c.GetHeader("admin")
		if adminHeader == os.Getenv("ADMIN") {
			var req models.AdminUser
			// convert req (JSON) to struct
			if jsonBindErr := c.ShouldBindJSON(&req); jsonBindErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "jsonBindErr",
					"error": jsonBindErr.Error(),
				})
				return
			}
		
			filter := bson.D{
				{Key: "username", Value: req.Username},
			}
			opts := bson.D{
				{Key: "_id", Value: 0},
				{Key: "password", Value: 0},
			}
		
			_, findErr := mongoAdmin.FindAdminOne(client, context.TODO(), filter, opts)
			if findErr != nil { // case not found user, created it!
				hash, hashErr := bcrypt.GenerateFromPassword([]byte(req.Password), 5)
				if hashErr != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "hashErr",
						"error": hashErr.Error(),
					})
					return
				}
		
				document := bson.D {
					{Key: "username", Value: req.Username},
					{Key: "password", Value: hash},
				}
		
				_, insertErr := mongoAdmin.InsertOne(client, context.TODO(), "admin", document)
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
			// in case of founded username, error response
			c.JSON(http.StatusConflict, gin.H{
				"message": "admin already exists",
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "authorized user only",
			})
		}
	}
}

func LogIn(client *mongo.Client) gin.HandlerFunc {
	return func (c *gin.Context) {
		var req models.AdminUser
	
		// convert req (JSON) to struct
		if jsonBindErr := c.ShouldBindJSON(&req); jsonBindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "jsonBindErr",
				"error": jsonBindErr.Error(),
			})
			return
		}

		filter := bson.D{
			{Key: "username", Value: req.Username},
		}
		opts := bson.D{
			{Key: "_id", Value: 0},
		}
	
		user, findErr := mongoAdmin.FindAdminOne(client, context.TODO(), filter, opts)
		if findErr != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "findErr",
				"error": findErr.Error(),
			})
			return
		}

		compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if compareErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "compareErr",
				"error": compareErr.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "log in successfully",
		})
	}
}
