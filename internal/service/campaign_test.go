package service

import (
	"testing"
	"time"

	"github.com/asymptoter/tonx-take-home-test/internal/repository"
	"github.com/asymptoter/tonx-take-home-test/internal/repository/mocks"
	"github.com/asymptoter/tonx-take-home-test/pkg/ctx"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const (
	mockCTX = mock.Anything
)

type campaignServiceSuite struct {
	suite.Suite
	ctx     ctx.CTX
	repo    *mocks.CampaignRepository
	service CampaignService
}

func (s *campaignServiceSuite) SetupSuite() {
	s.ctx = ctx.Background()

	s.repo = mocks.NewCampaignRepository(s.T())
	s.service = NewCampaignService(s.ctx, s.repo)
}

func (s *campaignServiceSuite) TearDownSuite() {
}

func (s *campaignServiceSuite) SetupTest() {
	loc, err := time.LoadLocation("Asia/Taipei")
	s.NoError(err)
	timeNow = func() time.Time {
		return time.Date(2024, 8, 26, 22, 55, 0, 0, loc)
	}
}

func (s *campaignServiceSuite) TestCreate() {
	campaignID := uint(1)
	now := time.Now().Unix()
	mockCampaign := &repository.Campaign{
		ID:      campaignID,
		Created: now,
	}
	s.repo.On("Create", mockCTX, repository.CreateCampaignInput{}).Return(mockCampaign, nil).Once()
	res, err := s.service.Create(s.ctx, CreateCampaignInput{})
	s.NoError(err)
	s.Equal(campaignID, res.ID)
	s.Equal(now, res.Created)
}

func (s *campaignServiceSuite) TestGetLatest() {
	campaignID := uint(1)
	now := time.Now().Unix()
	mockCampaign := &repository.Campaign{
		ID:      campaignID,
		Created: now,
	}

	s.repo.On("GetLatest", mockCTX, repository.GetLatestCampaignInput{}).Return(mockCampaign, nil).Once()

	res, err := s.service.GetLatest(s.ctx, GetLatestCampaignInput{})
	s.NoError(err)
	s.Equal(campaignID, res.ID)
	s.Equal(now, res.Created)
}

func (s *campaignServiceSuite) TestCreateCouponReservationWithNonemptyCouponCode() {
	campaignID := uint(1)
	userID := "user_id_4"
	mockCouponCode := "mock_coupon_code"
	newUUIDString = func() string {
		return mockCouponCode
	}
	couponReservation := &repository.CouponReservation{
		CampaignID: campaignID,
		UserID:     userID,
		CouponCode: mockCouponCode,
	}
	s.repo.On("CreateCouponReservation", mockCTX, repository.CreateCouponReservationInput{
		CampaignID: campaignID,
		UserID:     userID,
		CouponCode: mockCouponCode,
	}).Return(couponReservation, nil).Once()

	createCouponReservationInput := CreateCouponReservationInput{
		CampaignID: campaignID,
		UserID:     userID,
	}
	res, err := s.service.CreateCouponReservation(s.ctx, createCouponReservationInput)
	s.NoError(err)
	s.Equal(campaignID, res.CampaignID)
	s.Equal(userID, res.UserID)
	s.Equal(mockCouponCode, res.CouponCode)
}

func (s *campaignServiceSuite) TestCreateCouponReservationWithEmptyCouponCode() {
	campaignID := uint(1)
	userID := "user_id_1"
	mockCouponCode := ""
	newUUIDString = func() string {
		return mockCouponCode
	}
	couponReservation := &repository.CouponReservation{
		CampaignID: campaignID,
		UserID:     userID,
		CouponCode: mockCouponCode,
	}
	s.repo.On("CreateCouponReservation", mockCTX, repository.CreateCouponReservationInput{
		CampaignID: campaignID,
		UserID:     userID,
		CouponCode: mockCouponCode,
	}).Return(couponReservation, nil).Once()

	createCouponReservationInput := CreateCouponReservationInput{
		CampaignID: campaignID,
		UserID:     userID,
	}
	res, err := s.service.CreateCouponReservation(s.ctx, createCouponReservationInput)
	s.NoError(err)
	s.Equal(campaignID, res.CampaignID)
	s.Equal(userID, res.UserID)
	s.Equal(mockCouponCode, res.CouponCode)
}

func (s *campaignServiceSuite) TestCreateCouponReservationWithInvalidResercationTimeError() {
	loc, err := time.LoadLocation("Asia/Taipei")
	s.NoError(err)
	timeNow = func() time.Time {
		return time.Date(2024, 8, 26, 22, 54, 0, 0, loc)
	}
	campaignID := uint(1)
	userID := "user_id_4"

	createCouponReservationInput := CreateCouponReservationInput{
		CampaignID: campaignID,
		UserID:     userID,
	}
	_, err = s.service.CreateCouponReservation(s.ctx, createCouponReservationInput)
	s.Equal(ErrNotReservationTime, err)
}

func (s *campaignServiceSuite) TestGetCouponReservation() {
	campaignID := uint(1)
	userID := "user_id_1"
	mockCouponCode := "coupon_code"
	couponReservation := &repository.CouponReservation{
		CampaignID: campaignID,
		UserID:     userID,
		CouponCode: mockCouponCode,
	}
	s.repo.On("GetCouponReservation", mockCTX, repository.GetCouponReservationInput{
		CampaignID: campaignID,
		UserID:     userID,
	}).Return(couponReservation, nil).Once()

	getCouponReservationInput := GetCouponReservationInput{
		CampaignID: couponReservation.CampaignID,
		UserID:     couponReservation.UserID,
	}
	res, err := s.service.GetCouponReservation(s.ctx, getCouponReservationInput)
	s.NoError(err)
	s.Equal(campaignID, res.CampaignID)
	s.Equal(userID, res.UserID)
	s.Equal(mockCouponCode, res.CouponCode)
}

func TestCampaignServiceSuite(t *testing.T) {
	suite.Run(t, new(campaignServiceSuite))
}
