package ginhandler

type RegisterRequest struct {
	Username string `json:"username" example:"john_doe"`       // 用户名
	Password string `json:"password" example:"securepassword"` // 密码
	Email    string `json:"email" example:"john@example.com"`  // 邮箱
}

type RegisterResponse struct {
	Code    int32  `json:"code" example:"201"`                        // 响应码
	Message string `json:"message" example:"Registration successful"` // 响应消息
}

type LoginRequest struct {
	Username string `json:"username" example:"john_doe"`       // 用户名
	Password string `json:"password" example:"securepassword"` // 密码
}

type LoginResponse struct {
	Code    int32  `json:"code" example:"200"`                 // 响应码
	Token   string `json:"token" example:"some-token"`         // 认证令牌
	Message string `json:"message" example:"Login successful"` // 响应消息
}

type HelloRequest struct {
	Name string `json:"name" example:"John"` // 名字
}

type HelloResponse struct {
	Message string `json:"message" example:"Hello, John"` // 响应消息
}
