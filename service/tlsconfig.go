/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package builder implements the functions, types, and interfaces for the module.
package service

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/goexts/generic/settings"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type TLSOption = func(*tls.Config)

func SkipInsecure(b bool) TLSOption {
	return func(c *tls.Config) {
		c.InsecureSkipVerify = b
	}
}

func NewServerTLSConfig(cfg *configv1.TLSConfig, opts ...TLSOption) (*tls.Config, error) {
	if cfg == nil {
		return nil, nil
	}

	tlsCfg := settings.Apply(&tls.Config{}, opts)
	if cfg.GetFile() != nil {

	} else if cfg.GetPem() != nil {

	} else {

	}

	return tlsCfg, nil
}

func NewClientTLSConfig(cfg *configv1.TLSConfig, opts ...TLSOption) (*tls.Config, error) {
	if cfg == nil {
		return nil, nil
	}

	tlsCfg := settings.Apply(&tls.Config{}, opts)
	if cfg.GetFile() != nil {

	} else if cfg.GetPem() != nil {

	} else {

	}

	return tlsCfg, nil
}

// newRootCertFile creates x509 certPool with provided CA file
func newRootCertFile(filepath string) (*x509.CertPool, error) {
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
