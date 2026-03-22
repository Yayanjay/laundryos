package apiresponse

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Response struct {
	Timestamp   string      `json:"timestamp"`
	TraceID     string      `json:"trace_id"`
	ResponseKey string      `json:"response_key"`
	Message     Message     `json:"message"`
	Data        interface{} `json:"data"`
}

type Message struct {
	TitleIDN string `json:"title_idn"`
	TitleENG string `json:"title_eng"`
	DescIDN  string `json:"desc_idn"`
	DescENG  string `json:"desc_eng"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalPage   int `json:"total_page"`
	TotalItem   int `json:"total_item"`
}

type PaginatedData struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

func Success(data interface{}) Response {
	return Response{
		Timestamp:   formatTimestamp(),
		TraceID:     generateTraceID(),
		ResponseKey: "SUCCESS",
		Message:     Message{},
		Data:        data,
	}
}

func SuccessWithPagination(items interface{}, pagination Pagination) Response {
	return Response{
		Timestamp:   formatTimestamp(),
		TraceID:     generateTraceID(),
		ResponseKey: "SUCCESS",
		Message:     Message{},
		Data: PaginatedData{
			Items:      items,
			Pagination: pagination,
		},
	}
}

func Error(responseKey string, titleIDN, titleENG, descIDN, descENG string) Response {
	return Response{
		Timestamp:   formatTimestamp(),
		TraceID:     generateTraceID(),
		ResponseKey: responseKey,
		Message: Message{
			TitleIDN: titleIDN,
			TitleENG: titleENG,
			DescIDN:  descIDN,
			DescENG:  descENG,
		},
		Data: nil,
	}
}

func ValidationError(errors []ValidationErrorDetail) Response {
	return Response{
		Timestamp:   formatTimestamp(),
		TraceID:     generateTraceID(),
		ResponseKey: "VALIDATION_ERROR",
		Message: Message{
			TitleIDN: "Validasi Gagal",
			TitleENG: "Validation Failed",
			DescIDN:  "Mohon periksa input anda",
			DescENG:  "Please check your input",
		},
		Data: errors,
	}
}

type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func formatTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func generateTraceID() string {
	return uuid.New().String()
}

func CalculatePagination(page, limit, total int) Pagination {
	totalPage := total / limit
	if total%limit > 0 {
		totalPage++
	}
	return Pagination{
		CurrentPage: page,
		PageSize:    limit,
		TotalPage:   totalPage,
		TotalItem:   total,
	}
}

type ErrorCode string

const (
	ErrCodeSuccess      ErrorCode = "SUCCESS"
	ErrCodeError        ErrorCode = "ERROR"
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
)

func (e ErrorCode) String() string {
	return string(e)
}

func NewError(code ErrorCode, titleIDN, titleENG string) Response {
	return Error(code.String(), titleIDN, titleENG, "", "")
}

func BadRequest(titleIDN, titleENG string) Response {
	return Error(ErrCodeValidation.String(), titleIDN, titleENG, "", "")
}

func InternalError(titleIDN, titleENG string) Response {
	return Error(ErrCodeError.String(), titleIDN, titleENG, "", "")
}

func NotFound(resource string) Response {
	return Error(
		ErrCodeNotFound.String(),
		fmt.Sprintf("%s Tidak Ditemukan", resource),
		fmt.Sprintf("%s Not Found", resource),
		"",
		"",
	)
}

func Unauthorized() Response {
	return Error(
		ErrCodeUnauthorized.String(),
		"Unauthorized",
		"Unauthorized",
		"Anda tidak memiliki akses",
		"You do not have access",
	)
}

func Forbidden() Response {
	return Error(
		ErrCodeForbidden.String(),
		"Akses Ditolak",
		"Access Forbidden",
		"Anda tidak memiliki izin untuk mengakses resource ini",
		"You do not have permission to access this resource",
	)
}
