// Package users предоставляет примитивы для взаимодействия с API пользователей.
//
// Код сгенерирован с помощью github.com/deepmap/oapi-codegen версии v1.16.3, НЕ РЕДАКТИРОВАТЬ.
package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
)

// User определяет модель для пользователя.
type User struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"` 
    Password  string `json:"password"` 
}

type CreateUserRequest struct {
    Body *PostUsersJSONRequestBody `json:"body"`
}
// PostUsersJSONRequestBody определяет тело для PostUsers для ContentType application/json.
type PostUsersJSONRequestBody = User

// PatchUsersJSONRequestBody определяет тело для PatchUsers для ContentType application/json.
type PatchUsersJSONRequestBody = User

// DeleteUsersParams определяет параметры для DeleteUsers.
type DeleteUsersParams struct {
    Id uint `form:"id" json:"id"`
}

// ServerInterface представляет все серверные обработчики.
type ServerInterface interface {
    // Получить всех пользователей
    // (GET /users)
    GetUsers(ctx echo.Context) error
    // Создать нового пользователя
    // (POST /users)
    PostUsers(ctx echo.Context) error
    // Обновить существующего пользователя
    // (PATCH /users/{id})
    PatchUsers(ctx echo.Context, request PatchUsersRequestObject) error
    // Удалить пользователя
    // (DELETE /users/:id)
    DeleteUsers(ctx echo.Context) error // Измените здесь
}

// ServerInterfaceWrapper преобразует контексты echo в параметры.
type ServerInterfaceWrapper struct {
    Handler ServerInterface
}

// GetUsers преобразует контекст echo в параметры.
func (w *ServerInterfaceWrapper) GetUsers(ctx echo.Context) error {
    if err := w.Handler.GetUsers(ctx); err != nil {
        return err
    }
    return nil
}

// PostUsers преобразует контекст echo в параметры.
func (w *ServerInterfaceWrapper) PostUsers(ctx echo.Context) error {
    var body PostUsersJSONRequestBody
    if err := ctx.Bind(&body); err != nil {
        return err
    }
    
    // Теперь вызываем метод без передачи request
    err := w.Handler.PostUsers(ctx)
    return err
}

// PatchUsers преобразует контекст echo в параметры.
func (w *ServerInterfaceWrapper) PatchUsers(ctx echo.Context) error {
    var request PatchUsersRequestObject

    idParam := ctx.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
    }
    request.Id = uint(id)

    var body PatchUsersJSONRequestBody
    if err := ctx.Bind(&body); err != nil {
        return err
    }
    request.Body = &body

    err = w.Handler.PatchUsers(ctx, request)
    return err
}

// DeleteUsers преобразует контекст echo в параметры.
func (w *ServerInterfaceWrapper) DeleteUsers(ctx echo.Context) error {
    //var request DeleteUsersRequestObject

    idParam := ctx.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
    }
    _ = uint(id)

    err = w.Handler.DeleteUsers(ctx) // Измените здесь, убрав аргументы
    return err
}
// EchoRouter представляет интерфейс для добавления маршрутов Echo.
type EchoRouter interface {
    CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
    DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
    GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
    HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
    OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
    PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
    POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
    PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
    TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers добавляет каждый маршрут сервера к EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
    RegisterHandlersWithBaseURL(router, si, "")
}

// Регистрация обработчиков с добавлением базового URL.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {
    wrapper := ServerInterfaceWrapper{
        Handler: si,
    }

    router.GET(baseURL+"/users", wrapper.GetUsers)
    router.POST(baseURL+"/users", wrapper.PostUsers)
    router.PATCH(baseURL+"/users/:id", wrapper.PatchUsers)
    router.DELETE(baseURL+"/users/:id", wrapper.DeleteUsers)
}

type GetUsersRequestObject struct {
}

type GetUsersResponseObject interface {
    VisitGetUsersResponse(w http.ResponseWriter) error
}

type GetUsers200JSONResponse []User

func (response GetUsers200JSONResponse) VisitGetUsersResponse(w http.ResponseWriter) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)

    return json.NewEncoder(w).Encode(response)
}

type PostUsersRequestObject struct {
    Body *PostUsersJSONRequestBody
}

type PostUsersResponseObject interface {
    VisitPostUsersResponse(w http.ResponseWriter) error
}

type PostUsers201JSONResponse User

func (response PostUsers201JSONResponse) VisitPostUsersResponse(w http.ResponseWriter) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(201)

    return json.NewEncoder(w).Encode(response)
}

type PatchUsersRequestObject struct {
    Id   uint                        `param:"id"` // ID пользователя
    Body *PatchUsersJSONRequestBody `json:"body"` // Тело запроса
}

type PatchUsersResponseObject interface {
    VisitPatchUsersResponse(w http.ResponseWriter) error
}

