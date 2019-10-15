package app

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//Ошибки уровня бизнес логики
var (
	//Ошибки валидации
	ErrEmailIncorrect      				= status.Error(codes.InvalidArgument, "Некорректный E-mail")
	ErrPasswordIsEmpty      			= status.Error(codes.InvalidArgument, "Password не задан")
	ErrUserNameIsEmpty					= status.Error(codes.InvalidArgument, "E-mail не задан")
	ErrUserNameIsExist      			= status.Error(codes.AlreadyExists, "Пользователь уже зарегистрирован")

	ErrNotFound      					= status.Error(codes.NotFound, "Пользователь не найден")
	ErrIncorrectLoginOrPassword			= status.Error(codes.Unauthenticated,"Некорректный логин или пароль")

	//Ошибки CRUD
	ErrInsert      						= status.Error(codes.Internal, "Ошибка создания записи")
	ErrUpdate      						= status.Error(codes.Internal, "Ошибка сохранения записи")
	
)

//==================================================
// All gRPC err codes
//==================================================
// codes.OK - http.StatusOK
// codes.Canceled - http.StatusRequestTimeout
// codes.Unknown - http.StatusInternalServerError
// codes.InvalidArgument - http.StatusBadRequest
// codes.DeadlineExceeded - http.StatusGatewayTimeout
// codes.NotFound - http.StatusNotFound
// codes.AlreadyExists - http.StatusConflict
// codes.PermissionDenied - http.StatusForbidden
// codes.Unauthenticated - http.StatusUnauthorized
// codes.ResourceExhausted - http.StatusTooManyRequests
// codes.FailedPrecondition - http.StatusBadRequest
// codes.Aborted - http.StatusConflict
// codes.OutOfRange - http.StatusBadRequest
// codes.Unimplemented - http.StatusNotImplemented
// codes.Internal - http.StatusInternalServerError
// codes.Unavailable - http.StatusServiceUnavailable
// codes.DataLoss - http.StatusInternalServerError