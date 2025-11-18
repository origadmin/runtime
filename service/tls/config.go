/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package tls implements the functions, types, and interfaces for the module.
package tls

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/goexts/generic/configure"
	tlsv1 "github.com/origadmin/runtime/api/gen/go/config/transport/tls/v1"

	"github.com/origadmin/toolkits/errors"
)

var tlsVersionMap = map[string]uint16{
	"1.0": tls.VersionTLS10,
	"1.1": tls.VersionTLS11,
	"1.2": tls.VersionTLS12,
	"1.3": tls.VersionTLS13,
}

var cipherSuiteMap = map[string]uint16{
	"TLS_RSA_WITH_RC4_128_SHA":        tls.TLS_RSA_WITH_RC4_128_SHA,
	"TLS_RSA_WITH_3DES_EDE_CBC_SHA":   tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
	"TLS_RSA_WITH_AES_128_CBC_SHA":    tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	"TLS_RSA_WITH_AES_256_CBC_SHA":    tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	"TLS_RSA_WITH_AES_128_CBC_SHA256": tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
	//"TLS_RSA_WITH_AES_256_CBC_SHA256":         tls.TLS_RSA_WITH_AES_256_CBC_SHA256,
	"TLS_RSA_WITH_AES_128_GCM_SHA256":         tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_RSA_WITH_AES_256_GCM_SHA384":         tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":        tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_RC4_128_SHA":          tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
	"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":     tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":      tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":      tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
	//"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384": tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256": tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	//"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384":   tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384,
	"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384": tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":   tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305":    tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305":  tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	"TLS_AES_128_GCM_SHA256":                  tls.TLS_AES_128_GCM_SHA256,
	"TLS_AES_256_GCM_SHA384":                  tls.TLS_AES_256_GCM_SHA384,
	"TLS_CHACHA20_POLY1305_SHA256":            tls.TLS_CHACHA20_POLY1305_SHA256,
}

func parseCipherSuites(cipherSuites []string) []uint16 {
	var suites []uint16
	for _, s := range cipherSuites {
		if suite, ok := cipherSuiteMap[s]; ok {
			suites = append(suites, suite)
		}
	}
	return suites
}

func NewServerTLSConfig(cfg *tlsv1.TLSConfig, options ...Option) (*tls.Config, error) {
	if cfg == nil {
		return nil, nil
	}

	var err error
	var tlsCfg *tls.Config
	if cfg.GetFile() != nil {
		file := cfg.GetFile()
		if tlsCfg, err = NewServerTLSConfigFromFile(
			file.GetKey(),
			file.GetCert(),
			file.GetCa(),
			options...,
		); err != nil {
			return nil, err
		}
	} else if cfg.GetPem() != nil {
		pem := cfg.GetPem()
		if tlsCfg, err = NewServerTLSConfigFromPem(
			pem.GetKey(),
			pem.GetCert(),
			pem.GetCa(),
			options...,
		); err != nil {
			return nil, err
		}
	} else {
		// If no file or PEM config, create a default TLS config
		tlsCfg = configure.Apply(&tls.Config{}, options)
	}

	// Apply common TLS configurations from the proto
	if tlsCfg == nil {
		tlsCfg = configure.Apply(&tls.Config{}, options)
	}

	if version, ok := tlsVersionMap[cfg.GetMinVersion()]; ok {
		tlsCfg.MinVersion = version
	} else if cfg.GetMinVersion() != "" {
		// Default to TLS1.2 if not specified or invalid
		tlsCfg.MinVersion = tls.VersionTLS12
	}

	if len(cfg.GetCipherSuites()) > 0 {
		tlsCfg.CipherSuites = parseCipherSuites(cfg.GetCipherSuites())
	}

	if cfg.GetRequireClientCert() {
		tlsCfg.ClientAuth = tls.RequireAndVerifyClientCert
		if cfg.GetClientCaFile() != "" {
			cp, err := newRootCertWithFile(cfg.GetClientCaFile())
			if err != nil {
				return nil, errors.Wrap(err, "read client CA file error")
			}
			tlsCfg.ClientCAs = cp
		}
	} else {
		tlsCfg.ClientAuth = tls.NoClientCert
	}

	return tlsCfg, nil
}

func NewServerTLSConfigFromPem(key []byte, cert []byte, ca []byte, options ...Option) (*tls.Config, error) {
	if len(key) == 0 || len(cert) == 0 {
		return nil, fmt.Errorf("KeyPEMBlock and CertPEMBlock must both be present[key: %v, cert: %v]", key, cert)
	}

	cfg := configure.Apply(&tls.Config{}, options)
	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, errors.Wrap(err, "load x509 key pair error")
	}

	cfg.Certificates = []tls.Certificate{tlsCert}

	if len(ca) != 0 {
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
		cp, err := newRootCert(ca)
		if err != nil {
			return nil, errors.Wrap(err, "read cert PEM error")
		}

		cfg.RootCAs = cp
		cfg.ClientCAs = cp
	} else {
		cfg.ClientAuth = tls.NoClientCert
	}

	return cfg, nil
}

