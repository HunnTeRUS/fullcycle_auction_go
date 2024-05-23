package auction_controller

import (
	"context"
	"github.com/HunnTeRUS/fullcycle-auction-go/configuration/logger"
	"github.com/HunnTeRUS/fullcycle-auction-go/configuration/rest_err"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/usecase/auction_usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (sc *AuctionController) FindAuctions(c *gin.Context) {
	queryParams := struct {
		status      auction_usecase.AuctionStatus
		category    string
		productName string
	}{}

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		logger.Error("Error trying to validate query params",
			err,
		)
		errorMessage := rest_err.NewBadRequestError(
			"Query params are not valid",
		)

		c.JSON(errorMessage.Code, errorMessage)
		return
	}

	auctionsDomainList, err := sc.auctionUseCase.FindAuctions(
		context.Background(),
		queryParams.status,
		queryParams.category,
		queryParams.productName)
	if err != nil {
		errRest := rest_err.NewError(err.Message, err.Err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, auctionsDomainList)
}
