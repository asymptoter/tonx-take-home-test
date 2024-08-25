package repository

import (
	"github.com/asymptoter/tonx-take-home-test/pkg/ctx"
	"gorm.io/gorm"
)

type Campaign struct {
	ID      uint  `gorm:"primaryKey;autoIncrement:true"`
	Created int64 `gorm:"autoCreateTime"`
}

// CouponReservation represents a user's coupon reservation
type CouponReservation struct {
	CampaignID uint   `gorm:"primaryKey"`
	UserID     string `gorm:"primaryKey"`
	CouponCode string
	Campaign   Campaign `gorm:"foreignKey:CampaignID"`
}

type CreateCampaignInput struct {
}

type GetLatestCampaignInput struct {
}

type CreateCouponReservationInput struct {
	CampaignID uint
	UserID     string
	CouponCode string
}

type GetCouponReservationInput struct {
	CampaignID uint
	UserID     string
}

type CampaignRepository interface {
	Create(c ctx.CTX, p CreateCampaignInput) (*Campaign, error)
	GetLatest(c ctx.CTX, p GetLatestCampaignInput) (*Campaign, error)

	CreateCouponReservation(c ctx.CTX, p CreateCouponReservationInput) (*CouponReservation, error)
	GetCouponReservation(c ctx.CTX, p GetCouponReservationInput) (*CouponReservation, error)
}

type campaignRepository struct {
	db *gorm.DB
}

func NewCampaignRepository(c ctx.CTX, db *gorm.DB) CampaignRepository {
	if err := db.AutoMigrate(Campaign{}); err != nil {
		c.Fatal(err)
	}
	if err := db.AutoMigrate(CouponReservation{}); err != nil {
		c.Fatal(err)
	}
	return campaignRepository{
		db: db,
	}
}

func (r campaignRepository) Create(c ctx.CTX, p CreateCampaignInput) (*Campaign, error) {
	var res Campaign
	if err := r.db.Create(&res).Error; err != nil {
		c.Error(err)
		return nil, err
	}
	return &res, nil
}

func (r campaignRepository) GetLatest(c ctx.CTX, p GetLatestCampaignInput) (*Campaign, error) {
	var res Campaign
	if err := r.db.Last(&res).Error; err != nil {
		c.Error(err)
		return nil, err
	}
	return &res, nil
}

func (r campaignRepository) CreateCouponReservation(c ctx.CTX, p CreateCouponReservationInput) (*CouponReservation, error) {
	res := CouponReservation{
		CampaignID: p.CampaignID,
		UserID:     p.UserID,
		CouponCode: p.CouponCode,
	}
	if err := r.db.Create(&res).Error; err != nil {
		c.Error(err)
		return nil, err
	}
	return &res, nil
}

func (r campaignRepository) GetCouponReservation(c ctx.CTX, p GetCouponReservationInput) (*CouponReservation, error) {
	var res CouponReservation
	if err := r.db.First(&res, "campaign_id = ? AND user_id = ?", p.CampaignID, p.UserID).Error; err != nil {
		c.Error(err)
		return nil, err
	}
	return &res, nil
}
