package protobuf

import (
	"os"
	"log"
	"fmt"
	"context"
	"time"
	"strings"
	
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	
	app "../app"
	userService "./services/user/protobuf"
	categoryService "./services/category/protobuf"
	commentService "./services/comment/protobuf"
	
)

const (
	SRC_AVATAR="/images/author.jpg"
	SRC_COVER="/images/cover.png"
	SRC_POST="/images/house.jpg"
)

//=========================================================
//
//	Функции инициализации gRPC сервера
//
//=========================================================
type Server struct {
	Port string
	DbClient *mongo.Client

	UserService userService.UserServiceClient
	CategoryService categoryService.CategoryServiceClient
	CommentService commentService.CommentServiceClient

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

func (o *Server) UserConnect(){
	
	//Подключение к сервису User
	userServerAdress:=fmt.Sprintf("%s:%s",os.Getenv("USER_HOST"),os.Getenv("USER_PORT"))
	cnnUser, err := grpc.Dial(
		userServerAdress, 
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Printf("[IN] %s.UserConnect: ошибка подключения к сервису User", app.SERVICE_NAME)
	}
	
	o.UserService= userService.NewUserServiceClient(cnnUser)
	_, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
	defer cancel()
	
}

func (o *Server) CategoryConnect(){
	
	//Подключение к сервису User
	categoryServerAdress:=fmt.Sprintf("%s:%s",os.Getenv("CATEGORY_HOST"),os.Getenv("CATEGORY_PORT"))
	cnnCategory, err := grpc.Dial(
		categoryServerAdress, 
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Printf("[IN] %s.CategoryConnect: ошибка подключения к сервису Category", app.SERVICE_NAME)
	}

	o.CategoryService= categoryService.NewCategoryServiceClient(cnnCategory)
	_, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
	defer cancel()
	
}

func (o *Server) CommentConnect(){
	
	//Подключение к сервису User
	commentServerAdress:=fmt.Sprintf("%s:%s",os.Getenv("COMMENT_HOST"),os.Getenv("COMMENT_PORT"))
	cnnComment, err := grpc.Dial(
		commentServerAdress, 
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Printf("[IN] %s.CommentConnect: ошибка подключения к сервису Comment", app.SERVICE_NAME)
	}

	o.CommentService= commentService.NewCommentServiceClient(cnnComment)
	_, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
	defer cancel()
	
}



//----------------------------------------------------------------------------------------------------------------------
// Create
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) Create(ctx context.Context, in *CreatePostRequest) (*CreatePostResponse, error) {

	//Ответ
	out:=&CreatePostResponse{}

	md,_:=metadata.FromIncomingContext(ctx)
	var userId string
	if len(md["user-id"])>0{
		userId=md["user-id"][0]
	}

	//Проверка содержимого запроса перед выполнением
	//Проверка Title
	if in.Title==""{
		return out,app.ErrTitleIsEmpty
	}
	//Проверка SubTitle
	if in.SubTitle==""{
		return out,app.ErrSubTitleIsEmpty
	}
	//Проверка Content
	if in.Content==""{
		return out,app.ErrContentIsEmpty
	}
	
	collection := o.DbClient.Database("blog").Collection("posts")

	post:=&Post{
		Title:in.Title,
		SubTitle:in.SubTitle,
		Content:in.Content,
		Status:app.STATUS_NEW,
		UserId:userId,
		Categories:in.Categories,
	}
	
	insertResult, err := collection.InsertOne(context.TODO(), post)
	if err != nil {
		return out,err
	}
	if oid, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		post.Slug=fmt.Sprintf("%s",oid.Hex())
		
	}else {
		err:=app.ErrInsert
		return out,err
	}

	out.Post=post
	
	return out,nil
}

