package settings

import (
	"errors"
	"fmt"
	"iter"
	"net"
	"reflect"
	"regexp"
	"strings"
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
	errTimeoutShutdown = errors.New("shutdown timeout must be 0 or more")
	errPathSlash       = errors.New("path must start with /")
	errPathRoot        = errors.New("path / is reserved for root")
	errPathUnique      = errors.New("path must be unique")
	errPathSegment     = errors.New("path segments must be non-empty letters, digits or ._~-")
	errPathCount       = errors.New("path must have one {count} segment")
	errPathPrefix      = errors.New("path must end with /")
)

// countSegment is the message number wildcard.
const countSegment = "{count}"

var pathSegment = regexp.MustCompile(`^[A-Za-z0-9._~-]+$`)

// Validate returns an error if a setting is invalid.
func (s *Settings) Validate() error {
	if err := s.validateServer(); err != nil {
		return err
	}

	if err := s.validateAuth(); err != nil {
		return err
	}

	if err := s.validatePaths(); err != nil {
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

// validatePaths rejects invalid mux paths.
func (s *Settings) validatePaths() error {
	seen := make(map[string]string)
	for name, path := range s.named() {
		if path == "" {
			continue
		}

		if !strings.HasPrefix(path, "/") {
			return fmt.Errorf("%w - not %s %q",
				errPathSlash, name, path)
		}

		if err := pathShape(name, path); err != nil {
			return err
		}

		if path == "/" && name != "root" {
			return fmt.Errorf("%w - not %s", errPathRoot, name)
		}

		if other, ok := seen[path]; ok {
			return fmt.Errorf("%w - %s and %s are both %q",
				errPathUnique, other, name, path)
		}
		seen[path] = name
	}

	return nil
}

// pathShape returns an error for invalid mux path segments.
func pathShape(name, path string) error {
	counts := 0
	segs := strings.Split(path[1:], "/")
	for i, seg := range segs {
		switch {
		case seg == "" && i == len(segs)-1:
		case seg == countSegment &&
			(name == "msgGet" || name == "msgDelete"):
			counts++
		case seg == "." || seg == ".." || !pathSegment.MatchString(seg):
			return fmt.Errorf("%w - not %s %q",
				errPathSegment, name, path)
		}
	}

	switch name {
	case "msgGet", "msgDelete":
		if counts != 1 {
			return fmt.Errorf("%w - not %s %q",
				errPathCount, name, path)
		}
	case "assets", "download", "random", "root":
		if !strings.HasSuffix(path, "/") {
			return fmt.Errorf("%w - not %s %q",
				errPathPrefix, name, path)
		}
	}

	return nil
}

// named returns each path with its setting name and the JSON
// key it is configured by.
func (p *Paths) named() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for f, v := range reflect.ValueOf(p).Elem().Fields() {
			name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
			if !yield(name, v.String()) {
				return
			}
		}
	}
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
	if s.Idle.GetDuration() < 1 {
		return fmt.Errorf("%w - not %s",
			errTimeoutIdle, s.Idle.String())
	}

	if s.Read.GetDuration() < 0 {
		return fmt.Errorf("%w - not %s",
			errTimeoutRead, s.Read.String())
	}

	if s.ReadHeader.GetDuration() < 1 {
		return fmt.Errorf("%w - not %s",
			errTimeoutHeader, s.ReadHeader.String())
	}

	if s.Write.GetDuration() < 0 {
		return fmt.Errorf("%w - not %s",
			errTimeoutWrite, s.Write.String())
	}

	if s.Shutdown.GetDuration() < 0 {
		return fmt.Errorf("%w - not %s",
			errTimeoutShutdown, s.Shutdown.String())
	}

	return nil
}
