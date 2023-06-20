package parser

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	apperr "gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/error"
	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/internal/helper"
)

type WithPathID interface {
	SetID(int64)
}
type WithUserID interface {
	SetUserID(int64)
}
type WithPathIDAndUserID interface {
	SetID(int64)
	SetUserID(int64)
}

type BodyRequest interface{}

type Parser interface {
	ParserUserID(c *fiber.Ctx) (int64, error)
	ParserIntIDFromPathParams(c *fiber.Ctx) (int64, error)
	ParserBodyRequest(c *fiber.Ctx, req BodyRequest) error
	ParserBodyRequestWithUserID(c *fiber.Ctx, req WithUserID) error
	ParserBodyWithIntIDPathParams(c *fiber.Ctx, req WithPathID) error
	ParserBodyWithIntIDPathParamsAndUserID(c *fiber.Ctx, req WithPathIDAndUserID) error
}

type RequestParser struct {
	validator *validator.Validate
}

func NewParser() *RequestParser {
	return &RequestParser{
		validator: validator.New(),
	}
}

// Get User ID from Token
func (p *RequestParser) ParserUserID(c *fiber.Ctx) (int64, error) {
	userID := c.Locals("user_id").(int64)

	if userID == 0 {
		return 0, fmt.Errorf("EMPTY USER ID")
	}

	return userID, nil
}

// Get ID int64 from Path param
func (p *RequestParser) ParserIntIDFromPathParams(c *fiber.Ctx) (int64, error) {
	ID := c.Params("id")

	if ID == "" {
		return 0, fmt.Errorf("PATH PARAM ID EMPTY")
	}

	return helper.ToInt64(ID), nil
}

// Get request body and parse to struct
func (p *RequestParser) ParserBodyRequest(c *fiber.Ctx, req BodyRequest) error {
	body := c.Body()
	if err := json.Unmarshal(body, &req); err != nil {
		return apperr.ErrInvalidRequest()
	}

	return nil
}

// Get Request Body and ID int64 from request param
func (p *RequestParser) ParserBodyWithIntIDPathParams(c *fiber.Ctx, req WithPathID) error {
	if err := p.ParserBodyRequest(c, req); err != nil {
		return err
	}

	ID, err := p.ParserIntIDFromPathParams(c)

	if err != nil {
		return err
	}

	req.SetID(ID)

	return nil
}

// Get Request Body and User ID (from Token)
func (p *RequestParser) ParserBodyRequestWithUserID(c *fiber.Ctx, req WithUserID) error {
	if err := p.ParserBodyRequest(c, req); err != nil {
		return err
	}

	userID := c.Locals("user_id").(int64)

	req.SetUserID(userID)

	return nil
}

// Get Request Body, ID int64 from request param, and User ID (from Token)
func (p *RequestParser) ParserBodyWithIntIDPathParamsAndUserID(c *fiber.Ctx, req WithPathIDAndUserID) error {
	if err := p.ParserBodyWithIntIDPathParams(c, req); err != nil {
		return err
	}
	userID := c.Locals("user_id").(int64)

	req.SetUserID(userID)

	return nil
}
