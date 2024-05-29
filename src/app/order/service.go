package order

import (
	"errors"
	"net/http"
	"time"
	"tinder-clone/src/abstraction"
	"tinder-clone/src/dto"
	"tinder-clone/src/factory"
	"tinder-clone/src/models"
	"tinder-clone/src/repositories"
	"tinder-clone/utils/response"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type OrderService struct {
	*repositories.PremiumPackageRepository
	*repositories.UserOrderRepository
	*repositories.UserRepository
}

func NewOrderService(f *factory.Factory) *OrderService {
	return &OrderService{
		PremiumPackageRepository: f.PremiumPackageRepository,
		UserOrderRepository:      f.UserOrderRepository,
		UserRepository:           f.UserRepository,
	}
}

func (o *OrderService) CreateOrder(c *abstraction.Context, payload *dto.OrderRequestDto) (*dto.OrderResponseDto, error) {
	premiumPackage, err := o.PremiumPackageRepository.FindByID(c, payload.PremiumPackageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrorWrap(response.CustomError(http.StatusNotFound, err.Error()), err)
		}

		return nil, err
	}

	ordersCount, err := o.UserOrderRepository.CountOrders(c)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.FromString(c.AuthJwt.ID)
	if err != nil {
		return nil, err
	}

	userOrder := &models.UserOrder{
		Base: models.Base{
			ID: uuid.Must(uuid.NewV4()),
		},
		PremiumPackageID: premiumPackage.ID,
		UserID:           userID,
		TotalAmount:      premiumPackage.Amount,
		OrderTime:        time.Now(),
		Status:           "pending",
	}
	userOrder.GenerateInvoiceNumber(ordersCount)

	if err := o.UserOrderRepository.Create(c, userOrder); err != nil {
		return nil, response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
	}

	data := &dto.OrderResponseDto{
		PremiumPackage: &dto.PremiumPackageResponseDto{},
	}
	if err := copier.Copy(data, userOrder); err != nil {
		return nil, response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
	}

	if err := copier.Copy(data.PremiumPackage, premiumPackage); err != nil {
		return nil, response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
	}

	return data, nil
}

func (o *OrderService) Payment(c *abstraction.Context, ID string) (*dto.OrderResponseDto, error) {
	order, err := o.UserOrderRepository.FindByID(c, ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrorWrap(response.CustomError(http.StatusNotFound, err.Error()), err)
		}

		return nil, response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
	}

	now := time.Now()
	expire := now.AddDate(0, 1, 0)
	updatedOrder := &models.UserOrder{
		Base: models.Base{
			ID: order.ID,
		},
		PaymentTime: &now,
		ExpireTime:  &expire,
		Status:      "success",
	}

	if err := o.UserOrderRepository.Update(c, updatedOrder); err != nil {
		return nil, response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
	}

	updatedUser := &models.User{
		Base: models.Base{
			ID: order.UserID,
		},
	}
	if order.Package.Name == "No Swipe Quota" {
		_false := false
		updatedUser.HasSwipeLimit = &_false
	} else {
		_true := true
		updatedUser.IsVerified = &_true
	}

	if err := o.UserRepository.Update(c, updatedUser); err != nil {
		return nil, response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
	}

	data := &dto.OrderResponseDto{}
	if err := copier.Copy(data, order); err != nil {
		return nil, response.ErrorWrap(response.CustomError(http.StatusInternalServerError, err.Error()), err)
	}

	data.PaymentTime = updatedOrder.PaymentTime
	data.ExpireTime = updatedOrder.ExpireTime
	data.Status = updatedOrder.Status

	return data, nil
}
