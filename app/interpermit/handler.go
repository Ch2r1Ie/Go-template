package interpermit

import (
	"context"
	"net/http"
	"time"

	"github.com/3169a7e4c0eb100/go-simple-template/app"
	"github.com/3169a7e4c0eb100/go-simple-template/httpclient"
	"github.com/gin-gonic/gin"
)

type userByTokenService interface {
	UserByToken(ctx context.Context, token string) (httpclient.Response[GetUserByTokenResponse], error)
}

type userStorage interface {
	SaveUser(context.Context, GetUserByTokenResponse) error
}

type Handler struct {
	srv   userByTokenService
	store userStorage
}

func NewHandler(srv userByTokenService, store userStorage) *Handler {
	return &Handler{srv: srv, store: store}
}

func (handler *Handler) PermitTransaction(c *gin.Context) {
	var req SubmitTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, app.Response{
			Message: err.Error(),
		})
		return
	}

	resp, err := handler.srv.UserByToken(c.Request.Context(), req.TokenHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, app.Response{
			Message: err.Error(),
		})
		return
	}
	if resp.Code == http.StatusBadRequest || resp.Code == http.StatusInternalServerError {
		c.JSON(http.StatusInternalServerError, app.Response{
			Code:    resp.Response.Code,
			Message: resp.Response.Message,
		})
		return
	}

	user := resp.Response

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	if err := handler.store.SaveUser(ctx, user); err != nil {
		c.JSON(http.StatusInternalServerError, app.Response{
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

type SubmitTransactionRequest struct {
	TokenHash   string          `json:"tokenHash"`
	Dopa        SubmitTransDopa `json:"dopa"`
	Transaction Transaction     `json:"transaction"`
}

type SubmitTransDopa struct {
	LaserId       string `json:"laserId"`
	ThaiFirstName string `json:"thaiFirstName"`
	ThaiLastName  string `json:"thaiLastName"`
}

type Transaction struct {
	MobileNo         string       `json:"mobileNo"`
	Email            string       `json:"email"`
	InterLicenseType string       `json:"interLicenseType"`
	LicenseDesc      string       `json:"licenseDesc"`
	MaillingAddr     MaillingAddr `json:"maillingAddress"`
}

type MaillingAddr struct {
	Address     string `json:"address"`
	SubDistrict string `json:"subDistrict"`
	District    string `json:"district"`
	Province    string `json:"province"`
	PostalCode  string `json:"postalCode"`
	Country     string `json:"country"`
}
