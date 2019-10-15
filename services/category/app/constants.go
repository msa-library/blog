package app

const (
	//Service constants
	SERVICE_NAME	= "Category"
	
	//Таймаут подключения к сервисам GRPC (в милисекундах)
	CNN_TIMEOUT		=1000

	//
	STATUS_NEW			=1
	STATUS_PUBLISHED	=2
	STATUS_UNPUBLISHED	=3

	ACTION_STATUS_SUCCESS	=1
	ACTION_STATUS_FAIL		=2

	DB_NAME					="blog"
	DB_COLLECTION			="category"

)