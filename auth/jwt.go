package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ranggadablues/gosok/common"
	"google.golang.org/grpc/metadata"
)

var (
	accessSecret    = []byte(os.Getenv("ACCESS_SECRET"))  // load from env in real deployment
	refreshSecret   = []byte(os.Getenv("REFRESH_SECRET")) // separate key for refresh token
	ErrTokenExpired = errors.New("token is expired")
)

type Claims struct {
	UserInfo map[string]string `json:"userinfo"`
	jwt.RegisteredClaims
}

type contextKey string

const ClaimsContextKey contextKey = "jwt_claims"

// ---------------------------
// 🔸 Generate access + refresh pair
// ---------------------------
func GenerateTokenPair(userInfo map[string]string) (string, string, error) {
	// Access token expires fast
	accessClaims := &Claims{
		UserInfo: userInfo,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user-service",
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(accessSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh token lasts longer
	refreshClaims := &Claims{
		UserInfo: userInfo,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user-service",
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(refreshSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ---------------------------
// 🔸 Validate token (access or refresh)
// ---------------------------
func ValidateAccessToken(tokenStr string) (*Claims, error) {
	return validateToken(tokenStr)
}

func ValidateRefreshToken(tokenStr string) (*Claims, error) {
	return validateToken(tokenStr)
}

func validateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return accessSecret, nil
	})

	if err != nil {
		// Handle expiration separately
		if errors.Is(err, jwt.ErrTokenExpired) {
			return claims, jwt.ErrTokenExpired
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	tokenClaim := token.Claims.(*Claims)
	return tokenClaim, nil
}

// ---------------------------
// 🔸 Get claims from context
// ---------------------------
func GetClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(ClaimsContextKey).(*Claims)
	return claims, ok
}

func InjectToGRPCContext(ctx context.Context) context.Context {
	claims, ok := GetClaimsFromContext(ctx)
	if !ok {
		return ctx
	}
	md := metadata.New(claims.UserInfo)
	return metadata.NewOutgoingContext(ctx, md)
}

func IncomingContext(ctx context.Context, out interface{}) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("no metadata found")
	}

	auth := map[string]string{}
	for k := range md {
		val := md.Get(k)
		if len(val) > 0 {
			auth[k] = val[0]
		}
	}

	return common.MapToStruct(auth, out)
}
