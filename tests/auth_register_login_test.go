package tests

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	ssov1 "github.com/m1al04949/contracts/contracts/gen/go/sso"
	"github.com/m1al04949/sso-gRPC/internal/services/auth"
	"github.com/m1al04949/sso-gRPC/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID     = 0
	appID          = 1
	appSecret      = "test-secret"
	passDefaultLen = 10
	deltaSeconds   = 1
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePass()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err) // stop
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    appID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(toker *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))
	assert.InDelta(t, loginTime.Add(st.Cfg.JWT.TokenTTL).Unix(), int64(claims["exp"].(float64)), deltaSeconds)
}

func TestRegisterLogin_DuplicatedReg(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePass()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, auth.ErrUserExists.Error())
}

func TestRegisterLogin_InvalidCredentials(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Reguster with empty password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "Reguster with empty email",
			email:       "",
			password:    randomFakePass(),
			expectedErr: "email is required",
		},
		{
			name:        "Reguster with empty both credentials",
			email:       "",
			password:    "",
			expectedErr: "email is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int32
		expectedErr string
	}{
		{
			name:        "Login with empty password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       appID,
			expectedErr: "password is required",
		},
		{
			name:        "Login with empty email",
			email:       "",
			password:    randomFakePass(),
			appID:       appID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with empty both credentials",
			email:       "",
			password:    "",
			appID:       appID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with non-mathcing credentials",
			email:       gofakeit.Email(),
			password:    randomFakePass(),
			appID:       appID,
			expectedErr: "invalid email or password",
		}, {
			name:        "Login without appID",
			email:       gofakeit.Email(),
			password:    randomFakePass(),
			appID:       emptyAppID,
			expectedErr: "app_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    gofakeit.Email(),
				Password: randomFakePass(),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appID,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func randomFakePass() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
