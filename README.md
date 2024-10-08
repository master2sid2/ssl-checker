![golang logo](https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_LightBlue.png)

# Golang SSL certificate validation application

# Table of Contents

- [Overview](#overview)
- [Screenshots](#screenshots)
- [Run](#run)
- [Build](#build)
- [Bugs](#bugs)
- [What I want to do](#what-i-want-to-do)

# Overview
Golang version 1.23

The application provides an easy-to-use interface for tracking certificate expiration dates.
Features:
Implemented authorization with the ability to log in users with User and Admin roles
User role can only view the list of domains.
Admin role has the ability to add and remove domains as well as manage users
Implemented a background task that updates the cache every hourImplemented a background task that updates the cache every hour

# Screenshots
User based interface
![image](https://github.com/user-attachments/assets/42674360-cbd3-4e7a-93b7-331407b44ed9)

Admin based interface
![image](https://github.com/user-attachments/assets/080411a4-f86b-4baf-b4ae-5452eafef26e)

Admin manage users interface
![image](https://github.com/user-attachments/assets/1a319549-ec1f-4866-bd8a-65af3ef39cf0)


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

# What I want to do
* Add languages
* Add configs as command line arguments
* Add configs as config files
* Add prometheus metrics
* Prepare docker file
* Prepare Helm chart
* Maybe something else :)
