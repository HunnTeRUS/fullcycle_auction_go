package auction

import (
	"context"
	"errors"
	"github.com/HunnTeRUS/fullcycle-auction-go/configuration/logger"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/entity/auction_entity"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

func (repo *AuctionRepository) FindAuctionByAuctionId(ctx context.Context, auctionId string) (*auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{"_id": auctionId}

	var auctionEntity AuctionEntityMongo
	err := repo.Collection.FindOne(ctx, filter).Decode(&auctionEntity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error("Auction not found", err, zap.String("AuctionId", auctionId))
			return nil, internal_error.NewNotFoundError("User not found")
		}

		logger.Error("Error finding auction_usecase by AuctionId", err, zap.String("AuctionId", auctionId))
		return nil, internal_error.NewInternalServerError("Error finding auction_usecase by AuctionId")
	}

	auction := auction_entity.Auction{
		Id:          auctionEntity.ID,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Status:      auctionEntity.Status,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Timestamp:   time.Unix(auctionEntity.Timestamp, 0),
	}

	return &auction, nil
}

func (repo *AuctionRepository) FindAuctions(
	ctx context.Context,
	status auction_entity.AuctionStatus,
	category string,
	productName string) ([]auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}

	if category != "" {
		filter["category"] = category
	}

	if productName != "" {
		filter["productName"] = primitive.Regex{Pattern: productName, Options: "i"}
	}

	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error finding auctions", err)
		return nil, internal_error.NewInternalServerError("Error finding auctions")
	}
	defer cursor.Close(ctx)

	var auctionsMongo []AuctionEntityMongo
	if err := cursor.All(ctx, &auctionsMongo); err != nil {
		logger.Error("Error decoding auctions", err)
		return nil, internal_error.NewInternalServerError("Error decoding auctions")
	}

	var auctionsEntity []auction_entity.Auction
	for _, auction := range auctionsMongo {
		auctionsEntity = append(auctionsEntity, auction_entity.Auction{
			Id:          auction.ID,
			ProductName: auction.ProductName,
			Category:    auction.Category,
			Status:      auction.Status,
			Description: auction.Description,
			Condition:   auction.Condition,
			Timestamp:   time.Unix(auction.Timestamp, 0),
		})
	}

	return auctionsEntity, nil
}
