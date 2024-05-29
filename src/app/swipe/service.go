package swipe

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"tinder-clone/src/abstraction"
	"tinder-clone/src/dto"
	"tinder-clone/src/factory"
	"tinder-clone/src/models"
	"tinder-clone/src/repositories"
	"tinder-clone/utils/response"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
)

type SwipeService struct {
	*repositories.UserMatchRepository
	*repositories.UserRepository

	RedisClient *redis.Client
}

func NewSwipeService(f *factory.Factory) *SwipeService {
	return &SwipeService{
		UserMatchRepository: f.UserMatchRepository,
		UserRepository:      f.UserRepository,

		RedisClient: f.RedisClient,
	}
}
func (s *SwipeService) FindPotentialMatches(c *abstraction.Context, payload *dto.PotentialMatchRequestDto) (*dto.UserResponseDto, error) {
	user, err := s.UserRepository.FindByID(c, c.AuthJwt.ID)
	if err != nil {
		return nil, response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
	}

	if *user.IsExceedDailySwipeLimit {
		return nil, response.ErrorWrap(response.CustomError(http.StatusForbidden, "daily swipe limit exceeded"), nil)
	}

	minLatitude := payload.Latitude + 0.005
	maxLatitude := payload.Latitude
	minLongitude := payload.Longitude - 0.005
	maxLongitude := payload.Longitude
	potentialMatchKey := fmt.Sprintf("potential-matches:%s", c.AuthJwt.ID)

	if s.RedisClient.Exists(c, potentialMatchKey).Val() == 0 {
		matchedUserIDs, err := s.UserMatchRepository.FindMatchedUser(c, c.AuthJwt.ID)
		if err != nil {
			return nil, response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
		}

		gender := "Male"
		if user.Gender == "Male" {
			gender = "Female"
		}

		potentialUsers, err := s.UserRepository.Find(c, &dto.UserFilterDto{Age: user.Age, MatchedUserIDs: matchedUserIDs, Gender: gender})
		if err != nil {
			return nil, response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
		}

		for _, potentialUser := range potentialUsers {
			if potentialUser.ID.String() == c.AuthJwt.ID {
				continue
			}

			rangeLatitude := s.generateRange(minLatitude, maxLatitude)
			rangeLongitude := s.generateRange(minLongitude, maxLongitude)

			s.RedisClient.GeoAdd(c, potentialMatchKey, &redis.GeoLocation{
				Name:      potentialUser.ID.String(),
				Longitude: rangeLongitude,
				Latitude:  rangeLatitude,
			})
		}
		s.RedisClient.Expire(c, potentialMatchKey, 5*time.Minute)
	}

	radius := 100
	potentialMatches := []string{}
	for i := 0; i < 5; i++ {
		potentialMatches = s.RedisClient.GeoSearch(c, potentialMatchKey, &redis.GeoSearchQuery{
			Longitude:  payload.Longitude,
			Latitude:   payload.Latitude,
			Radius:     float64(radius),
			RadiusUnit: "m",
			Count:      15,
		}).Val()

		if len(potentialMatches) > 0 {
			break
		}

		radius += 200
	}

	if len(potentialMatches) == 0 {
		return nil, response.ErrorWrap(response.CustomError(http.StatusNotFound, "no potential matches found"), err)
	}

	potentialMatchID := ""
	userMatchesKey := fmt.Sprintf("user-matches:%s", c.AuthJwt.ID)
	if s.RedisClient.Exists(c, userMatchesKey).Val() == 0 {
		potentialMatchID = potentialMatches[0]
	} else {
		matchedUsers := s.RedisClient.SMembers(c, userMatchesKey).Val()
		mapMatchedUsers := make(map[string]bool)

		for _, matchedUser := range matchedUsers {
			mapMatchedUsers[matchedUser] = true
		}

		for _, potentialMatch := range potentialMatches {
			if _, ok := mapMatchedUsers[potentialMatch]; !ok {
				potentialMatchID = potentialMatch
				break
			}
		}
	}

	if potentialMatchID == "" {
		return nil, response.ErrorWrap(response.CustomError(http.StatusNotFound, "no potential matches found"), err)
	}
	potentialUser, err := s.UserRepository.FindByID(c, potentialMatchID)

	if err != nil {
		return nil, response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
	}

	data := &dto.UserResponseDto{}
	if err := copier.Copy(&data, potentialUser); err != nil {
		return nil, response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
	}

	return data, nil
}

func (s *SwipeService) generateRange(min, max float64) float64 {
	value := min + rand.Float64()*(max-min)
	valueString := fmt.Sprintf("%.6f", value)
	value, _ = strconv.ParseFloat(valueString, 64)

	return value
}

func (s *SwipeService) Match(c *abstraction.Context, payload *dto.MatchRequestDto) error {
	matchedUser, err := s.UserRepository.FindByID(c, payload.MatchedUserID)
	if err != nil {
		return response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
	}

	matchedUserModel := &models.UserMatch{
		Base: models.Base{
			ID: uuid.Must(uuid.NewV4()),
		},
		UserID:      uuid.Must(uuid.FromString(c.AuthJwt.ID)),
		MatchUserID: matchedUser.ID,
		ExpireTime:  time.Now().Add(24 * time.Hour),
		IsMatch:     true,
	}

	if payload.Direction == "left" {
		matchedUserModel.IsMatch = false
	}

	if err := s.UserMatchRepository.Create(c, matchedUserModel); err != nil {
		return response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
	}

	s.RedisClient.ZRem(c, fmt.Sprintf("potential-matches:%s", c.AuthJwt.ID), matchedUser.ID.String())
	s.RedisClient.SAdd(c, fmt.Sprintf("user-matches:%s", c.AuthJwt.ID), matchedUser.ID.String())
	s.RedisClient.Expire(c, fmt.Sprintf("user-matches:%s", c.AuthJwt.ID), 5*time.Minute)

	user, err := s.UserRepository.FindByID(c, c.AuthJwt.ID)
	if err != nil {
		return response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
	}

	if *user.HasSwipeLimit {
		if s.RedisClient.SCard(c, fmt.Sprintf("user-matches:%s", c.AuthJwt.ID)).Val() >= 10 {
			_true := true
			updatedUser := &models.User{
				Base: models.Base{
					ID: user.ID,
				},
				IsExceedDailySwipeLimit: &_true,
			}

			if err := s.UserRepository.Update(c, updatedUser); err != nil {
				return response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
			}
		}
	}

	return nil
}