type PatchUsers200JSONResponse User

func (response PatchUsers200JSONResponse) VisitPatchUsersResponse(w http.ResponseWriter) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)

    return json.NewEncoder(w).Encode(response)
}

type DeleteUsersRequestObject struct {
    Id uint `param:"id"` // ID пользователя
}

type DeleteUsersResponseObject interface {
    VisitDeleteUsersResponse(w http.ResponseWriter) error
}

type DeleteUsers204Response struct{}

func (response DeleteUsers204Response) VisitDeleteUsersResponse(w http.ResponseWriter) error {
    w.WriteHeader(204) // No Content
    return nil
}

// StrictServerInterface представляет все серверные обработчики.
type StrictServerInterface interface {
    GetUsers(ctx context.Context, request GetUsersRequestObject) (GetUsersResponseObject, error)
    PostUsers(ctx context.Context, request PostUsersRequestObject) (PostUsersResponseObject, error)
    PatchUsers(ctx context.Context, request PatchUsersRequestObject) (PatchUsersResponseObject, error)
    DeleteUsers(ctx context.Context, request DeleteUsersRequestObject) (DeleteUsersResponseObject, error) // Здесь все правильно
}

type StrictHandlerFunc = strictecho.StrictEchoHandlerFunc
type StrictMiddlewareFunc = strictecho.StrictEchoMiddlewareFunc

type strictHandler struct {
    ssi       StrictServerInterface
    middlewares []StrictMiddlewareFunc
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
    return &strictHandler{ssi: ssi, middlewares: middlewares}
}

// Реализация метода GetUsers
func (sh *strictHandler) GetUsers(ctx echo.Context) error {
    var request GetUsersRequestObject
    handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
        return sh.ssi.GetUsers(ctx.Request().Context(), request.(GetUsersRequestObject))
    }
    for _, middleware := range sh.middlewares {
        handler = middleware(handler, "GetUsers")
    }

    response, err := handler(ctx, request)
    if err != nil {
        return err
    } else if validResponse, ok := response.(GetUsersResponseObject); ok {
        return validResponse.VisitGetUsersResponse(ctx.Response())
    } else if response != nil {
        return fmt.Errorf("unexpected response type: %T", response)
    }
    return nil
}

// Реализация метода PostUsers
func (sh *strictHandler) PostUsers(ctx echo.Context) error {
    var request PostUsersRequestObject
    if err := ctx.Bind(&request.Body); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
    }

    handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
        return sh.ssi.PostUsers(ctx.Request().Context(), request.(PostUsersRequestObject))
    }

    for _, middleware := range sh.middlewares {
        handler = middleware(handler, "PostUsers")
    }

    response, err := handler(ctx, request)
    if err != nil {
        return err
    } else if validResponse, ok := response.(PostUsersResponseObject); ok {
        return validResponse.VisitPostUsersResponse(ctx.Response())
    } else if response != nil {
        return fmt.Errorf("unexpected response type: %T", response)
    }
    return nil
}

// Реализация метода PatchUsers
func (sh *strictHandler) PatchUsers(ctx echo.Context, request PatchUsersRequestObject) error {
    var body PatchUsersJSONRequestBody
    if err := ctx.Bind(&body); err != nil {
        return err
    }
    request.Body = &body

    handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
        return sh.ssi.PatchUsers(ctx.Request().Context(), request.(PatchUsersRequestObject))
    }
    for _, middleware := range sh.middlewares {
        handler = middleware(handler, "PatchUsers")
    }

    response, err := handler(ctx, request)
    if err != nil {
        return err
    } else if validResponse, ok := response.(PatchUsersResponseObject); ok {
        return validResponse.VisitPatchUsersResponse(ctx.Response())
    } else if response != nil {
        return fmt.Errorf("unexpected response type: %T", response)
    }
    return nil
}

// Реализация метода DeleteUsers
func (sh *strictHandler) DeleteUsers(ctx echo.Context) error {
    // Извлекаем ID из параметров URL
    idParam := ctx.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
    }
    
    // Создаем объект запроса
    request := DeleteUsersRequestObject{
        Id: uint(id), // Устанавливаем ID
    }

    handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
        return sh.ssi.DeleteUsers(ctx.Request().Context(), request.(DeleteUsersRequestObject))
    }
    for _, middleware := range sh.middlewares {
        handler = middleware(handler, "DeleteUsers")
    }

    response, err := handler(ctx, request)
    if err != nil {
        return err
    } else if validResponse, ok := response.(DeleteUsersResponseObject); ok {
        return validResponse.VisitDeleteUsersResponse(ctx.Response())
    } else if response != nil {
        return fmt.Errorf("unexpected response type: %T", response)
    }
    return nil
}