package settings

import (
	"errors"
	"fmt"
	"net"
)

var (
	errServerAddr      = errors.New("server address must be valid IP")
	errServerPort      = errors.New("server port must be between 1 and 65535")
	errTimeFormatEmpty = errors.New("timeFormat must not be an empty string")
	errTarpitDelay     = errors.New("tarpit delay must be 0 or more")
	errDownloads       = errors.New("downloads must be 1 or more")
	errDownloadsLimit  = errors.New("downloads limit must be 1 or more")
	errContentLimits   = errors.New("max content limits must be 1 or more")
	errFileSizeLimit   = errors.New("file size limit must be 1 or more")
	errTotalSizeLimit  = errors.New("total files size limit must be 1 or more")
	errReqsPerMinute   = errors.New("reqsPerMinute must be 1 or more")
	errExpiration      = errors.New("expiration duration must be 1 or more")
	errDurationLimit   = errors.New("duration limit must be 1 or more")
	errTimeoutIdle     = errors.New("idle timeout must be 1 or more")
	errTimeoutRead     = errors.New("read timeout must be 0 or more")
	errTimeoutHeader   = errors.New("readHeader timeout must be 1 or more")
	errTimeoutWrite    = errors.New("write timeout must be 0 or more")
)

// Validate returns an error if a setting is invalid.
func (s *Settings) Validate() error {
	if err := s.validateServer(); err != nil {
		return err
	}

	if err := s.validateAuth(); err != nil {
		return err
	}

	if err := s.validateContentLimits(); err != nil {
		return err
	}

	if err := s.validateTimeouts(); err != nil {
		return err
	}

	if s.TimeFormat == "" {
		return errTimeFormatEmpty
	}

	if s.Downloads < 1 {
		return fmt.Errorf("%w - not %d",
			errDownloads, s.Downloads)
	}

	if s.ReqsPerMinute < 1 {
		return fmt.Errorf("%w - not %d",
			errReqsPerMinute, s.ReqsPerMinute)
	}

	if s.Expiration.GetDuration() < 1 {
		return fmt.Errorf("%w - not %s",
			errExpiration, s.Expiration.Duration)
	}

	return nil
}

func (s *Settings) validateServer() error {
	if s.ServerAddr != "" && net.ParseIP(s.ServerAddr) == nil {
		return fmt.Errorf("%w - not %s",
			errServerAddr, s.ServerAddr)
	}

	if s.ServerPort < 1 || s.ServerPort > 65535 {
		return fmt.Errorf("%w - not %d",
			errServerPort, s.ServerPort)
	}

	return nil
}

func (s *Settings) validateAuth() error {
	if s.TarpitDelay.GetDuration() < 0 {
		return fmt.Errorf("%w - not %s",
			errTarpitDelay, s.TarpitDelay.String())
	}

	return nil
}

func (s *Settings) validateContentLimits() error {
	if s.FileLimits.MaxDownloads < 1 {
		return fmt.Errorf("%w - not %d",
			errDownloadsLimit, s.FileLimits.MaxDownloads)
	}

	if s.FileLimits.MaxDuration.GetDuration() < 1 {
		return fmt.Errorf("%w - not %s",
			errDurationLimit, s.FileLimits.MaxDuration.Duration)
	}

	if s.FileLimits.NameLength < 1 ||
		s.MessageLimits.LengthChars < 1 ||
		s.WallLimits.LengthChars < 1 {
		return errContentLimits
	}

	if s.FileLimits.SizeEachMb < 1 {
		return fmt.Errorf("%w - not %d",
			errFileSizeLimit, s.FileLimits.SizeEachMb)
	}

	if s.FileLimits.SizeTotalMb < 1 {
		return fmt.Errorf("%w - not %d",
			errTotalSizeLimit, s.FileLimits.SizeTotalMb)
	}

	return nil
}

func (s *Settings) validateTimeouts() error {
	if s.Timeout.Idle.GetDuration() < 1 {
		return fmt.Errorf("%w - not %s",
			errTimeoutIdle, s.Timeout.Idle.String())
	}

	if s.Timeout.Read.GetDuration() < 0 {
		return fmt.Errorf("%w - not %s",
			errTimeoutRead, s.Timeout.Read.String())
	}

	if s.Timeout.ReadHeader.GetDuration() < 1 {
		return fmt.Errorf("%w - not %s",
			errTimeoutHeader, s.Timeout.ReadHeader.String())
	}

	if s.Timeout.Write.GetDuration() < 0 {
		return fmt.Errorf("%w - not %s",
			errTimeoutWrite, s.Timeout.Write.String())
	}

	return nil
}
