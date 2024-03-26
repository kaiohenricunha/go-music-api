# go-music-k8s

```plaintext
.
├── Dockerfile
├── README.md
├── cli
│   └── main.go
├── go.mod
├── go.sum
├── helm
│   ├── Chart.yaml
│   ├── customvalues.yaml
│   ├── templates
│   │   ├── NOTES.txt
│   │   ├── _helpers.tpl
│   │   ├── configmap-musicapi.yaml
│   │   ├── deployment-musicapi.yaml
│   │   ├── deployment-mysql.yaml
│   │   ├── ns-musicapi.yaml
│   │   ├── ns-mysql.yaml
│   │   ├── pvc-mysql.yaml
│   │   ├── secret-musicapi.yaml
│   │   ├── secret-mysql.yaml
│   │   ├── service-musicapi.yaml
│   │   ├── service-mysql.yaml
│   │   └── tests
│   │       └── test-connection.yaml
│   └── values.yaml
├── manifests
│   ├── musicapi
│   │   ├── configmap-musicapi.yaml
│   │   ├── deployment-musicapi.yaml
│   │   ├── ns-musicapi.yaml
│   │   ├── secret-musicapi.yaml
│   │   └── service-musicapi.yaml
│   └── mysql
│       ├── deployment-mysql.yaml
│       ├── ns-mysql.yaml
│       ├── pvc-mysql.yaml
│       ├── secret-mysql.yaml
│       └── service-mysql.yaml
└── pkg
    ├── config
    │   ├── config.go
    │   └── env
    ├── db
    │   ├── db.go
    │   └── sql.sh
    ├── musicapp
    │   ├── init.go
    │   ├── music_dao.go
    │   ├── music_methods.go
    │   └── types.go
    └── server
        ├── handlers.go
        └── server.go
```

## Descrição

### Entrypoint

O entrypoint para nossa aplicação será `cli/main.go`.

Aqui, temos uma função `init()` que irá:

- inicializar a configuração do app
- inicializar o banco de dados

Confira a configuração do app. As mesmas credenciais de configuração serão fornecidas ao nosso app através do ConfigMap e Secret do k8s.

```bash
# env
CONFIG_DBHOST # host do banco de dados
CONFIG_DBNAME # nome do banco de dados
CONFIG_DBPASS # senha do banco de dados
CONFIG_DBUSER # usuário do banco de dados
CONFIG_SERVER_PORT # porta do servidor web
```

Em seguida, ele executará a função main() que irá iniciar o servidor web Golang, com a configuração inicializada:

```go
var cfg *config.Config
var err error

func init() {
	log.Print("Welcome to music api...")

	// get a config
	cfg, err = config.NewConfig()
	if err != nil {
		log.Fatal("Config init failed", err)
	}

	// migrate db
	if err = musicapp.DbInit(cfg.DB); err != nil {
		log.Fatal("DB migration failed...")
	}
}

func main() {
	server.Start(cfg)
}
```

### Database

Para simplificar, utilizei o GORM, um ORM para Go. Song será uma entidade em nosso banco de dados.

```go
type Song struct {
	gorm.Model
	Name   string `json:"name"`
	Artist string `json:"artist"`
}
```

Por padrão, gorm.Model inclui os seguintes campos:

```go
type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt DeletedAt `gorm:"index"`
}
```

Essa estrutura permite a manipulação eficiente de registros de músicas no banco de dados, aproveitando as funcionalidades do GORM para operações CRUD, além de gerenciamento automático de campos comuns como ID, CreatedAt, UpdatedAt, e DeletedAt.

### API

Por enquanto, manteremos 2 endpoints da API em nossa aplicação, definidos em server/server.go.

```go
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/music", GetSongHandler).Methods("GET")
	r.HandleFunc("/api/v1/music", PostSongHandler).Methods("POST")
```

- O método GET buscará todas as músicas do banco de dados.
- O método POST adicionará uma nova música ao banco de dados.

Ambos os endpoints utilizarão o mesmo caminho /api/v1/music.

O GetSongHandler consulta simplesmente todas as músicas e as retorna.

```go
func GetSongHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("GET /api/v1/music")

	// buscar músicas
	songs, err := musicapp.GetAllSongs(gcfg)
	if err != nil {
		log.Print("Falha na API de buscar músicas", err.Error())
		w.Write([]byte("Erro ao obter músicas!"))
		return
	}

	// resposta
	json.NewEncoder(w).Encode(songs)
}
```

