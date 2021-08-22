package skeletonkeyslx

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	passwordx "github.com/herebythere/passwordx/v0.1/golang"
	sclx "github.com/herebythere/supercachelx/v0.1/golang"
)

type KeyDetails struct {
	Password string   `json:"password"`
	Services []string `json:"services"`
}
type SkeletonKeyMap = map[string]KeyDetails
type AvailableServiceList = []string

const (
	expCache            = "EX"
	weekInSeconds       = 604800
	setCache            = "SET"
	getCache            = "GET"
	okCache             = "OK"
	trueAsStr           = "1"
	colonDelimiter      = ":"
	availableServices   = "available_services"
	saltedPasswordHash  = "salted_password_hash"
	skeletonKeyServices = "skeleton_key_services"
)

var (
	errAvailableServiceDoesNotExist   = errors.New("available service does not exist")
	errSkeletonKeysAreNil             = errors.New("skeleton keys are nil")
	errSkeletonKeyServiceDoesNotExist = errors.New("skeleton key service does not exist")
	errSetKeyUnsuccessful             = errors.New("set skeleton key was unsuccessful")
	errSetKeyServiceUnsuccessful      = errors.New("set skeleton key service was unsuccessful")
	errSetServiceUnsuccessful         = errors.New("set service was unsuccessful")
	errRequestFailedToResolve         = errors.New("request failed to resolve instructions")
)

/*
 * BUILD CACHE STORE FOR REGISTRY SKELETON KEYS
 * AND ASSOCIATED ROLES
 */

func getCacheSetID(categories ...string) string {
	return strings.Join(categories, colonDelimiter)
}

func setAvailableService(
	cacheAddress string,
	identifier string,
	service string,
) (
	bool,
	error,
) {
	setID := getCacheSetID(identifier, availableServices, service)
	instructions := []interface{}{
		setCache,
		setID,
		true,
		expCache,
		weekInSeconds,
	}

	respStr, errRespStr := sclx.ExecInstructionsAndParseString(
		cacheAddress,
		&instructions,
	)
	if errRespStr != nil {
		return false, errRespStr
	}

	if *respStr == okCache {
		return true, nil
	}

	return false, errRequestFailedToResolve
}

func getAvailableService(
	cacheAddress string,
	identifier string,
	service string,
) (
	bool,
	error,
) {
	setID := getCacheSetID(identifier, availableServices, service)
	instructions := []interface{}{getCache, setID}

	respStr, errRespStr := sclx.ExecInstructionsAndParseBase64(
		cacheAddress,
		&instructions,
	)
	if errRespStr != nil {
		return false, errRespStr
	}

	if *respStr == trueAsStr {
		return true, nil
	}

	return false, errRespStr
}

func parseAvailableServicesByFilepath(
	path string,
) (
	*AvailableServiceList,
	error,
) {
	servicesJSON, errServicesJSON := ioutil.ReadFile(path)
	if errServicesJSON != nil {
		return nil, errServicesJSON
	}

	var services AvailableServiceList
	errServices := json.Unmarshal(servicesJSON, &services)

	return &services, errServices
}

func setSkeletonKey(
	cacheAddress string,
	identifier string,
	username string,
	password string,
) (
	bool,
	error,
) {
	setID := getCacheSetID(identifier, saltedPasswordHash, username)
	hashResults, errHashResults := passwordx.HashPassword(
		password,
		&passwordx.DefaultHashParams,
	)
	if errHashResults != nil {
		return false, errHashResults
	}

	// marshal into json string
	hashResultsBytes, errHashResultsBytes := json.Marshal(hashResults)
	if errHashResultsBytes != nil {
		return false, errHashResultsBytes
	}

	// store hashed results as string
	hashResultsJSONStr := string(hashResultsBytes)
	instructions := []interface{}{
		setCache,
		setID,
		hashResultsJSONStr,
		expCache,
		weekInSeconds,
	}

	// setCache does not fail
	respStr, errRespStr := sclx.ExecInstructionsAndParseString(
		cacheAddress,
		&instructions,
	)
	if errRespStr != nil {
		return false, errRespStr
	}

	if *respStr == okCache {
		return true, nil
	}

	return false, errRequestFailedToResolve
}

func setSkeletonKeyService(
	cacheAddress string,
	identifier string,
	username string,
	service string,
) (
	bool,
	error,
) {
	setID := getCacheSetID(
		identifier,
		skeletonKeyServices,
		service,
	)
	instructions := []interface{}{
		setCache,
		setID,
		true,
		expCache,
		weekInSeconds,
	}

	respStr, errRespStr := sclx.ExecInstructionsAndParseString(
		cacheAddress,
		&instructions,
	)
	if errRespStr != nil {
		return false, errRespStr
	}
	if *respStr == okCache {
		return true, nil
	}

	return false, errRequestFailedToResolve
}

