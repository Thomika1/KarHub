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

func TestService_Priorities_NegativeCurrentStock(t *testing.T) {
	name := "Filtro de Óleo"
	currentStock := -10
	minimumStock := 20
	avgSales := 3
	leadTime := 7
	criticality := 3

	projected := currentStock - (avgSales * leadTime)
	urgency := (minimumStock - projected) * criticality

	repo := &mockRestockRepository{
		prioritiesFn: func(ctx context.Context, parameters crud.Query) ([]business.Product, error) {
			return []business.Product{
				{
					ProductData: business.ProductData{
						Name:              &name,
						Category:          productEnums.Engine,
						CurrentStock:      &currentStock,
						MinimumStock:      &minimumStock,
						AverageDailySales: &avgSales,
						LeadTimeDays:      &leadTime,
						CriticalityLevel:  &criticality,
					},
					ProjectedStock: &projected,
					UrgencyScore:   &urgency,
				},
			}, nil
		},
	}

	svc := &Service{domainRepository: repo}

	results, err := svc.Priorities(context.Background(), crud.Query{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedProjected := -31
	expectedUrgency := 153
	if *results[0].ProjectedStock != expectedProjected {
		t.Errorf("expected projected stock %d, got %d", expectedProjected, *results[0].ProjectedStock)
	}
	if *results[0].UrgencyScore != expectedUrgency {
		t.Errorf("expected urgency score %d, got %d", expectedUrgency, *results[0].UrgencyScore)
	}
}

func TestService_Priorities_ZeroDailySales(t *testing.T) {
	name := "Parabrisa Raro"
	currentStock := 15
	minimumStock := 20
	avgSales := 0
	leadTime := 7
	criticality := 2

	projected := currentStock - (avgSales * leadTime)
	urgency := (minimumStock - projected) * criticality

	repo := &mockRestockRepository{
		prioritiesFn: func(ctx context.Context, parameters crud.Query) ([]business.Product, error) {
			return []business.Product{
				{
					ProductData: business.ProductData{
						Name:              &name,
						Category:          productEnums.Cooling,
						CurrentStock:      &currentStock,
						MinimumStock:      &minimumStock,
						AverageDailySales: &avgSales,
						LeadTimeDays:      &leadTime,
						CriticalityLevel:  &criticality,
					},
					ProjectedStock: &projected,
					UrgencyScore:   &urgency,
				},
			}, nil
		},
	}

	svc := &Service{domainRepository: repo}

	results, err := svc.Priorities(context.Background(), crud.Query{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedProjected := 15
	expectedUrgency := 10
	if *results[0].ProjectedStock != expectedProjected {
		t.Errorf("expected projected stock %d, got %d", expectedProjected, *results[0].ProjectedStock)
	}
	if *results[0].UrgencyScore != expectedUrgency {
		t.Errorf("expected urgency score %d, got %d", expectedUrgency, *results[0].UrgencyScore)
	}
}

func TestService_Priorities_HighLeadTime(t *testing.T) {
	name := "Importado China"
	currentStock := 15
	minimumStock := 20
	avgSales := 3
	leadTime := 90
	criticality := 3

	projected := currentStock - (avgSales * leadTime)
	urgency := (minimumStock - projected) * criticality

	repo := &mockRestockRepository{
		prioritiesFn: func(ctx context.Context, parameters crud.Query) ([]business.Product, error) {
			return []business.Product{
				{
					ProductData: business.ProductData{
						Name:              &name,
						Category:          productEnums.Engine,
						CurrentStock:      &currentStock,
						MinimumStock:      &minimumStock,
						AverageDailySales: &avgSales,
						LeadTimeDays:      &leadTime,
						CriticalityLevel:  &criticality,
					},
					ProjectedStock: &projected,
					UrgencyScore:   &urgency,
				},
			}, nil
		},
	}

	svc := &Service{domainRepository: repo}

	results, err := svc.Priorities(context.Background(), crud.Query{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedProjected := -255
	expectedUrgency := 825
	if *results[0].ProjectedStock != expectedProjected {
		t.Errorf("expected projected stock %d, got %d", expectedProjected, *results[0].ProjectedStock)
	}
	if *results[0].UrgencyScore != expectedUrgency {
		t.Errorf("expected urgency score %d, got %d", expectedUrgency, *results[0].UrgencyScore)
	}
}

func TestService_Priorities_MaxCriticality_ZeroStock(t *testing.T) {
	name := "Peça Crítica"
	currentStock := 0
	minimumStock := 10
	avgSales := 5
	leadTime := 10
	criticality := 5

	projected := currentStock - (avgSales * leadTime)
	urgency := (minimumStock - projected) * criticality

	repo := &mockRestockRepository{
		prioritiesFn: func(ctx context.Context, parameters crud.Query) ([]business.Product, error) {
			return []business.Product{
				{
					ProductData: business.ProductData{
						Name:              &name,
						Category:          productEnums.Brakes,
						CurrentStock:      &currentStock,
						MinimumStock:      &minimumStock,
						AverageDailySales: &avgSales,
						LeadTimeDays:      &leadTime,
						CriticalityLevel:  &criticality,
					},
					ProjectedStock: &projected,
					UrgencyScore:   &urgency,
				},
			}, nil
		},
	}

	svc := &Service{domainRepository: repo}

	results, err := svc.Priorities(context.Background(), crud.Query{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedProjected := -50
	expectedUrgency := 300
	if *results[0].ProjectedStock != expectedProjected {
		t.Errorf("expected projected stock %d, got %d", expectedProjected, *results[0].ProjectedStock)
	}
	if *results[0].UrgencyScore != expectedUrgency {
		t.Errorf("expected urgency score %d, got %d", expectedUrgency, *results[0].UrgencyScore)
	}
}

func TestService_Priorities_HighSales_LowStock(t *testing.T) {
	name := "Alta Demanda"
	currentStock := 2
	minimumStock := 20
	avgSales := 50
	leadTime := 3
	criticality := 4

	projected := currentStock - (avgSales * leadTime)
	urgency := (minimumStock - projected) * criticality

	repo := &mockRestockRepository{
		prioritiesFn: func(ctx context.Context, parameters crud.Query) ([]business.Product, error) {
			return []business.Product{
				{
					ProductData: business.ProductData{
						Name:              &name,
						Category:          productEnums.Engine,
						CurrentStock:      &currentStock,
						MinimumStock:      &minimumStock,
						AverageDailySales: &avgSales,
						LeadTimeDays:      &leadTime,
						CriticalityLevel:  &criticality,
					},
					ProjectedStock: &projected,
					UrgencyScore:   &urgency,
				},
			}, nil
		},
	}

	svc := &Service{domainRepository: repo}

	results, err := svc.Priorities(context.Background(), crud.Query{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedProjected := -148
	expectedUrgency := 672
	if *results[0].ProjectedStock != expectedProjected {
		t.Errorf("expected projected stock %d, got %d", expectedProjected, *results[0].ProjectedStock)
	}
	if *results[0].UrgencyScore != expectedUrgency {
		t.Errorf("expected urgency score %d, got %d", expectedUrgency, *results[0].UrgencyScore)
	}
}

func TestService_Priorities_Ordering_UrgencyDesc(t *testing.T) {
	name1 := "Alta Urgência"
	name2 := "Média Urgência"
	name3 := "Baixa Urgência"

	current1 := 0
	current2 := 10
	current3 := 100

	min1 := 20
	min2 := 15
	min3 := 10

	avgSales1 := 5
	avgSales2 := 3
	avgSales3 := 1

	leadTime1 := 10
	leadTime2 := 5
	leadTime3 := 1

	crit1 := 5
	crit2 := 3
	crit3 := 1

	projected1 := current1 - (avgSales1 * leadTime1)
	projected2 := current2 - (avgSales2 * leadTime2)
	projected3 := current3 - (avgSales3 * leadTime3)

	urgency1 := (min1 - projected1) * crit1
	urgency2 := (min2 - projected2) * crit2
	urgency3 := (min3 - projected3) * crit3

	repo := &mockRestockRepository{
		prioritiesFn: func(ctx context.Context, parameters crud.Query) ([]business.Product, error) {
			return []business.Product{
				{
					ProductData: business.ProductData{
						Name: &name1, CurrentStock: &current1, MinimumStock: &min1,
						AverageDailySales: &avgSales1, LeadTimeDays: &leadTime1,
						CriticalityLevel: &crit1,
					},
					ProjectedStock: &projected1, UrgencyScore: &urgency1,
				},
				{
					ProductData: business.ProductData{
						Name: &name2, CurrentStock: &current2, MinimumStock: &min2,
						AverageDailySales: &avgSales2, LeadTimeDays: &leadTime2,
						CriticalityLevel: &crit2,
					},
					ProjectedStock: &projected2, UrgencyScore: &urgency2,
				},
				{
					ProductData: business.ProductData{
						Name: &name3, CurrentStock: &current3, MinimumStock: &min3,
						AverageDailySales: &avgSales3, LeadTimeDays: &leadTime3,
						CriticalityLevel: &crit3,
					},
					ProjectedStock: &projected3, UrgencyScore: &urgency3,
				},
			}, nil
		},
	}

	svc := &Service{domainRepository: repo}

	results, err := svc.Priorities(context.Background(), crud.Query{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 products, got %d", len(results))
	}

	if *results[0].UrgencyScore <= *results[1].UrgencyScore {
		t.Errorf("expected results ordered by urgency DESC, got %d then %d", *results[0].UrgencyScore, *results[1].UrgencyScore)
	}
	if *results[1].UrgencyScore <= *results[2].UrgencyScore {
		t.Errorf("expected results ordered by urgency DESC, got %d then %d", *results[1].UrgencyScore, *results[2].UrgencyScore)
	}
}
