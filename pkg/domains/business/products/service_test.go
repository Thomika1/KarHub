package products

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/Thomika1/KarHub/pkg/core/crud"
	productEnums "github.com/Thomika1/KarHub/pkg/domains/shared/enums/product"
	"github.com/Thomika1/KarHub/pkg/domains/shared/errors"
	"github.com/Thomika1/KarHub/pkg/domains/shared/models/business"
)

type mockRepository struct {
	createFn   func(ctx context.Context, entity any) error
	getFn      func(ctx context.Context, entity any, id any, preload bool) error
	listFn     func(ctx context.Context, entities any, filter crud.Query) error
	updateFn   func(ctx context.Context, entity any, id any) error
	updateByFn func(ctx context.Context, entity any, column string, value any) error
	deleteFn   func(ctx context.Context, entity any, id any) error
}

func (m *mockRepository) Create(ctx context.Context, entity any) error {
	if m.createFn != nil {
		return m.createFn(ctx, entity)
	}
	return nil
}

func (m *mockRepository) Get(ctx context.Context, entity any, id any, preload bool) error {
	if m.getFn != nil {
		return m.getFn(ctx, entity, id, preload)
	}
	return nil
}

func (m *mockRepository) GetBy(ctx context.Context, entity any, column string, value any, preload bool) error {
	return nil
}

func (m *mockRepository) List(ctx context.Context, entities any, filter crud.Query) error {
	if m.listFn != nil {
		return m.listFn(ctx, entities, filter)
	}
	return nil
}

func (m *mockRepository) Update(ctx context.Context, entity any, id any) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, entity, id)
	}
	return nil
}

func (m *mockRepository) UpdateBy(ctx context.Context, entity any, column string, value any) error {
	if m.updateByFn != nil {
		return m.updateByFn(ctx, entity, column, value)
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, entity any, id any) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, entity, id)
	}
	return nil
}

func (m *mockRepository) DeleteBy(ctx context.Context, entity any, column string, value any) error {
	return nil
}

func validProduct() *business.Product {
	name := "Filtro de Óleo"
	stock := 10
	minStock := 5
	avgSales := 3
	leadTime := 7
	unitCost := 1850
	critLevel := 3

	return &business.Product{
		ProductData: business.ProductData{
			Name:              &name,
			Category:          productEnums.Engine,
			CurrentStock:      &stock,
			MinimumStock:      &minStock,
			AverageDailySales: &avgSales,
			LeadTimeDays:      &leadTime,
			UnitCost:          &unitCost,
			CriticalityLevel:  &critLevel,
		},
	}
}

