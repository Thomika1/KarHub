package restock

import (
	"context"
	"errors"
	"testing"

	"github.com/Thomika1/KarHub/pkg/core/crud"
	productEnums "github.com/Thomika1/KarHub/pkg/domains/shared/enums/product"
	"github.com/Thomika1/KarHub/pkg/domains/shared/models/business"
)

type mockRestockRepository struct {
	prioritiesFn func(ctx context.Context, parameters crud.Query) ([]business.Product, error)
}

func (m *mockRestockRepository) Priorities(ctx context.Context, parameters crud.Query) ([]business.Product, error) {
	if m.prioritiesFn != nil {
		return m.prioritiesFn(ctx, parameters)
	}
	return nil, nil
}

func TestService_Priorities_Success(t *testing.T) {
	name1 := "Filtro de Óleo"
	name2 := "Pastilha de Freio"
	stock1 := 15
	stock2 := 8
	min1 := 20
	min2 := 10
	avgSales1 := 3
	avgSales2 := 5
	leadTime1 := 7
	leadTime2 := 4
	crit1 := 3
	crit2 := 4
	urgency1 := 45
	urgency2 := 36
	projected1 := -6
	projected2 := -12

	expected := []business.Product{
		{
			ProductData: business.ProductData{
				Name:              &name1,
				Category:          productEnums.Engine,
				CurrentStock:      &stock1,
				MinimumStock:      &min1,
				AverageDailySales: &avgSales1,
				LeadTimeDays:      &leadTime1,
				CriticalityLevel:  &crit1,
			},
			UrgencyScore:   &urgency1,
			ProjectedStock: &projected1,
		},
		{
			ProductData: business.ProductData{
				Name:              &name2,
				Category:          productEnums.Brakes,
				CurrentStock:      &stock2,
				MinimumStock:      &min2,
				AverageDailySales: &avgSales2,
				LeadTimeDays:      &leadTime2,
				CriticalityLevel:  &crit2,
			},
			UrgencyScore:   &urgency2,
			ProjectedStock: &projected2,
		},
	}

	repo := &mockRestockRepository{
		prioritiesFn: func(ctx context.Context, parameters crud.Query) ([]business.Product, error) {
			if parameters.Page != 0 {
				t.Errorf("expected page 0, got %d", parameters.Page)
			}
			if parameters.PageSize != 10 {
				t.Errorf("expected pageSize 10, got %d", parameters.PageSize)
			}
			return expected, nil
		},
	}

	svc := &Service{domainRepository: repo}

	results, err := svc.Priorities(context.Background(), crud.Query{Page: 0, PageSize: 10})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 products, got %d", len(results))
	}
	if *results[0].Name != name1 {
		t.Errorf("expected first product name %s, got %s", name1, *results[0].Name)
	}
	if *results[0].UrgencyScore != urgency1 {
		t.Errorf("expected first product urgency %d, got %d", urgency1, *results[0].UrgencyScore)
	}
	if *results[0].ProjectedStock != projected1 {
		t.Errorf("expected first product projected stock %d, got %d", projected1, *results[0].ProjectedStock)
	}
}

func TestService_Priorities_EmptyResult(t *testing.T) {
	repo := &mockRestockRepository{
		prioritiesFn: func(ctx context.Context, parameters crud.Query) ([]business.Product, error) {
			return []business.Product{}, nil
		},
	}

	svc := &Service{domainRepository: repo}

	results, err := svc.Priorities(context.Background(), crud.Query{})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 products, got %d", len(results))
	}
}

func TestService_Priorities_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database error")

	repo := &mockRestockRepository{
		prioritiesFn: func(ctx context.Context, parameters crud.Query) ([]business.Product, error) {
			return nil, expectedErr
		},
	}

	svc := &Service{domainRepository: repo}

	_, err := svc.Priorities(context.Background(), crud.Query{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}
