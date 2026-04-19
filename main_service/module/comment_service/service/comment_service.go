package comment_service

import (
	"context"
	"fmt"
	comment_dto "main_service/module/comment_service/dto"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ─── Interface ───────────────────────────────────────────────────────────────

type CommentService interface {
	Create(ctx context.Context, userID int64, req comment_dto.CreateCommentRequest) (*comment_dto.CommentResponse, error)
	GetByID(ctx context.Context, id int64) (*comment_dto.CommentResponse, error)
	Update(ctx context.Context, id, userID int64, req comment_dto.UpdateCommentRequest) (*comment_dto.CommentResponse, error)
	Delete(ctx context.Context, id, userID int64) error
	List(ctx context.Context, f comment_dto.CommentFilter, page, limit int) ([]*comment_dto.CommentResponse, int64, error)
}

// ─── Implementation ──────────────────────────────────────────────────────────

type commentService struct {
	db *pgxpool.Pool
}

func NewCommentService(db *pgxpool.Pool) CommentService {
	return &commentService{db: db}
}

// ─── Create ──────────────────────────────────────────────────────────────────

func (s *commentService) Create(ctx context.Context, userID int64, req comment_dto.CreateCommentRequest) (*comment_dto.CommentResponse, error) {
	var id int64
	err := s.db.QueryRow(ctx, `
		INSERT INTO comments (parent_id, user_id, vakansiya_id, resume_id, type, text)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, req.ParentID, userID, req.VakansiyaID, req.ResumeID, req.Type, req.Text).Scan(&id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

// ─── GetByID ─────────────────────────────────────────────────────────────────

func (s *commentService) GetByID(ctx context.Context, id int64) (*comment_dto.CommentResponse, error) {
	var c comment_dto.CommentResponse
	err := s.db.QueryRow(ctx, `
		SELECT id, parent_id, user_id, vakansiya_id, resume_id, type, text, created_at, updated_at, deleted_at
		FROM comments
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(
		&c.ID, &c.ParentID, &c.UserID, &c.VakansiyaID, &c.ResumeID,
		&c.Type, &c.Text, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
	)
	return &c, err
}

// ─── Update ──────────────────────────────────────────────────────────────────

func (s *commentService) Update(ctx context.Context, id, userID int64, req comment_dto.UpdateCommentRequest) (*comment_dto.CommentResponse, error) {
	setClauses := []string{"updated_at = NOW()"}
	args := []interface{}{}
	idx := 1

	if req.Type != nil {
		setClauses = append(setClauses, fmt.Sprintf("type = $%d", idx))
		args = append(args, *req.Type)
		idx++
	}
	if req.Text != nil {
		setClauses = append(setClauses, fmt.Sprintf("text = $%d", idx))
		args = append(args, *req.Text)
		idx++
	}

	args = append(args, id, userID)
	query := fmt.Sprintf(`
		UPDATE comments SET %s
		WHERE id = $%d AND user_id = $%d AND deleted_at IS NULL
		RETURNING id
	`, strings.Join(setClauses, ", "), idx, idx+1)

	var retID int64
	if err := s.db.QueryRow(ctx, query, args...).Scan(&retID); err != nil {
		return nil, fmt.Errorf("comment not found or access denied")
	}
	return s.GetByID(ctx, retID)
}

// ─── Delete ──────────────────────────────────────────────────────────────────

func (s *commentService) Delete(ctx context.Context, id, userID int64) error {
	tag, err := s.db.Exec(ctx, `
		UPDATE comments SET deleted_at = $1
		WHERE id = $2 AND user_id = $3 AND deleted_at IS NULL
	`, time.Now(), id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("comment not found or access denied")
	}
	return nil
}

// ─── List (threaded) ─────────────────────────────────────────────────────────

func (s *commentService) List(ctx context.Context, f comment_dto.CommentFilter, page, limit int) ([]*comment_dto.CommentResponse, int64, error) {
	args := []interface{}{}
	conditions := []string{"c.deleted_at IS NULL", "c.parent_id IS NULL"}
	idx := 1

	if f.VakansiyaID != nil {
		conditions = append(conditions, fmt.Sprintf("c.vakansiya_id = $%d", idx))
		args = append(args, *f.VakansiyaID)
		idx++
	}
	if f.ResumeID != nil {
		conditions = append(conditions, fmt.Sprintf("c.resume_id = $%d", idx))
		args = append(args, *f.ResumeID)
		idx++
	}
	if f.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("c.user_id = $%d", idx))
		args = append(args, *f.UserID)
		idx++
	}
	if f.Type != "" {
		conditions = append(conditions, fmt.Sprintf("c.type = $%d", idx))
		args = append(args, f.Type)
		idx++
	}

	where := strings.Join(conditions, " AND ")

	var total int64
	if err := s.db.QueryRow(ctx, "SELECT COUNT(*) FROM comments c WHERE "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rootArgs := append(args, limit, (page-1)*limit)
	rootRows, err := s.db.Query(ctx, fmt.Sprintf(
		"SELECT c.id FROM comments c WHERE %s ORDER BY c.id ASC LIMIT $%d OFFSET $%d",
		where, idx, idx+1,
	), rootArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rootRows.Close()

	var rootIDs []int64
	for rootRows.Next() {
		var id int64
		if err := rootRows.Scan(&id); err != nil {
			return nil, 0, err
		}
		rootIDs = append(rootIDs, id)
	}
	rootRows.Close()

	if len(rootIDs) == 0 {
		return []*comment_dto.CommentResponse{}, total, nil
	}

	allRows, err := s.db.Query(ctx, `
		WITH RECURSIVE tree AS (
			SELECT id, parent_id, user_id, vakansiya_id, resume_id, type, text, created_at, updated_at, deleted_at
			FROM comments
			WHERE id = ANY($1) AND deleted_at IS NULL
			UNION ALL
			SELECT c.id, c.parent_id, c.user_id, c.vakansiya_id, c.resume_id, c.type, c.text, c.created_at, c.updated_at, c.deleted_at
			FROM comments c
			JOIN tree t ON c.parent_id = t.id
			WHERE c.deleted_at IS NULL
		)
		SELECT id, parent_id, user_id, vakansiya_id, resume_id, type, text, created_at, updated_at, deleted_at
		FROM tree
		ORDER BY parent_id NULLS FIRST, id
	`, rootIDs)
	if err != nil {
		return nil, 0, err
	}
	defer allRows.Close()

	nodeMap := map[int64]*comment_dto.CommentResponse{}
	for allRows.Next() {
		var c comment_dto.CommentResponse
		if err := allRows.Scan(
			&c.ID, &c.ParentID, &c.UserID, &c.VakansiyaID, &c.ResumeID,
			&c.Type, &c.Text, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
		); err != nil {
			return nil, 0, err
		}
		cc := c
		nodeMap[c.ID] = &cc
	}
	if err := allRows.Err(); err != nil {
		return nil, 0, err
	}

	rootSet := make(map[int64]bool, len(rootIDs))
	for _, id := range rootIDs {
		rootSet[id] = true
	}

	for _, node := range nodeMap {
		if !rootSet[node.ID] && node.ParentID != nil {
			if parent, ok := nodeMap[*node.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	roots := make([]*comment_dto.CommentResponse, 0, len(rootIDs))
	for _, id := range rootIDs {
		if node, ok := nodeMap[id]; ok {
			roots = append(roots, node)
		}
	}
	return roots, total, nil
}