func TestService_Create_Success(t *testing.T) {
	repo := &mockRepository{
		listFn: func(ctx context.Context, entities any, filter crud.Query) error {
			*entities.(*[]business.Product) = []business.Product{}
			return nil
		},
	}
	svc := &Service{repository: repo}

	err := svc.Create(context.Background(), validProduct())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestService_Create_Validation_NameEmpty(t *testing.T) {
	repo := &mockRepository{}
	svc := &Service{repository: repo}

	product := validProduct()
	product.Name = nil

	err := svc.Create(context.Background(), product)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var ve errors.ValidationError
	if !stderrors.As(err, &ve) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestService_Create_Validation_CategoryEmpty(t *testing.T) {
	repo := &mockRepository{}
	svc := &Service{repository: repo}

	product := validProduct()
	product.Category = ""

	err := svc.Create(context.Background(), product)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var ve errors.ValidationError
	if !stderrors.As(err, &ve) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestService_Create_Validation_NegativeCurrentStock(t *testing.T) {
	repo := &mockRepository{}
	svc := &Service{repository: repo}

	product := validProduct()
	neg := -5
	product.CurrentStock = &neg

	err := svc.Create(context.Background(), product)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var ve errors.ValidationError
	if !stderrors.As(err, &ve) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestService_Create_Validation_NegativeMinimumStock(t *testing.T) {
	repo := &mockRepository{}
	svc := &Service{repository: repo}

	product := validProduct()
	neg := -1
	product.MinimumStock = &neg

	err := svc.Create(context.Background(), product)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var ve errors.ValidationError
	if !stderrors.As(err, &ve) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestService_Create_Validation_NegativeAverageDailySales(t *testing.T) {
	repo := &mockRepository{}
	svc := &Service{repository: repo}

	product := validProduct()
	neg := -2
	product.AverageDailySales = &neg

	err := svc.Create(context.Background(), product)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var ve errors.ValidationError
	if !stderrors.As(err, &ve) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestService_Create_Validation_NegativeLeadTimeDays(t *testing.T) {
	repo := &mockRepository{}
	svc := &Service{repository: repo}

	product := validProduct()
	neg := -3
	product.LeadTimeDays = &neg

	err := svc.Create(context.Background(), product)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var ve errors.ValidationError
	if !stderrors.As(err, &ve) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestService_Create_Validation_NegativeUnitCost(t *testing.T) {
	repo := &mockRepository{}
	svc := &Service{repository: repo}

	product := validProduct()
	neg := -100
	product.UnitCost = &neg

	err := svc.Create(context.Background(), product)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var ve errors.ValidationError
	if !stderrors.As(err, &ve) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestService_Create_Validation_CriticalityLevel_Zero(t *testing.T) {
	repo := &mockRepository{}
	svc := &Service{repository: repo}

	product := validProduct()
	zero := 0
	product.CriticalityLevel = &zero

	err := svc.Create(context.Background(), product)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var ve errors.ValidationError
	if !stderrors.As(err, &ve) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestService_Create_Validation_CriticalityLevel_TooHigh(t *testing.T) {
	repo := &mockRepository{}
	svc := &Service{repository: repo}

	product := validProduct()
	high := 6
	product.CriticalityLevel = &high

	err := svc.Create(context.Background(), product)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var ve errors.ValidationError
	if !stderrors.As(err, &ve) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestService_Create_DuplicateProduct(t *testing.T) {
	existing := validProduct()
	repo := &mockRepository{
		listFn: func(ctx context.Context, entities any, filter crud.Query) error {
			*entities.(*[]business.Product) = []business.Product{*existing}
			return nil
		},
	}
	svc := &Service{repository: repo}

	err := svc.Create(context.Background(), validProduct())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var ce errors.ConflictError
	if !stderrors.As(err, &ce) {
		t.Errorf("expected ConflictError, got %T", err)
	}
}

func TestService_Get_Success(t *testing.T) {
	expected := *validProduct()
	repo := &mockRepository{
		getFn: func(ctx context.Context, entity any, id any, preload bool) error {
			p := entity.(*business.Product)
			*p = expected
			return nil
		},
	}
	svc := &Service{repository: repo}

	result, err := svc.Get(context.Background(), "test-id")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Name == nil || *result.Name != *expected.Name {
		t.Errorf("expected name %s, got %v", *expected.Name, result.Name)
	}
}

func TestService_List_Success(t *testing.T) {
	expected := []business.Product{*validProduct()}
	repo := &mockRepository{
		listFn: func(ctx context.Context, entities any, filter crud.Query) error {
			*entities.(*[]business.Product) = expected
			return nil
		},
	}
	svc := &Service{repository: repo}

	results, err := svc.List(context.Background(), crud.Query{})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 product, got %d", len(results))
	}
}

func TestService_Update_Success(t *testing.T) {
	called := false
	repo := &mockRepository{
		updateByFn: func(ctx context.Context, entity any, column string, value any) error {
			called = true
			return nil
		},
	}
	svc := &Service{repository: repo}

	product := validProduct()
	err := svc.Update(context.Background(), "test-id", &product.ProductData)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !called {
		t.Error("expected repository.UpdateBy to be called")
	}
}

func TestService_Update_EmptyData(t *testing.T) {
	repo := &mockRepository{}
	svc := &Service{repository: repo}

	err := svc.Update(context.Background(), "test-id", &business.ProductData{})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestService_Delete_Success(t *testing.T) {
	called := false
	repo := &mockRepository{
		deleteFn: func(ctx context.Context, entity any, id any) error {
			called = true
			return nil
		},
	}
	svc := &Service{repository: repo}

	err := svc.Delete(context.Background(), "test-id")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !called {
		t.Error("expected repository.Delete to be called")
	}
}
