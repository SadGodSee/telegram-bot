package server

import (
	"github.com/SadGodSee/telegram-bot-one/pkg/repository"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
	"net/http"
	"strconv"
)

type AuthorizationServer struct {
	server          *http.Server
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURL     string
}

func NewAuthorizationServer(pocketClient *pocket.Client, tokenRepository repository.TokenRepository, redirectURL string) *AuthorizationServer {
	return &AuthorizationServer{pocketClient: pocketClient, tokenRepository: tokenRepository, redirectURL: redirectURL}
}

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":80", // делаем сервер на 80 порту, чтобы открывать в браузере
		Handler: s,
	}

	return s.server.ListenAndServe()
}

func (s *AuthorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	chatIDParam := r.URL.Query().Get("chat_id")
	if chatIDParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// теперь необходимо достать респонс токен из БД

	chatID, err := strconv.ParseInt(chatIDParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestToken, err := s.tokenRepository.Get(chatID, repository.RequestTokens)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authResp, err := s.pocketClient.Authorize(r.Context(), requestToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.tokenRepository.Save(chatID, authResp.AccessToken, repository.AccessTokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// напишем немного логов
	log.Printf("chat_id: %d\nrequest_token: %s\n access_token: %s\n", chatID, requestToken, authResp.AccessToken)

	w.Header().Add("Location", s.redirectURL)
	// браузер заредиректит, если статус код 301 и есть Location

	w.WriteHeader(http.StatusMovedPermanently)

}
