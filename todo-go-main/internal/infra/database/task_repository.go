package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const TasksTableName = "tasks"

type task struct {
	Id          uint64            `db:"id,omitempty"`
	UserId      uint64            `db:"user_id"`
	Title       string            `db:"title"`
	Description string            `db:"description"`
	Deadline    time.Time         `db:"deadline"`
	Status      domain.TaskStatus `db:"status"`
	CreatedDate time.Time         `db:"created_date"`
	UpdatedDate time.Time         `db:"updated_date"`
	DeletedDate *time.Time        `db:"deleted_date"`
}

type TaskRepository interface {
	Save(t domain.Task) (domain.Task, error)
	Read(id uint64) (domain.Task, error)
	Update(t domain.Task) (domain.Task, error)
	Delete(id uint64) error
}

type taskRepository struct {
	coll db.Collection
	sess db.Session
}

func NewTaskRepository(session db.Session) TaskRepository {
	return taskRepository{
		coll: session.Collection(TasksTableName),
		sess: session,
	}
}

func (r taskRepository) Save(t domain.Task) (domain.Task, error) {
	tsk := r.mapDomainToModel(t)
	tsk.CreatedDate = time.Now()
	tsk.UpdatedDate = time.Now()
	err := r.coll.InsertReturning(&tsk)
	if err != nil {
		return domain.Task{}, err
	}
	return r.mapModelToDomain(tsk), nil
}

func (r taskRepository) Read(id uint64) (domain.Task, error) {
	var tsk task
	err := r.coll.Find(db.Cond{"id": id}).One(&tsk)
	if err != nil {
		return domain.Task{}, err
	}
	return r.mapModelToDomain(tsk), nil
}

func (r taskRepository) Update(t domain.Task) (domain.Task, error) {
	tsk := r.mapDomainToModel(t)
	tsk.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": t.Id}).Update(&tsk)
	if err != nil {
		return domain.Task{}, err
	}
	return r.mapModelToDomain(tsk), nil
}

func (r taskRepository) Delete(id uint64) error {
	deletedDate := time.Now()
	err := r.coll.Find(db.Cond{"id": id}).Update(map[string]interface{}{
		"deleted_date": deletedDate,
	})
	return err
}

func (r taskRepository) mapDomainToModel(t domain.Task) task {
	return task{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Deadline:    t.Deadline,
		Status:      t.Status,
		CreatedDate: t.CreatedDate,
		UpdatedDate: t.UpdatedDate,
		DeletedDate: t.DeletedDate,
	}
}

func (r taskRepository) mapModelToDomain(t task) domain.Task {
	return domain.Task{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Deadline:    t.Deadline,
		Status:      t.Status,
		CreatedDate: t.CreatedDate,
		UpdatedDate: t.UpdatedDate,
		DeletedDate: t.DeletedDate,
	}
}
