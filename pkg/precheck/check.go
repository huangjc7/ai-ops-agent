package precheck

import (
	"ai-ops-agent/internal/version"
	"ai-ops-agent/pkg/i18n"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func CheckVar() error {
	if os.Getenv("BASE_URL") == "" {
		return fmt.Errorf(i18n.T("ErrBaseUrlMissing"))
	}

	if os.Getenv("MODEL") == "" {
		return fmt.Errorf(i18n.T("ErrModelMissing"))
	}

	if os.Getenv("API_KEY") == "" {
		return fmt.Errorf(i18n.T("ErrApiKeyMissing"))
	}

	return nil
}

func CheckVersion() (bool, error) {

	versionURL := "http://huangjc.cc/ai-ops-agent/latest"
	currentVersion := version.Version

	clinet := http.Client{Timeout: time.Second * 3}
	req, err := http.NewRequest("GET", versionURL, nil)
	if err != nil {
		return false, err
	}

	resp, err := clinet.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if strings.TrimSpace(currentVersion) == strings.TrimSpace(string(body)) {
		return true, nil
	}
	return false, nil
}
