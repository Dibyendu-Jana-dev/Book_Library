package middleware

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Payload struct {
	Role      string `json:"role"`
	Email     string `json:"email"`
	Source    string `json:"source"`
	AuthToken string `json:"authToken"`
	Jwt       string `json:"jwt"`
}

var contextKey = "contextMap"
var contextMap Payload
var INTERNAL = "INTERNAL"
var Secret string
var SecretKey []byte

func setSecret() {
	Secret = os.Getenv("SECRET_KEY")
	SecretKey = []byte(Secret)
}
func GenerateJWT(email, role string) (string, error) {
	setSecret()
	//var mySigningKey = []byte(os.Getenv("SECRET_KEY"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	log.Println("role", role)
	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	claims["iat"] = time.Now().Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				contextMap.AuthToken = "INVALID"
				ctx := context.WithValue(r.Context(), contextKey, contextMap)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
		}()
		setSecret()
		authToken := r.Header.Get("Authorization")
		contextMap.AuthToken = authToken
		log.Println("Authentication", authToken)
		split := strings.Split(authToken, " ")
		if len(split) > 1 {
			authToken = split[1]
		} else {
			authToken = split[0]
		}
		contextMap.Jwt = authToken
		if authToken == "" {
			contextMap.AuthToken = "INVALID"
			ctx := context.WithValue(r.Context(), contextKey, contextMap)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}
		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			//Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return SecretKey, nil
		})
		if err != nil {
			log.Println("err", err)
			contextMap.AuthToken = "INVALID"
			ctx := context.WithValue(r.Context(), contextKey, contextMap)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			claimKeyCheck := make(map[string]bool)
			for key := range claims {
				claimKeyCheck[key] = true
			}
			if claimKeyCheck["role"] {
				contextMap.Role = claims["role"].(string)
			}
			if claimKeyCheck["email"] {
				contextMap.Email = claims["email"].(string)
			}
			if claimKeyCheck["source"] {
				contextMap.Source = claims["source"].(string)
			} else {
				contextMap.Source = INTERNAL
			}
			ctx := context.WithValue(r.Context(), contextKey, contextMap)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else {
			contextMap.AuthToken = "INVALID"
			ctx := context.WithValue(r.Context(), contextKey, contextMap)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
	} )
}


func GetUserInfo(ctx context.Context) *Payload {
	value := ctx.Value(contextKey)
	raw := value.(Payload)
	return &raw
}
