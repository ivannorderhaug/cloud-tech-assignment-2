package db

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
)

var ctx context.Context
var client *firestore.Client

// InitializeFirestore Method for initializing the firestore client */
func InitializeFirestore() {
	// Firebase initialisation
	ctx = context.Background()

	sa := option.WithCredentialsFile("./service-account.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

// AddToFirestore Simple method to add data to the firestore database.
func AddToFirestore(collectionName string, documentID string, data interface{}) error {
	_, err := client.Collection(collectionName).Doc(documentID).Set(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

// CloseFirestore Method to gracefully close the firestore client
func CloseFirestore() {
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(client)
}

// GetAllDocumentsFromFirestore Method to retrieve documents from a collection by using an iterator
func GetAllDocumentsFromFirestore(collectionName string) ([]*firestore.DocumentSnapshot, error) {
	iterator := client.Collection(collectionName).Documents(ctx)
	documents, err := iterator.GetAll()
	if err != nil {
		return nil, err
	}

	//Returns a slice of document snapshots
	return documents, nil
}

// GetSingleDocumentFromFirestore Method to retrieve a single document from a collection by using the document id
func GetSingleDocumentFromFirestore(collectionName string, documentID string) (*firestore.DocumentSnapshot, error) {
	document, err := client.Collection(collectionName).Doc(documentID).Get(ctx)
	if err != nil {
		return nil, err
	}
	//Returns a single document snapshot
	return document, nil
}

// DeleteSingleDocumentFromFirestore Method to delete a single document from a collection by using the doc id
func DeleteSingleDocumentFromFirestore(collectionName string, documentID string) error {
	_, err := client.Collection(collectionName).Doc(documentID).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

// UpdateWebhook */
func UpdateWebhook(collectionName string, documentID string, path string, newValue interface{}) error {
	_, err := client.Collection(collectionName).Doc(documentID).Update(ctx, []firestore.Update{
		{
			Path:  path,
			Value: newValue,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
