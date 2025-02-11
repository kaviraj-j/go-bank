package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/kaviraj-j/go-bank/db/sqlc"
	"github.com/kaviraj-j/go-bank/token"
	"github.com/lib/pq"
)

type createAccountRequestParams struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD INR EUR CAD"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequestParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPaylod := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPaylod.Username,
		Balance:  "0",
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid owner")))
				return

			case "unique_violation":
				ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("account with the same currency exists")))
				return
			}

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)

}

type getAccountRequestParams struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequestParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPaylod := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if account.Owner != authPaylod.Username {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("you are not authorized to get this account")))
	}

	ctx.JSON(http.StatusOK, account)

}

type getAccountsRequestParams struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=20"`
}

func (server *Server) getAccounts(ctx *gin.Context) {
	var req getAccountsRequestParams
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPaylod := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.GetAccountsParams{
		Owner:  authPaylod.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	accounts, err := server.store.GetAccounts(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
