package main

import (
	familytree "github.com/bhborges/family-tree-api/internal"
	"github.com/bhborges/family-tree-api/pkg/db"
	"github.com/bhborges/family-tree-api/pkg/http"
	"github.com/bhborges/family-tree-api/pkg/log"
	"github.com/bhborges/family-tree-api/pkg/monitor"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		log.Module,
		db.PostgresModule,
		db.MigrateModule,
		http.RESTModule,
		monitor.APMModule(),
		familytree.APIModule(),
	).Run()
}
