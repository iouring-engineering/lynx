package main

import (
	"errors"
)

func validateExpiry(expType ExpiryType) error {
	if expType == EXPIRY_MINUTES || expType == EXPIRY_HOURS || expType == EXPIRY_DAYS {
		return nil
	}
	return errors.New("Invalid expiry type")
}

func validateAndroid(input MobileInputs) error {
	if LINK_DEFAULT == input.Type || LINK_DEEP == input.Type || LINK_WEB == input.Type {
		if LINK_DEEP == input.Type || LINK_WEB == input.Type {
			if input.Url == "" {
				return errors.New("Input url missing for android3")
			}
		}
		return nil
	}
	return errors.New("Invalid android type")
}

func validateIos(input MobileInputs) error {
	return nil
}

func validateDesktop(input DeskTopInput) error {
	if LINK_DEFAULT == input.Type || LINK_WEB == input.Type {
		if LINK_WEB == input.Type && input.Url == "" {
			return errors.New("Input url missing for desktop")
		}
		return nil
	}
	return errors.New("Invalid link type for desktop")
}

func validateCreateLink(req CreateShortLinkRequest) error {
	if err := validateExpiry(req.Expiry.Type); err != nil {
		return err
	}
	if err := validateAndroid(req.Android); err != nil {
		return err
	}
	if err := validateIos(req.Ios); err != nil {
		return err
	}
	if err := validateDesktop(req.Desktop); err != nil {
		return err
	}
	return nil
}
