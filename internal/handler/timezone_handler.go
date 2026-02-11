package handler

import (
	"fmt"
	"net/http"
	"time"

	"dove/pkg/response"
	"dove/pkg/timezone"

	"github.com/gin-gonic/gin"
)

// TimezoneHandler 时区处理器
type TimezoneHandler struct{}

// NewTimezoneHandler 创建时区处理器
func NewTimezoneHandler() *TimezoneHandler {
	return &TimezoneHandler{}
}

// GetCurrentTimezone 获取当前时区
// @Summary 获取当前时区
// @Description 获取系统当前设置的时区
// @Tags 时区
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Router /timezone/current [get]
func (h *TimezoneHandler) GetCurrentTimezone(c *gin.Context) {
	currentTimezone := timezone.GetCurrentTimezone()
	currentTime := timezone.GetCurrentTime()

	response.Success(c, gin.H{
		"timezone":         currentTimezone,
		"current_time":     currentTime.Format("2006-01-02 15:04:05"),
		"current_time_utc": currentTime.UTC().Format("2006-01-02 15:04:05"),
		"timestamp":        currentTime.Unix(),
	})
}

// GetAvailableTimezones 获取可用时区列表
// @Summary 获取可用时区列表
// @Description 获取系统支持的所有时区列表
// @Tags 时区
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]string}
// @Router /timezone/available [get]
func (h *TimezoneHandler) GetAvailableTimezones(c *gin.Context) {
	timezones := timezone.GetAvailableTimezones()
	response.Success(c, timezones)
}

// GetTimeInTimezone 获取指定时区的时间
// @Summary 获取指定时区的时间
// @Description 获取指定时区的当前时间
// @Tags 时区
// @Accept json
// @Produce json
// @Param timezone query string true "时区名称" example(Asia/Shanghai)
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 400 {object} response.Response
// @Router /timezone/time [get]
func (h *TimezoneHandler) GetTimeInTimezone(c *gin.Context) {
	tz := c.Query("timezone")
	if tz == "" {
		response.Error(c, http.StatusBadRequest, "timezone parameter is required")
		return
	}

	if !timezone.IsValidTimezone(tz) {
		response.Error(c, http.StatusBadRequest, "invalid timezone")
		return
	}

	currentTime, err := timezone.GetCurrentTimeInTimezone(tz)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get time in timezone")
		return
	}

	response.Success(c, gin.H{
		"timezone":         tz,
		"current_time":     currentTime.Format("2006-01-02 15:04:05"),
		"current_time_utc": currentTime.UTC().Format("2006-01-02 15:04:05"),
		"timestamp":        currentTime.Unix(),
	})
}

// ParseTime 解析时间字符串
// @Summary 解析时间字符串
// @Description 在指定时区解析时间字符串
// @Tags 时区
// @Accept json
// @Produce json
// @Param time query string true "时间字符串" example(2023-01-01 12:00:00)
// @Param timezone query string true "时区名称" example(Asia/Shanghai)
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 400 {object} response.Response
// @Router /timezone/parse [get]
func (h *TimezoneHandler) ParseTime(c *gin.Context) {
	timeStr := c.Query("time")
	tz := c.Query("timezone")

	if timeStr == "" {
		response.Error(c, http.StatusBadRequest, "time parameter is required")
		return
	}

	if tz == "" {
		response.Error(c, http.StatusBadRequest, "timezone parameter is required")
		return
	}

	if !timezone.IsValidTimezone(tz) {
		response.Error(c, http.StatusBadRequest, "invalid timezone")
		return
	}

	parsedTime, err := timezone.ParseTimeInTimezone(timeStr, tz)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to parse time: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"original_time":   timeStr,
		"timezone":        tz,
		"parsed_time":     parsedTime.Format("2006-01-02 15:04:05"),
		"parsed_time_utc": parsedTime.UTC().Format("2006-01-02 15:04:05"),
		"timestamp":       parsedTime.Unix(),
	})
}

