package jwt

import (
	"github.com/golang-jwt/jwt"
	"time"
)

var jwtSecret = []byte("woodwhale&sheepbotany")

const (
	AccessExpireTime  = 2 * time.Hour
	RefreshExpireTime = 7 * 24 * time.Hour
)

// Claims
// @Description: jwt claims
type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// GenAccessAndRefreshJWT
// @Description: 生成access token和refresh token
// @param id uint
// @param username string
// @param role string
// @return string
// @return string
// @return error
func GenAccessAndRefreshJWT(id int64, username string, role string) (string, string, error) {
	nowTime := time.Now()
	accessToken, err := GenAccessToken(nowTime, id, username, role)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := GenRefreshToken(nowTime, id, username, role)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func GenAccessToken(nowTime time.Time, id int64, username, role string) (string, error) {
	return genToken(nowTime, AccessExpireTime, id, username, role, "A1natawd")
}

func GenRefreshToken(nowTime time.Time, id int64, username, role string) (string, error) {
	return genToken(nowTime, RefreshExpireTime, id, username, role, "A1natawd")
}

func genToken(nowTime time.Time, expire time.Duration, id int64, username, role, issuer string) (string, error) {
	expireTime := nowTime.Add(expire)
	claims := Claims{
		ID:       id,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    issuer,
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

// ParseJWT
// @Description: 验证token
// @param token string
// @return *Claims
// @return error
func ParseJWT(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			// 只有token有效才返回claims
			return claims, nil
		}
	}
	return nil, err
}
