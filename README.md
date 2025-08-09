## Setup

**Create the .env file (Sample present: /.sampleenv )**

*Sample uses my DB password*
```
mv .sampleenv .env
```

**Create the database**
```
CREATE DATABASE IF NOT EXISTS kitchen;
USE kitchen;
```
**Execute the second migration file**

*path*
```
(migrations/000002_create_users_table.up.sql)
```

Now run the main go file
```
go run ./cmd/main.go
```

**Create the admin and chef users.**
```
curl -X POST http://localhost:3000/signup \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=admin&password=admingod&email=admin@god.com&phone=1234567890"

curl -X POST http://localhost:3000/signup \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=chef&password=chefgod&email=chef@god.com&phone=1234567890"
```
**Execute the third and forth migration file**

*path*
```
(migrations/000003_create_users_table.up.sql)
(migrations/000004_create_users_table.up.sql)
```

**Admin user:admin, passwd:admingod**

**Chef user:admin, passwd:chefgod**
