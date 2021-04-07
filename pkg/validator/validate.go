package validator

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/microlib/logger/pkg/multi"
)

// checkEnvars - private function, iterates through each item and checks the required field
func checkEnvar(item string, logger *multi.Logger) error {
	name := strings.Split(item, ",")[0]
	required, _ := strconv.ParseBool(strings.Split(item, ",")[1])
	logger.Trace(fmt.Sprintf("Input paramaters -> name %s : required %t", name, required))
	if os.Getenv(name) == "" {
		if required {
			logger.Errorf("%s envar is mandatory please set it", name)
			return errors.New(fmt.Sprintf("%s envar is mandatory please set it", name))
		}

		logger.Errorf(fmt.Sprintf("%s envar is empty please set it", name))
	}
	return nil
}

// ValidateEnvars : public call that groups all envar validations
// These envars are set via the openshift template
func ValidateEnvars(logger *multi.Logger) error {
	items := []string{
		"LOG_LEVEL,false",
		"NAME,false",
		"SERVER_PORT,true",
		"VERSION,true",
		"CICD_CONSOLE_DIR,true",
		"SONAR_TOKEN,true",
		"SONAR_URL,true",
		"SONARSCAN_PATH,true",
		"WORKDIR,true",
		"WEBHOOK_SECRET,true",
	}
	for x := range items {
		if err := checkEnvar(items[x], logger); err != nil {
			return err
		}
	}
	return nil
}
