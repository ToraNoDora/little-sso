package testing

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ssov1 "github.com/ToraNoDora/little-sso-protos/gen/go/sso"

	h "github.com/ToraNoDora/little-sso/sso/tests/helper"
	"github.com/ToraNoDora/little-sso/sso/tests/suite"
)

const (
	emptyAppID        = ""
	passDefaultLength = 10
)

var ra = h.GetRandomApp(suite.Cfg.Store)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.SignUp(ctx, &ssov1.SignUpRequest{
		Username: gofakeit.Username(),
		Email:    email,
		Password: pass,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.SignIn(ctx, &ssov1.SignInRequest{
		Email:    email,
		Password: pass,
		AppId:    ra.AppID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(ra.AppSecret), nil
		},
	)
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), claims["uid"].(string))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, ra.AppID, claims["app_id"].(string))

	const deltaSecond = 1
	assert.InDelta(
		t,
		loginTime.Add(suite.Cfg.TokenTtl).Unix(),
		claims["exp"].(float64),
		deltaSecond,
	)
}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	username := gofakeit.Username()
	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.SignUp(
		ctx,
		&ssov1.SignUpRequest{
			Username: username,
			Email:    email,
			Password: pass,
		},
	)

	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.SignUp(
		ctx,
		&ssov1.SignUpRequest{
			Username: username,
			Email:    email,
			Password: pass,
		},
	)

	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "internal error")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Register with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "Register with Empty Email",
			email:       "",
			password:    randomFakePassword(),
			expectedErr: "email is required",
		},
		{
			name:        "Register with Both Empty",
			email:       "",
			password:    "",
			expectedErr: "email is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.SignUp(ctx, &ssov1.SignUpRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)

		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       string
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       ra.AppID,
			expectedErr: "password is required",
		},
		{
			name:        "Login with Empty Email",
			email:       "",
			password:    randomFakePassword(),
			appID:       ra.AppID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with Both Empty Email and Password",
			email:       "",
			password:    "",
			appID:       ra.AppID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with Non-Matching Password",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       ra.AppID,
			expectedErr: "failed to login",
		},
		{
			name:        "Login without AppID",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       emptyAppID,
			expectedErr: "app_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.SignUp(ctx, &ssov1.SignUpRequest{
				Username: gofakeit.Username(),
				Email:    gofakeit.Email(),
				Password: randomFakePassword(),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.SignIn(ctx, &ssov1.SignInRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appID,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, true, passDefaultLength)
}