func getSkeletonKeyService(
	cacheAddress string,
	identifier string,
	service string,
) (
	bool,
	error,
) {
	setID := getCacheSetID(identifier, skeletonKeyServices, service)
	instructions := []interface{}{getCache, setID}

	respStr, errRespStr := sclx.ExecInstructionsAndParseBase64(
		cacheAddress,
		&instructions,
	)
	if errRespStr != nil {
		return false, errRespStr
	}
	if *respStr == trueAsStr {
		return true, nil
	}

	return false, errRequestFailedToResolve
}

func parseSkeletonKeysByFilepath(path string) (*SkeletonKeyMap, error) {
	skeletonKeysJSON, errSkeletonKeysJSON := ioutil.ReadFile(path)
	if errSkeletonKeysJSON != nil {
		return nil, errSkeletonKeysJSON
	}

	var skeletonKeys SkeletonKeyMap
	errSkeletonKeys := json.Unmarshal(skeletonKeysJSON, &skeletonKeys)

	return &skeletonKeys, errSkeletonKeys
}

/*
 * Exposed API
 *
 */

func ParseAndSetAvailableServices(
	cacheAddress string,
	identifier string,
	path string,
) error {
	availableServices, errAvailableServices := parseAvailableServicesByFilepath(path)
	if errAvailableServices != nil {
		return errAvailableServices
	}
	if availableServices == nil {
		return errAvailableServiceDoesNotExist
	}

	for _, service := range *availableServices {
		setSuccessful, errSetServices := setAvailableService(
			cacheAddress,
			identifier,
			service,
		)
		if errSetServices != nil {
			return errSetServices
		}
		if !setSuccessful {
			return errSetServiceUnsuccessful
		}
	}

	return nil
}

func ParseAndSetSkeletonKeys(
	cacheAddress string,
	identifier string,
	path string,
) error {
	skeletonKeys, errSkeletonKeys := parseSkeletonKeysByFilepath(path)
	if errSkeletonKeys != nil {
		return errSkeletonKeys
	}
	if skeletonKeys == nil {
		return errSkeletonKeysAreNil
	}

	for username, details := range *skeletonKeys {
		setKeySuccess, errSetKey := setSkeletonKey(
			cacheAddress,
			identifier,
			username,
			details.Password,
		)
		if errSetKey != errSetKey {
			return errSetKey
		}
		if !setKeySuccess {
			return errSetKeyUnsuccessful
		}

		for _, service := range details.Services {
			setServiceSuccess, errSetService := setSkeletonKeyService(
				cacheAddress,
				identifier,
				username,
				service,
			)
			if errSetService != nil {
				return errSetService
			}
			if !setServiceSuccess {
				return errSetKeyServiceUnsuccessful
			}
		}
	}

	return nil
}

func VerifySkeletonKey(
	cacheAddress string,
	identifier string,
	username string,
	password string,
) (
	bool,
	error,
) {
	setID := getCacheSetID(identifier, saltedPasswordHash, username)
	instructions := []interface{}{getCache, setID}

	respStr, errRespStr := sclx.ExecInstructionsAndParseBase64(
		cacheAddress,
		&instructions,
	)
	if errRespStr != nil {
		return false, errRespStr
	}

	var hashResults passwordx.HashResults
	errHashResults := json.Unmarshal([]byte(*respStr), &hashResults)
	if errHashResults != nil {
		return false, errHashResults
	}

	return passwordx.PasswordIsValid(password, &hashResults)
}

func VerifySkeletonKeyAndService(
	cacheAddress string,
	identifier string,
	service string,
	username string,
	password string,
) (
	bool,
	error,
) {
	skeletonKeyHasService, errSkeletonKeyService := getSkeletonKeyService(
		cacheAddress,
		identifier,
		service,
	)
	if errSkeletonKeyService != nil {
		return false, errSkeletonKeyService
	}
	if !skeletonKeyHasService {
		return false, errSkeletonKeyServiceDoesNotExist
	}

	serviceIsAvailable, errAvailableServices := getAvailableService(
		cacheAddress,
		identifier,
		service,
	)
	if errAvailableServices != nil {
		return false, errAvailableServices
	}
	if !serviceIsAvailable {
		return false, errAvailableServiceDoesNotExist
	}

	return VerifySkeletonKey(cacheAddress, identifier, username, password)
}

func SetupSkeletonKeysAndAvailableServices(
	cacheAddress string,
	identifier string,
	availableServicesPath string,
	skeletonKeysPath string,
) error {
	errPaseAvailableServices := ParseAndSetAvailableServices(
		cacheAddress,
		identifier,
		availableServicesPath,
	)
	if errPaseAvailableServices != nil {
		return errPaseAvailableServices
	}

	errParseSkeletonKeys := ParseAndSetSkeletonKeys(
		cacheAddress,
		identifier,
		skeletonKeysPath,
	)

	return errParseSkeletonKeys
}
