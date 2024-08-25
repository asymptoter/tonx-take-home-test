package service

import (
	"github.com/asymptoter/tonx-take-home-test/internal/repository"
	"github.com/asymptoter/tonx-take-home-test/pkg/ctx"
	"github.com/google/uuid"
)

var (
	newUUIDString = uuid.NewString
)

type Campaign struct {
	ID      uint
	Created int64
}

type CouponReservation struct {
	CampaignID uint
	UserID     string
	CouponCode string
}

type CreateCampaignInput struct {
}

type GetLatestCampaignInput struct {
}

type CreateCouponReservationInput struct {
	CampaignID uint
	UserID     string
}

type GetCouponReservationInput struct {
	CampaignID uint
	UserID     string
}

type CampaignService interface {
	Create(c ctx.CTX, p CreateCampaignInput) (*Campaign, error)
	GetLatest(c ctx.CTX, p GetLatestCampaignInput) (*Campaign, error)

	CreateCouponReservation(c ctx.CTX, p CreateCouponReservationInput) (*CouponReservation, error)
	GetCouponReservation(c ctx.CTX, p GetCouponReservationInput) (*CouponReservation, error)
}

type campaignService struct {
	repo repository.CampaignRepository
}

func NewCampaignService(c ctx.CTX, repo repository.CampaignRepository) CampaignService {
	return campaignService{
		repo: repo,
	}
}

func (s campaignService) Create(c ctx.CTX, p CreateCampaignInput) (*Campaign, error) {
	res, err := s.repo.Create(c, repository.CreateCampaignInput{})
	if err != nil {
		c.Error(err)
		return nil, err
	}
	return &Campaign{
		ID:      res.ID,
		Created: res.Created,
	}, nil
}

func (s campaignService) GetLatest(c ctx.CTX, p GetLatestCampaignInput) (*Campaign, error) {
	res, err := s.repo.GetLatest(c, repository.GetLatestCampaignInput{})
	if err != nil {
		c.Error(err)
		return nil, err
	}
	return &Campaign{
		ID:      res.ID,
		Created: res.Created,
	}, nil
}

func (s campaignService) CreateCouponReservation(c ctx.CTX, p CreateCouponReservationInput) (*CouponReservation, error) {
	// 根據 campaign_id 和 user_id 來決定 user 能不能拿到 coupon
	v := p.CampaignID
	for _, c := range p.UserID {
		v += uint(c)
	}

	couponCode := ""
	if v%5 == 0 { // 1/5 的機率可以拿到 coupon
		couponCode = newUUIDString()
	}

	input := repository.CreateCouponReservationInput{
		CampaignID: p.CampaignID,
		UserID:     p.UserID,
		CouponCode: couponCode,
	}
	res, err := s.repo.CreateCouponReservation(c, input)
	if err != nil {
		c.Error(err)
		return nil, err
	}

	return &CouponReservation{
		CampaignID: res.CampaignID,
		UserID:     res.UserID,
		CouponCode: res.CouponCode,
	}, nil
}

func (s campaignService) GetCouponReservation(c ctx.CTX, p GetCouponReservationInput) (*CouponReservation, error) {
	input := repository.GetCouponReservationInput{
		CampaignID: p.CampaignID,
		UserID:     p.UserID,
	}
	res, err := s.repo.GetCouponReservation(c, input)
	if err != nil {
		c.Error(err)
		return nil, err
	}

	return &CouponReservation{
		CampaignID: res.CampaignID,
		UserID:     res.UserID,
		CouponCode: res.CouponCode,
	}, nil
}
