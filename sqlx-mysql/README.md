## Prepare MySQL
Run a MySQL server in docker container:

```go
docker run -d --name mysql1 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7
```

Create a database:

```sql
CREATE DATABASE `test` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;
```

create a table：

```sql
CREATE TABLE test.Person (
	Id integer auto_increment NOT NULL,
	Name VARCHAR(30) NULL,
	City VARCHAR(50) NULL,
	AddTime DATETIME NOT NULL,
	UpdateTime DATETIME NOT NULL,
	CONSTRAINT Person_PK PRIMARY KEY (Id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_general_ci;
```

## Run program

    go run main.go