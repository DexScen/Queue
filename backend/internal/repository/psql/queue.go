package psql

import (
	"context"
	"database/sql"

	"github.com/DexScen/Queue/backend/internal/domain"
	"github.com/DexScen/Queue/backend/internal/errors"
	_ "github.com/lib/pq"
)

type Queues struct {
	db *sql.DB
}

func NewQueues(db *sql.DB) *Queues {
	return &Queues{db: db}
}

func (q *Queues) GetGameInfoByID(ctx context.Context, id int) (*domain.Game, error) {
    query := `
        SELECT 
            g.id,
            g.name,
            g.description,
            g.max_slots,
            g.duration_seconds,
            COALESCE(COUNT(q.id), 0) AS current_people
        FROM games g
        LEFT JOIN queue q 
            ON g.id = q.game_id AND q.status = 'waiting'
        WHERE g.id = $1
        GROUP BY g.id, g.name, g.description, g.max_slots, g.duration_seconds
    `

    var result domain.Game
    err := q.db.QueryRowContext(ctx, query, id).Scan(
        &result.ID,
        &result.Name,
        &result.Description,
        &result.Max_slots,
        &result.Duration_seconds,
        &result.Current_people,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.ErrGameNotFound
        }
        return nil, err
    }

    return &result, nil
}


func (q *Queues) GetAllGames(ctx context.Context, listGames *domain.ListGames) error {
    rows, err := q.db.QueryContext(ctx, `
        SELECT 
            g.id,
            g.name,
            g.description,
            g.max_slots,
            g.duration_seconds,
            COALESCE(COUNT(q.id), 0) AS current_people
        FROM games g
        LEFT JOIN queue q 
            ON g.id = q.game_id AND q.status = 'waiting'
        GROUP BY g.id, g.name, g.description, g.max_slots, g.duration_seconds
        ORDER BY g.id;
    `)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var game domain.Game
        if err := rows.Scan(
            &game.ID,
            &game.Name,
            &game.Description,
            &game.Max_slots,
            &game.Duration_seconds,
            &game.Current_people,
        ); err != nil {
            return err
        }

        *listGames = append(*listGames, game)
    }

    return rows.Err()
}

func (q *Queues) GetGamesByLogin(ctx context.Context, login string, listGames *domain.ListGames) error {
	rows, err := q.db.QueryContext(ctx, `
		SELECT 
			g.id,
			g.name,
			g.description,
			g.max_slots,
			g.duration_seconds,
			COALESCE(COUNT(q2.id), 0) AS current_people
		FROM users u
		JOIN queue q1 ON u.id = q1.user_id
		JOIN games g ON q1.game_id = g.id
		LEFT JOIN queue q2 ON g.id = q2.game_id AND q2.status = 'waiting'
		WHERE u.login = $1
		GROUP BY g.id, g.name, g.description, g.max_slots, g.duration_seconds
		ORDER BY g.id;
	`, login)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var game domain.Game
		if err := rows.Scan(
			&game.ID,
			&game.Name,
			&game.Description,
			&game.Max_slots,
			&game.Duration_seconds,
			&game.Current_people,
		); err != nil {
			return err
		}
		*listGames = append(*listGames, game)
	}

	return rows.Err()
}

func (q *Queues) GetPassword(ctx context.Context, login string) (string, error) {
	tr, err := q.db.Begin()
	if err != nil {
		return "", err
	}
	statement, err := tr.Prepare("SELECT password_hash FROM users WHERE login=$1")
	if err != nil {
		tr.Rollback()
		return "", err
	}
	defer statement.Close()

	var passwordHash string
	err = statement.QueryRow(login).Scan(&passwordHash)
	if err != nil {
		tr.Rollback()
		if err == sql.ErrNoRows {
			return "", errors.ErrUserNotFound
		}
		return "", err
	}

	if err := tr.Commit(); err != nil {
		return "", err
	}

	return passwordHash, nil
}

func (q *Queues) GetRole(ctx context.Context, login string) (string, error) {
	tr, err := q.db.Begin()
	if err != nil {
		return "", err
	}
	statement, err := tr.Prepare("SELECT role FROM users WHERE login=$1")
	if err != nil {
		tr.Rollback()
		return "", err
	}
	defer statement.Close()

	var role string
	err = statement.QueryRow(login).Scan(&role)
	if err != nil {
		tr.Rollback()
		if err == sql.ErrNoRows {
			return "", errors.ErrUserNotFound
		}
		return "", err
	}

	if err := tr.Commit(); err != nil {
		return "", err
	}

	return role, nil
}

func (q *Queues) UserExists(ctx context.Context, login string) (bool, error) {
	tr, err := q.db.Begin()
	if err != nil {
		return false, err
	}
	statement, err := tr.Prepare("SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)")
	if err != nil {
		tr.Rollback()
		return false, err
	}
	defer statement.Close()

	var exists bool
	err = statement.QueryRow(login).Scan(&exists)
	if err != nil {
		tr.Rollback()
		return false, err
	}

	if err := tr.Commit(); err != nil {
		return false, err
	}

	return exists, nil
}

func (q *Queues) Register(ctx context.Context, user *domain.User) error {
	tr, err := q.db.Begin()
	if err != nil {
		return err
	}
	statement, err := tr.Prepare("INSERT INTO users (login, password_hash) VALUES ($1, $2)")
	if err != nil {
		tr.Rollback()
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(user.Login, user.Password)
	if err != nil {
		tr.Rollback()
		return err
	}

	return tr.Commit()
}

func (q *Queues) RemovePlayerFromQueue(ctx context.Context, user_id, game_id int) error{
	tr, err := q.db.Begin()
	if err != nil {
		return err
	}
	statement, err := tr.Prepare("DELETE FROM queue WHERE user_id = $1 AND game_id = $2")
	if err != nil {
		tr.Rollback()
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(user_id, game_id)
	if err != nil {
		tr.Rollback()
		return err
	}

	return tr.Commit()
}

func (q *Queues) AddPlayerToQueue(ctx context.Context, userID, gameID int) (int, error) {
    var position int

    err := q.db.QueryRowContext(ctx, `
        SELECT COALESCE(MAX(position), 0) + 1
        FROM queue
        WHERE game_id = $1
    `, gameID).Scan(&position)
    if err != nil {
        return 0, err
    }

    // вставляем нового игрока
    _, err = q.db.ExecContext(ctx, `
        INSERT INTO queue (user_id, game_id, position, status)
        VALUES ($1, $2, $3, 'waiting')
    `, userID, gameID, position)
    if err != nil {
        return 0, err
    }

    return position, nil
}

func (q *Queues) GetIdByLogin(ctx context.Context, login string) (int, error) {
    var id int
    err := q.db.QueryRowContext(ctx,
        `SELECT id FROM users WHERE login = $1`,
        login,
    ).Scan(&id)

    if err != nil {
        if err == sql.ErrNoRows {
            return 0, errors.ErrUserNotFound
        }
        return 0, err
    }

    return id, nil
}

func (q *Queues) GetPlayersByGameID(ctx context.Context, gameID int, listUsers *domain.ListUsers) error {
    rows, err := q.db.QueryContext(ctx, `
        SELECT u.id, u.login
        FROM queue q
        JOIN users u ON q.user_id = u.id
        WHERE q.game_id = $1
    `, gameID)
    if err != nil {
        return err
    }
    defer rows.Close()

    var users domain.ListUsers
    for rows.Next() {
        var u domain.User
        if err := rows.Scan(&u.ID, &u.Login); err != nil {
            return err
        }
        users = append(users, u)
    }

    *listUsers = users
    return rows.Err()
}
