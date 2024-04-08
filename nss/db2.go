package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"github.com/joho/godotenv"
	. "github.com/protosam/go-libnss/structs"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

const (
	TOKEN_NOT_FOUND int = iota
)

type GroupMemberinfo struct {
	Username string `usernam:"type,omitempty"`
	// Attributes Attribute `attributes:"type,omitempty"`
}
type UserGroupInfo struct {
	Id   string `id:"type,omitempty"`
	Name string `name:"type,omitempty"`
	Path string `path:"type,omitempty"`
	// Attributes Attribute `attributes:"type,omitempty"`
}
type GroupIdInfo struct {
	 Id         string         `id:"type,omitempty"`
	 Name       string         `name:"type,omitempty"`
	 Path       string         `path:"type,omitempty"`
	Attributes GroupAttribute `attributes:"type,omitempty"`
}

type GroupAttribute struct {
	Gid []string `gid:"type,omitempty"`
	// Gid []string `json:"gid,omitempty"`
}

type UserInfo struct {
	Id         string    `id:"type,omitempty"`
	Username   string    `usernam:"type,omitempty"`
	Attributes Attribute `attributes:"type,omitempty"`
}



type Attribute struct {
	Uids     []string `json:"uid,omitempty"`
	Shells   []string `json:"shell,omitempty"`
	Homedirs []string `json:"homedir,omitempty"`
}

var Gid string
var Gid_uint uint
var Uid_uint uint
var Groupmember string

var dbtest_passwd []Passwd
var dbtest_group []Group
var dbtest_shadow []Shadow

func init() {
	godotenv.Load("/tmp/var.env")
	FQDN := os.Getenv("FQDN")
	auth_admin_user := os.Getenv("auth_admin_user")
	auth_admin_pass := os.Getenv("auth_admin_pass")
	currentUserDir := os.Getenv("HOME")
	dbPath := currentUserDir + "/sqlite-database.db"

	// var Username string
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	token_client := &http.Client{Transport: tr}
	openIDURL := "https://" + FQDN + "/auth/realms/master/protocol/openid-connect/token"
	// logrus.Info("OIDC Endpoint is" + openIDURL)
	formData := url.Values{
		"client_id":     {"admin-cli"},
		"grant_type":    {"password"},
		"username":      {auth_admin_user},
		"password":      {auth_admin_pass},
	}
	resp, err := token_client.PostForm(openIDURL, formData)
	if err != nil {
		fmt.Println(err.Error())

		// logrus.Warn(err.Error())
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	//fmt.Println(result)
	// logrus.Info(result)
	token := result["access_token"].(string)
	// client := gocloak.NewClient("https://ipa.tisafe.jump:8443", gocloak.SetAuthAdminRealms("admin/realms"), gocloak.SetAuthRealms("realms"))
	// ctx := context.Background()
	// token, err := client.LoginClient(ctx, "nss-client", "abcd", "master")
	if err != nil {
		panic("Login failed:" + err.Error())
	}

	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	user_client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", "https://" + FQDN + "/auth/admin/realms/master/users?briefRepresentation=false", nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", os.ExpandEnv("Bearer "+token))
	req.Header.Set("Content-Type", "application/json")
	resp, err = user_client.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	db_creation(dbPath)
	var uInfo []UserInfo
	Uid_contador:=5000
	err = json.NewDecoder(resp.Body).Decode(&uInfo)
	for _, s := range uInfo {
		// Group query end ###############################################################################
		inux_db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer inux_db.Close()

		//u64, err := strconv.ParseUint(s.Attributes.Uids[0], 10, 32)
		if err != nil {
			fmt.Println(err)
		}
		Uid_contador++
		insertUser(inux_db, s.Username, Uid_contador, Uid_contador, "/home/"+s.Username, "/bin/bash")
		insertGroup(inux_db, s.Username, Uid_contador, Groupmember)

	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	row, err := db.Query("SELECT * FROM users ")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var Username string
		var UID uint
		var GID uint
		var Dir string
		var Shell string
		row.Scan(&Username, &UID, &GID, &Dir, &Shell)
		// log.Println("Student: ", Username, " ", UID, " ", GID, " ", Dir, " ", Shell)
		dbtest_passwd = append(dbtest_passwd,
			Passwd{
				Username: Username,
				Password: "x",
				UID:      UID,
				GID:      GID,
				Gecos:    "Jump user",
				Dir:      Dir,
				Shell:    Shell,
			},
		)
		// fmt.Printf("%+v \n", dbtest_passwd)
	}

	row, err = db.Query("SELECT * FROM groups ")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var Groupname string
		var GID uint
		var Members string
		row.Scan(&Groupname, &GID, &Members)
		// fmt.Println(Members)
		split := strings.Split(Members, ",")
		split = delete_empty(split)
		// fmt.Println(split)
		Slack := []string{}
		for i := 0; i < len(split); i++ {
			// ustString := strings.Join(split[i], " ")
			// fmt.Println(len(split))
			// fmt.Println(split[i])
			Slack = append(Slack, split[i])
		}

		dbtest_group = append(dbtest_group,
			Group{
				Groupname: Groupname,
				Password:  "x",
				GID:       GID,
				Members:   Slack,
			},
		)
	}
	// fmt.Printf("%+v \n", dbtest_group)

	dbtest_shadow = append(dbtest_shadow,
		Shadow{
			Username:        "Srtyu",
			Password:        "$6$yZcX.DOY$7bgsJhILMYl3DfMZsYUwoObbVt5Sj9FuujuhVn05Vg9hk.2AXLNy6o1DcPNq0SIyaRZ5YBZer2rYaycuh3qtg1", // Password is "password"
			LastChange:      17920,
			MinChange:       0,
			MaxChange:       99999,
			PasswordWarn:    7,
			InactiveLockout: -1,
			ExpirationDate:  -1,
			Reserved:        -1,
		},
		Shadow{
			Username:        "web",
			Password:        "$6$yZcX.DOY$7bgsJhILMYl3DfMZsYUwoObbVt5Sj9FuujuhVn05Vg9hk.2AXLNy6o1DcPNq0SIyaRZ5YBZer2rYaycuh3qtg1", // Password is "password"
			LastChange:      17920,
			MinChange:       0,
			MaxChange:       99999,
			PasswordWarn:    7,
			InactiveLockout: 0,
			ExpirationDate:  0,
			Reserved:        -1,
		},
	)

	// fmt.Printf("%+v \n", dbtest_passwd)
	// fmt.Printf("%+v \n", dbtest_group)
	// fmt.Printf("%+v \n", dbtest_shadow)

}

