package bid

import (
	"context"
	"github.com/HunnTeRUS/fullcycle-auction-go/configuration/logger"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/entity/auction_entity"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/entity/bid_entity"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/infra/database/auction"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"sync"
	"time"
)

type BidEntityMongo struct {
	Id        string  `bson:"_id" copier:"Id"`
	UserId    string  `bson:"user_id" copier:"UserId"`
	AuctionId string  `bson:"auction_id" copier:"AuctionId"`
	Amount    float64 `bson:"amount" copier:"Amount"`
	Timestamp int64   `bson:"timestamp" copier:"Timestamp"`
}

type BidRepositoryMongo struct {
	Collection           *mongo.Collection
	AuctionRepository    *auction.AuctionRepository
	AuctionInterval      time.Duration
	AuctionStatusMap     map[string]auction_entity.AuctionStatus
	AuctionEndTimeMap    map[string]time.Time
	AuctionStatusMapMux  *sync.Mutex
	AuctionEndTimeMapMux *sync.Mutex
}

func NewBidRepository(database *mongo.Database) *BidRepositoryMongo {
	return &BidRepositoryMongo{
		Collection:           database.Collection("bids"),
		AuctionRepository:    auction.NewAuctionRepository(database),
		AuctionInterval:      getAuctionInterval(),
		AuctionStatusMap:     make(map[string]auction_entity.AuctionStatus),
		AuctionEndTimeMap:    make(map[string]time.Time),
		AuctionStatusMapMux:  &sync.Mutex{},
		AuctionEndTimeMapMux: &sync.Mutex{},
	}
}

func (repo *BidRepositoryMongo) CreateBids(ctx context.Context, bids []bid_entity.Bid) *internal_error.InternalError {
	var wg sync.WaitGroup

	for _, bidF := range bids {
		wg.Add(1)
		go func(bidValue bid_entity.Bid) {
			defer wg.Done()

			repo.AuctionStatusMapMux.Lock()
			auctionStatus, okStatus := repo.AuctionStatusMap[bidValue.AuctionId]
			repo.AuctionStatusMapMux.Unlock()

			repo.AuctionEndTimeMapMux.Lock()
			auctionEndTime, okEndTime := repo.AuctionEndTimeMap[bidValue.AuctionId]
			repo.AuctionEndTimeMapMux.Unlock()

			if okStatus && okEndTime {
				if auctionStatus != auction_entity.Active {
					// Leilão não está ativo, pule este lance
					return
				}

				now := time.Now()
				if now.After(auctionEndTime) {
					// Leilão está no estado "Completed", pule este lance
					return
				}

				repo.saveBid(ctx, bidValue)
			} else {
				auction, err := repo.AuctionRepository.FindAuctionByAuctionId(ctx, bidValue.AuctionId)
				if err != nil {
					logger.Error("Error finding auction_usecase", err)
					return
				}

				repo.AuctionStatusMapMux.Lock()
				repo.AuctionStatusMap[bidValue.AuctionId] = auction.Status
				repo.AuctionStatusMapMux.Unlock()

				repo.AuctionEndTimeMapMux.Lock()
				repo.AuctionEndTimeMap[bidValue.AuctionId] = auction.Timestamp.Add(repo.AuctionInterval)
				repo.AuctionEndTimeMapMux.Unlock()

				if auction.Status != auction_entity.Active {
					return
				}

				repo.saveBid(ctx, bidValue)
			}
		}(bidF)
	}

	wg.Wait()

	return nil
}

func (repo *BidRepositoryMongo) saveBid(ctx context.Context, bid bid_entity.Bid) *internal_error.InternalError {
	bidEntity := BidEntityMongo{
		Id:        bid.Id,
		UserId:    bid.UserId,
		AuctionId: bid.AuctionId,
		Amount:    bid.Amount,
		Timestamp: bid.Timestamp.Unix(),
	}

	_, err := repo.Collection.InsertOne(ctx, bidEntity)
	if err != nil {
		logger.Error("Error inserting bid_usecase into MongoDB", err)
		return internal_error.NewInternalServerError(err.Error())
	}

	return nil
}

func getAuctionInterval() time.Duration {
	batchInsertInterval := os.Getenv("AUCTION_INTERVAL")
	if duration, err := time.ParseDuration(batchInsertInterval); err == nil {
		return duration
	}

	return 15 * time.Second
}
