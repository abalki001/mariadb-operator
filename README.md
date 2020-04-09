# MariaDB server operator for Kubernetes

## Features
* Creates a maridb server and database from a CR


## CRs

### MariaDB
```yaml
apiVersion: mariadb.persistentsys/v1alpha1
kind: MariaDB
metadata:
  name: mariadb
spec:
  database: test-db
  password: db-user
  rootpwd: password
  size: 1
  username: db-user
```

This creates a database called `test-db`.
