package usecase

import (
    "github.com/Arlandaren/easyfund/internal/config"
    "github.com/Arlandaren/easyfund/pkg/banking"
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

type AuthUsecase struct {
    config      *config.Config
    vbankClient *banking.VBankClient
}

func NewAuthUsecase(config *config.Config, vbankClient *banking.VBankClient) *AuthUsecase {
    return &AuthUsecase{
        config:      config,
        vbankClient: vbankClient,
    }
}

func (a *AuthUsecase) GetRandomDemoClient() (*banking.RandomClientResponse, error) {
    return a.vbankClient.GetRandomDemoClient()
}

// LoginDemoClient теперь принимает DemoClientLoginRequest с personID и password,
// чтобы использовать правильный пароль для логина.
func (a *AuthUsecase) LoginDemoClient(req *banking.DemoClientLoginRequest) (*banking.AuthResponse, error) {
    loginResp, err := a.vbankClient.LoginClient(req.PersonID, req.Password)
    if err != nil {
        return nil, err
    }

    user := banking.User{
        ID:       uuid.New().String(),
        Name:     req.FullName,
        Email:    req.PersonID + "@demo.easyfund.ru",
        Role:     "borrower",
        PersonID: req.PersonID,
    }

    // Передаем токен VBank в JWT
    token, err := a.generateJWT(user, "demo", loginResp.AccessToken)
    if err != nil {
        return nil, err
    }

    return &banking.AuthResponse{
        Token:     token,
        TokenType: "Bearer",
        User:      user,
    }, nil
}

func (a *AuthUsecase) Login(req *banking.LoginRequest) (*banking.AuthResponse, error) {
    if req.Email == "admin@easyfund.ru" && req.Password == "password123" {
        user := banking.User{
            ID:    uuid.New().String(),
            Name:  "Administrator",
            Email: req.Email,
            Role:  "admin",
        }

        token, err := a.generateJWT(user, "regular", "")
        if err != nil {
            return nil, fmt.Errorf("failed to generate JWT: %w", err)
        }

        return &banking.AuthResponse{
            Token:     token,
            TokenType: "Bearer",
            User:      user,
        }, nil
    }

    return nil, fmt.Errorf("invalid credentials")
}

func (a *AuthUsecase) Register(req *banking.RegisterRequest) (*banking.AuthResponse, error) {
    user := banking.User{
        ID:    uuid.New().String(),
        Name:  req.Name,
        Email: req.Email,
        Role:  req.Role,
    }

    token, err := a.generateJWT(user, "regular", "")
    if err != nil {
        return nil, fmt.Errorf("failed to generate JWT: %w", err)
    }

    return &banking.AuthResponse{
        Token:     token,
        TokenType: "Bearer",
        User:      user,
    }, nil
}

func (a *AuthUsecase) ValidateJWT(tokenString string) (*banking.JWTClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(a.config.JWTSecret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        jwtClaims := &banking.JWTClaims{
            UserID: claims["user_id"].(string),
            Email:  claims["email"].(string),
            Mode:   claims["mode"].(string),
            Role:   claims["role"].(string),
        }

        if personID, exists := claims["person_id"]; exists && personID != nil {
            jwtClaims.PersonID = personID.(string)
        }

        if vbankToken, exists := claims["vbank_token"]; exists && vbankToken != nil {
            jwtClaims.VBankAccessToken = vbankToken.(string)
        }

        return jwtClaims, nil
    }

    return nil, fmt.Errorf("invalid token")
}

func (a *AuthUsecase) generateJWT(user banking.User, mode, vbankToken string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "role":    user.Role,
        "mode":    mode,
        "exp":     time.Now().Add(time.Hour * time.Duration(a.config.JWTExpiryHours)).Unix(),
        "iat":     time.Now().Unix(),
    }

    if user.PersonID != "" {
        claims["person_id"] = user.PersonID
    }

    if vbankToken != "" {
        claims["vbank_token"] = vbankToken
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(a.config.JWTSecret))
}
