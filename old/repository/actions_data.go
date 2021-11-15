package repository

import "context"

func (repo *Repository) GetBlacklistUnits() (map[string]int8, error) {
	countStmt := selectCountStmt(schema.blacklistUnit.name)
	colStmt := selectColStmt("name", schema.blacklistUnit.name)
	return getItems(countStmt, colStmt, repo)
}

func (repo *Repository) GetFruitsVeggies() (map[string]int8, error) {
	countStmt := selectCountStmt(schema.fruitVeggie.name)
	colStmt := selectColStmt("name", schema.fruitVeggie.name)
	return getItems(countStmt, colStmt, repo)
}

func getItems(countStmt string, colStmt string, repo *Repository) (map[string]int8, error) {
	ctx := context.Background()
	tx, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var count int64
	err = tx.QueryRow(countStmt).Scan(&count)
	if err != nil {
		return nil, err
	}

	items := make(map[string]int8, count)
	rows, err := db.Query(colStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items[name] = 0
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return items, err
}
