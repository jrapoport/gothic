package migration

import (
	"sort"
	"strings"
	"sync"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/jrapoport/gothic/store/drivers"
	"gorm.io/gorm"
)

// Plan represents a database migration plan as an ordered slice of migrations.
type Plan struct {
	migs []*Migration
	mu   sync.RWMutex
}

// NewPlan returns a new migration Plan.
func NewPlan() *Plan {
	return &Plan{migs: []*Migration{}}
}

// AddMigration add the Migration to the migration Plan.
func (plan *Plan) AddMigration(mg *Migration) {
	plan.AddMigrations([]*Migration{mg})
}

// AddMigrations adds a slice of Migration to the migration Plan.
func (plan *Plan) AddMigrations(mgs []*Migration) {
	plan.mu.Lock()
	defer plan.mu.Unlock()
	plan.migs = append(plan.migs, mgs...)
}

// Clear removes all migrations from the Plan.
func (plan *Plan) Clear() {
	plan.mu.Lock()
	defer plan.mu.Unlock()
	plan.migs = []*Migration{}
}

// Run executed all migrations in the Plan in a FIFO order. When sorted
// is true, the migrations are run in alphabetical order by id.
func (plan *Plan) Run(db *gorm.DB, sorted bool) error {
	plan.mu.RLock()
	defer plan.mu.RUnlock()
	name := db.Migrator().CurrentDatabase()
	log := db.Logger
	log.Info(nil, "migrating %s", name)
	migs := plan.migs
	if len(migs) <= 0 {
		log.Info(nil, "no migrations")
		return nil
	}
	planned := make([]*Migration, len(migs))
	for i, mig := range migs {
		planned[i] = mig
	}
	if sorted {
		sort.Slice(planned, func(i, j int) bool {
			return strings.Compare(planned[i].ID, planned[j].ID) == -1
		})
	}
	for idx, mg := range planned {
		log.Info(nil, "%s run migration %03d: %s", name, idx+1, mg.ID)
	}
	opts := *gormigrate.DefaultOptions
	n := opts.TableName
	opts.TableName = db.NamingStrategy.TableName(n)
	mgs := make([]*gormigrate.Migration, len(planned))
	for i, mg := range planned {
		mgs[i] = &mg.Migration
	}
	err := gormigrate.New(db, &opts, mgs).Migrate()
	if err != nil {
		return err
	}
	if db.Name() == string(drivers.SQLite) {
		// SQLite is the worst. Basically if we have a foreign key,
		// the indexes will get prefixed properly when the migrator
		// runs, but then re-added if the table is also an fk in a
		// different model. This results in the un-prefixed indexes
		// being re-added to the table. Later if another model comes
		// along that has the same indexes, all the work of pre-fixing
		// them in the first place will be undone, and it will fail.
		// This way we migrate the tables one at a time and commit
		// after each migration, which allows us to then go back and
		// clean up the indexes for *all* the tables each time again
		// in case they were (re)added by migrating an fk.
		err = db.Transaction(func(tx *gorm.DB) error {
			for _, p := range planned {
				err = migrateIndexes(db, p.model, p.indexes)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	log.Info(nil, "%s migration complete", name)
	return nil
}
