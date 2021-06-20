package constants

const DB_DRIVER = "sqlite3"
const DB_PATH = "./gitmonitor.db"
const INIT_PROJECT_TABLE = `CREATE TABLE IF NOT EXISTS project(
	"project_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"project_dir" TEXT
);`
const INIT_BRANCH_TABLE = `CREATE TABLE IF NOT EXISTS branch(
	"branch_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"project_id" INTEGER,
	"name" TEXT,
	"is_merge_target" INTEGER DEFAULT 0,
	"is_deleted" INTEGER DEFAULT 1,
	FOREIGN KEY (project_id)
	REFERENCES project (project_id)
		ON DELETE CASCADE 
		ON UPDATE NO ACTION
);`
const INIT_TASK_TABLE = `CREATE TABLE IF NOT EXISTS task(
	"task_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"project_id" INTEGER,
	"branch_id" INTEGER,
	"name" TEXT,
	"assignee_name" TEXT,
	"assignee_email" TEXT,
	"task_status" INT DEFAULT 0,
	"start_date" INTEGER,
	"end_date" INTEGER,
	FOREIGN KEY (project_id)
		REFERENCES project (project_id)
			ON DELETE CASCADE 
			ON UPDATE NO ACTION,
	FOREIGN KEY (branch_id)
		REFERENCES branch (branch_id)
			ON DELETE CASCADE 
			ON UPDATE NO ACTION
);`

type TaskStatus int

const (
	Waiting    TaskStatus = 0
	InProgress TaskStatus = 1
	Done       TaskStatus = 2
	Expired    TaskStatus = 3
)

var (
	TaskStatusMap = map[int]string{
		0: "Waiting",
		1: "In progress",
		2: "Done (Associated branch already deleted)",
		3: "Expired",
	}
)

const Separator = "^_^"
