package bid_controller

import (
	"context"
	"github.com/HunnTeRUS/fullcycle-auction-go/configuration/logger"
	"github.com/HunnTeRUS/fullcycle-auction-go/configuration/rest_err"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/usecase/bid_usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BidController struct {
	bidUseCase bid_usecase.BidUseCaseInterface
}

func NewBidController(useCase bid_usecase.BidUseCaseInterface) *BidController {
	return &BidController{useCase}
}

func (bc *BidController) CreateBid(c *gin.Context) {
	var bidRequest bid_usecase.BidInputDTO

	if err := c.ShouldBindJSON(&bidRequest); err != nil {
		errRest := rest_err.NewInternalServerError("error trying to parse bid_usecase data")
		logger.Error(errRest.Message, err)
		c.JSON(errRest.Code, errRest)
		return
	}

	if err := bc.bidUseCase.CreateBid(context.Background(), bidRequest); err != nil {
		errRest := rest_err.NewError(err.Message, err.Err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.Status(http.StatusCreated)
}
