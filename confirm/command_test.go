package confirm

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"testing"
	"time"
)

func TestOptions_errIsValidTimeout(t *testing.T) {
	type testCase struct {
		Description    string
		GivenTimeout   time.Duration
		GivenError     error
		ExpectedResult bool
	}

	cases := []testCase{
		{
			Description:    "timeout is positive, err is a timeout",
			GivenTimeout:   time.Second,
			GivenError:     huh.ErrTimeout,
			ExpectedResult: true,
		},
		{
			Description:    "timeout is zero, err is a timeout",
			GivenTimeout:   0,
			GivenError:     huh.ErrTimeout,
			ExpectedResult: false,
		},
		{
			Description:    "timeout is positive, err is not a timeout",
			GivenTimeout:   1,
			GivenError:     fmt.Errorf("i'm not a timeout"),
			ExpectedResult: false,
		},
		{
			Description:    "timeout is zero, err is not a timeout",
			GivenTimeout:   0,
			GivenError:     fmt.Errorf("i'm not a timeout"),
			ExpectedResult: false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.Description, func(t *testing.T) {
			sut := Options{Timeout: testCase.GivenTimeout}
			actualResult := sut.errIsValidTimeout(testCase.GivenError)
			if actualResult != testCase.ExpectedResult {
				t.Errorf("got: %v, want: %v", actualResult, testCase.ExpectedResult)
			}
		})
	}
}
