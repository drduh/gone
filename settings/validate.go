package settings

import (
	"errors"
	"fmt"
	"net"
)

var (
	errInvalidServerAddr     = errors.New("server address must be valid IP")
	errInvalidServerPort     = errors.New("server port must be 1..65535")
	errTimeFormatEmpty       = errors.New("timeFormat must not be empty")
	errMissingBasicToken     = errors.New("token must be set with basic auth")
	errInvalidDownloads      = errors.New("downloads must be >1")
	errInvalidDownloadsLimit = errors.New("downloads limit must be >1")
	errInvalidContentLimits  = errors.New("max content limits must be >1")
	errInvalidFileSizeLimit  = errors.New("file size limit must be >1")
	errInvalidTotalSizeLimit = errors.New("total files size limit must be >1")
	errInvalidReqsPerMinute  = errors.New("reqsPerMinute must be >1")
	errInvalidExpiration     = errors.New("expiration duration must be >1")
	errInvalidDurationLimit  = errors.New("duration limit must be >1")
)

// Validate returns an error if a setting is determined to be invalid.
func (s *Settings) Validate() error {
	if s.ServerAddr != "" && net.ParseIP(s.ServerAddr) == nil {
		return fmt.Errorf("%w: %s",
			errInvalidServerAddr, s.ServerAddr)
	}

	if s.ServerPort < 1 || s.ServerPort > 65535 {
		return fmt.Errorf("%w: %d",
			errInvalidServerPort, s.ServerPort)
	}

	if s.TimeFormat == "" {
		return errTimeFormatEmpty
	}

	if s.Basic.Field != "" && s.Basic.Token == "" {
		return errMissingBasicToken
	}

	if s.Downloads < 1 {
		return fmt.Errorf("%w: %d",
			errInvalidDownloads, s.Downloads)
	}

	if s.FileLimits.MaxDownloads < 1 {
		return fmt.Errorf("%w: %d",
			errInvalidDownloadsLimit, s.FileLimits.MaxDownloads)
	}

	if s.FileLimits.NameLength < 1 ||
		s.MessageLimits.LengthChars < 1 ||
		s.WallLimits.LengthChars < 1 {
		return errInvalidContentLimits
	}

	if s.FileLimits.SizeEachMb < 1 {
		return fmt.Errorf("%w: %d",
			errInvalidFileSizeLimit, s.FileLimits.SizeEachMb)
	}

	if s.FileLimits.SizeTotalMb < 1 {
		return fmt.Errorf("%w: %d",
			errInvalidTotalSizeLimit, s.FileLimits.SizeTotalMb)
	}

	if s.ReqsPerMinute < 1 {
		return fmt.Errorf("%w: %d",
			errInvalidReqsPerMinute, s.ReqsPerMinute)
	}

	if s.Expiration.GetDuration() < 1 {
		return fmt.Errorf("%w: %s",
			errInvalidExpiration, s.Expiration.Duration)
	}

	if s.FileLimits.MaxDuration.GetDuration() < 1 {
		return fmt.Errorf("%w: %s",
			errInvalidDurationLimit, s.FileLimits.MaxDuration.Duration)
	}

	return nil
}
