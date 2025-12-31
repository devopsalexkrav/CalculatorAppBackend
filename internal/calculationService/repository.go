package calculationService

import "gorm.io/gorm"

type CalculationRepository interface {
	CreateCalculation(cal Calculation) error
	//TODO: why []Calculation without * to return specific object no copy
	GetAllCalculations() ([]Calculation, error)
	GetCalculationById(id string) (Calculation, error)
	UpdateCalculation(cal Calculation) error
	DeleteCalculation(id string) error
}

type calcRepository struct {
	db *gorm.DB
}

func NewCalculationRepository(db *gorm.DB) CalculationRepository {
	return &calcRepository{db: db}
}	

func (r *calcRepository) CreateCalculation(cal Calculation) error {
	return r.db.Create(&cal).Error
}

func (r *calcRepository) GetAllCalculations() ([]Calculation, error) {
	var calculations []Calculation
	err := r.db.Find(&calculations).Error
	return calculations, err
}

func (r *calcRepository) GetCalculationById(id string) (Calculation, error) {
	var calculation Calculation
	err := r.db.First(&calculation, "id = ?", id).Error
	return calculation, err 
}	

func (r *calcRepository) UpdateCalculation(cal Calculation) error {
	return r.db.Save(&cal).Error
}

func (r *calcRepository) DeleteCalculation(id string) error {
	return r.db.Delete(&Calculation{}, "id = ?", id).Error
}