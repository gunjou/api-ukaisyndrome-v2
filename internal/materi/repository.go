package materi

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func (r *Repository) GetMateriByModul(ctx context.Context, userID int, modulID int, materiType *string) ([]MateriDTO, error) {

	query := `
		SELECT m.id_materi, m.id_modul, m.tipe_materi, m.judul, m.url_file, m.is_downloadable
		FROM materi m
		JOIN modul md ON md.id_modul = m.id_modul
		JOIN modulkelas mk ON mk.id_modul = md.id_modul
		JOIN paketkelas pk ON pk.id_paketkelas = mk.id_paketkelas
		JOIN pesertakelas p ON p.id_paketkelas = pk.id_paketkelas
		WHERE 
			m.id_modul = $2
			AND p.id_user = $1
			AND m.status = 1
			AND m.visibility = 'open'
	`

	args := []interface{}{userID, modulID}

	// 🔥 filter optional
	if materiType != nil {
		query += " AND m.tipe_materi = $3"
		args = append(args, *materiType)
	}

	query += " ORDER BY m.created_at ASC"

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []MateriDTO

	for rows.Next() {
		var m MateriDTO
		err := rows.Scan(
			&m.ID,
			&m.IDModul,
			&m.Type,
			&m.Title,
			&m.URL,
			&m.IsDownloadable,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil
}


func (r *Repository) GetMateriPrivateByUser(ctx context.Context, userID int, materiType *string) ([]MateriPrivateDTO, error) {

	query := `
		SELECT 
			mp.id_materi_private,
			mp.tipe_materi,
			mp.judul,
			mp.url_file,
			mp.is_downloadable
		FROM materi_private mp
		JOIN mentorship m ON m.id_mentorship = mp.id_mentorship
		WHERE 
			m.id_peserta = $1
			AND mp.status = 1
			AND mp.visibility = 'open'
	`

	args := []interface{}{userID}

	// optional filter
	if materiType != nil {
		query += " AND mp.tipe_materi = $2"
		args = append(args, *materiType)
	}

	query += " ORDER BY mp.created_at ASC"

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []MateriPrivateDTO

	for rows.Next() {
		var m MateriPrivateDTO
		err := rows.Scan(
			&m.ID,
			&m.Type,
			&m.Title,
			&m.URL,
			&m.IsDownloadable,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil
}


// ============================================================
// MATERI PROGRESS - UPSERT
// ============================================================

func (r *Repository) UpsertMateriProgress(ctx context.Context, userID, materiID int, openedAt time.Time) error {
	query := `
		INSERT INTO materi_progress (id_user, id_materi, first_opened_at, last_opened_at, created_at, updated_at)
		VALUES ($1, $2, $3, $3, $3, $3)
		ON CONFLICT (id_user, id_materi)
		DO UPDATE SET last_opened_at = EXCLUDED.last_opened_at, updated_at = EXCLUDED.updated_at
	`
	_, err := r.DB.Exec(ctx, query, userID, materiID, openedAt)
	return err
}

// ============================================================
// MATERI PROGRESS - RAW DATA PESERTA
// ============================================================

func (r *Repository) GetMateriProgress(ctx context.Context, userID int, modulID *int, page, limit int) ([]MateriProgressDTO, int, error) {
	offset := (page - 1) * limit

	countQuery := `
		SELECT COUNT(*)
		FROM materi_progress mp
		JOIN materi m ON m.id_materi = mp.id_materi
		WHERE mp.id_user = $1
	`
	countArgs := []interface{}{userID}

	if modulID != nil {
		countQuery += ` AND m.id_modul = $2`
		countArgs = append(countArgs, *modulID)
	}

	var total int
	if err := r.DB.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			mp.id_materi_progress, mp.id_materi, mp.first_opened_at, mp.last_opened_at
		FROM materi_progress mp
		JOIN materi m ON m.id_materi = mp.id_materi
		WHERE mp.id_user = $1
	`
	args := []interface{}{userID}
	paramIndex := 2

	if modulID != nil {
		query += ` AND m.id_modul = $` + strconv.Itoa(paramIndex)
		args = append(args, *modulID)
		paramIndex++
	}

	query += ` ORDER BY mp.last_opened_at DESC LIMIT $` + strconv.Itoa(paramIndex) + ` OFFSET $` + strconv.Itoa(paramIndex+1)
	args = append(args, limit, offset)

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	result := make([]MateriProgressDTO, 0)

	for rows.Next() {
		var item MateriProgressDTO
		var firstOpenedAt, lastOpenedAt time.Time

		if err := rows.Scan(&item.IDMateriProgress, &item.IDMateri, &firstOpenedAt, &lastOpenedAt); err != nil {
			return nil, 0, err
		}

		item.FirstOpenedAt = firstOpenedAt.Format(time.RFC3339)
		item.LastOpenedAt = lastOpenedAt.Format(time.RFC3339)
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

// ============================================================
// MATERI PROGRESS - MONITORING PESERTA
// ============================================================

func (r *Repository) GetMateriProgressMonitoring(ctx context.Context, userID int, modulID *int) (MateriProgressMonitoringSummaryDTO, []MateriProgressMonitoringDTO, error) {
	summaryQuery := `
		SELECT 
			COUNT(*) AS total_materi, COUNT(mp.id_materi_progress) AS materi_dibuka
		FROM materi m
		JOIN modul md ON md.id_modul = m.id_modul
		JOIN modulkelas mk ON mk.id_modul = md.id_modul
		JOIN paketkelas pk ON pk.id_paketkelas = mk.id_paketkelas
		JOIN pesertakelas p ON p.id_paketkelas = pk.id_paketkelas
		LEFT JOIN materi_progress mp ON mp.id_materi = m.id_materi AND mp.id_user = p.id_user
		WHERE p.id_user = $1 AND m.status = 1 AND m.visibility = 'open'
	`
	summaryArgs := []interface{}{userID}

	if modulID != nil {
		summaryQuery += ` AND m.id_modul = $2`
		summaryArgs = append(summaryArgs, *modulID)
	}

	var summary MateriProgressMonitoringSummaryDTO
	if err := r.DB.QueryRow(ctx, summaryQuery, summaryArgs...).Scan(&summary.TotalMateri, &summary.MateriDibuka); err != nil {
		return summary, nil, err
	}

	summary.MateriBelumDibuka = summary.TotalMateri - summary.MateriDibuka
	if summary.TotalMateri > 0 {
		summary.ProgressPercentage = float64(summary.MateriDibuka) / float64(summary.TotalMateri) * 100
	}

	moduleQuery := `
		SELECT 
			md.id_modul, md.judul AS nama_modul, COUNT(m.id_materi) AS total_materi, 
			COUNT(mp.id_materi_progress) AS materi_dibuka
		FROM materi m
		JOIN modul md ON md.id_modul = m.id_modul
		JOIN modulkelas mk ON mk.id_modul = md.id_modul
		JOIN paketkelas pk ON pk.id_paketkelas = mk.id_paketkelas
		JOIN pesertakelas p ON p.id_paketkelas = pk.id_paketkelas
		LEFT JOIN materi_progress mp ON mp.id_materi = m.id_materi AND mp.id_user = p.id_user
		WHERE p.id_user = $1 AND m.status = 1 AND m.visibility = 'open'
	`
	moduleArgs := []interface{}{userID}

	if modulID != nil {
		moduleQuery += ` AND m.id_modul = $2`
		moduleArgs = append(moduleArgs, *modulID)
	}

	moduleQuery += `
		GROUP BY md.id_modul, md.judul
		ORDER BY md.id_modul ASC
	`

	rows, err := r.DB.Query(ctx, moduleQuery, moduleArgs...)
	if err != nil {
		return summary, nil, err
	}
	defer rows.Close()

	modules := make([]MateriProgressMonitoringDTO, 0)

	for rows.Next() {
		var item MateriProgressMonitoringDTO

		if err := rows.Scan(&item.IDModul, &item.NamaModul, &item.TotalMateri, &item.MateriDibuka); err != nil {
			return summary, nil, err
		}

		item.MateriBelumDibuka = item.TotalMateri - item.MateriDibuka
		if item.TotalMateri > 0 {
			item.ProgressPercentage = float64(item.MateriDibuka) / float64(item.TotalMateri) * 100
		}

		modules = append(modules, item)
	}

	if err := rows.Err(); err != nil {
		return summary, nil, err
	}

	return summary, modules, nil
}