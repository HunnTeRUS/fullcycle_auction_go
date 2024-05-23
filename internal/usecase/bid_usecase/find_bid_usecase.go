package bid_usecase

import (
	"context"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/internal_error"
	"github.com/jinzhu/copier"
)

func (bs *BidUseCase) FindBidsByAuctionId(ctx context.Context, auctionId string) ([]BidOutputDTO, *internal_error.InternalError) {
	var bidOutput []BidOutputDTO

	bid, err := bs.BidRepository.FindBidsByAuctionId(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	if err := copier.Copy(&bidOutput, bid); err != nil {
		return nil, internal_error.NewInternalServerError(err.Error())
	}

	return bidOutput, nil
}

func (bs *BidUseCase) FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*BidOutputDTO, *internal_error.InternalError) {
	bid, err := bs.BidRepository.FindWinningBidByAuctionId(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	bidOutput := &BidOutputDTO{
		UserId:    bid.UserId,
		AuctionId: bid.AuctionId,
		Amount:    bid.Amount,
		Timestamp: bid.Timestamp,
	}

	return bidOutput, nil
}
