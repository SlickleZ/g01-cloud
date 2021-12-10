package reviewRoutes
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

func AddOne(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.Review

		if jsonBinderErr := c.ShouldBindJSON(&req); jsonBinderErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "jsonBinderErr",
				"error": jsonBinderErr.Error(),
			})
			return
		}

		document := bson.D{
			{Key: "review_id", Value: req.Review_Id},
			{Key: "reviewer", Value: req.Reviewer},
			{Key: "rating", Value: req.Rating},
			{Key: "comment", Value: req.Comment},
		}

		_, insertErr := mongoAdmin.InsertOne(client, context.TODO(), "reviews", document)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "insertErr",
				"error": insertErr.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "review added",
		})
	}
}

func GetAll(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		reviews, findErr := mongoAdmin.FindAllReviews(client, context.TODO())
		if findErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "findErr",
				"error": findErr.Error(),
			})
		}
		if reviews == nil {
			c.JSON(http.StatusNoContent, gin.H{
				"result": reviews,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"result": reviews,
			})
		}
	}
}