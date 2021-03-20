package global

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"github.com/dgraph-io/ristretto"
	"google.golang.org/api/option"
)

var (
	UserCache    *ristretto.Cache // user id -> account
	UserApiCache *ristretto.Cache // api key -> user
	FirebaseApp  *firebase.App
)

func Init() {
	var err error

	if UserCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M)
		MaxCost:     1 << 30, // max cost of cache (1 GB)
		BufferItems: 64,      // number of keys per Get buffer
	}); err != nil {
		fmt.Println("unable to initialize user cache", err.Error())
		panic(err)
	}

	if UserApiCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M)
		MaxCost:     1 << 30, // max cost of cache (1 GB)
		BufferItems: 64,      // number of keys per Get buffer
	}); err != nil {
		fmt.Println("unable to initialize user api cache", err.Error())
		panic(err)
	}

	opt := option.WithCredentialsFile("sample-domain-firebase.json")
	FirebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Println("error occurred while initiating firebase app")
		panic(err)
	}
}
