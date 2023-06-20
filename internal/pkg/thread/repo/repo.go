package repo

import (
	"context"
	"fmt"
	"github.com/IvanStukalov/DB_project/internal/models"
	"github.com/IvanStukalov/DB_project/internal/pkg/thread"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
	"strings"
	"time"
)

type repoPostgres struct {
	Conn *pgxpool.Pool
}

func NewRepoPostgres(Conn *pgxpool.Pool) thread.Repository {
	return &repoPostgres{Conn: Conn}
}

func (r *repoPostgres) GetForumByThread(ctx context.Context, id int) (string, error) {
	const selectForumByThread = `SELECT Forum 
															 FROM threads 
															 WHERE Id = $1;`

	row := r.Conn.QueryRow(ctx, selectForumByThread, id)
	var forumSlug string
	err := row.Scan(&forumSlug)
	if err != nil {
		return "", models.NotFound
	}
	return forumSlug, nil
}

func (r *repoPostgres) UpdateThread(ctx context.Context, slugOrId string, thread models.Thread) (models.Thread, error) {
	updateThread := `UPDATE threads 
									 SET Title=coalesce(nullif($1, ''), Title), Author=coalesce(nullif($2, ''), Author), 
											 Forum=coalesce(nullif($3, ''), Forum), Message=coalesce(nullif($4, ''), Message), 
										   Votes=coalesce(nullif($5, 0), Votes), Created=coalesce(nullif($6, make_timestamp(1, 1, 1, 0, 0, 0)), Created) `
	var row pgx.Row
	if id, err := strconv.Atoi(slugOrId); err == nil {
		updateThread += `	WHERE Id = $7 
											RETURNING Id, Title, Author, Forum, Message, Slug, Created;`

		row = r.Conn.QueryRow(ctx, updateThread, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Votes, thread.Created, id)
	} else {
		updateThread += ` WHERE Slug = $7 
											RETURNING Id, Title, Author, Forum, Message, Slug, Created;`

		row = r.Conn.QueryRow(ctx, updateThread, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Votes, thread.Created, slugOrId)
	}

	newThread := models.Thread{}
	err := row.Scan(&newThread.ID, &newThread.Title, &newThread.Author, &newThread.Forum, &newThread.Message, &newThread.Slug, &newThread.Created)
	if err != nil {
		return thread, models.NotFound
	}
	return newThread, nil
}

