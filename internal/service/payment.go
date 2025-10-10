package service

import (
	"event-management/entity"
	"event-management/internal/repository"
	"event-management/model"
	"event-management/pkg/database/mariadb"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type IPaymentService interface {
	CreatePayment(tx *gorm.DB, registrationID uuid.UUID, amount float64) (*model.CreatePaymentResponse, error)
	HandleNotification(notificationPayload map[string]interface{}, regService IRegistrationService) error
}

type PaymentService struct {
	db                *gorm.DB
	PaymentRepository repository.IPaymentRepository
	snapClient        snap.Client
	coreAPIClient     coreapi.Client
}

func NewPaymentService(paymentRepo repository.IPaymentRepository, snapClint snap.Client, coreAPIClient coreapi.Client) IPaymentService {
	return &PaymentService{
		db:                mariadb.Connection,
		PaymentRepository: paymentRepo,
		snapClient:        snapClint,
		coreAPIClient:     coreAPIClient,
	}
}

func (s *PaymentService) CreatePayment(tx *gorm.DB, registrationID uuid.UUID, amount float64) (*model.CreatePaymentResponse, error) {
	var result *model.CreatePaymentResponse

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  s.generateOrderID(registrationID.String()),
			GrossAmt: int64(amount),
		},
	}

	snapResp, midtransErr := s.snapClient.CreateTransaction(req)
	if midtransErr != nil {
		return nil, midtransErr
	}

	payment := &entity.Payment{
		OrderID:        req.TransactionDetails.OrderID,
		RegistrationID: registrationID,
		Amount:         amount,
		Status:         "pending",
		SnapURL:        snapResp.RedirectURL,
	}

	_, err := s.PaymentRepository.CreatePayment(tx, payment)
	if err != nil {
		return nil, err
	}

	result = &model.CreatePaymentResponse{
		SnapURL: snapResp,
	}

	return result, nil
}

func (s *PaymentService) HandleNotification(notificationPayload map[string]interface{}, regService IRegistrationService) error {
	orderID, exist := notificationPayload["order_id"].(string)
	if !exist {
		return fmt.Errorf("invalid notification payload: order_id not found")
	}

	transactionStatus, _ := s.coreAPIClient.CheckTransaction(orderID)
	if transactionStatus == nil {
		return fmt.Errorf("failed to get transaction status from midtrans")
	}

	if transactionStatus.StatusCode != "200" && transactionStatus.StatusCode != "201" {
		return fmt.Errorf("failed to check transaction status from midtrans")
	}

	tx := s.db.Begin()
	defer tx.Rollback()

	payment, err := s.PaymentRepository.GetPaymentByID(tx, orderID)
	if err != nil {
		return err
	}

	if payment.Status == "paid" {
		return nil
	}

	statusFromPayload, ok := notificationPayload["transaction_status"].(string)
	if !ok {
		return fmt.Errorf("transaction_status not found in notification payload")
	}

	switch statusFromPayload {
	case "capture", "settlement":
		payment.Status = "paid"
		paidAt := time.Now()
		payment.PaidAt = paidAt
	case "pending":
		payment.Status = "pending"
	default:
		payment.Status = "failed"
	}

	paymentType, ok := notificationPayload["payment_type"].(string)
	if ok {
		payment.PaymentType = paymentType
	}

	_, err = s.PaymentRepository.UpdatePayment(tx, payment)
	if err != nil {
		return err
	}

	if payment.Status == "paid" {
		err = regService.UpdateStatusAfterPayment(tx, payment.RegistrationID, "approved")
		if err != nil {
			return err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (s *PaymentService) generateOrderID(registrationID string) string {
	timestamp := time.Now().Unix()
	randomNum := rand.Intn(1000)

	return fmt.Sprintf("EVENT-%s-%d-%d", registrationID[:4], timestamp, randomNum)
}
