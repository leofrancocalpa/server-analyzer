package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/leofrancocalpa/server-analyzer/pkg/repository/model"
)

// Repository interface
type Repository interface {
	CreateServerInfo(serverInfo model.ServerInfo)
	FetchServersInfo() (model.ItemsInfo, error)
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
	fmt.Println(dB.Ping())
	return ServerInfoRepository{db: dB}
}

//CreateServerInfo persist in the DB the serverInfo queried
func (repo ServerInfoRepository) CreateServerInfo(serverInfo model.ServerInfo) {
	sqlstm0 := `SELECT count(*) FROM serverinfo WHERE dns='` + serverInfo.DNS + `'`
	rows, _ := repo.db.Query(sqlstm0)
	var count int
	rows.Next()
	rows.Scan(&count)

	if count < 1 {
		sqlstm := `INSERT INTO serverinfo (dns, data)
		VALUES ($1, $2)`
		sqlresult, err := repo.db.Exec(sqlstm, serverInfo.DNS, serverInfo.Data)
		if err != nil {
			fmt.Println("[ERROR] FAILURE CREATING")
			repo.Close()
			panic(err.Error())
		}
		fmt.Println("CREATE STATUS: ", sqlresult)
	} else {
		fmt.Println("[DB] The entry already exists, the data was not persisted")
	}

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
	fmt.Println(string(items[0].DNS))
	return itemsInfo, nil
}

//Close the connnection with DB
func (repo ServerInfoRepository) Close() {
	repo.db.Close()
}
