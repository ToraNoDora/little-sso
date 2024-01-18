package Helper

import (
	"fmt"
	"math/rand"

	_ "github.com/lib/pq" // draver pq

	"github.com/ToraNoDora/little-sso/sso/internal/src/domain/models"
	p "github.com/ToraNoDora/little-sso/sso/internal/src/store/postgres"
)

const (
	appsTable = "apps"
)

type RandomApp struct {
	AppID     string
	AppSecret string
}

func GetRandomApp(pCfg p.Config) RandomApp {
	db, err := p.NewPostgresDB(pCfg)
	if err != nil {
		panic(err)
	}

	var apps []models.App
	if err := db.Select(&apps, fmt.Sprintf("SELECT id, secret FROM %s", appsTable)); err != nil {
		panic(err)
	}

	ra := apps[rand.Intn(len(apps))]

	return RandomApp{
		AppID:     ra.ID,
		AppSecret: ra.Secret,
	}
}