O PostSongHandler adiciona uma nova música ao banco de dados.

```go
func PostSongHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("POST /api/v1/music")

	// adicionar música
	newsong := musicapp.Song{}
	json.NewDecoder(r.Body).Decode(&newsong)

	err := musicapp.PostSong(gcfg, &newsong)
	if err != nil {
		log.Print("Falha na API de adicionar música", err.Error())
		w.Write([]byte("Erro ao publicar música!"))
		return
	}

	// resposta
	w.Write([]byte("Nova música publicada!"))
}
```

Os métodos correspondentes do DAO de banco de dados estão definidos em music_dao.go.

```go
type musicDaoInterface interface {
	Get(db *gorm.DB) ([]*Song, error)
	Post(db *gorm.DB, song *Song) error
}

type musicApiOrm struct {
}

func NewMusicOrm() musicDaoInterface {
	return &musicApiOrm{}
}

func (m *musicApiOrm) Get(db *gorm.DB) ([]*Song, error) {
	songsArr := make([]*Song, 0)
	if err := db.Model(&Song{}).Find(&songsArr).Error; err != nil {
		return nil, err
	}
	log.Printf("%d linhas encontradas", len(songsArr))

	return songsArr, nil
}

func (m *musicApiOrm) Post(db *gorm.DB, song *Song) error {
	if err := db.Session(&gorm.Session{FullSaveAssociations: true, CreateBatchSize: CreateBatchSizeDefault}).Model(Song{}).Create(song).Error; err != nil {
		return err
	}

	return nil
}
```

## Unit Tests

Para testar a aplicação, utilizaremos o pacote `testing` do Go.

```go
package musicapp

import (
	"github.com/kaiohenricunha/go-music-k8s/pkg/config"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"testing"
)

func TestPostSong(t *testing.T) {
	db := SetupTestDB()     // Initialize and migrate your test database
	defer clearDatabase(db) // Ensure the database is cleared after the test runs.
	cfg := &config.Config{DB: db}

	song := &Song{Name: "Test Song", Artist: "Test Artist"}

	err := PostSong(cfg, song)
	assert.NoError(t, err, "Failed to post a new song")
}
```

In the pkg/musicapp/musicapp_test.go file, we have a test function for each of the operations in the musicapp API. To run the tests, execute the following command:

```bash
go test ./pkg/musicapp
```

You should see the following output:

```bash
go test
2024/03/25 00:57:26 2 rows found
PASS
ok      github.com/kaiohenricunha/go-music-k8s/pkg/musicapp     0.496s
```

## Deployment

