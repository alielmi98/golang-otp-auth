package handler

import (
	"net/http"
	"strconv"

	"github.com/alielmi98/golang-otp-auth/di"
	"github.com/alielmi98/golang-otp-auth/internal/user/api/dto"
	"github.com/alielmi98/golang-otp-auth/internal/user/usecase"
	"github.com/alielmi98/golang-otp-auth/pkg/config"
	"github.com/alielmi98/golang-otp-auth/pkg/helper"
	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	usecase    *usecase.UserUsecase
	otpUsecase *usecase.OtpUsecase
}

func NewUserHandler(cfg *config.Config) *UsersHandler {
	otpProvider := di.GetOtpProvider(cfg)
	userUsecase := usecase.NewUserUsecase(cfg, di.GetUserRepository(cfg), di.GetTokenProvider(cfg), otpProvider)
	otpUsecase := usecase.NewOtpUsecase(cfg, otpProvider)
	return &UsersHandler{usecase: userUsecase,
		otpUsecase: otpUsecase}
}

// RegisterLoginByMobileNumber godoc
// @Summary RegisterLoginByMobileNumber
// @Description RegisterLoginByMobileNumber
// @Tags Users
// @Accept  json
// @Produce  json
// @Param Request body dto.RegisterLoginByMobileRequest true "RegisterLoginByMobileRequest"
// @Success 201 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Failure 409 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/users/login-by-mobile [post]
func (h *UsersHandler) RegisterLoginByMobileNumber(c *gin.Context) {
	req := new(dto.RegisterLoginByMobileRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}
	token, err := h.usecase.RegisterAndLoginByMobileNumber(c, req.MobileNumber, req.Otp)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}

	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(token, true, helper.Success))
}

// SendOtp godoc
// @Summary Send otp to user
// @Description Send otp to user
// @Tags Users
// @Accept  json
// @Produce  json
// @Param mobile_number path string true "Mobile number"
// @Success 201 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Failure 409 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/users/send-otp/{mobile_number} [post]
func (h *UsersHandler) SendOtp(c *gin.Context) {
	mobileNumber := c.Param("mobile_number")
	if mobileNumber == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, nil))
		return
	}
	err := h.otpUsecase.SendOtp(mobileNumber)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	// TODO: Call internal SMS service
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(nil, true, helper.Success))
}

// GetUserByMobileNumber godoc
// @Summary Get user by mobile number
// @Description Get user by mobile number
// @Tags Users
// @Accept  json
// @Produce  json
// @Param mobile_number path string true "Mobile number"
// @Success 200 {object} dto.UserInfo "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Failure 409 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/users/{mobile_number} [get]
func (h *UsersHandler) GetUserByMobileNumber(c *gin.Context) {
	mobileNumber := c.Param("mobile_number")
	user, err := h.usecase.GetUserByMobileNumber(c, mobileNumber)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUsers godoc
// @Summary Get users
// @Description Get users
// @Tags Users
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param mobile_number query string false "Mobile number filter"
// @Success 200 {object} dto.UserList "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Failure 409 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/users [get]
func (h *UsersHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	mobileNumber := c.Query("mobile_number")

	users, err := h.usecase.GetAllUsers(c, page, pageSize, mobileNumber)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, users)
}
