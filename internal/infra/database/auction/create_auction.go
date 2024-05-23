package auction

import (
	"context"
	"github.com/HunnTeRUS/fullcycle-auction-go/configuration/logger"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/entity/auction_entity"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"os"
	"time"
)

type AuctionEntityMongo struct {
	ID          string                          `bson:"_id,omitempty" copier:"ID"`
	ProductName string                          `bson:"productName" copier:"ProductName"`
	Category    string                          `bson:"category" copier:"Category"`
	Status      auction_entity.AuctionStatus    `bson:"status" copier:"Status"`
	Description string                          `bson:"description" copier:"Description"`
	Condition   auction_entity.ProductCondition `bson:"condition" copier:"Condition"`
	Timestamp   int64                           `bson:"timestamp" copier:"Timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (repo *AuctionRepository) CreateAuction(ctx context.Context, auction *auction_entity.Auction) *internal_error.InternalError {
	auctionID, err := repo.createAuction(ctx, auction)
	if err != nil {
		logger.Error("error trying to create auction_usecase", err)
		return internal_error.NewInternalServerError("Error trying to create auction_usecase")
	}

	go func() {
		select {
		case <-time.After(getAuctionInterval()):
			update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}
			filter := bson.M{"_id": auctionID}

			_, err := repo.Collection.UpdateOne(context.Background(), filter, update)
			if err != nil {
				logger.Error("Error updating auction_usecase status", err, zap.String("AuctionId", auctionID))
			}
		case <-ctx.Done():
			return
		}
	}()

	return nil
}

func (repo *AuctionRepository) createAuction(ctx context.Context, auction *auction_entity.Auction) (string, *internal_error.InternalError) {
	auctionEntityMongo := AuctionEntityMongo{
		ID:          auction.Id,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Status:      auction.Status,
		Description: auction.Description,
		Condition:   auction.Condition,
		Timestamp:   auction.Timestamp.Unix(),
	}

	_, err := repo.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error inserting auction_usecase into MongoDB", err)
		return primitive.NilObjectID.Hex(), internal_error.NewInternalServerError(err.Error())
	}

	logger.Info("Auction created successfully", zap.String("AuctionID", auction.Id))
	return auctionEntityMongo.ID, nil
}

func getAuctionInterval() time.Duration {
	batchInsertInterval := os.Getenv("AUCTION_INTERVAL")
	if duration, err := time.ParseDuration(batchInsertInterval); err == nil {
		return duration
	}

	return 5 * time.Minute
}
