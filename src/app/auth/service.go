package auth

import (
	"errors"
	"net/http"
	"os"
	"time"
	"tinder-clone/src/abstraction"
	"tinder-clone/src/dto"
	"tinder-clone/src/factory"
	"tinder-clone/src/models"
	"tinder-clone/src/repositories"
	"tinder-clone/utils/response"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type AuthService struct {
	UserRepository *repositories.UserRepository
	JWTSecretKey   []byte
}

func NewAuthService(f *factory.Factory) *AuthService {
	return &AuthService{
		UserRepository: f.UserRepository,
		JWTSecretKey:   []byte(os.Getenv("JWT_SECRET_KEY")),
	}
}

func (a *AuthService) Login(c *abstraction.Context, payload *dto.LoginRequestDto) (*dto.LoginResponseDto, error) {
	user, err := a.UserRepository.FindByEmail(c, payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrorWrap(response.CustomError(http.StatusNotFound, "invalid credentials"), err)
		}

		return nil, response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
	}

	if !user.CheckPasswordHash(payload.Password) {
		return nil, response.ErrorWrap(response.CustomError(http.StatusUnauthorized, "invalid credentials"), err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	accessToken, err := token.SignedString(a.JWTSecretKey)
	if err != nil {
		return nil, response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
	}

	data := &dto.LoginResponseDto{
		AccessToken: accessToken,
		ExpiresIn:   time.Now().Add(time.Hour * 24).Unix(),
	}

	return data, nil
}

func (a *AuthService) Register(c *abstraction.Context, payload *dto.RegisterRequestDto) (*dto.RegisterResponseDto, error) {
	_, err := a.UserRepository.FindByEmail(c, payload.Email)
	if err == nil {
		return nil, response.ErrorWrap(response.CustomError(http.StatusConflict, "user already registered"), errors.New("user already registered"))
	}

	userID, _ := uuid.NewV4()
	user := &models.User{
		Base: models.Base{
			ID: userID,
		},
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
		Age:      payload.Age,
	}

	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	if err := a.UserRepository.Create(c, user); err != nil {
		return nil, err
	}

	data := &dto.RegisterResponseDto{}
	if err := copier.Copy(data, user); err != nil {
		return nil, err
	}

	return data, nil
}
