package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/mattn/go-sqlite3"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	var todo model.TODO
	if subject == "" {
		return &todo, sqlite3.ErrConstraint
	}

	stmt, _ := s.db.PrepareContext(ctx, insert)
	result, _ := stmt.ExecContext(ctx, subject, description)
	todo.ID, _ = result.LastInsertId()
	row := s.db.QueryRowContext(ctx, confirm, todo.ID)
	err := row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)

	return &todo, err
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	todos := make([]*model.TODO,0)
	if prevID == 0{
		rows,_ := s.db.QueryContext(ctx,read,size)
		for rows.Next(){
			var todo model.TODO
				err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
				if(err != nil){
					panic(err)
				}

			todos = append(todos, &todo)
		}
	}else{
		rows,_ := s.db.QueryContext(ctx,readWithID,prevID,size)
		for rows.Next(){
			var todo model.TODO
				err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
				if(err != nil){
					panic(err)
				}

			todos = append(todos, &todo)
		}
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	var todo model.TODO
	todo.ID = id

	if subject == ""{
		return &todo, sqlite3.ErrConstraint
	}

	stmt,err := s.db.PrepareContext(ctx,update)
	if(err !=nil){
		panic(err)
	}
	result,err :=stmt.ExecContext(ctx,subject,description,id)
	if(err !=nil){
		panic(err)
	}

	//ExecContext メソッドの戻り値から変更された Row の数を検査して、0 件だった場合は Station 11 で作成した ErrNotFound を返却するようにしよう
	rowCount, err2 := result.RowsAffected()
	if rowCount == 0{
		return &todo, &model.ErrNotFound{
			When: time.Now(),
			What: "Updated row not found",
		}
	}else if(err2 != nil){
		return &todo,err
	}

	row:= s.db.QueryRowContext(ctx,confirm,id)
	err = row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)

	if(err !=nil){
		panic(err)
	}

	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	if(len(ids) == 0){
		return nil
	}


	delete:= fmt.Sprintf(deleteFmt,strings.Repeat(",?",len(ids)-1))
	strIds := make([]interface{}, len(ids))
	for i,id := range ids{
		strIds[i] = fmt.Sprintf("%d",id);
	}

	stmt,err:= s.db.PrepareContext(ctx,delete)
	if(err !=nil){
		panic(err)
	}
	resp,err :=stmt.ExecContext(ctx,strIds...)
	if(err !=nil){
		panic(err)
	}

	deletedRowCount,err := resp.RowsAffected()
	if(err !=nil){
		panic(err)
	}else if (deletedRowCount == 0){
   return &model.ErrNotFound{
		When: time.Now(),
		What: "no items deleted",
	 }
	}






	
	return nil
}
