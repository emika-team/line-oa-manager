package firebase

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/joho/godotenv"
)

var App *firebase.App

var FirestoreClient *firestore.Client

func init() {
	var err error
	godotenv.Load()

	App, err = firebase.NewApp(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	FirestoreClient, err = App.Firestore(context.Background())
	if err != nil {
		panic(err)
	}
}
