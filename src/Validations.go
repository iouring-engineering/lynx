package main

import (
	"errors"
)

func validateAndroid(input MobileInputs) error {
	if input.Fbl == "" || input.WebUrl == "" {
		return errors.New("Fallback url and web url cannot be used at same time")
	}
	return nil
}

func validateIos(input MobileInputs) error {
	return nil
}

func validateCreateLink(req CreateShortLinkRequest) error {
	if err := validateAndroid(req.Android); err != nil {
		return err
	}
	if err := validateIos(req.Ios); err != nil {
		return err
	}
	return nil
}
