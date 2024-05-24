package domain

import "errors"

type Genre int

const (
	Action      = Genre(iota) // 0
	Adventure   = Genre(iota)
	Rpg         = Genre(iota)
	Strategy    = Genre(iota)
	Simulation  = Genre(iota)
	Puzzle      = Genre(iota)
	Sports      = Genre(iota)
	Racing      = Genre(iota)
	Fighting    = Genre(iota)
	Horror      = Genre(iota)
	Sandbox     = Genre(iota)
	Platformer  = Genre(iota)
	Shooter     = Genre(iota)
	Moba        = Genre(iota)
	Stealth     = Genre(iota)
	Survival    = Genre(iota)
	Educational = Genre(iota)
	Party       = Genre(iota)
	Casual      = Genre(iota)
)

var (
	ErrInvalidGenre = errors.New("genre_from_string: invalid genre")
)

func GenreFromString(game string) (Genre, error) {
	switch game {
	case "action":
		return Action, nil
	case "adventure":
		return Adventure, nil
	case "rpg":
		return Rpg, nil
	case "strategy":
		return Strategy, nil
	case "simulation":
		return Simulation, nil
	case "puzzle":
		return Puzzle, nil
	case "sports":
		return Sports, nil
	case "racing":
		return Racing, nil
	case "fighting":
		return Fighting, nil
	case "horror":
		return Horror, nil
	case "sandbox":
		return Sandbox, nil
	case "platformer":
		return Platformer, nil
	case "shooter":
		return Shooter, nil
	case "moba":
		return Moba, nil
	case "stealth":
		return Stealth, nil
	case "survival":
		return Survival, nil
	case "educational":
		return Educational, nil
	case "party":
		return Party, nil
	case "casual":
		return Casual, nil
	default:
		return Genre(0), ErrInvalidGenre
	}
}

func GenreToString(genre Genre) string {
	switch genre {
	case Action:
		return "action"
	case Adventure:
		return "adventure"
	case Rpg:
		return "rpg"
	case Strategy:
		return "strategy"
	case Simulation:
		return "simulation"
	case Puzzle:
		return "puzzle"
	case Sports:
		return "sports"
	case Racing:
		return "racing"
	case Fighting:
		return "fighting"
	case Horror:
		return "horror"
	case Sandbox:
		return "sandbox"
	case Platformer:
		return "platformer"
	case Shooter:
		return "shooter"
	case Moba:
		return "moba"
	case Stealth:
		return "stealth"
	case Survival:
		return "survival"
	case Educational:
		return "educational"
	case Party:
		return "party"
	case Casual:
		return "casual"
	default:
		return ""
	}
}
