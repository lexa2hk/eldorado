package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/romankravchuk/eldorado/internal/pkg/sl"
	"github.com/romankravchuk/eldorado/internal/server/http/api"
	"github.com/romankravchuk/eldorado/internal/server/http/api/response"
	"github.com/romankravchuk/eldorado/internal/services/auth/proto"
)

const (
	headerPrefix = "Bearer "
	tokenCookie  = "access_token"
)

func JWT(log *slog.Logger, client proto.AuthServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := log.With(
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			var token string

			rawToken := r.Header.Get(api.AuthorizationHeader)

			if strings.HasPrefix(rawToken, headerPrefix) {
				token = strings.TrimPrefix(rawToken, headerPrefix)
			} else {
				cookie, err := r.Cookie(tokenCookie)
				if err != nil {
					msg := "unauthorized"

					log.Error(msg, sl.Err(err))

					response.JSON(w, http.StatusUnauthorized, response.M{"error": msg})

					return
				}
				token = cookie.Value
			}

			if token == "" {
				msg := "unauthorized"

				log.Error(msg, slog.String("error", "token is empty"))

				response.JSON(w, http.StatusUnauthorized, response.M{"error": msg})

				return
			}

			resp, err := client.Verify(r.Context(), &proto.VerifyRequest{Token: token})
			if err != nil {
				msg := "internal server error"

				log.Error(msg, sl.Err(err), slog.String("token", token))

				response.JSON(w, http.StatusInternalServerError, response.M{"error": msg})

				return
			}
			if resp.Meta.Error != "" {
				msg := "forbidden"

				log.Error(msg, slog.String("error", resp.Meta.Error))

				response.JSON(w, int(resp.Meta.Status), response.M{"error": msg})

				return
			}

			log.Info("request verified", slog.String("token", token), slog.String("user_id", resp.UserID))

			ctx := context.WithValue(r.Context(), api.UserIDKey, resp.UserID)
			ctx = context.WithValue(ctx, api.TokenKey, token)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
