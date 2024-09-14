![golang logo](golang_logo.png)

# Golang SSL certificate validation application

# Table of Contents

- [Overview](#overview)
- [Run](#run)
- [Build](#build)
- [Bugs](#bugs)

# Overview
Golang version 1.23

The application provides an easy-to-use interface for tracking certificate expiration dates.
Features:
Implemented authorization with the ability to log in users with User and Admin roles
User role can only view the list of domains.
Admin role has the ability to add and remove domains as well as manage users
Implemented a background task that updates the cache every hourImplemented a background task that updates the cache every hour

# Run
For running application just clone repositiry and execute `go run main.go` command
```$bash
git clone https://github.com/master2sid2/ssl-checker.git
cd ssl-checker
go run main.go
```
The first time the application is launched, a default user (admin/admin) will be created

# Build
For build, run:
```$bash
go build
```

# Bugs
When adding a domain or user, the request is redirected to the page specified in the Action form instead of staying on the same page.
When adding a large number of domains, the next operations of adding or deleting can take a long time (15 minutes for 100+ domains).