func NewServerTLSConfigFromFile(keyFile, certFile, caFile string, options ...Option) (*tls.Config, error) {
	if keyFile == "" || certFile == "" {
		return nil, errors.Errorf("KeyFile and CertFile must both be present[key: %v, cert: %v]", keyFile, certFile)
	}

	cfg := configure.Apply(&tls.Config{}, options)
	tlsCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, errors.Wrap(err, "load x509 key pair error")
	}

	cfg.Certificates = []tls.Certificate{tlsCert}

	if caFile != "" {
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
		cp, err := newRootCertWithFile(caFile)
		if err != nil {
			return nil, errors.Wrap(err, "read cert file error")
		}

		cfg.RootCAs = cp
		cfg.ClientCAs = cp
	} else {
		cfg.ClientAuth = tls.NoClientCert
	}

	return cfg, nil
}

func NewClientTLSConfig(cfg *tlsv1.TLSConfig, options ...Option) (*tls.Config, error) {
	if cfg == nil {
		return nil, nil
	}

	var err error
	var tlsCfg *tls.Config
	if cfg.GetFile() != nil {
		file := cfg.GetFile()
		if tlsCfg, err = NewClientTLSConfigFromFile(
			file.GetKey(),
			file.GetCert(),
			file.GetCa(),
			options...,
		); err != nil {
			return nil, err
		}
	} else if cfg.GetPem() != nil {
		pem := cfg.GetPem()
		if tlsCfg, err = NewClientTLSConfigFromPem(
			pem.GetKey(),
			pem.GetCert(),
			pem.GetCa(),
			options...,
		); err != nil {
			return nil, err
		}
	} else {
		// If no file or PEM config, create a default TLS config
		tlsCfg = configure.Apply(&tls.Config{}, options)
	}

	// Apply common TLS configurations from the proto
	if tlsCfg == nil {
		tlsCfg = configure.Apply(&tls.Config{}, options)
	}

	if version, ok := tlsVersionMap[cfg.GetMinVersion()]; ok {
		tlsCfg.MinVersion = version
	} else if cfg.GetMinVersion() != "" {
		// Default to TLS1.2 if not specified or invalid
		tlsCfg.MinVersion = tls.VersionTLS12
	}

	if len(cfg.GetCipherSuites()) > 0 {
		tlsCfg.CipherSuites = parseCipherSuites(cfg.GetCipherSuites())
	}

	if cfg.GetClientCaFile() != "" {
		cp, err := newRootCertWithFile(cfg.GetClientCaFile())
		if err != nil {
			return nil, errors.Wrap(err, "read client CA file error")
		}
		tlsCfg.RootCAs = cp
	}

	return tlsCfg, nil
}

func NewClientTLSConfigFromPem(key []byte, cert []byte, ca []byte, options ...Option) (*tls.Config, error) {
	if len(key) == 0 || len(cert) == 0 {
		return nil, errors.Errorf("KeyPEMBlock and CertPEMBlock must both be present[key: %v, cert: %v]", key, cert)
	}

	cfg := configure.Apply(&tls.Config{}, options)
	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, errors.Wrap(err, "load x509 key pair error")
	}

	cfg.Certificates = []tls.Certificate{tlsCert}
	if len(ca) != 0 {
		cp, err := newRootCert(ca)
		if err != nil {
			return nil, errors.Wrap(err, "read cert PEM error")
		}
		cfg.RootCAs = cp
	}

	return cfg, nil
}

func NewClientTLSConfigFromFile(key string, cert string, ca string, options ...Option) (*tls.Config, error) {
	cfg := configure.Apply(&tls.Config{}, options)
	if key == "" || cert == "" {
		return cfg, nil
	}

	tlsCert, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, errors.Wrap(err, "load x509 key pair error")
	}

	cfg.Certificates = []tls.Certificate{tlsCert}

	if ca != "" {
		cp, err := newRootCertWithFile(ca)
		if err != nil {
			return nil, errors.Wrap(err, "read cert file error")
		}

		cfg.RootCAs = cp
	}

	return cfg, nil
}

// newRootCertWithFile creates x509 certPool with provided CA file
func newRootCertWithFile(filepath string) (*x509.CertPool, error) {
	rootPEM, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return newRootCert(rootPEM)
}

func newRootCert(rootPEM []byte) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	var err error
	block, _ := pem.Decode(rootPEM)
	if block == nil {
		return certPool, nil
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	certPool.AddCert(cert)
	return certPool, nil
}
