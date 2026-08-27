package course

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/iGuessImaDev/gocourse_domain/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, course *domain.Course) error
	GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error)
	Get(ctx context.Context, id string) (*domain.Course, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, name *string, startDate *string, endDate *string) error
	Count(ctx context.Context, filters Filters) (int, error)
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(ctx context.Context, course *domain.Course) error {
	if err := repo.db.WithContext(ctx).Create(course).Error; err != nil {
		repo.log.Println(err)
		return err
	}
	repo.log.Println("user created with id ", course.ID)
	return nil
}

func (repo *repo) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {
	var u []domain.Course

	tx := repo.db.WithContext(ctx).Model(&u)
	tx = applyFilters(tx, filters)
	tx = tx.Offset(offset).Limit(limit)
	result := tx.Order("created_at desc").Find(&u)

	if result.Error != nil {
		repo.log.Println(result.Error)
		return nil, result.Error
	}

	return u, nil
}

func (repo *repo) Get(ctx context.Context, id string) (*domain.Course, error) {
	course := domain.Course{ID: id}
	result := repo.db.WithContext(ctx).First(&course)
	if result.Error != nil {
		repo.log.Println(result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound{id}
		}
		return nil, result.Error
	}
	return &course, nil
}

func (repo *repo) Delete(ctx context.Context, id string) error {
	course := domain.Course{ID: id}

	result := repo.db.WithContext(ctx).Delete(&course)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		repo.log.Printf("course %s doesn't exist", id)
		return ErrNotFound{id}
	}

	repo.log.Println("course deleted with id ", course.ID)
	return nil
}

func (repo *repo) Update(ctx context.Context, id string, name *string, startDate *string, endDate *string) error {
	values := make(map[string]interface{})

	if name != nil {
		values["name"] = name
	}

	if startDate != nil {
		values["start_date"] = startDate
	}

	if endDate != nil {
		values["end_date"] = endDate
	}

	result := repo.db.WithContext(ctx).Model(&domain.Course{}).Where("id", id).Updates(values)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		repo.log.Printf("user %s doesn't exist", id)
		return ErrNotFound{id}
	}

	return nil
}

func (repo *repo) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := repo.db.WithContext(ctx).Model(domain.Course{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		repo.log.Println(err)
		return 0, err
	}
	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.Name != "" {
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("lower(name) like ?", filters.Name)
	}

	if filters.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", filters.StartDate)
		if err == nil {
			tx = tx.Where("start_date >= ?", startDate)
		}
	}

	if filters.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", filters.EndDate)
		if err == nil {
			tx = tx.Where("end_date < ?", endDate)
		}
	}
	return tx
}
