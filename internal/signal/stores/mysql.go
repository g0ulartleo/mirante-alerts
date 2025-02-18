//go:build mysql

package stores

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLSignalRepository struct {
	db *sql.DB
}

const (
	signalsDatabase = "mirante_signals"
)

func NewMySQLSignalRepository(cfg config.MySQLConfig) (signal.SignalRepository, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/", cfg.User, cfg.Password, cfg.Host, cfg.Port))
	if err != nil {
		return nil, err
	}
	repo := &MySQLSignalRepository{db: db}
	if err := repo.Init(); err != nil {
		return nil, err
	}
	db.Close()

	dsnWithDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.User, cfg.Password, cfg.Host, cfg.Port, signalsDatabase)
	conn, err := sql.Open("mysql", dsnWithDB)
	if err != nil {
		return nil, err
	}
	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)

	return &MySQLSignalRepository{db: conn}, nil
}

func (r *MySQLSignalRepository) Save(signal signal.Signal) error {
	query := `INSERT INTO signals (alert_id, status, message, created_at) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, signal.AlertID, signal.Status, signal.Message, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (r *MySQLSignalRepository) GetAlertLatestSignals(alertID string, limit int) ([]signal.Signal, error) {
	query := `
		SELECT alert_id, status, message, created_at 
		FROM signals WHERE alert_id = ? ORDER BY created_at DESC LIMIT ?`
	rows, err := r.db.Query(query, alertID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	signals := make([]signal.Signal, 0)
	for rows.Next() {
		var s signal.Signal
		err := rows.Scan(&s.AlertID, &s.Status, &s.Message, &s.Timestamp)
		if err != nil {
			return nil, err
		}
		signals = append(signals, s)
	}
	return signals, nil
}

func (r *MySQLSignalRepository) GetAlertHealth(alertID string) (signal.Status, error) {
	query := `
		SELECT status 
		FROM signals 
		WHERE alert_id = ? 
		ORDER BY created_at DESC 
		LIMIT 1`
	rows, err := r.db.Query(query, alertID)
	if err != nil {
		return signal.StatusUnknown, err
	}
	defer rows.Close()
	if rows.Next() {
		var status signal.Status
		err := rows.Scan(&status)
		if err != nil {
			return signal.StatusUnknown, err
		}
		return status, nil
	}
	return signal.StatusUnknown, nil
}

func (r *MySQLSignalRepository) Init() error {
	query := `CREATE DATABASE IF NOT EXISTS ` + signalsDatabase
	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}
	query = `
		CREATE TABLE IF NOT EXISTS ` + signalsDatabase + `.signals (
			alert_id VARCHAR(255) NOT NULL,
			status VARCHAR(255) NOT NULL,
			message VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_alert_created (alert_id, created_at)
		)`
	_, err = r.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (r *MySQLSignalRepository) CleanOldSignals() error {
	query := `DELETE FROM signals WHERE created_at < NOW() - INTERVAL 14 DAY`
	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (r *MySQLSignalRepository) Close() error {
	return r.db.Close()
}
