package cmd

import (
	"barbar/pkg/mongodb"
	"context"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func InitIndexMongo() *cobra.Command {
	return &cobra.Command{
		Use:   "init-index-mongo",
		Short: "use init-index-mongo",
		Long:  `Init Index for setup index on mongo database`,
		Run: func(cmd *cobra.Command, args []string) {
			runScriptMongo()
		},
	}
}

func runScriptMongo() {

	timeoutCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// init mongo
	mongoUser, err := mongodb.Connect(timeoutCtx, "mongodb://root:rootpassword@localhost:27017/?authSource=admin", "users")
	if err != nil {
		log.Fatal(err)
	}
	defer mongoUser.Client().Disconnect(timeoutCtx)

	mongoAuth, err := mongodb.Connect(timeoutCtx, "mongodb://root:rootpassword@localhost:27017/?authSource=admin", "auth")
	if err != nil {
		log.Fatal(err)
	}
	defer mongoAuth.Client().Disconnect(timeoutCtx)

	logger.Info("Starting Mongo Script.....")

	// create index uid and email on table user
	logger.Info("set index for database users")
	modUser := []mongo.IndexModel{
		{
			Keys: bson.M{
				"uid":   -1,
			}, Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{
				"email":   1,
			}, Options: options.Index().SetUnique(true),
		},
	}

	colUser := mongoUser.Collection("users")

	// Declare an options object
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err = colUser.Indexes().CreateMany(timeoutCtx, modUser, opts)
	// Check for the options errors
	if err != nil {
		fmt.Println("Indexes().CreateIndexes() User ERROR:", err)
		os.Exit(1) // exit in case of error
	} else {
		fmt.Println("CreateIndexes User () opts:", opts)
	}

	// create index uid db auth
	logger.Info("set index for database users")
	modAuth := []mongo.IndexModel{
		{
			Keys: bson.M{
				"uid":   -1,
			}, Options: options.Index().SetUnique(true),
		},
	}

	colAuth := mongoAuth.Collection("auth")

	// Declare an options object
	_, err = colAuth.Indexes().CreateMany(timeoutCtx, modAuth, opts)
	// Check for the options errors
	if err != nil {
		fmt.Println("Indexes().CreateIndexes() Auth ERROR:", err)
		os.Exit(1) // exit in case of error
	} else {
		fmt.Println("CreateIndexes Auth () opts:", opts)
	}

	logger.Info("Create Index Finish")
}
