# NotiSync-Server
a simple implement for the app called [Notisync](https://github.com/SquareFong/Notisync)

# dependencies

this program relies on [this project](https://github.com/go-sql-driver/mysql). Please install it before use.

Simple install the package to your $GOPATH with the go tool from shell:

```bash
$ go get -u github.com/go-sql-driver/mysql
```

# how to boot

this project use MySQL or MariaDB. Before you run this program, I recommend you to prepare a database and a account which has all properties to this database.

For example, you could run these command to initialize the database:

```bash
$ sudo mysql -u root
MariaDB [(none)]> CREATE DATABASE NotiSync;
MariaDB [(none)]> CREATE USER 'notisync'@'localhost' IDENTIFIED BY '666666';
MariaDB [(none)]> GRANT ALL PRIVILEGES ON NotiSync.* TO 'notisync'@'localhost';
MariaDB [(none)]> FLUSH PRIVILEGES;
```

then write a config file to /etc/notisync/config.json

```
{
    "UserName":"notisync",
    "Password":"666666",
    "DBName":"NotiSync"
}
```

then you should create users table manually

```bash
$ mysql -u notisync
MariaDB [(none)]> create table if not exists Users (id INTEGER key AUTO_INCREMENT, uuid TEXT);
```

finally, add user manually, for example

```bash
MariaDB [(none)]> insert Users(uuid) values('7517e18a-40a6-4902-a7c9-23bd0ef7f00f');
```

# TODO

- [ ] improve database performance
- [ ] add user check before service
- [ ] make it to add user automatically
- [ ] add lock to improve multithread performance
