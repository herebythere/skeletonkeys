package skeletonkeyslx

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	passwordx "github.com/herebythere/passwordx/v0.1/golang"
)

type KeyDetails struct {
	Password string   `json:"password"`
	Services []string `json:"services"`
}
type SkeletonKeyMap = map[string]KeyDetails
type AvailableServiceList = []string

const (
	hset      = "HSET"
	hget      = "HGET"
	trueAsStr = "1"

	SkeletonKey         = "skeleton_key"
	colonDelimiter      = ":"
	AvailableServices   = "available_services"
	SaltedPasswordHash  = "salted_password_hash"
	SkeletonKeyServices = "skeleton_key_services"
	applicationJSON     = "application/json"
)

var (
	localCacheAddress = os.Getenv("LOCAL_CACHE_ADDRESS")

	errSkeletonKeyDoesNotExist        = errors.New("skeleton key does not exist")
	errNilEntry                       = errors.New("nil entry was provided")
	errAvailableServiceDoesNotExist   = errors.New("available service does not exist")
	errSkeletonKeysAreNil             = errors.New("skeleton keys are nil")
	errSkeletonKeyServiceDoesNotExist = errors.New("skeleton key service does not exist")

	errSetKeyUnsuccessful        = errors.New("set skeleton key was unsuccessful")
	errSetKeyServiceUnsuccessful = errors.New("set skeleton key service was unsuccessful")
	errSetServiceUnsuccessful    = errors.New("set service was unsuccessful")
)

/*
 * BUILD CACHE STORE FOR REGISTRY SKELETON KEYS
 * AND ASSOCIATED ROLES
 */

func getCacheSetID(categories ...string) string {
	return strings.Join(categories, colonDelimiter)
}

func parseCachedString(resp *http.Response) (*string, error) {
	var serviceAsBase64 string
	errJSONResponse := json.NewDecoder(resp.Body).Decode(&serviceAsBase64)
	if errJSONResponse != nil {

		return nil, errJSONResponse
	}

	serviceAsBytes, errServiceAsBytes := base64.URLEncoding.DecodeString(
		serviceAsBase64,
	)
	if errServiceAsBytes != nil {

		return nil, errServiceAsBytes
	}

	serviceReturned := string(serviceAsBytes)
	return &serviceReturned, nil
}

func postJSONRequest(
	instructions []interface{},
) (
	*http.Response,
	error,
) {
	if instructions == nil {
		return nil, errNilEntry
	}

	instructionsAsJSON, errInstructionsAsJSON := json.Marshal(instructions)
	if errInstructionsAsJSON != nil {

		return nil, errInstructionsAsJSON
	}

	requestBody := bytes.NewBuffer(instructionsAsJSON)

	return http.Post(localCacheAddress, applicationJSON, requestBody)
}