// FormatTime 格式化时间
// @Summary 格式化时间
// @Description 将时间戳格式化为指定时区的时间字符串
// @Tags 时区
// @Accept json
// @Produce json
// @Param timestamp query int64 true "时间戳" example(1672531200)
// @Param timezone query string true "时区名称" example(Asia/Shanghai)
// @Param format query string false "时间格式" example(2006-01-02 15:04:05)
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 400 {object} response.Response
// @Router /timezone/format [get]
func (h *TimezoneHandler) FormatTime(c *gin.Context) {
	timestampStr := c.Query("timestamp")
	tz := c.Query("timezone")
	format := c.Query("format")

	if format == "" {
		format = "2006-01-02 15:04:05"
	}

	if timestampStr == "" {
		response.Error(c, http.StatusBadRequest, "timestamp parameter is required")
		return
	}

	if tz == "" {
		response.Error(c, http.StatusBadRequest, "timezone parameter is required")
		return
	}

	if !timezone.IsValidTimezone(tz) {
		response.Error(c, http.StatusBadRequest, "invalid timezone")
		return
	}

	// 解析时间戳
	var timestamp int64
	if _, err := fmt.Sscanf(timestampStr, "%d", &timestamp); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid timestamp")
		return
	}

	timeObj := time.Unix(timestamp, 0)
	formattedTime, err := timezone.FormatTimeInTimezone(timeObj, tz, format)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to format time")
		return
	}

	response.Success(c, gin.H{
		"timestamp":      timestamp,
		"timezone":       tz,
		"format":         format,
		"formatted_time": formattedTime,
		"utc_time":       timeObj.UTC().Format(format),
	})
}

// ConvertTime 转换时区
// @Summary 转换时区
// @Description 将时间从一个时区转换到另一个时区
// @Tags 时区
// @Accept json
// @Produce json
// @Param time query string true "时间字符串" example(2023-01-01 12:00:00)
// @Param from_timezone query string true "源时区" example(Asia/Shanghai)
// @Param to_timezone query string true "目标时区" example(America/New_York)
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 400 {object} response.Response
// @Router /timezone/convert [get]
func (h *TimezoneHandler) ConvertTime(c *gin.Context) {
	timeStr := c.Query("time")
	fromTz := c.Query("from_timezone")
	toTz := c.Query("to_timezone")

	if timeStr == "" {
		response.Error(c, http.StatusBadRequest, "time parameter is required")
		return
	}

	if fromTz == "" {
		response.Error(c, http.StatusBadRequest, "from_timezone parameter is required")
		return
	}

	if toTz == "" {
		response.Error(c, http.StatusBadRequest, "to_timezone parameter is required")
		return
	}

	if !timezone.IsValidTimezone(fromTz) {
		response.Error(c, http.StatusBadRequest, "invalid from_timezone")
		return
	}

	if !timezone.IsValidTimezone(toTz) {
		response.Error(c, http.StatusBadRequest, "invalid to_timezone")
		return
	}

	// 解析源时区的时间
	parsedTime, err := timezone.ParseTimeInTimezone(timeStr, fromTz)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to parse time: "+err.Error())
		return
	}

	// 转换到目标时区
	convertedTime, err := timezone.GetCurrentTimeInTimezone(toTz)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to convert timezone")
		return
	}

	// 计算时间差
	timeDiff := convertedTime.Sub(parsedTime)

	response.Success(c, gin.H{
		"original_time":         timeStr,
		"from_timezone":         fromTz,
		"to_timezone":           toTz,
		"converted_time":        convertedTime.Format("2006-01-02 15:04:05"),
		"converted_time_utc":    convertedTime.UTC().Format("2006-01-02 15:04:05"),
		"timestamp":             convertedTime.Unix(),
		"time_difference_hours": timeDiff.Hours(),
	})
}
