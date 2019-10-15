package protobuf

import (
	"os"
	"log"
	"fmt"
	"time"
	"context"
	
	"net/mail"
	"crypto/md5"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"

	app "../app"
)

//=========================================================
//
//	Функции инициализации gRPC сервера
//
//=========================================================
type Server struct {
	Port string
	DbClient *mongo.Client
}

type MongoItem struct {
	ID      primitive.ObjectID  `json:"_id" bson:"_id"`
}

//----------------------------------------------------------------------------------------------------------------------
// Midelware
//----------------------------------------------------------------------------------------------------------------------
func AccessLogInterceptor(ctx context.Context,req interface{},info *grpc.UnaryServerInfo,handler grpc.UnaryHandler,) (interface{}, error) {

	start:=time.Now()
	md,_:=metadata.FromIncomingContext(ctx)

	// Calls the handler
	reply, err := handler(ctx, req)

	var traceId,userId,userRole string
	if len(md["trace-id"])>0{
		traceId=md["trace-id"][0]
	}
	if len(md["user-id"])>0{
		userId=md["user-id"][0]
	}
	if len(md["user-role"])>0{
		userRole=md["user-role"][0]
	}
	

	msg:=fmt.Sprintf("Call:%v, traceId: %v, userId: %v, userRole: %v, time: %v", info.FullMethod,traceId,userId,userRole,time.Since(start))
	app.AccesLog(msg)

	return reply, err
	
}

//----------------------------------------------------------------------------------------------------------------------
// SignUp
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) SignUp(ctx context.Context, in *SignUpRequest) (*SignUpResponse, error) {
	
	//Ответ, по умолчанию STATUS_FAIL
	out:=&SignUpResponse{}

	
	//Проверка содержимого запроса перед выполнением
	//Проверка Username
	err:=checkUserName(in.Username)
	if err!=nil{
	 	log.Printf("[ERR] %s.SignUp, %v", app.SERVICE_NAME,err)
	 	return out,err
	}

	//Проверка Username на дубль
	err=o.checkUserNameExist(in.Username)
	if err!=nil{
		log.Printf("[ERR] %s.SignUp, %v", app.SERVICE_NAME,err)
		return out,err
	}
	
	//Проверка Password
	err=checkPassword(in.Password)
	if err!=nil{
		log.Printf("[ERR] %s.SignUp, %v", app.SERVICE_NAME,err)
		return out,err
	}
	
	
	user:=&User{
		Username:in.Username,
		FirstName:in.FirstName,
		LastName:in.LastName,
		Password:getMD5(in.Password),
	}

	var slug string
	collection:= o.DbClient.Database("blog").Collection("users")
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return out,err
	}
	if oid, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		slug=fmt.Sprintf("%s",oid.Hex())
	}else {
		err:=app.ErrInsert
		return out,err
	}

	out.Slug=slug
	out.Username=in.Username
	out.Role=app.ROLE_USER
	return out,nil
}

//----------------------------------------------------------------------------------------------------------------------
// SignIn
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) SignIn(ctx context.Context, in *SignInRequest) (*SignInResponse, error) {

	//Ответ, по умолчанию STATUS_FAIL
	out:=&SignInResponse{}

	collection := o.DbClient.Database("blog").Collection("users")
	filter := bson.D{{"username",in.Username}}
	user:=&User{}
	result:=collection.FindOne(context.TODO(), filter)
	
	//Получаю ID записи
	item:=&MongoItem{}
	err:= result.Decode(item)
	if err != nil {
		err=app.ErrIncorrectLoginOrPassword
		return out,err
	}
	slug:=fmt.Sprintf("%s",item.ID.Hex())
	
	//Получаю атрибуты записи
	err= result.Decode(user)
	if err != nil {
		err=app.ErrIncorrectLoginOrPassword
		return out,err
	}

	//Проверяю пароль
	if user.Password!=getMD5(in.Password){
	 	err=app.ErrIncorrectLoginOrPassword
	 	log.Printf("[ERR] %s.SignIn, некорректный пароль", app.SERVICE_NAME)
	 	return out, err
	}

	out.Slug=slug
	out.Username=in.Username
	out.Role=app.ROLE_USER

	return out,nil
}

