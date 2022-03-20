package gl_tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	go_http "net/http"
	go_time "time"
)

func VerifyHostCertificate(host string, port int) error {
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), nil)
	if err != nil {
		return fmt.Errorf("server does not support tls certificate: %s", err.Error())
	}

	err = conn.VerifyHostname(host)
	if err != nil {
		return fmt.Errorf("hostname does not match with the certificate: %s", err.Error())
	}

	expiry := conn.ConnectionState().PeerCertificates[0].NotAfter
	if go_time.Now().UTC().After(expiry) {
		certExpiryInfo := fmt.Sprintf("Issuer: %s\nExpiry: %v\n", conn.ConnectionState().PeerCertificates[0].Issuer, expiry.Format(go_time.RFC850))
		return fmt.Errorf("certificate is expired:\n%s", certExpiryInfo)
	}
	return nil
}

func VerifySelfSignedCertificate(hostAddr string, certPem []byte) error {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	ok := rootCAs.AppendCertsFromPEM(certPem)
	if !ok {
		return fmt.Errorf("unable to append certificate")
	}

	config := &tls.Config{
		RootCAs: rootCAs,
	}
	transport := &go_http.Transport{TLSClientConfig: config}
	httpCli := &go_http.Client{Transport: transport}

	req, _ := go_http.NewRequest(go_http.MethodGet, hostAddr, nil)
	_, err := httpCli.Do(req)
	if err != nil {
		return fmt.Errorf("call to %s failed: %s", hostAddr, err.Error())
	}
	return nil
}
