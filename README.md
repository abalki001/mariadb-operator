# MariaDB server operator for Kubernetes

## Features
* Setup a MariaDB server. Database server version can be configured in the CR file.
* Creates a new custom Database along with a user credential set for the custom database.
* Operator uses Persistent Volume where MariaDB can write its data files.
* Seamless upgrades of MariaDB is possible without loosing data.


## CRs

### MariaDB
```yaml
apiVersion: mariadb.persistentsys/v1alpha1
kind: MariaDB
metadata:
  name: mariadb
spec:
  # Keep this parameter value unchanged.
  size: 1
  
  # Root user password
  rootpwd: password

  # New Database name
  database: test-db
  # Database additional user details (base64 encoded)
  username: db-user 
  password: db-user 

  # Image name with version
  image: "mariadb/server:10.3"
```

This Operator with create a database called `test-db`, along with user credentials.
The Server image name is mentioned in "image" parameter.

## Setup Instructions
MariaDB Database uses external location on host to store all DB files. This location is default set to "/mnt/data" in YAML file: deploy/mariadb_pv.yaml.
Ensure that the path exists for DB files and has all necessary permissions.

### Start operator and create all resources
Run the following make command to start all resources:
```
# make install
```
By default, all resources will be created in a namespace called "mariadb"

### Stop Operator and delete all resources
Run the command to delete all resources:
```
# make install
```

Verify list of pods. One Operator and One Server pod should be created.
```
# kubectl get pods -n mariadb
NAME                              READY   STATUS    RESTARTS   AGE
mariadb-operator-78c95468-m824g   1/1     Running   0          118s
mariadb-server-778b9b7cb5-nt6n5   1/1     Running   0          109s
```

Verify that "mariadb-service" is created.
```
# kubectl get svc -n mariadb
NAME                       TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)             AGE
mariadb-operator-metrics   ClusterIP   10.109.149.148   <none>        8383/TCP,8686/TCP   2m
mariadb-service            NodePort    10.100.40.161    <none>        80:30685/TCP        119s
```
Service "mariadb-service" is a NodePort Service that exposes mariadb service on port 3306
```
# kubectl describe service/mariadb-service -n mariadb
Name:                     mariadb-service
Namespace:                mariadb
Labels:                   MariaDB_cr=mariadb
                          app=MariaDB
                          tier=mariadb
Annotations:              <none>
Selector:                 MariaDB_cr=mariadb,app=MariaDB,tier=mariadb
Type:                     NodePort
IP:                       10.100.40.161
Port:                     <unset>  80/TCP
TargetPort:               3306/TCP
NodePort:                 <unset>  30685/TCP
Endpoints:                172.17.0.5:3306
Session Affinity:         None
External Traffic Policy:  Cluster
Events:                   <none>
```

Test connectivity to MariaDB Server:
```
# mysql -h 172.17.0.5 -P 3306 -u db-user -pdb-user
```
where 172.17.0.5 is the IP of mariadb-server pod and 3306 is the Target Port.
If everything is correct, mysql prompt will be presented to the user.

## Upgrade MariaDB server version
Mariadb server image is mentioned in CR key "image".
To upgrade the Server version, change value of "image" key and reapply the CR YAML file.
```
# kubectl apply -f deploy/crds/mariadb.persistentsys_v1alpha1_mariadb_cr.yaml
```




