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

### musicapi

WIP

### mysql

WIP
