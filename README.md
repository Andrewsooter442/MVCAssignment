## Setup

**Create the .env file (Sample present: /.sampleenv )**

*Sample uses my DB password*
```
mv .sampleenv .env
```

**Create the database**
```
 make migrate-up
```



**Create the admin and chef users.**
```
make run dev
```
```

curl -X POST http://localhost:3000/signup \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=admin&password=admingod&email=admin@god.com&phone=1234567890"

curl -X POST http://localhost:3000/signup \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=chef&password=chefgod&email=chef@god.com&phone=1234567890"

```

**Stop the server**


**Give admin and chef access to the admin and chef users.**
```
make create-admin
```

**Start the server**

```
make build
make run
```



**Admin user:admin, passwd:admingod**

**Chef user:admin, passwd:chefgod**
