package pkg

import (
	"log"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	initConfig("../")

	viper.Set("db.type", "sqlite3")
	viper.Set("db.path", "test.db")

	initDb()

	retCode := m.Run()

	log.Print("Dropping tables...")
	db.DropTable("apps")
	db.DropTable("histories")
	log.Print("Tables dropped !")

	os.Remove("test.db")
	os.Exit(retCode)
}
