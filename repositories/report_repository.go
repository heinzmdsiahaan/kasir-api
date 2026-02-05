package repositories

import (
	"database/sql"
	models "kasir-api/model"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetReport(startDate, endDate string) (*models.ReportResponse, error) {
	var result models.ReportResponse

	// 1️⃣ Total revenue & total transaksi
	err := repo.db.QueryRow(`
		SELECT 
			COALESCE(SUM(total_amount), 0) AS total_revenue,
			COUNT(*) AS total_transaksi
		FROM transactions
		WHERE created_at BETWEEN $1 AND $2
	`, startDate, endDate).Scan(
		&result.TotalRevenue,
		&result.TotalTransaksi,
	)
	if err != nil {
		return nil, err
	}

	// 2️⃣ Produk terlaris
	err = repo.db.QueryRow(`
		SELECT 
			p.name,
			SUM(td.quantity) AS qty_terjual
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		JOIN transactions t ON t.id = td.transaction_id
		WHERE t.created_at BETWEEN $1 AND $2
		GROUP BY p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`, startDate, endDate).Scan(
		&result.ProdukTerlaris.Nama,
		&result.ProdukTerlaris.QtyTerjual,
	)

	// kalau belum ada transaksi
	if err == sql.ErrNoRows {
		result.ProdukTerlaris = models.BestSellingProduct{}
		return &result, nil
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}
