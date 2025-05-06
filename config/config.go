package config

type Config struct {
	MONGO_URI             string
	MOVIE_COLLECTION_NAME string
	USER_COLLECTION_NAME  string
	DB_NAME               string
	JWT_SECRET            string
}

var MyConfig *Config
