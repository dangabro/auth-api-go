package data

import (
	"fmt"
	"strings"
)

func CheckUserData(data UpdateUserData) error {
	err := checkStringRange(data.Id, "id", true, 1, true, 64, true)
	if err != nil {
		return err
	}

	err = checkStringRange(data.Name, "name", true, 1, true, 128, true)
	if err != nil {
		return err
	}

	err = checkStringRange(data.Login, "login", true, 1, true, 64, true)
	if err != nil {
		return err
	}

	return nil
}

func checkStringRange(str string, name string, trim bool, minLength int, checkMinimumLength bool, maxLength int, checkMaximumLength bool) error {
	strData := str

	// see if we trim
	if trim {
		strData = strings.TrimSpace(strData)
	}

	// min length
	if checkMinimumLength {
		if len(strData) < minLength {
			return fmt.Errorf("%s too short, must be at least %d in length", name, minLength)
		}
	}

	// max length
	if checkMaximumLength {
		if len(strData) > maxLength {
			return fmt.Errorf("%s is too long, must have maximum %d characters in length", name, maxLength)
		}
	}

	return nil
}

// CheckValidRights all the rights in the string must exist in the published list
func CheckValidRights(rights []string, allRights []RightData) error {
	rightMap := make(map[string]bool)
	for _, rData := range allRights {
		rightMap[rData.Cd] = true
	}

	var invalidRights []string
	for _, right := range rights {
		_, ok := rightMap[right]

		if !ok {
			invalidRights = append(invalidRights, right)
		}
	}

	if len(invalidRights) > 0 {
		// error, there are invalid rights
		strRights := strings.Join(invalidRights, ", ")

		return fmt.Errorf("Invalid rights: %s", strRights)
	}

	return nil
}

// CheckAuthFullConnect check auth admin and connect are part of the list
func CheckAuthFullConnect(rights []string) error {
	rightMap := make(map[string]bool)
	for _, right := range rights {
		rightMap[right] = true
	}

	checkedRights := []string{AUTH_ADMIN, AUTH_CONNECT}
	for _, right := range checkedRights {
		_, ok := rightMap[right]

		if !ok {
			return fmt.Errorf("when changing rights for yourself, you cannot remove AUTH_CONNECT or AUTH_ADMIN")
		}
	}

	return nil
}
