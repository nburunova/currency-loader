package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/nburunova/currency-loader/src/services/log"
	"github.com/pkg/errors"
)

var (
	ErrInvalidPageSize   = errors.New("invalid page size")
	ErrInvalidValuteCode = errors.New("invalid valute code")
	tokenAuth            = jwtauth.New("HS256", []byte("secret"), nil)
	user                 = map[string]interface{}{"user_id": 123}
)

type Storage interface {
	TotalPages(ctx context.Context, pageSize int) (*int, error)
	ValutesOnPage(ctx context.Context, pageNum, pageSize int) (Valutes, error)
	ValuteValueByCode(ctx context.Context, code string) (*float64, error)
}

type CurrencyAPI struct {
	strg     Storage
	pageSize int
	logger   log.Logger
}

func NewCurrencyAPI(strg Storage, pageSize int, logger log.Logger) (*CurrencyAPI, error) {
	_, tokenString, _ := tokenAuth.Encode(user)
	logger.Debugf("DEBUG: a sample jwt is %s", tokenString)
	if pageSize <= 0 {
		return nil, errors.Wrap(ErrInvalidPageSize, "page size must be greater then 0")
	}
	return &CurrencyAPI{
		strg:     strg,
		pageSize: pageSize,
		logger:   logger,
	}, nil
}

func (api *CurrencyAPI) Routers() *chi.Mux {
	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(jwtauth.Authenticator)
	r.Get("/currencies", api.currenciesHandler)
	r.Get("/currency/{code}", api.valuteValueHandler)
	return r
}

func protected(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			render.Render(w, r, ErrorRenderer(err))
		}
		f(w, r)
	}
}

func (api *CurrencyAPI) currenciesHandler(w http.ResponseWriter, r *http.Request) {
	pageNum, err := pageNum(r)
	if err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	totalPages, err := api.strg.TotalPages(ctx, api.pageSize)
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
	}
	valutesOnPage, err := api.strg.ValutesOnPage(ctx, pageNum, api.pageSize)
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
	}
	render.JSON(w, r, NewCurrenciesResponse(pageNum, *totalPages, valutesOnPage))
}

func (api *CurrencyAPI) valuteValueHandler(w http.ResponseWriter, r *http.Request) {
	code := valuteCode(r)
	if code == "" {
		render.Render(w, r, ErrorRenderer(ErrInvalidValuteCode))
	}
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	value, err := api.strg.ValuteValueByCode(ctx, code)
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
	}
	render.JSON(w, r, NewValuteResponse(NewValute(code, *value)))
}

func pageNum(r *http.Request) (int, error) {
	if r.URL.Query().Get("page") == "" {
		return 1, nil
	}
	return strconv.Atoi(r.URL.Query().Get("page"))
}

func valuteCode(r *http.Request) string {
	return chi.URLParam(r, "code")
}