//----------------------------------------------------------------------------------------------------------------------
// LogOut
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) LogOut(ctx context.Context, in *LogOutRequest) (*LogOutResponse, error) {

	//Ответ, по умолчанию STATUS_FAIL
	out:=&LogOutResponse{
		Status:app.STATUS_FAIL,
	}
	
	 md, ok := metadata.FromIncomingContext(ctx)
	 if !ok {
	 	fmt.Println("md: %v",md)
	 }

	header:= metadata.Pairs("authorization", "")
	grpc.SendHeader(ctx, header)
	
	out.Status = app.STATUS_SUCCESS
	return out,nil
}

//----------------------------------------------------------------------------------------------------------------------
// Create
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) Create(ctx context.Context, in *CreateUserRequest) (*CreateUserResponse, error) {
	
	//Ответ
	out:=&CreateUserResponse{}

	
	// //Проверка содержимого запроса перед выполнением
	//Проверка Username
	err:=checkUserName(in.Username)
	if err!=nil{
	 	log.Printf("[ERR] %s.SignUp, %v", app.SERVICE_NAME,err)
	 	return out,err
	}

	//Проверка Username на дубль
	err=o.checkUserNameExist(in.Username)
	if err!=nil{
		log.Printf("[ERR] %s.SignUp, %v", app.SERVICE_NAME,err)
		return out,err
	}
	
	//Проверка Password
	err=checkPassword(in.Password)
	if err!=nil{
		log.Printf("[ERR] %s.SignUp, %v", app.SERVICE_NAME,err)
		return out,err
	}

	//Проверка Email
	err=checkEmail(in.Email)
	if err!=nil{
		log.Printf("[ERR] %s.Create, %v", app.SERVICE_NAME,err)
		return out,err
	}
	
	user:=&User{
		Username:in.Username,
		Email:in.Email,
		Password:getMD5(in.Password),
		FirstName:in.FirstName,
		LastName:in.LastName,
		Phone:in.Phone,
	}

	var slug string
	collection := o.DbClient.Database("blog").Collection("users")
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return out,err
	}
	if oid, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		slug=fmt.Sprintf("%s",oid.Hex())
	}else {
		err:=app.ErrInsert
		return out,err
	}
	user.Slug=slug
	out.User=user
	return out,nil
}

//----------------------------------------------------------------------------------------------------------------------
// Update
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) Update(ctx context.Context, in *UpdateUserRequest) (*UpdateUserResponse, error) {
	
	//Ответ
	out:=&UpdateUserResponse{}

	//Проверка содержимого запроса перед выполнением
	//Проверка Email
	err:=checkEmail(in.Email)
	if err!=nil{
		log.Printf("[ERR] %s.Update, %v", app.SERVICE_NAME,err)
		return out,err
	}
	
	id, err := primitive.ObjectIDFromHex(in.Slug)
	if err != nil {
		log.Printf("[ERR] %s.Update: %v", app.SERVICE_NAME, err)
		err=app.ErrNotFound
		return out,err
	}
	filter:= bson.M{"_id": id}
	
	user:=&User{}
	collection:= o.DbClient.Database("blog").Collection("users")
	err= collection.FindOne(context.TODO(), filter).Decode(user)
	if err!= nil {
		log.Printf("[ERR] %s.Update: %v", app.SERVICE_NAME, err)
		err=app.ErrNotFound
		return out,err
	}

	password:=user.Password
	if in.Password!="" {
		password=getMD5(in.Password)
	}

	
	update:= bson.M{"$set": bson.M{
		"email": in.Email,
		"password": password,
		"firstname": in.FirstName,
		"lastname":in.LastName,
		"phone":in.Phone,
	}}

	_, err= collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("[ERR] %s.Update: %v", app.SERVICE_NAME, err)
		err=app.ErrUpdate
		return out,err
	}

	out.Status=app.STATUS_SUCCESS
	return out,nil
}

