package main

import (
	"context"
	"github.com/HunnTeRUS/fullcycle-auction-go/configuration/database/mongodb"
	"github.com/HunnTeRUS/fullcycle-auction-go/configuration/logger"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/infra/api/web/auction_controller"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/infra/api/web/bid_controller"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/infra/api/web/user_controller"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/infra/database/auction"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/infra/database/bid"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/infra/database/user"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/usecase/auction_usecase"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/usecase/bid_usecase"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/usecase/user_usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	databaseConnection, err := mongodb.NewMongoDBConnection(context.Background())
	if err != nil {
		log.Fatal("Error initiating mongodb databaseConnection")
	}

	router := gin.Default()

	user, bid, auction := initUserDependencies(databaseConnection)

	router.GET("/auction", auction.FindAuctions)
	router.POST("/auction", auction.CreateAuctionController)
	router.POST("/bid", bid.CreateBid)
	router.GET("/bid/:auctionId", bid.FindBidsByAuctionIdController)
	router.GET("/bid/winner/:auctionId", bid.FindWinningBidByAuctionIdController)
	router.GET("/user/:userId", user.FindUserById)
	logger.Info("aaaa")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func initUserDependencies(db *mongo.Database) (
	userDependencies *user_controller.UserController,
	bidDependencies *bid_controller.BidController,
	auctionDependencies *auction_controller.AuctionController) {

	bidRepository := bid.NewBidRepository(db)
	auctionRepository := auction.NewAuctionRepository(db)

	userDependencies = user_controller.NewUserController(
		user_usecase.NewUserService(user.NewUserRepository(db)))

	bidDependencies = bid_controller.NewBidController(
		bid_usecase.NewBidService(bidRepository))

	auctionDependencies = auction_controller.NewAuctionController(
		auction_usecase.NewAuctionService(auctionRepository, bidRepository))

	return
}
