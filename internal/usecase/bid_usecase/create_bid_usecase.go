package bid_usecase

import (
	"context"
	"github.com/HunnTeRUS/fullcycle-auction-go/configuration/logger"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/entity/bid_entity"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/internal_error"
	"os"
	"strconv"
	"time"
)

type BidUseCase struct {
	BidRepository bid_entity.BidRepositoryInterface

	timer               *time.Timer
	maxBatchSize        int
	batchInsertInterval time.Duration
}

func NewBidService(bidRepository bid_entity.BidRepositoryInterface) BidUseCaseInterface {
	return &BidUseCase{
		BidRepository:       bidRepository,
		timer:               time.NewTimer(getBatchSizeInterval()),
		maxBatchSize:        getMaxBatchSize(),
		batchInsertInterval: getBatchSizeInterval(),
	}
}

type BidUseCaseInterface interface {
	CreateBid(ctx context.Context, bid BidInputDTO) *internal_error.InternalError
	FindBidsByAuctionId(ctx context.Context, auctionId string) ([]BidOutputDTO, *internal_error.InternalError)
	FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*BidOutputDTO, *internal_error.InternalError)
}

type BidInputDTO struct {
	UserId    string  `json:"user_id"`
	AuctionId string  `json:"auction_id"`
	Amount    float64 `json:"amount"`
}

type BidOutputDTO struct {
	UserId    string    `json:"user_id" copier:"UserId"`
	AuctionId string    `json:"auction_id" copier:"AuctionId"`
	Amount    float64   `json:"amount" copier:"Amount"`
	Timestamp time.Time `json:"timestamp" time_format:"2006-01-02 15:04:05" copier:"Timestamp"`
}

var bidBatch []bid_entity.Bid

func (bs *BidUseCase) CreateBid(ctx context.Context, bidDTO BidInputDTO) *internal_error.InternalError {
	bidChannel := make(chan bid_entity.Bid, bs.maxBatchSize)

	go func() {
		defer close(bidChannel)

		for {
			select {
			case bid, ok := <-bidChannel:
				if !ok {
					if len(bidBatch) > 0 {
						err := bs.BidRepository.CreateBids(ctx, bidBatch)
						if err != nil {
							logger.Error("error trying to process the bid", err)
						}
					}
					return
				}

				bidBatch = append(bidBatch, bid)

				select {
				case <-bs.timer.C:
					err := bs.BidRepository.CreateBids(ctx, bidBatch)
					if err != nil {
						logger.Error("error trying to process the bid", err)
					}
					bidBatch = nil
					bs.timer.Reset(bs.batchInsertInterval)
				default:
					if len(bidBatch) >= bs.maxBatchSize {
						err := bs.BidRepository.CreateBids(ctx, bidBatch)
						if err != nil {
							logger.Error("error trying to process the bid", err)
						}
						bidBatch = nil
						bs.timer.Reset(bs.batchInsertInterval)
					}
				}
			case <-bs.timer.C:
				err := bs.BidRepository.CreateBids(ctx, bidBatch)
				if err != nil {
					logger.Error("error trying to process the bid", err)
				}
				bidBatch = nil
				bs.timer.Reset(bs.batchInsertInterval)
			}
		}
	}()

	bidEntity, err := bid_entity.CreateBid(bidDTO.UserId, bidDTO.AuctionId, bidDTO.Amount)
	if err != nil {
		return err
	}

	bidChannel <- *bidEntity

	return nil
}

func getBatchSizeInterval() time.Duration {
	batchInsertInterval := os.Getenv("BATCH_INSERT_INTERVAL")
	if duration, err := time.ParseDuration(batchInsertInterval); err == nil {
		return duration
	}

	return 3 * time.Minute
}

func getMaxBatchSize() int {
	if batchSize, err := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE")); err == nil {
		return batchSize
	}

	return 5
}
