package model

import "testing"

func TestCertificateName(t *testing.T) {
	cert := Certificate{
		Domain: "*.example.com",
	}

	expectedSecretName := "core_SSL_CERTIFICATE___example_com_BUNDLE"
	secretName := cert.SecretName()
	if expectedSecretName != secretName {
		t.Errorf("unexpected secret name %s (expected %s input %s)", secretName, expectedSecretName, cert.Domain)
	}

}
