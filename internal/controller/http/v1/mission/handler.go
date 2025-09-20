package mission

import (
	"backend/config"
	request "backend/internal/controller/http/request/cat"
	"backend/internal/controller/http/response"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h handler) getMissions(c *gin.Context) {
	missions, svcCode, err := h.svc.GetMissions(c.Request.Context())
	if err != nil {
		c.JSON(config.CodeToHttpStatus(svcCode), response.NewErr(svcCode, err))
		return
	}

	c.JSON(http.StatusOK, response.New(config.CodeOK).AddKey("missions", missions))
}

func (h handler) getMission(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("mission_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.NewErr(config.CodeBadRequest, fmt.Errorf("invalid mission_id: %v", err)))
		return
	}

	mission, svcCode, err := h.svc.GetMission(c.Request.Context(), uint(missionID))
	if err != nil {
		c.JSON(config.CodeToHttpStatus(svcCode), response.NewErr(svcCode, err))
		return
	}

	c.JSON(http.StatusOK, response.New(config.CodeOK).AddKey("mission", mission))
}

func (h handler) createMission(c *gin.Context) {
	var body request.Mission
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErr(config.CodeBadRequest, err))
		return
	}

	if valid, err := h.validator.ValidateStruct(body); !valid || err != nil {
		c.JSON(http.StatusBadRequest, response.NewErr(config.CodeBadRequest, err))
		return
	}

	if err := body.ValidateTargetsLen(); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErr(config.CodeBadRequest, err))
		return
	}

	mission, svcCode, err := h.svc.CreateMission(c.Request.Context(), body)
	if err != nil {
		c.JSON(config.CodeToHttpStatus(svcCode), response.NewErr(svcCode, err))
		return
	}

	c.JSON(http.StatusOK, response.New(svcCode).
		AddKey("mission", mission).
		SetMessage("mission created"))
}

func (h handler) assignCat(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("mission_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.NewErr(config.CodeBadRequest, fmt.Errorf("invalid mission_id: %v", err)))
		return
	}

	catID, err := strconv.ParseUint(c.Param("cat_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.NewErr(config.CodeBadRequest, fmt.Errorf("invalid cat_id: %v", err)))
		return
	}

	svcCode, err := h.svc.AssignCat(c.Request.Context(), uint(missionID), uint(catID))
	if err != nil {
		c.JSON(config.CodeToHttpStatus(svcCode), response.NewErr(svcCode, err))
		return
	}

	c.JSON(http.StatusOK, response.New(svcCode).SetMessage("cat assigned to mission"))
}

func (h handler) completeMission(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("mission_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.NewErr(config.CodeBadRequest, fmt.Errorf("invalid mission_id: %v", err)))
		return
	}

	svcCode, err := h.svc.CompleteMission(c.Request.Context(), uint(missionID))
	if err != nil {
		c.JSON(config.CodeToHttpStatus(svcCode), response.NewErr(svcCode, err))
		return
	}

	c.JSON(http.StatusOK, response.New(svcCode).SetMessage("mission completed"))
}

func (h handler) deleteMission(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("mission_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.NewErr(config.CodeBadRequest, fmt.Errorf("invalid mission_id: %v", err)))
		return
	}

	svcCode, err := h.svc.DeleteMission(c.Request.Context(), uint(missionID))
	if err != nil {
		c.JSON(config.CodeToHttpStatus(svcCode), response.NewErr(svcCode, err))
		return
	}

	c.JSON(http.StatusOK, response.New(svcCode).SetMessage("mission deleted"))
}

func (h handler) createTarget(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("mission_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.NewErr(config.CodeBadRequest, fmt.Errorf("invalid mission_id: %v", err)))
		return
	}

	var body request.Target
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErr(config.CodeBadRequest, err))
		return
	}

	if valid, err := h.validator.ValidateStruct(body); !valid || err != nil {
		c.JSON(http.StatusBadRequest, response.NewErr(config.CodeBadRequest, err))
		return
	}

	svcCode, err := h.svc.CreateTarget(c.Request.Context(), body, uint(missionID))
	if err != nil {
		c.JSON(config.CodeToHttpStatus(svcCode), response.NewErr(svcCode, err))
		return
	}

	c.JSON(http.StatusOK, response.New(svcCode).SetMessage("target created"))
}

func (h handler) updateTarget(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("mission_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.NewErr(config.CodeBadRequest, fmt.Errorf("invalid mission_id: %v", err)))
		return
	}

	targetID, err := strconv.ParseUint(c.Param("target_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.NewErr(config.CodeBadRequest, fmt.Errorf("invalid target_id: %v", err)))
		return
	}

	var body request.UpdateTarget
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErr(config.CodeBadRequest, err))
		return
	}

	if valid, err := h.validator.ValidateStruct(body); !valid || err != nil {
		c.JSON(http.StatusBadRequest, response.NewErr(config.CodeBadRequest, err))
		return
	}

	svcCode, err := h.svc.UpdateTarget(c.Request.Context(), body, uint(targetID), uint(missionID))
	if err != nil {
		c.JSON(config.CodeToHttpStatus(svcCode), response.NewErr(svcCode, err))
		return
	}

	c.JSON(http.StatusOK, response.New(svcCode).SetMessage("target updated"))
}

func (h handler) completeTarget(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("mission_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.NewErr(config.CodeBadRequest, fmt.Errorf("invalid mission_id: %v", err)))
		return
	}

	targetID, err := strconv.ParseUint(c.Param("target_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.NewErr(config.CodeBadRequest, fmt.Errorf("invalid target_id: %v", err)))
		return
	}

	svcCode, err := h.svc.CompleteTarget(c.Request.Context(), uint(targetID), uint(missionID))
	if err != nil {
		c.JSON(config.CodeToHttpStatus(svcCode), response.NewErr(svcCode, err))
		return
	}

	c.JSON(http.StatusOK, response.New(svcCode).SetMessage("target completed"))
}

func (h handler) deleteTarget(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("mission_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.NewErr(config.CodeBadRequest, fmt.Errorf("invalid mission_id: %v", err)))
		return
	}

	targetID, err := strconv.ParseUint(c.Param("target_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.NewErr(config.CodeBadRequest, fmt.Errorf("invalid target_id: %v", err)))
		return
	}

	svcCode, err := h.svc.DeleteTarget(c.Request.Context(), uint(targetID), uint(missionID))
	if err != nil {
		c.JSON(config.CodeToHttpStatus(svcCode), response.NewErr(svcCode, err))
		return
	}

	c.JSON(http.StatusOK, response.New(svcCode).SetMessage("target deleted"))
}
