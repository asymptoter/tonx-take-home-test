package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asymptoter/tonx-take-home-test/internal/service"
	"github.com/asymptoter/tonx-take-home-test/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	mockCTX = mock.Anything
)

type handlerSuite struct {
	suite.Suite
	router      *gin.Engine
	mockService *mocks.CampaignService
}

func (s *handlerSuite) SetupSuite() {
	s.mockService = mocks.NewCampaignService(s.T())
	gin.SetMode(gin.TestMode)
	s.router = gin.Default()
	RegisterHTTPHandler(s.router, s.mockService)
}

func (s *handlerSuite) request(method, path string, res any) (int, error) {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	if res != nil {
		if err := json.Unmarshal([]byte(w.Body.String()), res); err != nil {
			return 0, err
		}
	}

	return w.Code, nil
}

func (s *handlerSuite) TestGetLatestCampaign_Success() {
	campaignID := uint(1)
	s.mockService.On("GetLatest", mockCTX, service.GetLatestCampaignInput{}).Return(&service.Campaign{ID: campaignID}, nil).Once()

	var res getLatestCampaignResponse
	code, err := s.request(http.MethodGet, "/campaigns/latest", &res)
	s.NoError(err)
	s.Equal(http.StatusOK, code)
	s.Equal(campaignID, res.LatestCampaignID)
}

func (s *handlerSuite) TestGetLatestCampaign_UnexpectedError() {
	s.mockService.On("GetLatest", mockCTX, service.GetLatestCampaignInput{}).Return(nil, errors.New("error")).Once()

	code, err := s.request(http.MethodGet, "/campaigns/latest", nil)
	s.NoError(err)
	s.Equal(http.StatusInternalServerError, code)
}

func (s *handlerSuite) TestCreateCouponReservation_Success() {
	mockUserID := "mock_user_id"
	getUserID = func() (string, error) {
		return mockUserID, nil
	}

	createCouponReservationInput := service.CreateCouponReservationInput{
		CampaignID: 1,
		UserID:     mockUserID,
	}
	s.mockService.On("CreateCouponReservation", mockCTX, createCouponReservationInput).Return(nil, nil).Once()

	code, err := s.request(http.MethodPost, "/campaigns/1/reservations", nil)
	s.NoError(err)
	s.Equal(http.StatusNoContent, code)
}

func (s *handlerSuite) TestCreateCouponReservation_InvalidTime() {
	mockUserID := "mock_user_id"
	getUserID = func() (string, error) {
		return mockUserID, nil
	}

	createCouponReservationInput := service.CreateCouponReservationInput{
		CampaignID: 1,
		UserID:     mockUserID,
	}
	s.mockService.On("CreateCouponReservation", mockCTX, createCouponReservationInput).Return(nil, service.ErrNotReservationTime).Once()

	code, err := s.request(http.MethodPost, "/campaigns/1/reservations", nil)
	s.NoError(err)
	s.Equal(http.StatusForbidden, code)
}

func (s *handlerSuite) TestCreateCouponReservation_InvalidCampaignID() {
	code, err := s.request(http.MethodPost, "/campaigns/-1/reservations", nil)
	s.NoError(err)
	s.Equal(http.StatusBadRequest, code)
}

func (s *handlerSuite) TestGetCouponReservation_Success() {
	mockUserID := "mock_user_id"
	getUserID = func() (string, error) {
		return mockUserID, nil
	}

	couponCode := "coupon_code"
	getCouponReservationInput := service.GetCouponReservationInput{
		CampaignID: 1,
		UserID:     mockUserID,
	}
	s.mockService.On("GetCouponReservation", mockCTX, getCouponReservationInput).Return(&service.CouponReservation{CouponCode: couponCode}, nil).Once()

	var res getCouponReservationResponse
	code, err := s.request(http.MethodGet, "/campaigns/1/reservations", &res)
	s.NoError(err)
	s.Equal(http.StatusOK, code)
	s.Equal(couponCode, res.CouponCode)
}

func (s *handlerSuite) TestGetCouponReservation_UnexpectedError() {
	mockUserID := "mock_user_id"
	getUserID = func() (string, error) {
		return mockUserID, nil
	}

	getCouponReservationInput := service.GetCouponReservationInput{
		CampaignID: 1,
		UserID:     mockUserID,
	}
	s.mockService.On("GetCouponReservation", mockCTX, getCouponReservationInput).Return(nil, errors.New("")).Once()

	code, err := s.request(http.MethodGet, "/campaigns/1/reservations", nil)
	s.NoError(err)
	s.Equal(http.StatusInternalServerError, code)
}

// Test Suite Runner
func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(handlerSuite))
}