func delete_empty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func db_creation(dbPath string) {
	// os.Remove(dbPath)
	// fmt.Println("Creating sqlite-database.db...")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// O arquivo não existe, então podemos criá-lo
		fmt.Println("Creating sqlite-database.db...")
		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatal(err.Error())
		}
	file.Close()
	}
	// fmt.Println("/tmp/sqlite-database.db created")
	linux_db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer linux_db.Close()

	linux_db_drop := `DROP TABLE IF EXISTS users; DROP TABLE IF EXISTS groups;` // SQL Statement for Drop Table
	statementDrop, err := linux_db.Prepare(linux_db_drop) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statementDrop.Exec() // Execute SQL Statements
	// log.Println("tables dropped")	

	linux_db_table := `CREATE TABLE users ( Username string PRIMARY KEY, UID INT, GID INT NOT NULL,Dir string NOT NULL,Shell string NOT NULL);` // SQL Statement for Create Table
	statementUsers, err := linux_db.Prepare(linux_db_table) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statementUsers.Exec() // Execute SQL Statements
	// log.Println("user table created")

	linux_db_group_table := `CREATE TABLE groups ( Groupname string PRIMARY KEY,GID INT NOT NULL,Members string NOT NULL);` // SQL Statement for Create Table
	statementGroups, err = linux_db.Prepare(linux_db_group_table) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statementGroups.Exec() // Execute SQL Statements
	// log.Println("group table created")

}

func insertUser(db *sql.DB, username string,uid int ,gid int, dir string, shell string) {
	// log.Println("Inserting user record ...")
	//insertuser := `INSERT INTO users(username, uid, gid,dir,shell) VALUES (?,?,?,?,?)`
	insertuser := `INSERT INTO users(username,uid,gid,dir,shell) VALUES (?,?,?,?,?)`
	statement, err := db.Prepare(insertuser) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(username, uid, gid, dir, shell)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
func insertGroup(db *sql.DB, groupname string, gid int, members string) {

	// log.Println("Inserting user record ...")
	insertgroup := `INSERT INTO groups(groupname,gid,members) VALUES (?,?,?)`

	statement, err := db.Prepare(insertgroup) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = statement.Exec(groupname, gid, members)
	if err != nil {
		log.Fatalln(err.Error())
	}
}