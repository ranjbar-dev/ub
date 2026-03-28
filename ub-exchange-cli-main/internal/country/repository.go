package country

import (
	"context"
	"database/sql"
)

type Country struct {
	ID                            int64
	Name                          sql.NullString
	FullName                      sql.NullString
	Code                          sql.NullString
	ImagePath                     sql.NullString
	Iso31661alpha3                sql.NullString
	Iso31661alpha2                sql.NullString
	Iso4217currencyAlphabeticCode sql.NullString
	UnitermEnglishFormal          sql.NullString
	RegionName                    sql.NullString
	Languages                     sql.NullString
}

// Repository provides data access for country records.
type Repository interface {
	// All returns every country in the database.
	All(ctx context.Context) []Country
	// GetCountryByID looks up a single country by its database ID.
	GetCountryByID(id int64, country *Country) error
}
