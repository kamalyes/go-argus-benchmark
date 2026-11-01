/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-16 11:11:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-18 11:58:15
 * @FilePath: \go-argus-benchmark\models.go
 * @Description: 测试模型
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package bench

// SimpleUser 简单用户模型
// EN SimpleUser model
type SimpleUser struct {
	Name  string `json:"name" validate:"required,min=2,max=50"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"required,gte=0,lte=130"`
}

// ComplexOrder 复杂订单模型
// EN ComplexOrder model
type ComplexOrder struct {
	OrderID     string       `json:"order_id" validate:"required,uuid"`
	CustomerID  string       `json:"customer_id" validate:"required,uuid"`
	Email       string       `json:"email" validate:"required,email"`
	Phone       string       `json:"phone" validate:"required,e164"`
	Currency    string       `json:"currency" validate:"required,alpha,len=3"`
	TotalAmount float64      `json:"total_amount" validate:"required,gt=0"`
	Status      string       `json:"status" validate:"required,oneof=pending confirmed shipped delivered cancelled"`
	Items       []OrderItem  `json:"items" validate:"required,min=1,dive"`
	Shipping    ShippingInfo `json:"shipping" validate:"required"`
}

// OrderItem 订单项模型
// EN OrderItem model
type OrderItem struct {
	ProductID string  `json:"product_id" validate:"required,uuid"`
	SKU       string  `json:"sku" validate:"required,alphanum"`
	Quantity  int     `json:"quantity" validate:"required,gte=1,lte=999"`
	UnitPrice float64 `json:"unit_price" validate:"required,gt=0"`
}

// ShippingInfo 配送信息模型
// EN ShippingInfo model
type ShippingInfo struct {
	Country string `json:"country" validate:"required,alpha,len=2"`
	State   string `json:"state" validate:"required,min=1,max=100"`
	City    string `json:"city" validate:"required,min=1,max=100"`
	Street  string `json:"street" validate:"required,min=1,max=200"`
	ZipCode string `json:"zip_code" validate:"required,alphanum,min=3,max=10"`
}

// NestedWorkspace 嵌套工作空间模型
// EN NestedWorkspace model
type NestedWorkspace struct {
	WorkspaceID string              `json:"workspace_id" validate:"required,uuid"`
	Code        string              `json:"code" validate:"required,alphanum,min=2,max=20"`
	DisplayName string              `json:"display_name" validate:"required,min=2,max=100"`
	AdminEmail  string              `json:"admin_email" validate:"required,email"`
	Envs        []EnvironmentConfig `json:"envs" validate:"required,min=1,dive"`
}

// EnvironmentConfig 环境配置模型
// EN EnvironmentConfig model
type EnvironmentConfig struct {
	EnvID            string   `json:"env_id" validate:"required,uuid"`
	DisplayName      string   `json:"display_name" validate:"required,min=1,max=50"`
	Region           string   `json:"region" validate:"required,oneof=US EU AP SA ME AF"`
	SupportedLocales []string `json:"supported_locales" validate:"required,min=1,dive,oneof=zh en ja ko es pt ru ar hi"`
	Theme            string   `json:"theme" validate:"required,oneof=dark light auto"`
}

// SoftwareRelease 软件发布模型
// EN SoftwareRelease model
type SoftwareRelease struct {
	Version    string `json:"version" validate:"required,semver"`
	LicenseKey string `json:"license_key" validate:"required,alphanum,min=8,max=32"`
	ISBN       string `json:"isbn" validate:"required,isbn13"`
	ISSN       string `json:"issn" validate:"required,issn"`
	Schedule   string `json:"schedule" validate:"required,cron"`
	Locale     string `json:"locale" validate:"required,bcp47"`
}

// CryptoWallet 加密钱包模型
// EN CryptoWallet model
type CryptoWallet struct {
	WalletID  string `json:"wallet_id" validate:"required,uuid"`
	EthAddr   string `json:"eth_addr" validate:"required,eth_addr"`
	BtcAddr   string `json:"btc_addr" validate:"required,btc_addr"`
	BIC       string `json:"bic" validate:"required,bic"`
	AvatarURI string `json:"avatar_uri" validate:"required,datauri"`
}

// CompatProfile 兼容配置文件模型
// EN CompatProfile model
type CompatProfile struct {
	Name       string            `json:"name" validate:"required,min=2,max=16,alphanumunicode"`
	Email      string            `json:"email" validate:"required,email"`
	Age        int               `json:"age" validate:"gte=18,lte=120"`
	Password   string            `json:"password" validate:"required,min=8"`
	Confirm    string            `json:"confirm" validate:"required,eqfield=password"`
	Website    string            `json:"website" validate:"omitempty,http_url"`
	Role       string            `json:"role" validate:"oneof=admin member guest"`
	Tags       []string          `json:"tags" validate:"min=1,dive,required,lowercase"`
	Meta       map[string]string `json:"meta" validate:"omitempty,dive,required"`
	TraceID    string            `json:"trace_id" validate:"omitempty,uuid4"`
	RemoteIP   string            `json:"remote_ip" validate:"omitempty,ip"`
	RemoteCIDR string            `json:"remote_cidr" validate:"omitempty,cidr"`
}
