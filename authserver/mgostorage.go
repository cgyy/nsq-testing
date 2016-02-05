package authserver

import (
	"log"
	"time"

	"github.com/RangelReale/osin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// collection names for the entities
const (
	CLIENT_COL    = "clients"
	AUTHORIZE_COL = "authorizations"
	ACCESS_COL    = "accesses"
)

const REFRESHTOKEN = "refreshtoken"

//type OAuthMongoStorage struct{}

//func NewOAuthMongoStorage() *OAuthMongoStorage {
//	return &OAuthMongoStorage{}
//}

type AuthorizeData struct {
	ClientId    string
	Code        string
	ExpiresIn   int32
	Scope       string
	RedirectUri string
	State       string
	CreatedAt   time.Time
	UserData    interface{}
}

func AuthorizeDataFromOSIN(data *osin.AuthorizeData) *AuthorizeData {
	return &AuthorizeData{
		ClientId:    data.Client.GetId(),
		Code:        data.Code,
		ExpiresIn:   data.ExpiresIn,
		Scope:       data.Scope,
		RedirectUri: data.RedirectUri,
		State:       data.State,
		CreatedAt:   data.CreatedAt,
		UserData:    data.UserData,
	}
}

func (data *AuthorizeData) AuthorizeDataToOSIN(s *MongoStorage) (*osin.AuthorizeData, error) {
	client, err := s.GetClient(data.ClientId)
	if err != nil {
		return nil, err
	}
	return &osin.AuthorizeData{
		Client:      client,
		Code:        data.Code,
		ExpiresIn:   data.ExpiresIn,
		Scope:       data.Scope,
		RedirectUri: data.RedirectUri,
		State:       data.State,
		CreatedAt:   data.CreatedAt,
		UserData:    data.UserData,
	}, nil
}

type AccessData struct {
	ClientId string
	// current authorizeData, Get from Code
	//Code string
	// previous AccessData, Get from Prev AccessToken
	//PrevAccessToken string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int32
	Scope        string
	RedirectUri  string
	CreatedAt    time.Time
	UserData     interface{}
}

func AccessDataFromOSIN(data *osin.AccessData) *AccessData {
	/*
		var prev_access_token string
		if data.AccessData == nil {
			prev_access_token = ""
		} else {
			prev_access_token = data.AccessData.AccessToken
		}
	*/
	return &AccessData{
		ClientId: data.Client.GetId(),
		//Code:            data.AuthorizeData.Code,
		//PrevAccessToken: prev_access_token,
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		ExpiresIn:    data.ExpiresIn,
		Scope:        data.Scope,
		RedirectUri:  data.RedirectUri,
		CreatedAt:    data.CreatedAt,
		UserData:     data.UserData,
	}
}

func (data *AccessData) AccessDataToOSIN(s *MongoStorage) (*osin.AccessData, error) {
	log.Println("AccessDataToOSIN come?")
	client, err := s.GetClient(data.ClientId)
	if err != nil {
		log.Println("Get Client Error")
		return nil, err
	}
	/*
		var authorize_data *osin.AuthorizeData
		log.Println("data.Code =", data.Code)
		if data.Code == "" {
			authorize_data = nil
		} else {
			authorize_data, err = s.LoadAuthorize(data.Code)
			if err != nil {
				log.Println("loading authorize_data error ==> not error just skip")
				//return nil, err
			}
		}
		var access_data *osin.AccessData
		if data.PrevAccessToken == "" {
			access_data = nil
		} else {
			access_data, err = s.LoadAccess(data.PrevAccessToken)
			if err != nil {
				log.Println("loading access_data error")
				return nil, err
			}
		}
	*/
	return &osin.AccessData{
		Client: client,
		//AuthorizeData: authorize_data,
		//AccessData:    access_data,
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		ExpiresIn:    data.ExpiresIn,
		Scope:        data.Scope,
		RedirectUri:  data.RedirectUri,
		CreatedAt:    data.CreatedAt,
		UserData:     data.UserData,
	}, nil
}

type MongoStorage struct {
	DbName  string
	Session *mgo.Session
}

func NewMongoStorage(session *mgo.Session, dbName string) *MongoStorage {
	storage := &MongoStorage{dbName, session}
	index := mgo.Index{
		Key:        []string{REFRESHTOKEN},
		Unique:     false, // refreshtoken is sometimes empty
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}
	accesses := storage.Session.DB(dbName).C(ACCESS_COL)
	err := accesses.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	return storage
}

func (store *MongoStorage) GetClient(id string) (osin.Client, error) {
	session := store.Session.Copy()
	defer session.Close()
	clients := session.DB(store.DbName).C(CLIENT_COL)
	client := new(osin.DefaultClient)
	err := clients.FindId(id).One(client)
	return client, err
}

func (store *MongoStorage) SetClient(id string, client osin.Client) error {
	session := store.Session.Copy()
	defer session.Close()
	clients := session.DB(store.DbName).C(CLIENT_COL)
	_, err := clients.UpsertId(id, client)
	return err
}

func (store *MongoStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	log.Println("mongo SaveAuthorize - come")
	session := store.Session.Copy()
	defer session.Close()
	authorizations := session.DB(store.DbName).C(AUTHORIZE_COL)
	log.Println("SaveAuthorize : ", data.Code)
	_, err := authorizations.UpsertId(data.Code, AuthorizeDataFromOSIN(data))
	return err
}

func (store *MongoStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	log.Println("mongo LoadAuthorize - come")
	session := store.Session.Copy()
	defer session.Close()
	authorizations := session.DB(store.DbName).C(AUTHORIZE_COL)
	authData := new(AuthorizeData)
	err := authorizations.FindId(code).One(authData)
	if err != nil {
		return nil, err
	}
	return authData.AuthorizeDataToOSIN(store)
}

func (store *MongoStorage) RemoveAuthorize(code string) error {
	log.Println("mongo RemoveAuthorize - come")
	session := store.Session.Copy()
	defer session.Close()
	authorizations := session.DB(store.DbName).C(AUTHORIZE_COL)
	return authorizations.RemoveId(code)
}

func (store *MongoStorage) SaveAccess(data *osin.AccessData) error {
	session := store.Session.Copy()
	defer session.Close()
	accesses := session.DB(store.DbName).C(ACCESS_COL)
	_, err := accesses.UpsertId(data.AccessToken, AccessDataFromOSIN(data))
	return err
}

func (store *MongoStorage) LoadAccess(token string) (*osin.AccessData, error) {
	log.Println("mongo LoadAccess - come")
	session := store.Session.Copy()
	defer session.Close()
	accesses := session.DB(store.DbName).C(ACCESS_COL)
	accData := new(AccessData)
	//log.Println("Token = ", token)
	err := accesses.FindId(token).One(accData)
	if err != nil {
		log.Println("find id error")
		return nil, err
	}
	return accData.AccessDataToOSIN(store)
}

func (store *MongoStorage) RemoveAccess(token string) error {
	log.Println("mongo RemoveAccess - come")
	session := store.Session.Copy()
	defer session.Close()
	accesses := session.DB(store.DbName).C(ACCESS_COL)
	return accesses.RemoveId(token)
}

func (store *MongoStorage) LoadRefresh(token string) (*osin.AccessData, error) {
	log.Println("mongo LoadRefresh - come")
	session := store.Session.Copy()
	defer session.Close()
	accesses := session.DB(store.DbName).C(ACCESS_COL)
	accData := new(AccessData)
	err := accesses.Find(bson.M{REFRESHTOKEN: token}).One(accData)
	if err != nil {
		return nil, err
	}
	return accData.AccessDataToOSIN(store)
}

func (store *MongoStorage) RemoveRefresh(token string) error {
	log.Println("mongo RemoveRefresh - come")
	session := store.Session.Copy()
	defer session.Close()
	accesses := session.DB(store.DbName).C(ACCESS_COL)
	return accesses.Update(bson.M{REFRESHTOKEN: token}, bson.M{
		"$unset": bson.M{
			REFRESHTOKEN: 1,
		}})
}

func (store *MongoStorage) Clone() osin.Storage {
	return store
}

// implementation ?
func (store *MongoStorage) Close() {

}
