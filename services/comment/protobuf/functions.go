package protobuf

import (
	"os"
	"log"
	"fmt"
	"context"
	"time"
	
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"

	
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	
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

//----------------------------------------------------------------------------------------------------------------------
// Create
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) Create(ctx context.Context, in *CreateCommentRequest) (*CreateCommentResponse, error) {

	//Ответ
	out:=&CreateCommentResponse{}
	
	
	// md,_:=metadata.FromIncomingContext(ctx)
	// var userId string
	// if len(md["user-id"])>0{
	// 	userId=md["user-id"][0]
	// }
	

	//Проверка содержимого запроса перед выполнением
	//Проверка Content
	if in.Content==""{
		return out,app.ErrContentIsEmpty
	}

	collection := o.DbClient.Database(app.DB_NAME).Collection(app.DB_COLLECTION)
	
	comment:=&Comment{
		//ParentId:in.ParentId,
		//UserId:userId,
		PostId:in.PostId,
		Content:in.Content,
		Status:app.STATUS_PUBLISHED,
	}
	
	insertResult, err := collection.InsertOne(context.TODO(), comment)
	if err != nil {
		return out,err
	}
	
	if oid, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		comment.Slug=fmt.Sprintf("%s",oid.Hex())
	}else {
		err:=app.ErrInsert
		
		return out,err
	}

	out.Comment=comment
	
	return out,nil
}

//----------------------------------------------------------------------------------------------------------------------
// Update
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) Update(ctx context.Context, in *UpdateCommentRequest) (*UpdateCommentResponse, error) {
	
	//Ответ
	out:=&UpdateCommentResponse{}

	//Проверка содержимого запроса перед выполнением
	//Проверка Title
	if in.Content==""{
		return out,app.ErrContentIsEmpty
	}
	
	collection := o.DbClient.Database(app.DB_NAME).Collection(app.DB_COLLECTION)

	id, err := primitive.ObjectIDFromHex(in.Slug)
	if err != nil {
		return out,err
	}
	filter := bson.M{"_id": id}
	
	update:= bson.M{"$set": bson.M{
		"Content": in.Content,
		"Status": in.Status,
	}}

	_, err= collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return out,err
	}

	
	out.Status=app.ACTION_STATUS_SUCCESS
	return out,nil
}

//----------------------------------------------------------------------------------------------------------------------
// Delete
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) Delete(ctx context.Context, in *DeleteCommentRequest) (*DeleteCommentResponse, error) {
	
	//Ответ
	out:=&DeleteCommentResponse{}

	
	collection := o.DbClient.Database(app.DB_NAME).Collection(app.DB_COLLECTION)

	id, err := primitive.ObjectIDFromHex(in.Slug)
	if err != nil {
		return out,err
	}

	filter:= bson.M{"_id": id}
	_, err= collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		return out,err
	}

	
	out.Status=app.ACTION_STATUS_SUCCESS
	return out,nil
}

//----------------------------------------------------------------------------------------------------------------------
// Get
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) Get(ctx context.Context, in *GetCommentRequest) (*GetCommentResponse, error) {
	//Ответ
	out:=&GetCommentResponse{}

	collection := o.DbClient.Database(app.DB_NAME).Collection(app.DB_COLLECTION)

	comment:=&Comment{}
	
	id, err := primitive.ObjectIDFromHex(in.Slug)
	if err != nil {
		return out,err
	}

	filter:= bson.M{"_id": id}

	err= collection.FindOne(context.TODO(), filter).Decode(comment)
	if err != nil {
		return out,err
	}
	
	
	out.Comment=comment	
	return out,nil
}

//----------------------------------------------------------------------------------------------------------------------
// Find
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) Find(ctx context.Context, in *FindCommentRequest) (*FindCommentResponse, error) {

	out:=&FindCommentResponse{}
	
	collection := o.DbClient.Database(app.DB_NAME).Collection(app.DB_COLLECTION)

	options:= options.Find()
	//options.SetLimit(2)
	
	
	filter := bson.D{{"postid", in.PostId}}
	//filter := bson.D{}
	var results []*Comment
	
	
	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return out,err
	}

	for cur.Next(context.TODO()) {

		
		comment:= &Comment{}
		
		err := cur.Decode(comment)
		if err != nil {
			return out,err
		}
		
		type Item struct {
			ID      primitive.ObjectID  `json:"_id" bson:"_id"`
		}
		item:=&Item{}

		err= cur.Decode(item)
		if err != nil {
			return out,err
		}
		
		comment.Slug=fmt.Sprintf("%s",item.ID.Hex())
		results = append(results, comment)
	}
	
	if err:= cur.Err(); err != nil {
		return out,err
	}
	
	// Close the cursor once finished
	cur.Close(context.TODO())

	out.Comments=results
	return out,nil
}
