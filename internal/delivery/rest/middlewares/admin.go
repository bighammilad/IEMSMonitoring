package middlewares

import (
	"fmt"
	"time"

	. "monitoring/internal/globals"
	hashpass "monitoring/pkg/hashPass"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// echo middleware to check jwt token
func IsAdminUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Get("user") != nil {
			tokenString := c.Get("user").(*jwt.Token)
			token, err := jwt.Parse(tokenString.Raw, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte("secret"), nil
			})
			if err != nil {
				return echo.ErrInternalServerError
			}
			jwtToken := token.Claims.(jwt.MapClaims)
			if _, ok := jwtToken["role"]; !ok {
				return echo.ErrForbidden
			}
			//check if jwtToken["role"] is float64
			if _, ok := jwtToken["role"].(float64); !ok {
				return echo.ErrForbidden
			}
			role := int(jwtToken["role"].(float64))
			if role == 1 {
				return next(c)
			} else {
				return echo.ErrForbidden
			}
		} else {
			return echo.ErrForbidden
		}
	}
}

func MakeAdminUser(c echo.Context) (err error) {
	var count int
	row, err := GlobalPG.QueryContext(c.Request().Context(), "select count(*) as count from Users")
	if err != nil {
		panic("Failed to execute query to check user count")
	}
	defer row.Close()

	for row.Next() {
		err = row.Scan(&count)
		if err != nil {
			panic("Failed to scan user count")
		}
	}
	if count == 0 {
		adminUsername := "admin"
		adminPassword := "admin_password"
		role := 1
		// hash password
		hashedPass, err := hashpass.HashPassword(adminPassword)
		if err != nil {
			panic("Failed to hash password")
		}
		// insert admin user
		_, err = GlobalPG.QueryContext(c.Request().Context(), "INSERT INTO users (username, password, role) VALUES ($1, $2, $3)", adminUsername, hashedPass, role)
		if err != nil {
			panic("Failed to insert admin user")
		}
		// create token
		token, err := MakeToken(adminUsername, role, 7)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Admin user created with username: %s and password: %s and the bearer token: %s\n", adminUsername, adminPassword, token)
	} else {
		fmt.Println("Admin user already exists")
	}
	return nil
}

func MakeToken(usname string, role int, expiredays int64) (tokenString string, err error) {
	type MyCustomClaims struct {
		Sm   string `json:"sm"`
		Role int    `json:"role"`
		jwt.RegisteredClaims
	}

	claims := MyCustomClaims{
		usname,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * time.Duration(expiredays))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err = token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
