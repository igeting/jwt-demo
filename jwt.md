# jwt

## jwt structure

### 1.header
```
{
  "alg": "HS256",
  "typ": "JWT"
}
```

### 2.payload
- iss：发行人
- exp：到期时间
- sub：主题
- aud：用户
- nbf：在此之前不可用
- iat：发布时间
- jti：JWT ID用于标识该JWT

```
{
  "sub": "1234567890",
  "name": "Helen",
  "admin": true
}
```

### 3.signature
```
JWTString=Base64(Header).Base64(Payload).HMACSHA256(base64UrlEncode(header)+"."+base64UrlEncode(payload),secret)
```

## effect

> header和payload可以直接利用base64解码出原文，从header中获取哈希签名的算法，从payload中获取有效数据

> signature由于使用了不可逆的加密算法，无法解码出原文，它的作用是校验token有没有被篡改。服务端获取header中的加密算法之后，利用该算法加上secretKey对header、payload进行加密，比对加密后的数据和客户端发送过来的是否一致。注意secretKey只能保存在服务端，而且对于不同的加密算法其含义有所不同，一般对于MD5类型的摘要加密算法，secretKey实际上代表的是盐值
