package domain

import "errors"

type Platform int

const (
	Pc          = Platform(iota)
	Mobile      = Platform(iota)
	PlayStation = Platform(iota)
	Xbox        = Platform(iota)
	Nintendo    = Platform(iota)
)

var (
	ErrInvalidPlatform = errors.New("platform_from_string: invalid platform")
)

func PlatformFromString(platform string) (Platform, error) {
	switch platform {
	case "pc":
		return Pc, nil
	case "Mobile":
		return Mobile, nil
	case "PlayStation":
		return PlayStation, nil
	case "Xbox":
		return Xbox, nil
	case "Nintendo":
		return Nintendo, nil
	default:
		return Platform(0), ErrInvalidPlatform
	}
}

func PlatformToString(platform Platform) string {
	switch platform {
	case Pc:
		return "pc"
	case Mobile:
		return "Mobile"
	case PlayStation:
		return "PlayStation"
	case Xbox:
		return "Xbox"
	case Nintendo:
		return "Nintendo"
	default:
		return ""
	}
}
