package main

func validateAndroid(input MobileInputs) error {
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
