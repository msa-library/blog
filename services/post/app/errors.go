package app

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//Ошибки уровня бизнес логики
var (
	//Ошибки валидации
	ErrTitleIsEmpty      			= status.Error(codes.InvalidArgument, "Не заполнено поле заголовок")
	ErrSubTitleIsEmpty      		= status.Error(codes.InvalidArgument, "Не заполнено поле подзаголовок")
	ErrContentIsEmpty      			= status.Error(codes.InvalidArgument, "Не заполнено поле текст")

	ErrRetrievingMetadata			= status.Errorf(codes.InvalidArgument, "Ошибка чтения метаданных")
	ErrAuthorization				= status.Errorf(codes.Unauthenticated, "Не авторизованные запросы не поддерживаются")
	
	//Ошибки CRUD
	ErrInsert      					= status.Error(codes.Internal, "Ошибка создания записи")
	
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