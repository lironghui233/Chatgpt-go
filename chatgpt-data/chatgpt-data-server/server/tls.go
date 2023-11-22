package server

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/*
param：cert（服务器的证书文件路径），key（密钥文件路径）
*/
func GetTlsOpt(cert, key string) grpc.ServerOption {
	creds, err := credentials.NewServerTLSFromFile(cert, key) //从给定的文件路径加载服务器的TLS证书和密钥，并返回一个credentials.TransportCredentials对象
	if err != nil {
		panic(err)
	}
	return grpc.Creds(creds)
}

/*
param： clientCaCert（客户端CA证书），certFile（服务器证书文件），keyFile（服务器私钥文件）
*/
func GetMTlsOpt(clientCaCert, certFile, keyFile string) grpc.ServerOption {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile) //加载服务器证书和私钥，用于在服务器上验证自己
	if err != nil {
		panic(err)
	}
	ca := x509.NewCertPool()                //创建一个新的证书池，用于存储客户端需要验证的CA证书
	bytes, err := os.ReadFile(clientCaCert) //读取客户端CA证书文件的内容
	if err != nil {
		panic(err)
	}
	ok := ca.AppendCertsFromPEM(bytes) //将客户端CA证书添加到之前创建的证书池中
	if !ok {
		panic("append cert failed")
	}
	tlsConfig := &tls.Config{ //创建一个新的TLS配置
		ClientAuth:   tls.RequireAndVerifyClientCert, //设置服务器需要客户端提供证书，并且会对客户端证书进行验证
		Certificates: []tls.Certificate{cert},        //设置服务器使用的证书，这里只使用了一个证书，即之前加载的服务器证书
		ClientCAs:    ca,                             //设置服务器要验证的客户端CA证书池，这里使用的是之前创建的证书池
	}
	return grpc.Creds(credentials.NewTLS(tlsConfig))
}