//----------------------------------------------------------------------------------------------------------------------
// Update
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) Update(ctx context.Context, in *UpdatePostRequest) (*UpdatePostResponse, error) {
	
	//Ответ
	out:=&UpdatePostResponse{}

	//Проверка содержимого запроса перед выполнением
	//Проверка Title
	if in.Title==""{
		return out,app.ErrTitleIsEmpty
	}
	//Проверка SubTitle
	if in.SubTitle==""{
		return out,app.ErrSubTitleIsEmpty
	}
	//Проверка Content
	if in.Content==""{
		return out,app.ErrContentIsEmpty
	}

	collection := o.DbClient.Database("blog").Collection("posts")

	id, err := primitive.ObjectIDFromHex(in.Slug)
	if err != nil {
		return out,err
	}
	filter := bson.M{"_id": id}
	
	update:= bson.M{"$set": bson.M{
		"title": in.Title,
		"subtitle": in.SubTitle,
		"content": in.Content,
		"categories":in.Categories,
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
func (o *Server) Delete(ctx context.Context, in *DeletePostRequest) (*DeletePostResponse, error) {
	
	//Ответ
	out:=&DeletePostResponse{}

	
	collection := o.DbClient.Database("blog").Collection("posts")

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
func (o *Server) Get(ctx context.Context, in *GetPostRequest) (*GetPostResponse, error) {
	//Ответ
	out:=&GetPostResponse{}

	
	collection := o.DbClient.Database("blog").Collection("posts")

	post:=&Post{}

	
	id, err := primitive.ObjectIDFromHex(in.Slug)
	if err != nil {
		return out,err
	}

	filter:= bson.M{"_id": id}

	err= collection.FindOne(context.TODO(), filter).Decode(post)
	if err != nil {
		return out,err
	}
	post.Slug=in.Slug
	post.Src=SRC_POST //TODO - заглушка

	
	//Запрос к сервису User
	var header, trailer metadata.MD
	resp, err := o.UserService.Get(
		getCallContext(ctx),
		&userService.GetUserRequest{Slug:post.UserId},
		grpc.Header(&header), //метадата со стороны сервера в начале запоса
		grpc.Trailer(&trailer), //метадата со стороны сервера в коне запоса
	)


	if err != nil {
		return out,err
	}
	
	author:=&Author{
		Slug:resp.User.Slug,
		FirstName:resp.User.FirstName,
		LastName:resp.User.LastName,
		SrcAvatar:SRC_AVATAR, //TODO - заглушка
		SrcCover:SRC_COVER,   //TODO - заглушка
	}
	post.Author=author


	//Запрос к сервису Category, JOIN category
	respCategory,err:= o.CategoryService.Find(getCallContext(ctx), &categoryService.FindCategoryRequest{})
	if err != nil {
		return out,err
	}
	for _, category:= range respCategory.Categories {
		for _, category_slug:= range strings.Split(post.Categories,",") {
			if category.Slug==category_slug{
				postCategor:=&PostCategory{
					Slug:category.Slug,
					Name:category.Name,
				}
				post.PostCategories=append(post.PostCategories,postCategor)
			}
		}
	}

	//Запрос к сервису Comments, JOIN comments
	respComment,err:= o.CommentService.Find(getCallContext(ctx), &commentService.FindCommentRequest{PostId:in.Slug})
	if err != nil {
		return out,err
	}
	for _, comment:= range respComment.Comments {
		postComment:=&PostComment{
			Slug:comment.Slug,
			Content:comment.Content,
		}
		post.PostComments=append(post.PostComments,postComment)
	}
	
	out.Post=post	

	//fmt.Printf("Found a single document: %+v\n", resp)

	return out,nil
}

//----------------------------------------------------------------------------------------------------------------------
// Find - список всех постов
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) Find(ctx context.Context, in *FindPostRequest) (*FindPostResponse, error) {

	out:=&FindPostResponse{}
	
	collection := o.DbClient.Database("blog").Collection("posts")

	options:= options.Find()
	//options.SetLimit(2)
	
	
	//filter := bson.D{{"title", "1111"}}
	filter := bson.D{}
	var results []*Post
	
	
	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return out,err
	}

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		post:= &Post{}
		
		
		err := cur.Decode(post)
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
		//fmt.Printf("raw: %v \n",item)
		post.Slug=fmt.Sprintf("%s",item.ID.Hex())
		post.Src=SRC_POST //TODO - заглушка

		results = append(results, post)
	}
	
	if err:= cur.Err(); err != nil {
		return out,err
	}
	
	// Close the cursor once finished
	cur.Close(context.TODO())

	out.Posts=results
	return out,nil
}

//----------------------------------------------------------------------------------------------------------------------
// GET AUTHOR & related Posts
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) GetAuthor(ctx context.Context, in *GetAuthorRequest) (*GetAuthorResponse, error) {
	out:=&GetAuthorResponse{}

	//Запрос к сервису User
	resp, err := o.UserService.Get(getCallContext(ctx), &userService.GetUserRequest{Slug:in.Slug})
	if err != nil {
		return out,err
	}
	
	out.Author=&Author{
		Slug:resp.User.Slug,
		FirstName:resp.User.FirstName,
		LastName:resp.User.LastName,
		SrcAvatar:SRC_AVATAR,  //TODO - заглушка
		SrcCover:SRC_COVER,    //TODO - заглушка
	}

	var posts []*Post

	collection := o.DbClient.Database("blog").Collection("posts")
	options:= options.Find()
	filter := bson.D{{"userid",in.Slug}}
	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return out,err
	}
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		post:= &Post{}
		err := cur.Decode(post)
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
		post.Slug=fmt.Sprintf("%s",item.ID.Hex())
		post.Src=SRC_POST //TODO - заглушка
		posts = append(posts, post)
	}
	if err:= cur.Err(); err != nil {
		return out,err
	}
	// Close the cursor once finished
	cur.Close(context.TODO())

	out.Author.Posts=posts
	return out,nil
}

