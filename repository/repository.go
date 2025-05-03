package repository

import "main/database"

var Client database.Client

func GetClient() database.Client {
	return Client
}

func SetClient(client database.Client) {
	Client = client
}