func (r *repoPostgres) GetThread(ctx context.Context, slugOrId string) (models.Thread, error) {
	var row pgx.Row
	if id, err := strconv.Atoi(slugOrId); err == nil {
		const selectThreadBySlug = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created 
																FROM threads 	
																WHERE $1 = Id;`

		row = r.Conn.QueryRow(ctx, selectThreadBySlug, id)
	} else {
		const selectThreadBySlug = `SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created 
																FROM threads 
																WHERE $1 = Slug;`

		row = r.Conn.QueryRow(ctx, selectThreadBySlug, slugOrId)
	}
	finalThread := models.Thread{}
	err := row.Scan(&finalThread.ID, &finalThread.Title, &finalThread.Author, &finalThread.Forum, &finalThread.Message, &finalThread.Votes, &finalThread.Slug, &finalThread.Created)
	if err != nil {
		return finalThread, models.NotFound
	}
	return finalThread, nil
}

func (r *repoPostgres) CreatePosts(ctx context.Context, thread int, forum string, posts []models.Post) ([]models.Post, error) {
	var insertPost = `INSERT INTO posts (Author, Created, Forum, IsEdited, Message, Parent, Thread) VALUES `
	values := make([]interface{}, 0)

	createdDate := time.Now()
	var finalPosts []models.Post
	for i, post := range posts {
		insertPost += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d),", i*7+1, i*7+2, i*7+3, i*7+4, i*7+5, i*7+6, i*7+7)
		values = append(values, post.Author, createdDate, forum, post.IsEdited, post.Message, post.Parent, thread)

		if post.Parent != 0 {
			foundThread := 0
			selectThread := `SELECT Thread FROM posts WHERE Id = $1;`
			row := r.Conn.QueryRow(ctx, selectThread, post.Parent)
			err := row.Scan(&foundThread)

			if err != nil || foundThread != thread {
				return []models.Post{}, models.Conflict
			}
		}

		selectAuthor := `SELECT Nickname FROM users WHERE Nickname = $1;`
		row := r.Conn.QueryRow(ctx, selectAuthor, post.Author)
		var author string
		err := row.Scan(&author)
		if err != nil {
			return []models.Post{}, models.NotFound
		}
	}

	insertPost = strings.TrimSuffix(insertPost, ",")
	insertPost += ` RETURNING Id, Author, Created, Forum, IsEdited, Message, Parent, Thread, Path;`
	rows, err := r.Conn.Query(ctx, insertPost, values...)

	if err != nil {
		if pqError, ok := err.(*pgconn.PgError); ok {
			switch pqError.Code {
			case "23503": // foreign key violation
				return posts, models.NotFound
			default:
				return posts, models.InternalError
			}
		}
	}
	defer rows.Close()

	for range posts {
		if rows.Next() {
			foundPost := models.Post{}
			err = rows.Scan(&foundPost.ID, &foundPost.Author, &foundPost.Created, &foundPost.Forum, &foundPost.IsEdited, &foundPost.Message, &foundPost.Parent, &foundPost.Thread, &foundPost.Path)
			if err != nil {
				return posts, models.InternalError
			}
			finalPosts = append(finalPosts, foundPost)
		}
	}

	return finalPosts, nil
}

func (r *repoPostgres) CreateVote(ctx context.Context, thread int, vote models.Vote) error {
	const insertVote = `INSERT INTO votes (Author, Voice, Thread) 
											VALUES ($1, $2, $3);`

	_, err := r.Conn.Exec(ctx, insertVote, vote.Nickname, vote.Voice, thread)
	if err != nil {
		if pqError, ok := err.(*pgconn.PgError); ok {
			switch pqError.Code {
			case "23505": // unique violation
				return models.Conflict
			default:
				return models.InternalError
			}
		}
	}
	return nil
}

func (r *repoPostgres) ChangeVote(ctx context.Context, thread int, vote models.Vote) error {
	const selectVote = `SELECT Voice 
											FROM votes 
             					WHERE Author = $1 AND Thread = $2;`

	row := r.Conn.QueryRow(ctx, selectVote, vote.Nickname, thread)
	var voice int
	err := row.Scan(&voice)
	if err != nil {
		return models.InternalError
	}

	if voice == vote.Voice {
		return nil
	}

	const updateVote = `UPDATE votes 
											SET Voice = $1 
             					WHERE Author = $2 AND Thread = $3;`

	_, err = r.Conn.Exec(ctx, updateVote, vote.Voice, vote.Nickname, thread)
	if err != nil {
		return models.InternalError
	}
	return nil
}

func (r *repoPostgres) GetPostsFlat(ctx context.Context, thread int, limit string, since string, desc string) ([]models.Post, error) {
	var rows pgx.Rows
	var errQuery error
	selectPosts := `SELECT Id, Author, Created, Forum, IsEdited, Message, Parent, Thread 
									FROM posts 
									WHERE Thread = $1`

	if limit == "" {
		if since != "" && desc == "true" {
			selectPosts += ` AND Id < $2 ORDER BY Id DESC`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, since)
		}
		if since == "" && desc == "true" {
			selectPosts += ` ORDER BY Id DESC`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread)
		}
		if since != "" && desc != "true" {
			selectPosts += ` AND Id > $2 ORDER BY Id ASC`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, since)
		}
		if since == "" && desc != "true" {
			selectPosts += ` ORDER BY Id ASC`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread)
		}
	} else {
		if since != "" && desc == "true" {
			selectPosts += ` AND Id < $2 ORDER BY Id DESC LIMIT $3`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, since, limit)
		}
		if since == "" && desc == "true" {
			selectPosts += ` ORDER BY Id DESC LIMIT $2`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, limit)
		}
		if since != "" && desc != "true" {
			selectPosts += ` AND Id > $2 ORDER BY Id ASC LIMIT $3`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, since, limit)
		}
		if since == "" && desc != "true" {
			selectPosts += ` ORDER BY Id ASC LIMIT $2`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, limit)
		}
	}
	if errQuery != nil {
		return []models.Post{}, models.InternalError
	}
	defer rows.Close()

	posts := make([]models.Post, 0)
	for rows.Next() {
		postOne := models.Post{}
		err := rows.Scan(&postOne.ID, &postOne.Author, &postOne.Created, &postOne.Forum, &postOne.IsEdited, &postOne.Message, &postOne.Parent, &postOne.Thread)
		if err != nil {
			return []models.Post{}, models.InternalError
		}
		posts = append(posts, postOne)
	}
	return posts, nil
}

func (r *repoPostgres) GetPostsTree(ctx context.Context, thread int, limit string, since string, desc string) ([]models.Post, error) {
	var rows pgx.Rows
	var errQuery error
	selectPosts := `SELECT posts.Id, posts.Author, posts.Created, posts.Forum, posts.IsEdited, posts.Message, posts.Parent, posts.Thread, posts.Path
									FROM posts`

	if limit == "" {
		if since != "" && desc == "true" {
			selectPosts += ` JOIN posts parent ON parent.id = $2 
											 WHERE posts.path < parent.path AND posts.thread = $1 
											 ORDER BY posts.path DESC, posts.id DESC`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, since)
		}
		if since == "" && desc == "true" {
			selectPosts += ` WHERE posts.Thread = $1 
										 	 ORDER BY posts.Path DESC, posts.Id DESC`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread)
		}
		if since != "" && desc != "true" {
			selectPosts += ` JOIN posts parent ON parent.id = $2 
											 WHERE posts.path > parent.path AND posts.thread = $1 
											 ORDER BY posts.path ASC, posts.id ASC`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, since)
		}
		if since == "" && desc != "true" {
			selectPosts += ` WHERE posts.Thread = $1 
										   ORDER BY posts.Path ASC, posts.Id ASC`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread)
		}
	} else {
		if since != "" && desc == "true" {
			selectPosts += ` JOIN posts parent ON parent.id = $2 
											 WHERE posts.path < parent.path AND posts.thread = $1 
											 ORDER BY posts.path DESC, posts.id DESC LIMIT $3`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, since, limit)
		}
		if since == "" && desc == "true" {
			selectPosts += ` WHERE posts.Thread = $1 
										   ORDER BY posts.Path DESC, posts.Id DESC LIMIT $2`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, limit)
		}
		if since != "" && desc != "true" {
			selectPosts += ` JOIN posts parent ON parent.id = $2 
											 WHERE posts.path > parent.path AND posts.thread = $1 
											 ORDER BY posts.path ASC, posts.id ASC LIMIT $3`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, since, limit)
		}
		if since == "" && desc != "true" {
			selectPosts += ` WHERE posts.Thread = $1 
										   ORDER BY posts.Path ASC, posts.Id ASC LIMIT $2`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, limit)
		}
	}
	if errQuery != nil {
		return []models.Post{}, models.InternalError
	}
	defer rows.Close()

	posts := make([]models.Post, 0)
	for rows.Next() {
		postOne := models.Post{}
		err := rows.Scan(&postOne.ID, &postOne.Author, &postOne.Created, &postOne.Forum, &postOne.IsEdited, &postOne.Message, &postOne.Parent, &postOne.Thread, &postOne.Path)

		if err != nil {
			return []models.Post{}, models.InternalError
		}
		posts = append(posts, postOne)
	}
	return posts, nil
}

func (r *repoPostgres) GetPostsParentTree(ctx context.Context, thread int, limit string, since string, desc string) ([]models.Post, error) {
	selectPostParents := fmt.Sprintf(`SELECT Id FROM posts WHERE Thread = %d AND Parent = 0 `, thread)

	if since == "" {
		if desc == "true" {
			selectPostParents += ` ORDER BY Id DESC `
		} else {
			selectPostParents += ` ORDER BY Id ASC `
		}
	} else {
		if desc == "true" {
			selectPostParents += fmt.Sprintf(` AND Path[1] < (SELECT Path[1] FROM posts WHERE Id = %s) ORDER BY Id DESC `, since)
		} else {
			selectPostParents += fmt.Sprintf(` AND Path[1] > (SELECT Path[1] FROM posts WHERE Id = %s) ORDER BY Id ASC `, since)
		}
	}

	if limit != "" {
		selectPostParents += " LIMIT " + limit
	}

	selectPosts := fmt.Sprintf(`SELECT Id, Author, Created, Forum, IsEdited, Message, Parent, Thread FROM posts WHERE Path[1] = ANY (%s) `, selectPostParents)

	if desc == "true" {
		selectPosts += ` ORDER BY Path[1] DESC, Path, Id `
	} else {
		selectPosts += ` ORDER BY Path[1] ASC, Path, Id `
	}

	rows, _ := r.Conn.Query(ctx, selectPosts)
	defer rows.Close()
	posts := make([]models.Post, 0)
	for rows.Next() {
		onePost := models.Post{}
		err := rows.Scan(&onePost.ID, &onePost.Author, &onePost.Created, &onePost.Forum, &onePost.IsEdited, &onePost.Message, &onePost.Parent, &onePost.Thread)
		if err != nil {
			return posts, models.InternalError
		}
		posts = append(posts, onePost)
	}

	return posts, nil
}