//----------------------------------------------------------------------------------------------------------------------
// Delete
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) Delete(ctx context.Context, in *DeleteUserRequest) (*DeleteUserResponse, error) {
	
	//Ответ
	out:=&DeleteUserResponse{}

	
	collection := o.DbClient.Database("blog").Collection("users")
	id, err := primitive.ObjectIDFromHex(in.Slug)
	if err != nil {
		return out,err
	}

	filter:= bson.M{"_id": id}
	_, err= collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		log.Printf("[ERR] %s.Delete, %v", app.SERVICE_NAME,err)
		return out,err
	}

	
	out.Status=app.STATUS_SUCCESS
	return out,nil
}

//----------------------------------------------------------------------------------------------------------------------
// Get
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) Get(ctx context.Context, in *GetUserRequest) (*GetUserResponse, error) {
	
	
	out:=&GetUserResponse{}

	id, err := primitive.ObjectIDFromHex(in.Slug)
	if err != nil {
		log.Printf("[ERR] %s.Get, %v", app.SERVICE_NAME,err)
		err=app.ErrNotFound
		return out,err
	}
	filter:= bson.M{"_id": id}
	user:=&User{}
	collection := o.DbClient.Database("blog").Collection("users")
	err= collection.FindOne(context.TODO(), filter).Decode(user)
	if err!= nil {
		log.Printf("[ERR] %s.Get, %v", app.SERVICE_NAME,err)
		err=app.ErrNotFound
		return out,err
	}
	user.Slug=in.Slug
	out.User=user
	return out,nil
}


//----------------------------------------------------------------------------------------------------------------------
// Find
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) Find(ctx context.Context, in *FindUserRequest) (*FindUserResponse, error) {
	
	out:=&FindUserResponse{}
	collection := o.DbClient.Database("blog").Collection("users")
	options:= options.Find()
	//options.SetLimit(2)
	
	filter := bson.D{}
	var users []*User
	
	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return out,err
	}

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		user:= &User{}
		
		
		err := cur.Decode(user)
		if err != nil {
			return out,err
		}
		
		item:=&MongoItem{}

		err= cur.Decode(item)
		if err != nil {
			return out,err
		}
		//fmt.Printf("raw: %v \n",item)
		user.Slug=fmt.Sprintf("%s",item.ID.Hex())
		//user.Src=SRC_POST //TODO - заглушка

		users = append(users, user)
	}
	
	if err:= cur.Err(); err != nil {
		return out,err
	}
	
	// Close the cursor once finished
	cur.Close(context.TODO())

	out.Users=users
	return out,nil
}

//----------------------------------------------------------------------------------------------------------------------
// Подключение/Отключение к БД
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) DbConnect() error {
	var client *mongo.Client
	
	// Create client
	strURI:=fmt.Sprintf("mongodb://%s:%s@%s:%s",os.Getenv("MONGO_USER"),os.Getenv("MONGO_PASS"),os.Getenv("MONGO_HOST"),os.Getenv("MONGO_PORT"))
	client, err:= mongo.NewClient(options.Client().ApplyURI(strURI))
	if err != nil {
		return err
	}

	// Create connect
	err = client.Connect(context.TODO())
	if err != nil {
		return err
	}
	o.DbClient=client
	return nil
}
func (o *Server) DbDisconnect() {
	err:= o.DbClient.Disconnect(context.TODO())

	if err != nil {
		log.Printf("Ошибка завершения соединения БД: %v",err)
	}
}

func getMD5(value string) string {
	slugData:= []byte(value)
    return fmt.Sprintf("%x",md5.Sum(slugData))
}



func checkUserName(username string) error{
	if username==""{
		return app.ErrUserNameIsEmpty
	}
	return nil
}
func (o *Server) checkUserNameExist(username string) error{
	
	collection := o.DbClient.Database("blog").Collection("users")
	
	filter := bson.D{{"username",username}}
	user:=&User{}
	err:= collection.FindOne(context.TODO(), filter).Decode(user)

	if err!=nil && err!=mongo.ErrNoDocuments{
		return err
	}else if user.Username==username{
		return app.ErrUserNameIsExist
	}
	return nil
}

func checkEmail(email string) error{
	if email!=""{
		//Проверка E-mail формата
		_, err := mail.ParseAddress(email)
		if err!=nil{
			return app.ErrEmailIncorrect
		}
	}
	return nil
}

func checkPassword(password string) error{
	if password==""{
		return app.ErrPasswordIsEmpty
	}
	return nil
}