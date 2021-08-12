package skeletonkeyslx

import (
	"fmt"
	"os"
	"testing"
)

const (
	testIncrement  = "INCR"
	testIncrSet    = "test_increment_test"
	testExecSetStr = "test_exec_set_gateway"

	skeletonKeysTest = "test_skeleton_keys_test"
	testService      = "test_service"
	testUser         = "test_user"
	testFailUser     = "test_failuser"
	testPassword     = "test_password"
	testBadPassword  = "test_badpassword"
)

var (
	localCacheAddress = os.Getenv("LOCAL_CACHE_ADDRESS")
	// localCacheAddress = "http://10.88.0.1:1234"

	testAvailableServicesPath = "./available_services.json.example"
	testSkeletonKeysPath      = "./skeleton_keys.json.example"
)

func TestExecGetInstructionsAndParseString(t *testing.T) {
	setInstructions := []interface{}{setCache, testExecSetStr, "hello_world"}
	execInstructionsAndParseString(localCacheAddress, &setInstructions)

	getInstructions := []interface{}{getCache, testExecSetStr}
	parsedStr, errParsedStr := execInstructionsAndParseString(
		localCacheAddress,
		&getInstructions,
	)
	if errParsedStr != nil {
		t.Fail()
		t.Logf(errParsedStr.Error())
	}
	if parsedStr == nil {
		t.Fail()
		t.Logf("increment was not successfuul")
	}
	if parsedStr != nil && *parsedStr == "" {
		t.Fail()
		t.Logf("integer less than one returned")
	}
}

func TestSetAvailableService(t *testing.T) {
	setSuccessful, errSetSuccessful := setAvailableService(
		localCacheAddress,
		skeletonKeysTest,
		testService,
	)
	if !setSuccessful {
		t.Fail()
		t.Logf(fmt.Sprint("set service was not successfuul"))
	}
	if errSetSuccessful != nil {
		t.Fail()
		t.Logf(errSetSuccessful.Error())
	}
}

func TestGetAvailableService(t *testing.T) {
	getSuccessful, errGetSuccessful := getAvailableService(
		localCacheAddress,
		skeletonKeysTest,
		testService,
	)
	if !getSuccessful {
		t.Fail()
		t.Logf(fmt.Sprint("get service was not successfuul"))
	}
	if errGetSuccessful != nil {
		t.Fail()
		t.Logf(errGetSuccessful.Error())
	}
}

func TestParseAvailableService(t *testing.T) {
	availableServices, errAvailableServices := parseAvailableServicesByFilepath(
		testAvailableServicesPath,
	)
	if availableServices == nil {
		t.Fail()
		t.Logf(fmt.Sprint("parsing resulted in nil available services"))
	}
	if errAvailableServices != nil {
		t.Fail()
		t.Logf(errAvailableServices.Error())
	}
}

func TestParseAndSetAvailableServices(t *testing.T) {
	errParseAndSet := ParseAndSetAvailableServices(
		localCacheAddress,
		skeletonKeysTest,
		testAvailableServicesPath,
	)
	if errParseAndSet != nil {
		t.Fail()
		t.Logf(errParseAndSet.Error())
	}
}

func TestSetSkeletonKey(t *testing.T) {
	setSuccessful, errSetSuccessful := setSkeletonKey(
		localCacheAddress,
		skeletonKeysTest,
		testUser,
		testPassword,
	)
	if !setSuccessful {
		t.Fail()
		t.Logf(fmt.Sprint("set skeleton key was not successfuul"))
	}
	if errSetSuccessful != nil {
		t.Fail()
		t.Logf(errSetSuccessful.Error())
	}
}

func TestSetSkeletonKeyService(t *testing.T) {
	setSuccessful, errSetSuccessful := setSkeletonKeyService(
		localCacheAddress,
		skeletonKeysTest,
		testUser,
		testService,
	)
	if !setSuccessful {
		t.Fail()
		t.Logf(fmt.Sprint("set service was not successfuul"))
	}
	if errSetSuccessful != nil {
		t.Fail()
		t.Logf(errSetSuccessful.Error())
	}
}

func TestVerifySkeletonKey(t *testing.T) {
	verified, errVerified := VerifySkeletonKey(
		localCacheAddress,
		skeletonKeysTest,
		testUser,
		testPassword,
	)
	if !verified {
		t.Fail()
		t.Logf(fmt.Sprint("verify skeleton key was not successfuul"))
	}
	if errVerified != nil {
		t.Fail()
		t.Logf(errVerified.Error())
	}
}

func TestBadSkeletonKeyWontVerify(t *testing.T) {
	verified, errVerified := VerifySkeletonKey(
		localCacheAddress,
		skeletonKeysTest,
		testFailUser,
		testPassword,
	)
	if verified {
		t.Fail()
		t.Logf(fmt.Sprint("unidentified skeleton key was successfuul"))
	}
	if errVerified == nil {
		t.Fail()
		t.Logf(fmt.Sprint("an error should have been returned on bad skeleton key"))
	}
}

func TestBadSkeletonKeyPasswordWontVerify(t *testing.T) {
	verified, errVerified := VerifySkeletonKey(
		localCacheAddress,
		skeletonKeysTest,
		testUser,
		testBadPassword,
	)
	if verified {
		t.Fail()
		t.Logf(fmt.Sprint("badd password successfuul"))
	}
	if errVerified != nil {
		t.Fail()
		t.Logf(errVerified.Error())
	}
}

func TestGetSkeletonKeyService(t *testing.T) {
	getSuccessful, errGetSuccessful := getSkeletonKeyService(
		localCacheAddress,
		skeletonKeysTest,
		testService,
	)
	if !getSuccessful {
		t.Fail()
		t.Logf(fmt.Sprint("set service was not successfuul"))
	}
	if errGetSuccessful != nil {
		t.Fail()
		t.Logf(errGetSuccessful.Error())
	}
}

func TestParseSkeletonKeys(t *testing.T) {
	skeletonKeys, errParseSkeletonKeys := parseSkeletonKeysByFilepath(
		testSkeletonKeysPath,
	)
	if skeletonKeys == nil {
		t.Fail()
		t.Logf(fmt.Sprint("parsing resulted in nil skeleton keys"))
	}
	if errParseSkeletonKeys != nil {
		t.Fail()
		t.Logf(errParseSkeletonKeys.Error())
	}
}

func TestParseAndSetSkeletonKeys(t *testing.T) {
	errParseAndSet := ParseAndSetSkeletonKeys(
		localCacheAddress,
		skeletonKeysTest,
		testSkeletonKeysPath,
	)
	if errParseAndSet != nil {
		t.Fail()
		t.Logf(errParseAndSet.Error())
	}
}

func TestVerifySkeletonKeyAndService(t *testing.T) {
	verified, errVerified := VerifySkeletonKeyAndService(
		localCacheAddress,
		skeletonKeysTest,
		testService,
		testUser,
		testPassword,
	)
	if !verified {
		t.Fail()
		t.Logf(fmt.Sprint("verify skeleton key was not successfuul"))
	}
	if errVerified != nil {
		t.Fail()
		t.Logf(errVerified.Error())
	}
}

func TestSetupSkeletonKeysAndAvailableServices(t *testing.T) {
	errSetupSkeletonKeys := SetupSkeletonKeysAndAvailableServices(
		localCacheAddress,
		skeletonKeysTest,
		testAvailableServicesPath,
		testSkeletonKeysPath,
	)
	if errSetupSkeletonKeys != nil {
		t.Fail()
		t.Logf(errSetupSkeletonKeys.Error())
	}
}
