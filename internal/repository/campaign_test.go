package repository

import (
	"testing"

	"github.com/asymptoter/tonx-take-home-test/pkg/ctx"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type campaignRepositorySuite struct {
	suite.Suite
	ctx  ctx.CTX
	db   *gorm.DB
	repo CampaignRepository
}

func (s *campaignRepositorySuite) SetupSuite() {
	s.ctx = ctx.Background()
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.ctx.Fatal(err)
	}

	s.repo = NewCampaignRepository(s.ctx, s.db)
}

func (s *campaignRepositorySuite) TearDownSuite() {
	db, err := s.db.DB()
	if err != nil {
		s.ctx.Fatal(err)
	}
	if err := db.Close(); err != nil {
		s.ctx.Fatal(err)
	}
}

func (s *campaignRepositorySuite) SetupTest() {
	// Clear the campaigns table before each test
	if err := s.db.Exec("DELETE FROM campaigns").Error; err != nil {
		s.ctx.Fatal(err)
	}

	// Clear the reservations table before each test
	if err := s.db.Exec("DELETE FROM coupon_reservations").Error; err != nil {
		s.ctx.Fatal(err)
	}
}

func (s *campaignRepositorySuite) TestCreate() {
	campaign, err := s.repo.Create(s.ctx, CreateCampaignInput{})
	s.NoError(err)
	s.Equal(uint(1), campaign.ID)
}

func (s *campaignRepositorySuite) TestGetLatest() {
	_, err := s.repo.Create(s.ctx, CreateCampaignInput{})
	s.NoError(err)

	campaign2, err := s.repo.Create(s.ctx, CreateCampaignInput{})
	s.NoError(err)

	res, err := s.repo.GetLatest(s.ctx, GetLatestCampaignInput{})
	s.NoError(err)
	s.Equal(campaign2.ID, res.ID)
}

func (s *campaignRepositorySuite) TestCreateCouponReservation() {
	createCouponReservationInput := CreateCouponReservationInput{
		CampaignID: 1,
		UserID:     "user_id_1",
		CouponCode: "coupon_code_1",
	}
	_, err := s.repo.CreateCouponReservation(s.ctx, createCouponReservationInput)
	s.NoError(err)
}

func (s *campaignRepositorySuite) TestGetCouponReservation() {
	createCouponReservationInput := CreateCouponReservationInput{
		CampaignID: 1,
		UserID:     "user_id_1",
		CouponCode: "coupon_code_1",
	}
	couponReservation, err := s.repo.CreateCouponReservation(s.ctx, createCouponReservationInput)
	s.NoError(err)

	getCouponReservationInput := GetCouponReservationInput{
		CampaignID: couponReservation.CampaignID,
		UserID:     couponReservation.UserID,
	}
	res, err := s.repo.GetCouponReservation(s.ctx, getCouponReservationInput)
	s.NoError(err)
	s.Equal(createCouponReservationInput.CouponCode, res.CouponCode)
}

func TestCampaignRepositorySuite(t *testing.T) {
	suite.Run(t, new(campaignRepositorySuite))
}
