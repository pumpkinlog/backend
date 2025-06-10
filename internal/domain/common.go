package domain

import (
	"regexp"
)

type RegionID string

func (rid RegionID) Validate() error {
	if rid == "" {
		return ValidationError("region ID is required")
	}

	if !regionIDRegex.MatchString(string(rid)) {
		return ValidationError("region ID must match regular expression: %s", regionIDRegex)
	}

	return nil
}

type Code string

func (c Code) Validate() error {
	if c == "" {
		return ValidationError("code is required")
	}

	if len(c) > 128 {
		return ValidationError("code must be less than 128 characters")
	}

	if !codeRegex.MatchString(string(c)) {
		return ValidationError("code must match regular expression: %s", codeRegex)
	}

	return nil
}

var (
	regionIDRegex = regexp.MustCompile(`^[A-Z]{2}(?:-[A-Z]{2})?$`)
	codeRegex     = regexp.MustCompile(`^[A-Z0-9_-]+$`)
)
