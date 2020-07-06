package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/leofrancocalpa/server-analyzer/pkg/repository/model"
)

// Repository interface
type Repository interface {
	CreateServerInfo(serverInfo model.ServerInfo) bool
	UpdateEntry(serverInfo model.ServerInfo)
	FetchServersInfo() (model.ItemsInfo, error)
	GetServerInfoFromDB(dns string) (time.Time, string, []byte)
	ExistEntry(dns string) bool
	Close()
}

// ServerInfoRepository instance
type ServerInfoRepository struct {
	db *sql.DB
}

//NewRepository creates an instance of ServerInfoRepository
func NewRepository() Repository {
	dB, err := createDBConnection()
	if err != nil {
		fmt.Println("Error conecting DB")
		panic("Exiting...")
	}
	return ServerInfoRepository{db: dB}
}

//CreateServerInfo persist in the DB the serverInfo queried
func (repo ServerInfoRepository) CreateServerInfo(serverInfo model.ServerInfo) bool {

	sqlstm := `INSERT INTO serverinfo (dns, data)
		VALUES ($1, $2)`
	sqlresult, err := repo.db.Exec(sqlstm, serverInfo.DNS, serverInfo.Data)
	if err != nil {
		fmt.Println("[ERROR] FAILURE CREATING")
		repo.Close()
		return false
	}
	fmt.Println("CREATE STATUS: ", sqlresult)
	return true

}

//UpdateEntry update an entry with the given dns
func (repo ServerInfoRepository) UpdateEntry(serverInfo model.ServerInfo) {
	sqlstm := `UPDATE serverinfo SET (last_updated, data) = (now(), $1) WHERE dns = $2`
	_, err := repo.db.Exec(sqlstm, serverInfo.Data, serverInfo.DNS)
	if err != nil {
		fmt.Println("[ERROR DB] Failure Updating")
		fmt.Println(err.Error())
	}
	fmt.Println("[UPDATED] successful")
}

//FetchServersInfo obtains all servers info queried before
func (repo ServerInfoRepository) FetchServersInfo() (model.ItemsInfo, error) {
	fmt.Println("Fetching")
	sqlstm := `SELECT dns, data FROM serverinfo`
	rows, err := repo.db.Query(sqlstm)
	if err != nil {
		fmt.Print("[ERROR] FAILURE fetching data")
	}
	defer rows.Close()

	var items []model.ServerInfo
	for rows.Next() {
		var servInfo model.ServerInfo
		if err := rows.Scan(&servInfo.DNS, &servInfo.Data); err != nil {
			log.Println(err)
			continue
		}
		items = append(items, servInfo)
	}
	itemsInfo := model.ItemsInfo{
		Items: items,
	}
	fmt.Println("Fetching finished")
	return itemsInfo, nil
}

//GetServerInfoFromDB Fff
func (repo ServerInfoRepository) GetServerInfoFromDB(dns string) (time.Time, string, []byte) {

	sqlstm := `SELECT last_updated, data->'ssl_grade', data->'servers' FROM serverinfo WHERE dns='` + dns + `'`
	rows, err := repo.db.Query(sqlstm)
	if err != nil {
		fmt.Println("[ERROR DB] Failure while trying to get serverInfo (last_updated, sslgrade, data)")
	}
	var timestmp time.Time
	var sslgrade string
	var servs []byte
	rows.Next()
	if err := rows.Scan(&timestmp, &sslgrade, &servs); err != nil {
		fmt.Println("FF " + err.Error())
	}
	fmt.Println(timestmp)

	return time.Time{}, strings.Replace(sslgrade, `"`, ``, 2), servs

}

//ExistEntry return true if exist an entry with a given dns
func (repo ServerInfoRepository) ExistEntry(dns string) bool {
	sqlstm := `SELECT count(*) FROM serverinfo WHERE dns='` + dns + `'`
	rows, _ := repo.db.Query(sqlstm)
	var count int
	rows.Next()
	rows.Scan(&count)
	if count == 1 {
		return true
	} else {
		return false
	}
}

//Close the connnection with DB
func (repo ServerInfoRepository) Close() {
	repo.db.Close()
}