O ambiente escolhido foi [Kubernetes](https://kubernetes.io/) através do [minikube](https://minikube.sigs.k8s.io/docs/start/).

```
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-amd64
sudo install minikube-darwin-amd64 /usr/local/bin/minikube
minikube start
```

A partir de então, podemos interagir com a aplicação através do [kubectl](https://kubernetes.io/docs/tasks/tools/).

### mysql

O banco de dados MySQL ficará em um namespace chamado `db-ns`.

![mysql](./images/mysql.png)

Todos os components do MySQL, como o Deployment, Service, PersistentVolumeClaim, ConfigMap e Secret, serão definidos em um arquivo YAML usando dry-run, como nos exemplos abaixo:

```bash
# manifests/mysql 
> kubectl create ns db-ns --dry-run -oyaml > ns.yaml
> kubectl create deployment mysql --image=mysql:5.7 --dry-run -oyaml > deployment.yaml
```

Com os arquivos YAML prontos, podemos implantar o MySQL no cluster minikube.

```bash
> kubectl apply -f manifests/mysql
```

Para resumir, o banco de dados MySQL precisa dos seguintes componentes:

- Namespace: db-ns
- Secret com credenciais de banco de dados
- Um Service para expor o MySQL. Funciona como um ponto de acesso para o MySQL(DNS)
- Um PersistentVolumeClaim para armazenamento persistente dos dados

É importante observar qual a storage class padrão do cluster utilizado. No caso do Minikube:

```bash
kubectl get storageclasses.storage.k8s.io
NAME                 PROVISIONER                RECLAIMPOLICY   VOLUMEBINDINGMODE   ALLOWVOLUMEEXPANSION   AGE
standard (default)   k8s.io/minikube-hostpath   Delete          Immediate           false                  10m
```

Com essa informação, podemos definir o PersistentVolumeClaim para o MySQL.

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-pv-claim
  namespace: db-ns
spec:
  storageClassName: standard
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 20Gi
```

Com tudo pronto, podemos verificar se o MySQL foi implantado com sucesso.

```bash
kubectl get pods -A                                
NAMESPACE     NAME                               READY   STATUS    RESTARTS      AGE
db-ns         mysql-bdccddf9b-98qtz              1/1     Running   0             80s
kube-system   coredns-5dd5756b68-nhwfg           1/1     Running   0             18m
kube-system   etcd-minikube                      1/1     Running   0             19m
kube-system   kube-apiserver-minikube            1/1     Running   0             19m
kube-system   kube-controller-manager-minikube   1/1     Running   1 (19m ago)   19m
kube-system   kube-proxy-m6jbs                   1/1     Running   0             18m
kube-system   kube-scheduler-minikube            1/1     Running   0             19m
kube-system   storage-provisioner                1/1     Running   1 (18m ago)   18m
```

Como podemos ver, o pod mysql-bdccddf9b-98qtz está em execução no namespace db-ns, junto com outros pods do sistema.

### Backend em Golang

A API de música ficará em um namespace chamado `music-ns`.

![musicapi](./images/musicapi.png)

Para os componentes de Kubernetes puro, como namespaces, utilizaremos uma abordagem semelhante à do MySQL para criar os arquivos YAML, usando dry-run.

```bash
> kubectl create ns music-ns --dry-run -oyaml > ns.yaml
```

A API também precisará de um Secret, que carregará as credenciais do banco de dados.

```yaml
apiVersion: v1
data:
  rootpassword: Z3JlZW4=
kind: Secret
metadata:
  creationTimestamp: null
  name: mysql-password
  namespace: music-ns
```

É possível confirmar se o conteúdo do Secret foi criado corretamente usando base64.

```bash
echo -n "senha-escolhida" | base64
```

Precisaremos também de uma forma de passar todas as variáveis de ambiente para o nosso app. Para isso, utilizaremos um ConfigMap.

```bash
> kubectl create configmap music-cm -n music-ns 
--from-literal serverport=8081 
--from-literal dbuser=root 
--from-literal dbname=infnet_music_db 
--dry-run -oyaml > manifests/mysql/configmap.yaml
```

Finalmente, podemos implantar a aplicação de música no cluster minikube com a imagem `kaiohenricunha/musicapi:1.0.0` que geramos.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: musicapi
  name: musicapi
  namespace: music-ns
spec:
  replicas: 1
  selector:
    matchLabels:
      app: musicapi
  template:
    metadata:
      labels:
        app: musicapi
    spec:
      containers:
      - image: kaiohenricunha/go-music-k8s:1.0.0
        name: go-music-k8s
        resources:
          limits:
            cpu: 100m
            memory: 100Mi
          requests:
            cpu: 100m
            memory: 100Mi
        env:
          - name: CONFIG_DBPASS
            valueFrom:
              secretKeyRef:
                key: rootpassword
                name: mysql-password
          - name: CONFIG_DBNAME
            valueFrom:
              configMapKeyRef:
                key: dbname
                name: music-cm
          - name: CONFIG_DBUSER
            valueFrom:
              configMapKeyRef:
                key: dbuser
                name: music-cm
          - name: CONFIG_SERVER_PORT
            valueFrom:
              configMapKeyRef:
                key: serverport
                name: music-cm
          - name: CONFIG_DBHOST
            valueFrom:
              configMapKeyRef:
                key: dbhost
                name: music-cm
```

É o deployment que não só cria o pod com a aplicação, mas também a associa com todos os componentes necessários, como ConfigMap, Secret e Service.

Com tudo pronto, podemos verificar se a aplicação de música foi implantada com sucesso.

```bash
kubectl apply -f manifests/musicapi
```

```bash
kubectl get pods -A
NAMESPACE     NAME                               READY   STATUS             RESTARTS       AGE
db-ns         mysql-698ff8f95d-4qvfw             1/1     Running            0              67s
kube-system   coredns-5dd5756b68-nhwfg           1/1     Running            0              40m
kube-system   etcd-minikube                      1/1     Running            0              40m
kube-system   kube-apiserver-minikube            1/1     Running            0              40m
kube-system   kube-controller-manager-minikube   1/1     Running            1 (40m ago)    40m
kube-system   kube-proxy-m6jbs                   1/1     Running            0              40m
kube-system   kube-scheduler-minikube            1/1     Running            0              40m
kube-system   storage-provisioner                1/1     Running            1 (39m ago)    40m
music-ns      musicapi-78b9f69c4d-qvs9z          0/1     CrashLoopBackOff   5 (98s ago)    5m26s
music-ns      musicapi-7d4c877f5d-54m4j          0/1     CrashLoopBackOff   5 (110s ago)   6m29s
```

Como podemos ver, o pod musicapi-78b9f69c4d-qvs9z está em CrashLoopBackOff. Ao investigar o log do pod, podemos ver que o motivo é que o pod não consegue se conectar ao banco de dados.

```bash
kubectl logs musicapi-78b9f69c4d-qvs9z -n music-ns
2024/03/25 23:07:02 Welcome to music api...
2024/03/25 23:07:02 failed to connect to databaseError 1049 (42000): Unknown database 'infnet_music_db'

2024/03/25 23:07:02 /go/pkg/mod/github.com/!jana!sabuj/music-api-k8s@v0.0.0-20230401103529-67db1fbe644c/pkg/db/db.go:11
[error] failed to initialize database, got error Error 1049 (42000): Unknown database 'infnet_music_db'
```

Isso ocorre porque o banco de dados não foi criado explicitamente. Para corrigir isso, podemos acessar o pod mysql e criar o banco de dados manualmente. Para tal, é necessário "entrar" no pod mysql.

```bash
kubectl get pods -n db-ns 
NAME                     READY   STATUS    RESTARTS   AGE
mysql-698ff8f95d-4qvfw   1/1     Running   0          11m
```

```bash
kubectl exec -it -n db-ns mysql-698ff8f95d-4qvfw -- bash
bash-4.4# ls
bin  boot  dev  docker-entrypoint-initdb.d  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
bash-4.4# 
```

Dentro do pod, podemos acessar o MySQL com o comando `mysql -u root -p` e criar o banco de dados.

```bash
bash-4.4# mysql -u root -p
Enter password: 
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 14
Server version: 8.3.0 MySQL Community Server - GPL

Copyright (c) 2000, 2024, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> 
```

Finalmente, devemos executar os comandos SQL para criar o banco de dados e a tabela.

```sql
mysql> create database infnet_music_db;
Query OK, 1 row affected (0.02 sec)

mysql> use infnet_music_db;
Database changed

mysql> INSERT INTO songs (created_at, updated_at, name, artist) VALUES (now(), now(), "Borbulhas de Amor", "Fagner");
Query OK, 1 row affected (0.04 sec)
```

```sql
mysql> select * from songs;
+----+---------------------+---------------------+------------+-------------------+---------------------+
| id | created_at          | updated_at          | deleted_at | name              | artist              |
+----+---------------------+---------------------+------------+-------------------+---------------------+
|  1 | 2024-03-25 23:17:32 | 2024-03-25 23:17:32 | NULL       | Burn It Down      | Linkin Park         |
|  2 | 2024-03-25 23:17:40 | 2024-03-25 23:17:40 | NULL       | Earth Song        | Michael Jackson     |
|  3 | 2024-03-25 23:17:48 | 2024-03-25 23:17:48 | NULL       | Hey Jude          | The Beatles         |
|  4 | 2024-03-25 23:17:55 | 2024-03-25 23:17:55 | NULL       | Sound of Silence  | Simon and Garfunkel |
|  5 | 2024-03-25 23:18:02 | 2024-03-25 23:18:02 | NULL       | Hotel California  | The Eagles          |
|  6 | 2024-03-25 23:18:13 | 2024-03-25 23:18:13 | NULL       | Comfortably Numb  | Pink Floyd          |
|  7 | 2024-03-25 23:18:38 | 2024-03-25 23:18:38 | NULL       | Borbulhas de Amor | Fagner              |
+----+---------------------+---------------------+------------+-------------------+---------------------+
7 rows in set (0.00 sec)

mysql> 
```

Com o banco de dados criado, o pod musicapi deve ser reiniciado automaticamente e, desta vez, ele deve ser capaz de se conectar ao banco de dados.

```bash
kubectl logs -f musicapi-78b9f69c4d-qvs9z -n music-ns
2024/03/25 23:17:10 Welcome to music api...
```

## Frontend em React

WIP
