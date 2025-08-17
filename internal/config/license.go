package config

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v3/log"
	"github.com/keygen-sh/keygen-go/v3"
	"github.com/keygen-sh/machineid"
)

var (
	keygenAccount   string
	keygenProduct   string
	keygenPublicKey string
)

const (
	editionPersonal  = "personal"
	editionLite      = "lite"
	editionAirGapped = "air-gapped"
)

func Limits(ctx context.Context, licenseKeySigned, licenseFile string) (PackageLimits, func() error, error) {
	if licenseKeySigned == "" {
		if licenseFile == "" {
			return editionLimits(editionPersonal), nil, nil
		}
		return editionLimits(editionPersonal), nil, errors.New("license key is required. only online validation is supported")
	}

	fingerprint, err := machineid.ProtectedID(keygen.Product)
	if err != nil {
		return editionLimits(editionPersonal), nil, err
	}

	license, err := licenseValidate(ctx, licenseKeySigned, fingerprint)

	deactivationFunc := func() error {
		license.Deactivate(ctx, fingerprint)
		return nil
	}

	if err != nil {
		return editionLimits(editionPersonal), deactivationFunc, err
	}

	return editionLimits(editionLite), deactivationFunc, nil
}

func licenseValidate(ctx context.Context, licenseKeySigned string, fingerprint string) (*keygen.License, error) {
	keygen.Account = keygenAccount
	keygen.Product = keygenProduct
	keygen.PublicKey = keygenPublicKey
	keygen.LicenseKey = licenseKeySigned

	license, err := keygen.Validate(ctx, fingerprint)
	log.Info("license validated")
	switch {
	case err == keygen.ErrLicenseNotActivated:
		_, err := license.Activate(ctx, fingerprint)
		switch {
		case err == keygen.ErrMachineLimitExceeded:
			return nil, errors.New("machine limit has been exceeded")
		case err != nil:
			return nil, errors.New("machine activation failed")
		}
	case err == keygen.ErrLicenseExpired:
		return nil, errors.New("license is expired")
	case err != nil:
		return nil, err
	}
	log.Info("license activated")
	return license, nil
}