func setAvailableService(serverName string, service string) (bool, error) {
	setID := getCacheSetID(serverName, AvailableServices)
	instructions := []interface{}{hset, setID, service, true}

	// HSET just produced an update of the value, 0 is returned,
	// otherwise if a new field is created 1 is returned.
	// HSET does not fail
	resp, errResponse := postJSONRequest(instructions)
	if errResponse != nil {
		return false, errResponse
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

func getAvailableService(serverName string, service string) (bool, error) {
	setID := getCacheSetID(serverName, AvailableServices)
	instructions := []interface{}{hget, setID, service}

	resp, errResponse := postJSONRequest(instructions)
	if errResponse != nil {
		return false, errResponse
	}
	defer resp.Body.Close()

	serviceReturned, errServiceReturned := parseCachedString(resp)
	if errServiceReturned != nil {
		return false, errServiceReturned
	}

	return *serviceReturned == trueAsStr, nil
}

func parseAvailableServicesByFilepath(path string) (*AvailableServiceList, error) {
	servicesJSON, errServicesJSON := ioutil.ReadFile(path)
	if errServicesJSON != nil {
		return nil, errServicesJSON
	}

	var services AvailableServiceList
	errServices := json.Unmarshal(servicesJSON, &services)

	return &services, errServices
}

func parseAndSetAvailableServices(serverName string, path string, err error) error {
	availableServices, errAvailableServices := parseAvailableServicesByFilepath(path)
	if errAvailableServices != nil {
		return errAvailableServices
	}
	if availableServices == nil {
		return errAvailableServiceDoesNotExist
	}

	for _, service := range *availableServices {
		setSuccessful, errSetServices := setAvailableService(serverName, service)
		if errSetServices != nil {
			return errSetServices
		}
		if !setSuccessful {
			return errSetServiceUnsuccessful
		}
	}

	return nil
}

func setSkeletonKey(serverName string, username string, password string) (bool, error) {
	setID := getCacheSetID(serverName, SaltedPasswordHash)
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
	instructions := []interface{}{hset, setID, username, hashResultsJSONStr}

	// HSET does not fail
	resp, errResponse := postJSONRequest(instructions)
	if errResponse != nil {
		return false, errResponse
	}
	defer resp.Body.Close()

	if resp != nil && resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

func setSkeletonKeyService(serverName string, username string, service string) (bool, error) {
	setID := getCacheSetID(serverName, SkeletonKeyServices)
	serviceID := getCacheSetID(username, service)
	instructions := []interface{}{hset, setID, serviceID, true}

	resp, errResponse := postJSONRequest(instructions)
	if errResponse != nil {
		return false, errResponse
	}
	defer resp.Body.Close()

	if resp != nil && resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

func getSkeletonKeyService(serverName string, service string) (bool, error) {
	setID := getCacheSetID(serverName, SkeletonKeyServices)
	instructions := []interface{}{hget, setID, service}

	resp, errResponse := postJSONRequest(instructions)
	if errResponse != nil {
		return false, errResponse
	}
	defer resp.Body.Close()

	serviceReturned, errServiceReturned := parseCachedString(resp)
	if errServiceReturned != nil {
		return false, errServiceReturned
	}

	return *serviceReturned == trueAsStr, nil
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

func parseAndSetSkeletonKeys(serverName string, path string, err error) error {
	if err != nil {
		return err
	}

	skeletonKeys, errSkeletonKeys := parseSkeletonKeysByFilepath(path)
	if errSkeletonKeys != nil {
		return errSkeletonKeys
	}
	if skeletonKeys == nil {
		return errSkeletonKeysAreNil
	}

	for username, details := range *skeletonKeys {
		setKeySuccess, errSetKey := setSkeletonKey(serverName, username, details.Password)
		if errSetKey != errSetKey {
			return errSetKey
		}
		if !setKeySuccess {
			return errSetKeyUnsuccessful
		}

		for _, service := range details.Services {
			setServiceSuccess, errSetService := setSkeletonKeyService(serverName, username, service)
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

func VerifySkeletonKey(serverName string, username string, password string) (bool, error) {
	setID := getCacheSetID(serverName, SaltedPasswordHash)
	instructions := []interface{}{hget, setID, username}

	resp, errResponse := postJSONRequest(instructions)
	if errResponse != nil {
		return false, errResponse
	}
	defer resp.Body.Close()

	hashResultsDecoded, errHashResultsDecoded := parseCachedString(resp)
	if errHashResultsDecoded != nil {
		return false, errHashResultsDecoded
	}
	if len(*hashResultsDecoded) == 0 {
		return false, errSkeletonKeyDoesNotExist
	}

	var hashResults passwordx.HashResults
	errHashResults := json.Unmarshal([]byte(*hashResultsDecoded), &hashResults)
	if errHashResults != nil {
		return false, errHashResults
	}

	return passwordx.PasswordIsValid(password, &hashResults)
}

func VerifySkeletonKeyAndService(
	serverName string,
	service string,
	username string,
	password string,
) (bool, error) {
	skeletonKeyHasService, errSkeletonKeyService := getSkeletonKeyService(serverName, service)
	if errSkeletonKeyService != nil {
		return false, errSkeletonKeyService
	}
	if !skeletonKeyHasService {
		return false, errSkeletonKeyServiceDoesNotExist
	}

	serviceIsAvailable, errAvailableServices := getAvailableService(serverName, service)
	if errAvailableServices != nil {
		return false, errAvailableServices
	}
	if !serviceIsAvailable {
		return false, errAvailableServiceDoesNotExist
	}

	return VerifySkeletonKey(serverName, username, password)
}

func SetupSkeletonKeysAndAvailableServices(
	serverName string,
	availableServicesPath string,
	skeletonKeysPath string,
) error {
	errPaseAvailableServices := parseAndSetAvailableServices(serverName, availableServicesPath, nil)
	errParseSkeletonKeys := parseAndSetSkeletonKeys(serverName, skeletonKeysPath, errPaseAvailableServices)

	return errParseSkeletonKeys
}
