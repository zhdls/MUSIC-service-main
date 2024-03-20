package app

//组合其提供的 API，设计我们的鉴权场景

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/travel_study/blog-service/global"
	"github.com/travel_study/blog-service/pkg/util"
	"time"
)

type Claims struct {
	//嵌入的 AppKey 和 AppSecret，用于我们自定义的认证信息
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`

	//jwt.StandardClaims 结构体，它是 jwt-go 库中预定义的，也是 JWT 的规范
	//对应Payload的相关字段
	jwt.StandardClaims
}

// GetJWTSecret 获取该项目的 JWT Secret
func GetJWTSecret() []byte {
	return []byte(global.JWTSetting.Secret)
}

// GenerateToken 生成 JWT Token
//根据客户端传入的 AppKey 和 AppSecret 以及在项目配置中所设置的签发者（Issuer）和过期时间（ExpiresAt），
//根据指定的算法生成签名后的 Token
func GenerateToken(appKey, appSecret string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(global.JWTSetting.Expire) //过期时间
	claims := Claims{
		AppKey:    util.EncodeMD5(appKey),
		AppSecret: util.EncodeMD5(appSecret),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    global.JWTSetting.Issuer,
		},
	}

	//根据 Claims 结构体创建 Token 实例
	//第一个参数：SigningMethod，其包含 SigningMethodHS256、SigningMethodHS384、SigningMethodHS512 三种 crypto.Hash 加密算法的方案
	//第二个参数：Claims，主要是用于传递用户所预定义的一些权利要求，便于后续的加密、校验等行为
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//生成签名字符串，根据所传入 Secret 不同，进行签名并返回标准的 Token
	//使用指定的secret签名并获得完整的编码后的字符串token
	token, err := tokenClaims.SignedString(GetJWTSecret())
	return token, err
}

// ParseToken 解析和校验 Token
//解析传入的 Token，然后根据 Claims 的相关属性要求进行校验
func ParseToken(token string) (*Claims, error) {
	//解析鉴权的声明	方法内部主要是具体的解码和校验的过程，最终返回 *Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil { //如果有声明在令牌中
		//Valid：验证基于时间的声明；如：过期时间（ExpiresAt）、签发者（Issuer）、生效时间（Not Before）
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	//如果没有任何声明在令牌中，仍然会被认为是有效的

	return nil, err
}
