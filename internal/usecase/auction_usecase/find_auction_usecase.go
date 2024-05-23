package auction_usecase

import (
	"context"
	"fmt"
	"github.com/HunnTeRUS/fullcycle-auction-go/configuration/logger"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/entity/auction_entity"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/internal_error"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/usecase/bid_usecase"
	"github.com/jinzhu/copier"
)

func (bs *AuctionUseCase) FindAuctions(
	ctx context.Context,
	status AuctionStatus,
	category string,
	productName string) ([]AuctionOutputDTO, *internal_error.InternalError) {
	auctionsEntity, err := bs.auctionRepository.FindAuctions(
		ctx, auction_entity.AuctionStatus(status), category, productName)
	if err != nil {
		return nil, err
	}

	var auctionOutputDTO []AuctionOutputDTO
	if err := copier.Copy(&auctionOutputDTO, auctionsEntity); err != nil {
		return nil, internal_error.NewInternalServerError(err.Error())
	}

	return auctionOutputDTO, nil
}

func (bs *AuctionUseCase) FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*WinningInfoOutputDTO, *internal_error.InternalError) {
	auctionDomain, err := bs.auctionRepository.FindAuctionByAuctionId(ctx, auctionId)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying to find the auction_usecase with id %s", auctionId), err)
		return nil, err
	}

	var auctionOutputDTO AuctionOutputDTO
	if err := copier.Copy(&auctionOutputDTO, auctionDomain); err != nil {
		return nil, internal_error.NewInternalServerError(err.Error())
	}

	winningBid, err := bs.bidRepository.FindWinningBidByAuctionId(ctx, auctionId)
	if err != nil {
		return &WinningInfoOutputDTO{
			Auction: auctionOutputDTO,
			Bid:     nil,
		}, nil
	}

	var bidOutputDTO bid_usecase.BidOutputDTO
	if err := copier.Copy(&bidOutputDTO, winningBid); err != nil {
		return nil, internal_error.NewInternalServerError(err.Error())
	}

	return &WinningInfoOutputDTO{
		Auction: auctionOutputDTO,
		Bid:     &bidOutputDTO,
	}, nil
}