//----------------------------------------------------------------------------------------------------------------------
// FIND AUTHORS
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) FindAuthors(ctx context.Context, in *FindAuthorRequest) (*FindAuthorResponse, error) {
	out:=&FindAuthorResponse{}

	var authors []*Author
	
	//Запрос к сервису User
	resp, err := o.UserService.Find(getCallContext(ctx), &userService.FindUserRequest{})
	if err != nil {
		return out,err
	}
	for _, user:= range resp.Users {
		author:=&Author{
			Slug:user.Slug,
			FirstName:user.FirstName,
			LastName:user.LastName,
			SrcAvatar:SRC_AVATAR,  //TODO - заглушка
			SrcCover:SRC_COVER,    //TODO - заглушка
		}

		authors=append(authors,author)
	}
	out.Authors=authors
	return out,nil
}

//----------------------------------------------------------------------------------------------------------------------
// GET PostCategory & related Posts
//----------------------------------------------------------------------------------------------------------------------
func (o *Server) GetPostCategory(ctx context.Context, in *GetPostCategoryRequest) (*GetPostCategoryResponse, error) {
	out:=&GetPostCategoryResponse{}

	//Запрос к сервису Category
	resp, err := o.CategoryService.Get(getCallContext(ctx), &categoryService.GetCategoryRequest{Slug:in.Slug})
	if err != nil {
		return out,err
	}
	
	out.Category=&PostCategory{
		Name:resp.Category.Name,
	}

	var posts []*Post

	collection := o.DbClient.Database("blog").Collection("posts")
	options:= options.Find()
	//filter := bson.D{{"categories",in.Slug}}
	query:=fmt.Sprintf(".*%s.*",in.Slug)
	//filter:= bson.M{{"categories" : {$regex : query}}
	filter := bson.D{{"categories", primitive.Regex{Pattern: query, Options: ""}}}
	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return out,err
	}
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		post:= &Post{}
		err := cur.Decode(post)
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
		post.Slug=fmt.Sprintf("%s",item.ID.Hex())
		post.Src=SRC_POST //TODO - заглушка
		posts = append(posts, post)
	}
	if err:= cur.Err(); err != nil {
		return out,err
	}
	// Close the cursor once finished
	cur.Close(context.TODO())

	out.Category.Posts=posts
	return out,nil
}

//----------------------------------------------------------------------------------------------------------------------
// Формироание контекста для вызова связанных сервисов
//----------------------------------------------------------------------------------------------------------------------
func getCallContext(ctx context.Context) context.Context{
	callContext:=context.Background()
	mdIn,_:=metadata.FromIncomingContext(ctx)
	var traceId,userId,userRole string
	if len(mdIn["trace-id"])>0{
		traceId=mdIn["trace-id"][0]
	}
	if len(mdIn["user-id"])>0{
		userId=mdIn["user-id"][0]
	}
	if len(mdIn["user-role"])>0{
		userRole=mdIn["user-role"][0]
	}

	mdOut:=metadata.Pairs(
		"trace-id",traceId,
		"user-id",userId,
		"user-role",userRole,
	)
	callContext=metadata.NewOutgoingContext(callContext,mdOut)
	return callContext
}



