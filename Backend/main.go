package main

import (
	// "fmt"
	"context"
	"errors"
	"g01API/adminRoutes"
	"g01API/gcsAdmin"
	"g01API/mongoAdmin"
	"g01API/productRoutes"
	"g01API/reviewRoutes"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var gcsManager *gcsAdmin.ClientUploader

const (
	authorizationKey = "Authorization"
	authorizationType = "bearer"
)

// function that run at the first time
func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", os.Getenv("CREDENTIALS_PATH"))
	client = mongoAdmin.Connect(os.Getenv("MONGO_URL"))
	Newclient, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	gcsManager = &gcsAdmin.ClientUploader{
		Client: Newclient,
		BucketName: os.Getenv("GCS_BCK"),
		ProjectID: os.Getenv("GCS_PROJ_ID"),
		UploadPath: "product-pics/",
	}

	gin.SetMode(gin.ReleaseMode)
}

func main() {
	// create gin router
	router := gin.Default()

	authRoute := router.Group("/").Use(authMiddleWare())

	// declare routes
	authRoute.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello Gopher",
		})
	})

	// admin routes
	authRoute.POST("api/admin/add", adminRoutes.AddAdmin(client))
	authRoute.POST("api/admin/login", adminRoutes.LogIn(client))

	// product routes
	authRoute.GET("api/products/list", productRoutes.GetAll(client))
	authRoute.GET("api/products/:prod_id", productRoutes.GetOne(client))
	authRoute.POST("api/products/add", productRoutes.AddOne(client))
	authRoute.PUT("api/products/edit", productRoutes.UpdateOne(client))
	authRoute.DELETE("api/products/del", productRoutes.DeleteOne(client))
	authRoute.POST("api/products/upload-image", uploadFile())
	authRoute.POST("api/products/delete-image", deleteFile())

	// review routes
	authRoute.POST("/api/reviews/add", reviewRoutes.AddOne(client))
	authRoute.GET("/api/reviews/list", reviewRoutes.GetAll(client))

	router.Run()
}

func authMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationKey)

		if len(authHeader) == 0 {
			err := errors.New("authorization header is required")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != authorizationType {
			err := errors.New("unsupported authorization type")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		token := fields[1]
		if token != os.Getenv("API_KEY") {
			err := errors.New("authorization user only")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		} else {
			c.Next()
		}
	}
}

func deleteFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		idHeader := c.GetHeader("prod_id")

		err := gcsManager.DeleteFile(idHeader+".png")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "object deleted",
		})
	}
}

func uploadFile() gin.HandlerFunc {
	return func(c *gin.Context) {

		idHeader := c.GetHeader("prod_id")
		formFile, file_err := c.FormFile("file_input")
		if file_err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": file_err.Error(),
			})
			return
		}

		blobFile, err := formFile.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = gcsManager.UploadFile(blobFile, idHeader+".png")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	}
}