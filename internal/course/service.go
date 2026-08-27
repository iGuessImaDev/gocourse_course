package course

import (
	"context"
	"log"
	"time"

	"github.com/iGuessImaDev/gocourse_domain/domain"
)

type (
	Service interface {
		Create(ctx context.Context, name, startDate, endDate string) (*domain.Course, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error)
		Get(ctx context.Context, id string) (*domain.Course, error)
		Update(ctx context.Context, id string, name *string, startDate *string, endDate *string) error
		Delete(ctx context.Context, id string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}

	Filters struct {
		Name      string
		StartDate string
		EndDate   string
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, name, startDate, endDate string) (*domain.Course, error) {

	startDateParsed, err := time.Parse("2006-01-01", startDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	endDateParsed, err := time.Parse("2006-01-01", endDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	course := domain.Course{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}
	if err := s.repo.Create(ctx, &course); err != nil {
		return nil, err
	}
	return &course, nil
}

func (s service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {
	s.log.Println("getall course service")
	courses, err := s.repo.GetAll(ctx, filters, offset, limit)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}
	return courses, nil
}

func (s service) Get(ctx context.Context, id string) (*domain.Course, error) {
	courses, err := s.repo.Get(ctx, id)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}
	return courses, nil
}

func (s service) Update(ctx context.Context, id string, name *string, startDate *string, endDate *string) error {
	return s.repo.Update(ctx, id, name, startDate, endDate)
}

func (s service) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s service) Count(ctx context.Context, filters Filters) (int, error) {
	return s.repo.Count(ctx, filters)
}
