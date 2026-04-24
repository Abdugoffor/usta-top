package language_service

const allQuery = `
	SELECT
		l.id,
		l.name,
		COALESCE(l.description, ''),
		l.is_active,
		l.created_at,
		l.updated_at
	FROM languages l
	WHERE l.deleted_at IS NULL
	  AND ($1 = '' OR l.name ILIKE '%' || $1 || '%')
	  AND ($2::boolean IS NULL OR l.is_active = $2)
	  AND ($3 = 0 OR l.id < $3)
	ORDER BY l.id DESC
	LIMIT $4
`

const showQuery = `
	SELECT
		l.id,
		l.name,
		COALESCE(l.description, ''),
		l.is_active,
		l.created_at,
		l.updated_at
	FROM languages l
	WHERE l.id = $1
	  AND l.deleted_at IS NULL
`

const createQuery = `
	INSERT INTO languages (name, description, is_active)
	VALUES ($1, $2, $3)
	RETURNING
		id,
		name,
		COALESCE(description, ''),
		is_active,
		created_at,
		updated_at
`

const updateQuery = `
	UPDATE languages SET
		name        = COALESCE($1, name),
		description = COALESCE($2, description),
		is_active   = COALESCE($3, is_active),
		updated_at  = NOW()
	WHERE id = $4
	  AND deleted_at IS NULL
	RETURNING
		id,
		name,
		COALESCE(description, ''),
		is_active,
		created_at,
		updated_at
`

const deleteQuery = `
	UPDATE languages
	SET deleted_at = $1
	WHERE id = $2
	  AND deleted_at IS NULL
`
