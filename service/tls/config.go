/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package service implements the functions, types, and interfaces for the module.
package tls

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/goexts/generic/settings"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/toolkits/errors"
)

func NewServerTLSConfig(cfg *configv1.TLSConfig, options ...Option) (*tls.Config, error) {
	if cfg == nil {
		return nil, nil
	}

	var err error
	var tlsCfg *tls.Config
	if cfg.File != nil {
		if tlsCfg, err = NewServerTLSConfigFromFile(
			cfg.File.GetKey(),
			cfg.File.GetCert(),
			cfg.File.GetCa(),
			options...,
		); err != nil {
			return nil, err
		}
	} else if cfg.Pem != nil {
		if tlsCfg, err = NewServerTLSConfigFromPem(
			cfg.Pem.GetKey(),
			cfg.Pem.GetCert(),
			cfg.Pem.GetCa(),
			options...,
		); err != nil {
			return nil, err
		}
	} else {

	}
	return tlsCfg, nil
}

func NewServerTLSConfigFromPem(key []byte, cert []byte, ca []byte, options ...Option) (*tls.Config, error) {
	if len(key) == 0 || len(cert) == 0 {
		return nil, fmt.Errorf("KeyPEMBlock and CertPEMBlock must both be present[key: %v, cert: %v]", key, cert)
	}

	cfg := settings.Apply(&tls.Config{}, options)
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

	cfg := settings.Apply(&tls.Config{}, options)
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

func NewClientTLSConfig(cfg *configv1.TLSConfig, options ...Option) (*tls.Config, error) {
	if cfg == nil {
		return nil, nil
	}

	var err error
	var tlsCfg *tls.Config
	if cfg.File != nil {
		if tlsCfg, err = NewClientTLSConfigFromFile(
			cfg.File.GetKey(),
			cfg.File.GetCert(),
			cfg.File.GetCa(),
			options...,
		); err != nil {
			return nil, err
		}
	} else if cfg.Pem != nil {
		if tlsCfg, err = NewClientTLSConfigFromPem(
			cfg.Pem.GetKey(),
			cfg.Pem.GetCert(),
			cfg.Pem.GetCa(),
			options...,
		); err != nil {
			return nil, err
		}
	} else {

	}

	return tlsCfg, nil
}

func NewClientTLSConfigFromPem(key []byte, cert []byte, ca []byte, options ...Option) (*tls.Config, error) {
	if len(key) == 0 || len(cert) == 0 {
		return nil, errors.Errorf("KeyPEMBlock and CertPEMBlock must both be present[key: %v, cert: %v]", key, cert)
	}

	cfg := settings.Apply(&tls.Config{}, options)
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
	cfg := settings.Apply(&tls.Config{}, options)
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
