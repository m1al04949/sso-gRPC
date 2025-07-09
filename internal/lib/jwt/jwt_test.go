package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/m1al04949/sso-gRPC/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewToken_Happy(t *testing.T) {
	user := models.User{
		ID:    1,
		Email: "test@example.com",
	}

	app := models.App{
		ID:     1,
		Secret: "test-secret",
	}

	ttl := 1 * time.Hour

	t.Run("successful token generation", func(t *testing.T) {
		token, err := NewToken(user, app, ttl)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(app.Secret), nil
		})
		require.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		require.True(t, ok)

		assert.Equal(t, float64(user.ID), claims["uid"])
		assert.Equal(t, user.Email, claims["email"])
		assert.Equal(t, float64(app.ID), claims["app_id"])

		expectedExp := time.Now().Add(ttl).Unix()
		actualExp := int64(claims["exp"].(float64))
		assert.InDelta(t, expectedExp, actualExp, 1)
	})
}
