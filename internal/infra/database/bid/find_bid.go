package bid

import (
	"context"
	"github.com/HunnTeRUS/fullcycle-auction-go/configuration/logger"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/entity/bid_entity"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

func (repo *BidRepositoryMongo) FindBidsByAuctionId(ctx context.Context, auctionId string) ([]bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"auction_id": auctionId}

	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error finding bids by AuctionId", err, zap.String("AuctionId", auctionId))
		return nil, internal_error.NewInternalServerError("Error finding bids by AuctionId")
	}
	defer cursor.Close(ctx)

	var bids []BidEntityMongo
	if err := cursor.All(ctx, &bids); err != nil {
		logger.Error("Error decoding bids", err)
		return nil, internal_error.NewInternalServerError("Error decoding winning bid_usecase")
	}

	var bidsEntity []bid_entity.Bid
	for _, bid := range bids {
		bidsEntity = append(bidsEntity, bid_entity.Bid{
			Id:        bid.Id,
			UserId:    bid.UserId,
			AuctionId: bid.AuctionId,
			Amount:    bid.Amount,
			Timestamp: time.Unix(bid.Timestamp, 0),
		})
	}

	return bidsEntity, nil
}

func (repo *BidRepositoryMongo) FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*bid_entity.Bid, *internal_error.InternalError) {
	var winningBid BidEntityMongo

	opts := options.FindOne().SetSort(bson.D{{"amount", -1}})
	filter := bson.M{"auction_id": auctionId}

	err := repo.Collection.FindOne(ctx, filter, opts).Decode(&winningBid)
	if err != nil {
		logger.Error("Error finding winning bid_usecase by AuctionId", err, zap.String("AuctionId", auctionId))
		return nil, internal_error.NewInternalServerError("Error finding winning bid_usecase by AuctionId")
	}

	bidEntity := bid_entity.Bid{
		Id:        winningBid.Id,
		UserId:    winningBid.UserId,
		AuctionId: winningBid.AuctionId,
		Amount:    winningBid.Amount,
		Timestamp: time.Unix(winningBid.Timestamp, 0),
	}

	return &bidEntity, nil
}
