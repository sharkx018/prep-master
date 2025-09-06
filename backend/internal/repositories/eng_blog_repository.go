package repositories

import (
	"database/sql"
	"fmt"
	"strconv"

	"interview-prep-app/internal/models"
)

// EngBlogRepository handles database operations for engineering blogs
type EngBlogRepository struct {
	db *sql.DB
}

// NewEngBlogRepository creates a new engineering blog repository
func NewEngBlogRepository(db *sql.DB) *EngBlogRepository {
	return &EngBlogRepository{db: db}
}

// GetAll retrieves all engineering blogs with their articles
func (r *EngBlogRepository) GetAll(limit, offset int) ([]models.EngBlog, int, error) {
	// First get the total count
	var total int
	countQuery := `SELECT COUNT(*) FROM eng_blogs`
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Build the main query
	query := `
		SELECT 
			eb.id, eb.name, eb.link, eb.order_idx,
			eba.id, eba.title, eba.order_idx, eba.external_link
		FROM eng_blogs eb
		LEFT JOIN eng_blog_articles eba ON eb.id = eba.blog_id
		ORDER BY eb.order_idx ASC, eba.order_idx ASC`

	// Add pagination if specified
	args := []interface{}{}
	if limit > 0 {
		query += ` LIMIT $1`
		args = append(args, limit)
		if offset > 0 {
			query += ` OFFSET $2`
			args = append(args, offset)
		}
	} else if offset > 0 {
		query += ` OFFSET $1`
		args = append(args, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query engineering blogs: %w", err)
	}
	defer rows.Close()

	// Map to store blogs by ID
	blogMap := make(map[int]*models.EngBlog)
	var blogOrder []int

	for rows.Next() {
		var (
			blogID       int
			blogName     string
			blogLink     string
			blogOrderIdx int
			articleID    sql.NullInt64
			articleTitle sql.NullString
			articleOrder sql.NullInt64
			articleLink  sql.NullString
		)

		err := rows.Scan(
			&blogID, &blogName, &blogLink, &blogOrderIdx,
			&articleID, &articleTitle, &articleOrder, &articleLink,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan row: %w", err)
		}

		// Get or create blog
		blog, exists := blogMap[blogID]
		if !exists {
			blog = &models.EngBlog{
				ID:               strconv.Itoa(blogID),
				Name:             blogName,
				Link:             blogLink,
				OrderIdx:         blogOrderIdx,
				PracticeProblems: []models.EngBlogProblem{},
			}
			blogMap[blogID] = blog
			blogOrder = append(blogOrder, blogID)
		}

		// Add article if it exists
		if articleID.Valid {
			article := models.EngBlogProblem{
				ID:           strconv.FormatInt(articleID.Int64, 10),
				Title:        articleTitle.String,
				OrderIdx:     int(articleOrder.Int64),
				ExternalLink: articleLink.String,
			}
			blog.PracticeProblems = append(blog.PracticeProblems, article)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate rows: %w", err)
	}

	// Convert map to slice maintaining order
	blogs := make([]models.EngBlog, 0, len(blogMap))
	for _, blogID := range blogOrder {
		blogs = append(blogs, *blogMap[blogID])
	}

	return blogs, total, nil
}

// GetByID retrieves a specific engineering blog by ID
func (r *EngBlogRepository) GetByID(id string) (*models.EngBlog, error) {
	blogID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid blog ID: %w", err)
	}

	query := `
		SELECT 
			eb.id, eb.name, eb.link, eb.order_idx,
			eba.id, eba.title, eba.order_idx, eba.external_link
		FROM eng_blogs eb
		LEFT JOIN eng_blog_articles eba ON eb.id = eba.blog_id
		WHERE eb.id = $1
		ORDER BY eba.order_idx ASC`

	rows, err := r.db.Query(query, blogID)
	if err != nil {
		return nil, fmt.Errorf("failed to query engineering blog: %w", err)
	}
	defer rows.Close()

	var blog *models.EngBlog
	for rows.Next() {
		var (
			blogName     string
			blogLink     string
			blogOrderIdx int
			articleID    sql.NullInt64
			articleTitle sql.NullString
			articleOrder sql.NullInt64
			articleLink  sql.NullString
		)

		err := rows.Scan(
			&blogID, &blogName, &blogLink, &blogOrderIdx,
			&articleID, &articleTitle, &articleOrder, &articleLink,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Initialize blog on first row
		if blog == nil {
			blog = &models.EngBlog{
				ID:               strconv.Itoa(blogID),
				Name:             blogName,
				Link:             blogLink,
				OrderIdx:         blogOrderIdx,
				PracticeProblems: []models.EngBlogProblem{},
			}
		}

		// Add article if it exists
		if articleID.Valid {
			article := models.EngBlogProblem{
				ID:           strconv.FormatInt(articleID.Int64, 10),
				Title:        articleTitle.String,
				OrderIdx:     int(articleOrder.Int64),
				ExternalLink: articleLink.String,
			}
			blog.PracticeProblems = append(blog.PracticeProblems, article)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	if blog == nil {
		return nil, fmt.Errorf("engineering blog not found")
	}

	return blog, nil
}

// CreateBlog creates a new engineering blog
func (r *EngBlogRepository) CreateBlog(name, link string, orderIdx int) (*models.EngBlogDB, error) {
	query := `
		INSERT INTO eng_blogs (name, link, order_idx) 
		VALUES ($1, $2, $3) 
		RETURNING id, name, link, order_idx, created_at, updated_at`

	var blog models.EngBlogDB
	err := r.db.QueryRow(query, name, link, orderIdx).Scan(
		&blog.ID, &blog.Name, &blog.Link, &blog.OrderIdx,
		&blog.CreatedAt, &blog.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create engineering blog: %w", err)
	}

	return &blog, nil
}

// CreateArticle creates a new article for an engineering blog
func (r *EngBlogRepository) CreateArticle(blogID int, title, externalLink string, orderIdx int) (*models.EngBlogArticleDB, error) {
	query := `
		INSERT INTO eng_blog_articles (blog_id, title, external_link, order_idx) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, blog_id, title, external_link, order_idx, created_at, updated_at`

	var article models.EngBlogArticleDB
	err := r.db.QueryRow(query, blogID, title, externalLink, orderIdx).Scan(
		&article.ID, &article.BlogID, &article.Title, &article.ExternalLink, &article.OrderIdx,
		&article.CreatedAt, &article.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create engineering blog article: %w", err)
	}

	return &article, nil
}
