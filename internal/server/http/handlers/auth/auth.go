package auth

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/romankravchuk/eldorado/internal/pkg/sl"
	"github.com/romankravchuk/eldorado/internal/pkg/validator"
	"github.com/romankravchuk/eldorado/internal/server/http/api"
	"github.com/romankravchuk/eldorado/internal/server/http/api/response"
	"github.com/romankravchuk/eldorado/internal/services/auth/proto"
)

func HandleRegister(log *slog.Logger, client proto.AuthServiceClient) api.APIFunc {
	const op = "server.http.handlers.auth.Register"

	type req struct {
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"required,gte=3,lte=20"`
		Password string `json:"password" validate:"required,gte=8,alphanum,lte=20"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		input := new(req)
		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err))

			return response.APIError{
				Status:  http.StatusBadRequest,
				Message: msg,
			}
		}

		if err := validator.ValidateStruct(*input); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err))

			return response.APIError{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			}
		}

		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		resp, err := client.SignUp(ctx, &proto.SignUpRequest{
			Email:    input.Email,
			Username: input.Username,
			Password: input.Password,
		})
		if err != nil {
			msg := "internal server error"

			log.Error(msg, sl.Err(err), slog.Any("request_body", input))

			return response.APIError{
				Status:  http.StatusInternalServerError,
				Message: msg,
			}
		}
		if resp.Error != "" {
			msg := "invalid request"

			log.Error(msg, slog.Any("request_body", input), slog.Any("response", resp))

			return response.APIError{
				Status:  int(resp.Status),
				Message: msg,
			}
		}

		return response.JSON(w, http.StatusOK, response.M{"message": "ok"})
	}
}

func HandleGetTokenPairs(log *slog.Logger, client proto.AuthServiceClient) api.APIFunc {
	const op = "server.http.handlers.auth.GetTokenPairs"

	type req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,gte=8,alphanum,lte=20"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		input := new(req)
		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err))

			return response.APIError{
				Status:  http.StatusBadRequest,
				Message: msg,
			}
		}

		if err := validator.ValidateStruct(*input); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err))

			return response.APIError{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			}
		}

		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		resp, err := client.Token(ctx, &proto.TokenRequest{
			Email:    input.Email,
			Password: input.Password,
		})
		if err != nil {
			msg := "internal server error"

			log.Error(msg, sl.Err(err), slog.Any("request_body", input))

			return response.APIError{
				Status:  http.StatusInternalServerError,
				Message: msg,
			}
		}
		if resp.Meta.Error != "" {
			msg := "invalid request"

			log.Error(msg, slog.Any("request_body", input), slog.Any("response", resp))

			return response.APIError{
				Status:  int(resp.Meta.Status),
				Message: msg,
			}
		}

		return response.JSON(w, http.StatusOK, response.M{
			"access_token":  resp.AccessToken,
			"refresh_token": resp.RefreshToken,
		})
	}
}

func HandleRefreshToken(log *slog.Logger, client proto.AuthServiceClient) api.APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		token, err := r.Cookie("refresh_token")
		if err != nil {
			msg := "refresh token cookie not found"

			log.Error(msg, sl.Err(err))

			return response.APIError{
				Status:  http.StatusForbidden,
				Message: msg,
			}
		}

		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		resp, err := client.Refresh(ctx, &proto.RefreshRequest{
			Refresh: token.Value,
		})
		if err != nil {
			msg := "internal server error"

			log.Error(msg, sl.Err(err), slog.Any("refresh_token", token.Value))

			return response.APIError{
				Status:  http.StatusInternalServerError,
				Message: msg,
			}
		}
		if resp.Meta.Error != "" {
			msg := "invalid request"

			log.Error(msg, slog.Any("refresh_token", token.Value), slog.Any("response", resp))

			return response.APIError{
				Status:  int(resp.Meta.Status),
				Message: msg,
			}
		}

		return response.JSON(w, http.StatusOK, response.M{
			"access_token": resp.AccessToken,
		})
	}
}
