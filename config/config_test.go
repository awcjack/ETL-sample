package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/spf13/viper"
)

func TestGetStringConfigWithDefault(t *testing.T) {
	type testcase struct {
		testcase       string
		key            string
		keyEnv         string
		defaultValue   string
		expectedResult string
	}

	testcases := []testcase{
		{
			testcase:       "Normal",
			key:            "testKey",
			keyEnv:         "test",
			defaultValue:   "testDefault",
			expectedResult: "test",
		},
		{
			testcase:       "Empty env",
			key:            "testKey",
			keyEnv:         "",
			defaultValue:   "testDefault",
			expectedResult: "testDefault",
		},
		{
			testcase:       "Both empty",
			key:            "testKey",
			keyEnv:         "",
			defaultValue:   "",
			expectedResult: "",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.testcase, func(t *testing.T) {
			viper.BindEnv(tc.key, tc.key)
			os.Setenv(tc.key, tc.keyEnv)
			v := getStringConfigWithDefault(tc.key, tc.defaultValue)

			if v != tc.expectedResult {
				t.Errorf("expected %v, but got %v", tc.expectedResult, v)
			}
		})
	}
}

func TestGetIntConfigWithDefault(t *testing.T) {
	type testcase struct {
		testcase       string
		key            string
		keyEnv         int
		defaultValue   int
		expectedResult int
	}

	testcases := []testcase{
		{
			testcase:       "Normal",
			key:            "testKey",
			keyEnv:         1000,
			defaultValue:   100,
			expectedResult: 1000,
		},
		{
			testcase:       "Empty env",
			key:            "testKey",
			keyEnv:         0,
			defaultValue:   100,
			expectedResult: 100,
		},
		{
			testcase:       "Both empty",
			key:            "testKey",
			keyEnv:         0,
			defaultValue:   0,
			expectedResult: 0,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.testcase, func(t *testing.T) {
			viper.BindEnv(tc.key, tc.key)
			os.Setenv(tc.key, strconv.Itoa(tc.keyEnv))
			v := getIntConfigWithDefault(tc.key, tc.defaultValue)

			if v != tc.expectedResult {
				t.Errorf("expected %v, but got %v", tc.expectedResult, v)
			}
		})
	}
}
