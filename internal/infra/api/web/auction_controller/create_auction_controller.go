package auction_controller

import (
	"context"
	"github.com/HunnTeRUS/fullcycle-auction-go/configuration/logger"
	"github.com/HunnTeRUS/fullcycle-auction-go/configuration/rest_err"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/usecase/auction_usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuctionController struct {
	auctionUseCase auction_usecase.AuctionUseCaseInterface
}

func NewAuctionController(auctionUseCase auction_usecase.AuctionUseCaseInterface) *AuctionController {
	return &AuctionController{auctionUseCase: auctionUseCase}
}

func (sc *AuctionController) CreateAuctionController(c *gin.Context) {
	var auctionRequest auction_usecase.AuctionInputDTO

	if err := c.ShouldBindJSON(&auctionRequest); err != nil {
		errRest := rest_err.NewInternalServerError("error trying to parse auction_usecase data")
		logger.Error(errRest.Message, err)
		c.JSON(errRest.Code, errRest)
		return
	}

	if err := sc.auctionUseCase.CreateAuction(context.Background(), auctionRequest); err != nil {
		errRest := rest_err.NewError(err.Message, err.Err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.Status(http.StatusCreated)
}
