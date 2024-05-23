package auction_entity

import (
	"context"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/entity/bid_entity"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/entity/user_entity"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/internal_error"
	"github.com/google/uuid"
	"time"
)

type Auction struct {
	Id          string           `copier:"ID"`
	ProductName string           `copier:"ProductName"`
	Category    string           `copier:"Category"`
	Status      AuctionStatus    `copier:"Status"`
	Description string           `copier:"Description"`
	Condition   ProductCondition `copier:"Condition"`
	Timestamp   time.Time        `copier:"Timestamp"`
}

type WinningInfo struct {
	User    *user_entity.User
	Auction Auction
	Bid     *bid_entity.Bid
}

func CreateAuction(productName, category, description string, condition ProductCondition) (*Auction, *internal_error.InternalError) {
	auction := &Auction{
		Id:          uuid.New().String(),
		ProductName: productName,
		Category:    category,
		Status:      Active,
		Description: description,
		Condition:   condition,
		Timestamp:   time.Now(),
	}

	if err := auction.Validate(); err != nil {
		return nil, err
	}

	return auction, nil
}

func (auc *Auction) Validate() *internal_error.InternalError {
	if len(auc.ProductName) <= 1 &&
		len(auc.Category) <= 1 &&
		len(auc.Description) <= 10 && (auc.Condition != New &&
		auc.Condition != Used &&
		auc.Condition != Refurbished) {
		return internal_error.NewBadRequestError("invalid auction object, data is incorrect")
	}

	return nil
}

type AuctionRepositoryInterface interface {
	CreateAuction(ctx context.Context, auction *Auction) *internal_error.InternalError

	FindAuctions(
		ctx context.Context,
		status AuctionStatus,
		category string,
		productName string) ([]Auction, *internal_error.InternalError)

	FindAuctionByAuctionId(
		ctx context.Context,
		auctionId string) (*Auction, *internal_error.InternalError)
}

type ProductCondition int
type AuctionStatus int

const (
	New ProductCondition = iota
	Used
	Refurbished
)

const (
	Active AuctionStatus = iota
	Completed
)
