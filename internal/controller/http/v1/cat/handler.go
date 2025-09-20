package cat

import (
	"backend/config"
	request "backend/internal/controller/http/request/cat"
	"backend/internal/controller/http/response"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h handler) getCats(c *gin.Context) {
	breed := c.Query("breed")

	cats, svcCode, err := h.svc.GetCats(c.Request.Context(), breed)
	if err != nil {
		c.JSON(config.CodeToHttpStatus(svcCode), response.NewErr(svcCode, err))
		return
	}

	c.JSON(http.StatusOK, response.New(config.CodeOK).AddKey("cats", cats))
}

func (h handler) getCatByID(c *gin.Context) {
	catID, err := strconv.ParseUint(c.Param("cat_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.NewErr(config.CodeBadRequest, fmt.Errorf("invalid cat_id: %v", err)))
		return
	}

	cat, svcCode, err := h.svc.GetCatByID(c.Request.Context(), uint(catID))
	if err != nil {
		c.JSON(config.CodeToHttpStatus(svcCode), response.NewErr(svcCode, err))
		return
	}

	c.JSON(http.StatusOK, response.New(config.CodeOK).AddKey("cat", cat))
}

func (h handler) createCat(c *gin.Context) {
	var body request.Cat
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErr(config.CodeBadRequest, err))
		return
	}

	if valid, err := h.validator.ValidateStruct(body); !valid || err != nil {
		c.JSON(http.StatusBadRequest, response.NewErr(config.CodeBadRequest, err))
		return
	}

	cat, svcCode, err := h.svc.CreateCat(c.Request.Context(), body)
	if err != nil {
		c.JSON(config.CodeToHttpStatus(svcCode), response.NewErr(svcCode, err))
		return
	}

	c.JSON(http.StatusOK, response.New(svcCode).
		AddKey("cat", cat).
		SetMessage("record created"))
}

func (h handler) updateCat(c *gin.Context) {
	catID, err := strconv.ParseUint(c.Param("cat_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.NewErr(config.CodeBadRequest, fmt.Errorf("invalid cat_id: %v", err)))
		return
	}

	var body request.UpdateCat
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErr(config.CodeBadRequest, err))
		return
	}

	if valid, err := h.validator.ValidateStruct(body); !valid || err != nil {
		c.JSON(http.StatusBadRequest, response.NewErr(config.CodeBadRequest, err))
		return
	}

	svcCode, err := h.svc.UpdateCat(c.Request.Context(), body, uint(catID))
	if err != nil {
		c.JSON(config.CodeToHttpStatus(svcCode), response.NewErr(svcCode, err))
		return
	}

	c.JSON(http.StatusOK, response.New(svcCode).SetMessage("record updated"))
}

func (h handler) deleteCat(c *gin.Context) {
	catID, err := strconv.ParseUint(c.Param("cat_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.NewErr(config.CodeBadRequest, fmt.Errorf("invalid cat_id: %v", err)))
		return
	}

	svcCode, err := h.svc.DeleteCat(c.Request.Context(), uint(catID))
	if err != nil {
		c.JSON(config.CodeToHttpStatus(svcCode), response.NewErr(svcCode, err))
		return
	}

	c.JSON(http.StatusOK, response.New(svcCode).SetMessage("record deleted"))
}
