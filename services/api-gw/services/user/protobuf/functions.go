package protobuf

import (
	"net/http"
	"fmt"
	"time"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"github.com/dgrijalva/jwt-go"
)

const JWT_SECRET	="eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9"

func forwardSignIn(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, resp proto.Message, opts ...func(context.Context, http.ResponseWriter, proto.Message) error) {
	
	//Преобразую proto.Message в SignInResponse
	signInResponse:=&SignInResponse{}
	signInResponse.XXX_Merge(resp)

	token,err:=GetJWTToken(signInResponse.Slug,signInResponse.Role)
	if err!=nil{
		http.Error(w, fmt.Sprintf("%v",err), http.StatusUnauthorized)
		return
	}
	w.Header().Set("authorization", token)
	runtime.ForwardResponseMessage(ctx, mux, marshaler, w, req, resp, opts...)
}

func forwardSignUp(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, resp proto.Message, opts ...func(context.Context, http.ResponseWriter, proto.Message) error) {
	
	//Преобразую proto.Message в SignUpResponse
	signUpResponse:=&SignUpResponse{}
	signUpResponse.XXX_Merge(resp)

	token,err:=GetJWTToken(signUpResponse.Slug,signUpResponse.Role)
	if err!=nil{
		http.Error(w, fmt.Sprintf("%v",err), http.StatusUnauthorized)
		return
	}
	w.Header().Set("authorization", token)
	runtime.ForwardResponseMessage(ctx, mux, marshaler, w, req, resp, opts...)
}


func init() {
	//Переопределяю обработку ответа для вызовова SignIn
	forward_UserService_SignIn_0 = forwardSignIn

	//Переопределяю обработку ответа для вызовова SignUp
	forward_UserService_SignUp_0 = forwardSignUp
}

func GetJWTToken(userid string, role string) (string,error) {
    // Создаем новый токен

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "UserID": userid,
        "UserRole": role,
        "nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
    })

    // Подписываем токен нашим секретным ключем
    tokenString, err := token.SignedString([]byte(JWT_SECRET))

    // Отдаем токен клиенту
    return tokenString, err
}

func CheckGetJWTToken(tokenString string) (error,jwt.MapClaims){
    
    // Parse takes the token string and a function for looking up the key. The latter is especially
    // useful if you use multiple keys for your application.  The standard is to use 'kid' in the
    // head of the token to identify which key to use, but the parsed token (head and claims) is provided
    // to the callback, providing flexibility.
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Don't forget to validate the alg is what you expect:
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }
        
        // hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
        return []byte(JWT_SECRET), nil
    })

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        //fmt.Println(claims["role"])
        return nil,claims
    } else {
        //fmt.Println(err)
        return err, nil
    }
}
