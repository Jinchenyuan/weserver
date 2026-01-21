package ginhandler

type RegisterRequest struct {
	Account  string `json:"account" example:"john"`            // 用户名
	Name     string `json:"name" example:"John Doe"`           // 姓名
	Password string `json:"password" example:"securepassword"` // 密码
	Email    string `json:"email" example:"john@example.com"`  // 邮箱
}

type RegisterResponse struct {
	Code    int32  `json:"code" example:"201"`                        // 响应码
	Message string `json:"message" example:"Registration successful"` // 响应消息
}

type LoginRequest struct {
	Account  string `json:"account" example:"john_doe"`        // 账号
	Password string `json:"password" example:"securepassword"` // 密码
}

type LoginResponse struct {
	Code      int32  `json:"code" example:"200"`                 // 响应码
	AccountId uint32 `json:"account_id" example:"123"`           // 账号ID
	Token     string `json:"token" example:"some-token"`         // 认证令牌
	Message   string `json:"message" example:"Login successful"` // 响应消息
}

type HelloRequest struct {
	Name string `json:"name" example:"John"` // 名字
}

type HelloResponse struct {
	Message string `json:"message" example:"Hello, John"` // 响应消息
}
