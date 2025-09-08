package adminapi

import (
	"context"
	"github/wildwind123/shop/internal/provider"
	"github/wildwind123/shop/pkg/db"
	"github/wildwind123/shop/pkg/ogenapi"
	"log/slog"

	"github.com/go-faster/errors"
	"github.com/huandu/go-sqlbuilder"
	"github.com/wildwind123/slogger"
)

// pH *ProductHandler github/wildwind123/shop/pkg/ogenapi.ProductHandler
type ProductHandler struct {
	Provider *provider.Provider
}

// ProductGet implements GET /product operation.
// Get product list.
//
// GET /product
func (pH *ProductHandler) ProductGet(ctx context.Context, params ogenapi.ProductGetParams) (ogenapi.ProductGetRes, error) {
	logger := slogger.FromCtx(ctx)

	res, err := db.QueryTableRows[db.Product](pH.Provider.DB.Sqlx(), db.ReqQueryTableRow{
		Tags:  []string{"db"},
		Table: "product",
		Builder: func(sb *sqlbuilder.SelectBuilder) *sqlbuilder.SelectBuilder {
			return sb.Limit(params.Limit.Value).Offset(params.Offset.Value)
		},
	})

	if err != nil {
		logger.Error("cant product get", slog.Any("err", err))
		return nil, errors.Wrap(err, "cant product get")
	}

	return &ogenapi.ResGetProduct{
		Data:   convertDBProductsToOgenAPI(res.Item),
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
		Count:  int(db.ItemCount(pH.Provider.DB.Sqlx(), res.Sb, logger)),
	}, nil
}

// ProductPost implements POST /product operation.
//
// Receive the exact message you've sent.
//
// POST /product
func (pH *ProductHandler) ProductPost(ctx context.Context, req *ogenapi.ProductPostReq) (ogenapi.ProductPostRes, error) {
	logger := slogger.FromCtx(ctx)
	res, err := pH.Provider.DB.Sqlx().NamedExec(`
	insert into product(name, description) values( :name, :description)
	`, map[string]interface{}{
		"name":        req.Name,
		"description": "",
	})
	if err != nil {
		logger.Error("cant insert new item", slog.Any("err", err))
		return &ogenapi.Error{
			Message: err.Error(),
		}, nil
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		logger.Error("cant get last insert", slog.Any("err", err))
		return &ogenapi.Error{
			Message: err.Error(),
		}, nil
	}

	return &ogenapi.ResponseId{
		ID: int(lastInsertID),
	}, nil
}
