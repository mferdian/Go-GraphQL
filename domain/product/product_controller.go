package product

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mferdian/Go-GraphQL/constants"
	"github.com/mferdian/Go-GraphQL/logging"
	"github.com/mferdian/Go-GraphQL/utils"
)

type (
	IProductController interface {
		CreateProduct(ctx *gin.Context)
		UpdateProduct(ctx *gin.Context)
		DeleteProduct(ctx *gin.Context)
	}

	ProductController struct {
		productService IProductService
	}
)

func NewUserController(productService IProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

func (pc *ProductController) CreateProduct(ctx *gin.Context) {
	var payload CreateProductRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := pc.productService.CreateProduct(ctx.Request.Context(), payload)
	if err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_CREATE_PRODUCT)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_CREATE_PRODUCT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_CREATE_PRODUCT+": %s", result.Name)
	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_CREATE_PRODUCT, result)
	ctx.JSON(http.StatusCreated, res)
}
func (pc *ProductController) UpdateProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if _, err := uuid.Parse(idParam); err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_UUID_FORMAT)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var payload UpdateProductRequest
	payload.ID = idParam

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := pc.productService.UpdateProduct(ctx.Request.Context(), payload)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_UPDATE_PRODUCT)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UPDATE_PRODUCT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_UPDATE_PRODUCT+": %s", result.ID)
	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_UPDATE_PRODUCT, result)
	ctx.JSON(http.StatusOK, res)
}
func (pc *ProductController) DeleteProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if _, err := uuid.Parse(idParam); err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_UUID_FORMAT)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	payload := DeleteProductRequest{ProductID: idParam}

	result, err := pc.productService.DeleteProduct(ctx.Request.Context(), payload)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_DELETE_USER)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_DELETE_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_DELETE_USER+": %s", idParam)
	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_DELETE_USER, result)
	ctx.JSON(http.StatusOK, res)
}
