package service

import (
	"context"

	"github.com/alifmufthi91/ecommerce-system/services/product/config"
	warehouseservice "github.com/alifmufthi91/ecommerce-system/services/product/external/warehouse_service"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg/observ"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/product/payload"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/product/repository"
	"go.opentelemetry.io/otel/codes"
)

//go:generate mockery --name=ProductService --case underscore
type ProductService interface {
	CreateProduct(ctx context.Context, req payload.CreateProductReq) error
	GetProducts(ctx context.Context, token string) ([]payload.GetProductsResp, error)
	GetProductByID(ctx context.Context, productID string) (model.Product, error)
}

type productService struct {
	config       *config.Config
	warehouseSvc warehouseservice.IWarehouseSvc
	productRepo  repository.ProductRepository
}

func NewProductService(config *config.Config, whSvc warehouseservice.IWarehouseSvc, productRepo repository.ProductRepository) ProductService {
	return &productService{
		config:       config,
		warehouseSvc: whSvc,
		productRepo:  productRepo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req payload.CreateProductReq) (err error) {
	ctx, span := observ.GetTracer().Start(ctx, "productService.CreateProduct")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ShopID:      req.ShopID,
	}

	err = s.productRepo.CreateProduct(ctx, product)

	return err
}

func (s *productService) GetProducts(ctx context.Context, token string) (result []payload.GetProductsResp, err error) {
	ctx, span := observ.GetTracer().Start(ctx, "productService.GetProducts")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	products, err := s.productRepo.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	var productIDs []string
	var productIndexMap = make(map[string]int)
	for i, product := range products {
		productIDs = append(productIDs, product.ID.String())
		productIndexMap[product.ID.String()] = i
	}

	availableStocks, err := s.warehouseSvc.GetStockAvailables(ctx, warehouseservice.GetStockAvailablesReq{
		ProductIDIN: productIDs,
		Token:       token,
	})
	if err != nil {
		return nil, err
	}

	for _, stock := range availableStocks.Data {
		if index, exists := productIndexMap[stock.ProductID]; exists {
			result = append(result, payload.GetProductsResp{
				ID:             products[index].ID,
				Name:           products[index].Name,
				Description:    products[index].Description,
				Price:          products[index].Price,
				ShopID:         products[index].ShopID,
				AvailableStock: stock.AvailableStock,
				CreatedAt:      products[index].CreatedAt,
				UpdatedAt:      products[index].UpdatedAt,
			})
		}
	}

	return result, nil
}

func (s *productService) GetProductByID(ctx context.Context, productID string) (model.Product, error) {
	ctx, span := observ.GetTracer().Start(ctx, "productService.GetProductByID")
	defer span.End()

	product, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return model.Product{}, err
	}

	return product, nil
}
