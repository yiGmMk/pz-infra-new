package base

import (
	"io/ioutil"
	"strings"

	"github.com/gyf841010/pz-infra-new/confUtil"
	"github.com/gyf841010/pz-infra-new/database"
	"github.com/gyf841010/pz-infra-new/log"
	. "github.com/gyf841010/pz-infra-new/logging"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/ziutek/mymysql/native" // Native engine
)

var initializedConf = false

// fot *_test.go files to read conf file via relative path to intialize db connection
func InitConfigFile(configFileName string) {
	log.Debug("########InitConfigFile##########")

	if initializedConf {
		return
	}
	configFilePath := confUtil.FileHierarchyFind(configFileName)
	log.Debugf("configFilePath is %s", configFilePath)
	beego.LoadAppConfig("ini", configFilePath)
	//dbString := "root:b553e6e21a8ff@tcp(localhost:3306)/pz_base?charset=utf8mb4&parseTime=True"
	dbString := beego.AppConfig.String("db")
	log.Debug("DB String is %s", dbString)
	database.InitDB(dbString, nil)
	if Log == nil {
		InitLogger("testing")
	}
	confUtil.LoadLocaleConfig()
	// Initialized language type list.

	initializedConf = true
}

func InitDBWithFile(sqlFile string, dbs ...*gorm.DB) {
	if len(sqlFile) == 0 {
		return
	}
	bytes, _ := ioutil.ReadFile(sqlFile)

	db := database.GetNonTransactionDatabases(dbs)
	sqls := strings.Split(string(bytes), ";")
	for _, sql := range sqls {
		err := db.Exec(sql).Error
		if err != nil && !(strings.Contains(string(err.Error()), "Query was empty")) {
			Log.Error("execute setup sql failed. error is:", WithError(err))
			panic(err)
		}
	}
}
