package response

import "github.com/labstack/echo/v4"

type successResponse interface{}

type Success struct {
	Response successResponse `json:"response"`
	Code     int             `json:"code"`
}

type responseData struct {
	Data interface{} `json:"data"`
}

func SuccessBuilder(res *Success, data interface{}) *Success {
	res.Response = nil
	if data != nil {
		res.Response = data
	}

	return res
}

func SuccessResponse(data interface{}) *Success {
	return SuccessBuilder(&Success{Code: 200}, data)
}

func (s *Success) Send(c echo.Context) error {
	return c.JSON(s.Code, responseData{Data: s.Response})
}
