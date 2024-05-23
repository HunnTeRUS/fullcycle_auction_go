package auction_usecase

import (
	"context"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/entity/auction_entity"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/entity/bid_entity"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/internal_error"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/usecase/bid_usecase"
	"time"
)

type AuctionInputDTO struct {
	ProductName string           `json:"product_name" binding:"required,min=4"`
	Category    string           `json:"category" binding:"required,min=4"`
	Description string           `json:"description" binding:"required,min=10"`
	Condition   ProductCondition `json:"condition"`
}

type AuctionOutputDTO struct {
	Id          string           `json:"id" copier:"Id"`
	ProductName string           `json:"product_name" copier:"ProductName"`
	Category    string           `json:"category" copier:"Category"`
	Status      AuctionStatus    `json:"status" copier:"Status"`
	Description string           `json:"description" copier:"Description"`
	Condition   ProductCondition `json:"condition" copier:"Condition"`
	Timestamp   time.Time        `json:"timestamp" copier:"Timestamp" time_format:"2006-01-02 15:04:05"`
}

type WinningInfoOutputDTO struct {
	Auction AuctionOutputDTO          `json:"auction"`
	Bid     *bid_usecase.BidOutputDTO `json:"bid"`
}

type AuctionUseCase struct {
	auctionRepository auction_entity.AuctionRepositoryInterface
	bidRepository     bid_entity.BidRepositoryInterface
}

type ProductCondition int
type AuctionStatus int

type AuctionUseCaseInterface interface {
	CreateAuction(
		ctx context.Context,
		auction AuctionInputDTO) *internal_error.InternalError

	FindAuctions(
		ctx context.Context,
		status AuctionStatus,
		category string,
		productName string) ([]AuctionOutputDTO, *internal_error.InternalError)
}

func (bs *AuctionUseCase) CreateAuction(ctx context.Context, auctionInput AuctionInputDTO) *internal_error.InternalError {
	auction, err := auction_entity.CreateAuction(
		auctionInput.ProductName,
		auctionInput.Category,
		auctionInput.Description,
		auction_entity.ProductCondition(auctionInput.Condition))
	if err != nil {
		return nil
	}

	return bs.auctionRepository.CreateAuction(ctx, auction)
}

func NewAuctionService(
	auctionRepository auction_entity.AuctionRepositoryInterface,
	bidRepository bid_entity.BidRepositoryInterface) *AuctionUseCase {
	return &AuctionUseCase{
		auctionRepository, bidRepository,
	}
}
