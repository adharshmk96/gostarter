package pgstorage

import (
	"context"
	"database/sql"
	"gostarter/infra"
	"gostarter/internals/domain"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type accountRepository struct {
	conn   *sql.DB
	logger *slog.Logger
	tracer trace.Tracer
}

func NewAccountRepository(container *infra.Container) domain.AccountRepository {
	return &accountRepository{
		conn:   container.DbConn,
		logger: container.Logger,
		tracer: container.Tracer,
	}
}

const (
	createAccountQuery = `
		INSERT INTO gostarter_account (username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	getRoleIDByNameQuery = `
		SELECT id FROM gostarter_role WHERE name = $1`

	createRoleQuery = `
		INSERT INTO gostarter_role (name, created_at, updated_at)
		VALUES ($1, $2, $3)
		RETURNING id`

	assignRoleToAccountQuery = `
		INSERT INTO gostarter_account_role (account_id, role_id, created_at)
		VALUES ($1, $2, $3)`

	getRolesByAccountIDQuery = `
		SELECT r.name
		FROM gostarter_role r
		JOIN gostarter_account_role ar ON r.id = ar.role_id
		WHERE ar.account_id = $1`

	getAccountByIDQuery = `
		SELECT a.id, a.username, a.email, a.password, a.created_at, a.updated_at 
		FROM gostarter_account a
		WHERE a.id = $1
		GROUP BY a.id`

	getAccountByEmailQuery = `
		SELECT a.id, a.username, a.email, a.password, a.created_at, a.updated_at
		FROM gostarter_account a
		WHERE a.email = $1
		GROUP BY a.id`

	getAccountByUsernameQuery = `
		SELECT a.id, a.username, a.email, a.password, a.created_at, a.updated_at
		FROM gostarter_account a
		WHERE a.username = $1
		GROUP BY a.id`

	updateAccountQuery = `
		UPDATE gostarter_account
		SET username = $1, email = $2, password = $3, updated_at = $4
		WHERE id = $5`

	deleteAccountQuery = `
		DELETE FROM gostarter_account WHERE id = $1`

	listAccountsQuery = `
		SELECT a.id, a.username, a.email, a.password, a.created_at, a.updated_at,
			   array_agg(r.name) as roles
		FROM gostarter_account a
		LEFT JOIN gostarter_account_role ar ON a.id = ar.account_id
		LEFT JOIN gostarter_role r ON ar.role_id = r.id
		GROUP BY a.id
		ORDER BY a.id
		LIMIT $1 OFFSET $2`

	totalAccountsQuery = `
		SELECT COUNT(id) FROM gostarter_account
		`
)

func (a *accountRepository) CreateAccount(ctx context.Context, account *domain.Account) error {
	ctx, span := a.tracer.Start(ctx, "AccountRepository.CreateAccount")
	defer span.End()

	// Start transaction
	tx, err := a.conn.BeginTx(ctx, nil)
	if err != nil {
		a.logger.Error("failed to begin transaction", "error", err)
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.logger.Error("failed to rollback transaction", "error", rbErr)
			}
		}
	}()

	now := time.Now()
	if account.Username == "" {
		account.Username = account.Email
	}

	// Insert account
	err = tx.QueryRowContext(
		ctx,
		createAccountQuery,
		account.Username,
		account.Email,
		account.Password,
		now,
		now,
	).Scan(&account.Id)

	if err != nil {
		a.logger.Error("failed to create account", "error", err)
		return err
	}

	// Assign roles
	for _, roleName := range account.Roles {
		// Try to get existing role Id
		var roleID int
		err = tx.QueryRowContext(ctx, getRoleIDByNameQuery, roleName).Scan(&roleID)

		if err == sql.ErrNoRows {
			// Role doesn't exist, create it
			err = tx.QueryRowContext(
				ctx,
				createRoleQuery,
				roleName,
				now,
				now,
			).Scan(&roleID)

			if err != nil {
				a.logger.Error("failed to create role", "error", err, "role", roleName)
				return err
			}
		} else if err != nil {
			a.logger.Error("failed to get role id", "error", err, "role", roleName)
			return err
		}

		// Assign role to account
		_, err = tx.ExecContext(
			ctx,
			assignRoleToAccountQuery,
			account.Id,
			roleID,
			now,
		)

		if err != nil {
			a.logger.Error("failed to assign role to account",
				"error", err,
				"accountId", account.Id,
				"roleId", roleID)
			return err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		a.logger.Error("failed to commit transaction", "error", err)
		return err
	}

	account.CreatedAt = now
	account.UpdatedAt = now
	return nil
}

func (a *accountRepository) GetAccountByID(ctx context.Context, id int) (*domain.Account, error) {
	ctx, span := a.tracer.Start(ctx, "AccountRepository.GetAccountByID")
	defer span.End()

	account := &domain.Account{}
	var roles []string

	err := a.conn.QueryRowContext(ctx, getAccountByIDQuery, id).Scan(
		&account.Id,
		&account.Username,
		&account.Email,
		&account.Password,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrAccountNotFound
	}

	if err != nil {
		a.logger.Error("failed to get account by id", "error", err)
		return nil, err
	}

	err = a.conn.QueryRowContext(ctx, getRolesByAccountIDQuery, id).Scan(&roles)
	if err != nil {
		a.logger.Error("failed to get account roles", "error", err)
		return nil, err
	}

	account.Roles = roles
	return account, nil
}

func (a *accountRepository) GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	ctx, span := a.tracer.Start(ctx, "AccountRepository.GetAccountByEmail")
	defer span.End()

	var account domain.Account
	var roles []string

	err := a.conn.QueryRowContext(ctx, getAccountByEmailQuery, email).Scan(
		&account.Id,
		&account.Username,
		&account.Email,
		&account.Password,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrAccountNotFound
	}

	if err != nil {
		a.logger.Error("failed to get account by email", "error", err)
		return nil, err
	}

	rows, err := a.conn.QueryContext(ctx, getRolesByAccountIDQuery, account.Id)
	if err != nil {
		a.logger.Error("failed to get account roles", "error", err)
		return nil, err
	}

	for rows.Next() {
		var role string
		err := rows.Scan(&role)
		if err != nil {
			a.logger.Error("failed to scan role row", "error", err)
			return nil, err
		}
		roles = append(roles, role)
	}

	account.Roles = roles
	return &account, nil
}

func (a *accountRepository) UpdateAccount(ctx context.Context, account *domain.Account) error {
	ctx, span := a.tracer.Start(ctx, "AccountRepository.UpdateAccount")
	defer span.End()

	res, err := a.conn.ExecContext(ctx, updateAccountQuery,
		account.Username,
		account.Email,
		account.Password,
		time.Now(),
		account.Id,
	)

	if err != nil {
		a.logger.Error("failed to update account", "error", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrAccountNotFound
	}

	return nil
}

func (a *accountRepository) DeleteAccount(ctx context.Context, id int) error {
	ctx, span := a.tracer.Start(ctx, "AccountRepository.DeleteAccount")
	defer span.End()

	res, err := a.conn.ExecContext(ctx, deleteAccountQuery, id)
	if err != nil {
		a.logger.Error("failed to delete account", "error", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrAccountNotFound
	}

	return nil
}

func (a *accountRepository) ListAccounts(ctx context.Context, pagination *domain.Pagination) ([]*domain.Account, error) {
	ctx, span := a.tracer.Start(ctx, "AccountRepository.ListAccounts")
	defer span.End()

	var total int

	err := a.conn.QueryRowContext(ctx, totalAccountsQuery).Scan(&total)
	if err != nil {
		a.logger.Error("failed to get total accounts", "error", err)
		return nil, err
	}

	pagination.SetTotal(total)

	rows, err := a.conn.QueryContext(ctx, listAccountsQuery,
		pagination.Size,
		pagination.GetOffset(),
	)
	if err != nil {
		a.logger.Error("failed to list accounts", "error", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			a.logger.Error("failed to close rows", slog.String("error", err.Error()))
		}
	}(rows)

	var accounts []*domain.Account

	for rows.Next() {
		account := &domain.Account{}
		var roles []string

		err := rows.Scan(
			&account.Id,
			&account.Username,
			&account.Email,
			&account.Password,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			a.logger.Error("failed to scan account row", "error", err)
			return nil, err
		}

		err = a.conn.QueryRowContext(ctx, getRolesByAccountIDQuery, account.Id).Scan(&roles)
		if err != nil {
			a.logger.Error("failed to get account roles", "error", err)
			return nil, err
		}

		account.Roles = roles
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}
