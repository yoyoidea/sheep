package huobi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"sort"
)

// 构造签名
// mapParams: 送进来参与签名的参数, Map类型
// strMethod: 请求的方法 GET, POST......
// strHostUrl: 请求的主机
// strRequestPath: 请求的路由路径
// strSecretKey: 进行签名的密钥
func createSign(mapParams map[string]string, strMethod, strHostUrl, strRequestPath, strSecretKey string) string {
	// 参数处理, 按API要求, 参数名应按ASCII码进行排序(使用UTF-8编码, 其进行URI编码, 16进制字符必须大写)
	sortedParams := MapSortByKey(mapParams)
	encodeParams := mapValueEncodeURI(sortedParams)
	strParams := map2UrlQuery(encodeParams)

	strPayload := strMethod + "\n" + strHostUrl + "\n" + strRequestPath + "\n" + strParams

	return computeHmac256(strPayload, strSecretKey)
}

// 对Map的值进行URI编码
// mapParams: 需要进行URI编码的map
// return: 编码后的map
func mapValueEncodeURI(mapValue map[string]string) map[string]string {
	for key, value := range mapValue {
		valueEncodeURI := url.QueryEscape(value)
		mapValue[key] = valueEncodeURI
	}

	return mapValue
}

// 将map格式的请求参数转换为字符串格式的
// mapParams: map格式的参数键值对
// return: 查询字符串
func map2UrlQuery(mapParams map[string]string) string {
	var keys []string
	for key := range mapParams {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var strParams string
	for _, key := range keys {
		strParams += (key + "=" + mapParams[key] + "&")
	}

	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}

	return strParams
}

// HMAC SHA256加密
// strMessage: 需要加密的信息
// strSecret: 密钥
// return: BASE64编码的密文
func computeHmac256(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// 对Map按着ASCII码进行排序
// mapValue: 需要进行排序的map
// return: 排序后的map
func MapSortByKey(mapValue map[string]string) map[string]string {
	var keys []string
	for key := range mapValue {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	mapReturn := make(map[string]string)
	for _, key := range keys {
		mapReturn[key] = mapValue[key]
	}

	return mapReturn
}
