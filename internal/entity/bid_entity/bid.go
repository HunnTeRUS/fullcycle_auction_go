package bid_entity

import (
	"context"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/internal_error"
	"github.com/google/uuid"
	"time"
)

type Bid struct {
	Id        string    `copier:"Id"`
	UserId    string    `copier:"UserId"`
	AuctionId string    `copier:"AuctionId"`
	Amount    float64   `copier:"Amount"`
	Timestamp time.Time `copier:"Timestamp"`
}

func CreateBid(userId, auctionId string, amount float64) (*Bid, *internal_error.InternalError) {
	entityBid := &Bid{
		Id:        uuid.New().String(),
		UserId:    userId,
		AuctionId: auctionId,
		Amount:    amount,
		Timestamp: time.Now(),
	}

	if err := entityBid.Validate(); err != nil {
		return nil, err
	}

	return entityBid, nil
}

type BidRepositoryInterface interface {
	CreateBids(ctx context.Context, bids []Bid) *internal_error.InternalError

	FindBidsByAuctionId(context.Context, string) ([]Bid, *internal_error.InternalError)
	FindWinningBidByAuctionId(context.Context, string) (*Bid, *internal_error.InternalError)
}

func (bd *Bid) Validate() *internal_error.InternalError {
	if uuid.Validate(bd.UserId) != nil {
		return internal_error.NewBadRequestError("userId is not a valid identifier")
	} else if uuid.Validate(bd.AuctionId) != nil {
		return internal_error.NewBadRequestError("auctionId is not a valid identifier")
	} else if bd.Amount <= 0 {
		return internal_error.NewBadRequestError("amount is not a valid price value")
	}

	return nil
}
