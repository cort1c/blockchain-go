package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
)

type (
	NewTransactionRequest struct {
		Sender    string `validate:"required"`
		Recipient string `validate:"required"`
		Amount    int    `validate:"required"`
	}

	NewTransactionResponse struct {
		Message string
	}

	FullChainResponse struct {
		Chain  *[]Block
		Length int
	}

	BlockchainContext struct {
		echo.Context
		Blockchain *Blockchain
	}

	CustomValidator struct {
		Validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func main() {
	blockchain := New()
	e := echo.New()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &BlockchainContext{c, blockchain}
			return h(cc)
		}
	})
	e.Validator = &CustomValidator{validator.New()}
	e.GET("/", rootHandler)
	e.POST("/transactions/new", newTransactionHandler)
	e.GET("/chain", fullChainHandler)
	e.Logger.Fatal(e.Start(":1323"))
}

func rootHandler(c echo.Context) error {
	return c.HTML(http.StatusOK, "<b>Hello, World!</b>")
}

func newTransactionHandler(c echo.Context) error {
	request := &NewTransactionRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}
	if err := c.Validate(request); err != nil {
		return err
	}
	cc := c.(*BlockchainContext)
	index := cc.Blockchain.NewTransaction(request.Sender, request.Recipient, request.Amount)
	message := fmt.Sprintf("Transaction will be added to Block %v.", index)
	return c.JSON(http.StatusOK, &NewTransactionResponse{message})
}

func fullChainHandler(c echo.Context) error {
	cc := c.(*BlockchainContext)
	return c.JSON(http.StatusOK, &FullChainResponse{&cc.Blockchain.Chain, len(cc.Blockchain.Chain)})
}
