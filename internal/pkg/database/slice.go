package database

func GetSlice(db *Database, query string) (any, error) {

	// open connection
	sess, err := db.NewSession(false)
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	rows, err := sess.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]any, 0, 10)

	for rows.Next() {

		var item any
		if err := rows.Scan(&item); err != nil {
			return nil, err
		}

		result = append(result, item)

	}

	return result, nil

}
