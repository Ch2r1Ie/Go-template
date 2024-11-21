package interpermit

import (
	"context"
	"net/http"
	"time"

	"github.com/3169a7e4c0eb100/go-simple-template/app"
	"github.com/3169a7e4c0eb100/go-simple-template/httpclient"
)

type userService struct {
	url    string
	client *http.Client
}

func NewUserService(client *http.Client) *userService {
	return &userService{client: client}
}

func (a *userService) UserByToken(ctx context.Context, token string) (httpclient.Response[GetUserByTokenResponse], error) {
	req := SubmitTransactionRequest{
		TokenHash: token,
	}
	return httpclient.Post[SubmitTransactionRequest, GetUserByTokenResponse](ctx, a.client, a.url, req)
}

type GetUserByTokenResponse struct {
	app.Response
	Data *UserProfile `json:"data,omitempty"`
}

type UserProfile struct {
	CifNo                string    `json:"cifNo"`
	CdiToken             string    `json:"cdiToken"`
	CitizenID            string    `json:"citizenId"`
	MobileNo             string    `json:"mobileNo"`
	Token                string    `json:"token"`
	TokenExpiredAt       time.Time `json:"tokenExpiredAt"`
	DopaCount            int16     `json:"dopaCount"`
	FaceCount            int16     `json:"faceCount"`
	IdentityCountResetAt time.Time `json:"identityCountResetAt"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}
