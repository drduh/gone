package settings

import "fmt"

// Validate returns an error if a setting is determined to be invalid.
func (s *Settings) Validate() error {
	if s.Port < 1 || s.Port > 65535 {
		return fmt.Errorf("port must be 1..65535, got %d",
			s.Port)
	}

	if s.TimeFormat == "" {
		return fmt.Errorf("timeFormat must not be empty")
	}

	if s.Basic.Field != "" && s.Basic.Token == "" {
		return fmt.Errorf("token must be set with basic auth")
	}

	if s.Downloads < 1 {
		return fmt.Errorf("downloads must be >1, got %d",
			s.Downloads)
	}

	if s.MaxDownloads < 1 {
		return fmt.Errorf("downloads limit must be >1, got %d",
			s.MaxDownloads)
	}

	if s.MaxSizeMsg < 1 || s.MaxSizeName < 1 || s.MaxSizeWall < 1 {
		return fmt.Errorf("max content limits must be >1")
	}

	if s.MaxSizeFileMb < 1 {
		return fmt.Errorf("file size limit must be >1, got %d",
			s.MaxSizeFileMb)
	}

	if s.ReqsPerMinute < 1 {
		return fmt.Errorf("reqsPerMinute must be >1, got %d",
			s.ReqsPerMinute)
	}

	if s.Expiration.GetDuration() < 1 {
		return fmt.Errorf("expiration duration must be >1, got %s",
			s.Expiration.Duration)
	}

	if s.MaxDuration.GetDuration() < 1 {
		return fmt.Errorf("duration limit must be >1, got %s",
			s.MaxDuration.Duration)
	}

	return nil
}
