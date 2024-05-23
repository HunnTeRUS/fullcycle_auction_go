package bid_controller

import (
	"context"
	"github.com/HunnTeRUS/fullcycle-auction-go/configuration/logger"
	"github.com/HunnTeRUS/fullcycle-auction-go/configuration/rest_err"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (bc *BidController) FindBidsByAuctionIdController(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		logger.Error("Error trying to validate auctionId",
			err,
		)
		errorMessage := rest_err.NewBadRequestError(
			"AuctionId is not a valid id",
		)

		c.JSON(errorMessage.Code, errorMessage)
		return
	}

	auctions, err := bc.bidUseCase.FindBidsByAuctionId(context.Background(), auctionId)
	if err != nil {
		errRest := rest_err.NewError(err.Message, err.Err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, auctions)
}

func (bc *BidController) FindWinningBidByAuctionIdController(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		logger.Error("Error trying to validate auctionId",
			err,
		)
		errorMessage := rest_err.NewBadRequestError(
			"AuctionId is not a valid id",
		)

		c.JSON(errorMessage.Code, errorMessage)
		return
	}

	auction, err := bc.bidUseCase.FindWinningBidByAuctionId(context.Background(), auctionId)
	if err != nil {
		errRest := rest_err.NewError(err.Message, err.Err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, auction)
}
