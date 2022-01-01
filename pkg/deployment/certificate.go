package deployment

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"math/big"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func keyToPem(k *rsa.PrivateKey) ([]byte, error) {
	data := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	var buf bytes.Buffer
	if err := pem.Encode(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func certToPem(cert []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
	return buf.Bytes(), err
}

func GenerateCertificate(service, namespace string) ([]byte, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	sub := pkix.Name{
		CommonName:   fmt.Sprintf("%s.%s.svc.cluster.local", service, namespace),
		Organization: []string{"cmsnr"},
	}

	template := x509.Certificate{
		SerialNumber:       big.NewInt(1),
		Subject:            sub,
		IsCA:               false,
		SignatureAlgorithm: x509.SHA256WithRSA,
		NotBefore:          time.Now(),
		NotAfter:           time.Now().Add(time.Hour * 24 * 365),
		KeyUsage:           x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames: []string{
			fmt.Sprintf("%s.%s.svc", service, namespace),
			fmt.Sprintf("%s.%s", service, namespace),
			fmt.Sprintf("%s.%s.svc.cluster.local", service, namespace),
		},
	}

	certificate, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return nil, nil, err
	}

	keyPem, err := keyToPem(key)
	if err != nil {
		return nil, nil, err
	}

	certPem, err := certToPem(certificate)

	return certPem, keyPem, err
}

func GetCertificate(name, namespace string) ([]byte, error) {
	ctx := context.Background()
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	secret, err := clientSet.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return secret.Data["tls.crt"], nil

}

func CertAsSecret(cert, key []byte, name, namespace string) *corev1.Secret {

	data := map[string][]byte{
		"tls.key": key,
		"tls.crt": cert,
	}

	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: corev1.SecretTypeTLS,
		Data: data,
	}

	return secret
}
