package handler

import (
	"net/http"
	"strconv"

	"github.com/asymptoter/tonx-take-home-test/internal/service"
	"github.com/asymptoter/tonx-take-home-test/pkg/ctx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	getUserID = func() (string, error) {
		return uuid.NewString(), nil
	}
)

type handler struct {
	campaignService service.CampaignService
}

func RegisterHTTPHandler(r *gin.Engine, campaignService service.CampaignService) {
	h := handler{
		campaignService: campaignService,
	}

	// Get latest campaign id
	r.GET("/campaigns/latest", h.GetLatestCampaign)
	// Create reservation
	r.POST("/campaigns/:id/reservations", h.CreateCouponReservation)
	// Get coupon code
	r.GET("/campaigns/:id/reservations", h.GetCouponReservation)
}

type getLatestCampaignResponse struct {
	LatestCampaignID uint `json:"latest_campaign_id"`
}

func (h handler) GetLatestCampaign(c *gin.Context) {
	ctx := ctx.Background()
	userID, err := getUserID()
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	ctx = ctx.With("user_id", userID)

	campaign, err := h.campaignService.GetLatest(ctx, service.GetLatestCampaignInput{})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, getLatestCampaignResponse{
		LatestCampaignID: campaign.ID,
	})
}

func (h handler) CreateCouponReservation(c *gin.Context) {
	ctx := ctx.Background()
	userID, err := getUserID()
	if err != nil {
		ctx.Error(err)
		c.Status(http.StatusUnauthorized)
		return
	}
	ctx = ctx.With("user_id", userID)

	campaignID, err := strconv.Atoi(c.Param("id"))
	if err != nil || campaignID < 0 {
		ctx.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid campaign id",
		})
		return
	}

	input := service.CreateCouponReservationInput{
		CampaignID: uint(campaignID),
		UserID:     userID,
	}
	_, err = h.campaignService.CreateCouponReservation(ctx, input)
	if err == service.ErrNotReservationTime {
		c.Status(http.StatusForbidden)
		return
	} else if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

type getCouponReservationResponse struct {
	CouponCode string `json:"coupon_code"`
}

func (h handler) GetCouponReservation(c *gin.Context) {
	ctx := ctx.Background()

	userID, err := getUserID()
	if err != nil {
		ctx.Error(err)
		c.Status(http.StatusUnauthorized)
		return
	}
	ctx = ctx.With("user_id", userID)

	campaignID, err := strconv.Atoi(c.Param("id"))
	if err != nil || campaignID < 0 {
		ctx.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid campaign id",
		})
		return
	}

	input := service.GetCouponReservationInput{
		CampaignID: uint(campaignID),
		UserID:     userID,
	}
	reservation, err := h.campaignService.GetCouponReservation(ctx, input)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, getCouponReservationResponse{
		CouponCode: reservation.CouponCode,
	})
}
