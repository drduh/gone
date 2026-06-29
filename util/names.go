package util

import (
	"bufio"
	"io"
	"os"
	"strings"
)

var adjectives = []string{
	"amber", "ancient", "autumn",
	"bright", "breezy", "blissful",
	"clear", "cool", "crimson",
	"dewy", "dappled", "dreamy",
	"earthy", "emerald", "evening",
	"fresh", "feathered", "frosted",
	"gentle", "golden", "glowing",
	"hazy", "hollow", "hushed",
	"icy", "indigo", "ivory",
	"jolly", "jade", "joyful",
	"kind", "keen", "kissed",
	"lucid", "lilting", "luminous",
	"misty", "merry", "moonlit",
	"narrow", "nimble", "noble",
	"opal", "open", "olden",
	"pale", "peaceful", "petaled",
	"quiet", "quick", "quivering",
	"rustic", "rosy", "rainy",
	"soft", "silver", "sunlit",
	"tidy", "twilight", "tranquil",
	"umber", "uplifted", "untamed",
	"verdant", "velvet", "vivid",
	"wild", "willowy", "windy",
	"xanthic", "xenial", "xylic",
	"yearning", "yellow", "yawning",
	"zesty", "zephyr", "zenith",
}

var nouns = []string{
	"acorn", "aurora", "aspen",
	"brook", "bloom", "briar",
	"cloud", "cedar", "cove",
	"dell", "dawn", "drift",
	"ember", "elm", "estuary",
	"field", "fern", "feather",
	"grove", "glade", "garden",
	"harbor", "hill", "horizon",
	"isle", "inlet", "iris",
	"juniper", "jewel", "journey",
	"knoll", "kite", "kingfisher",
	"lagoon", "leaf", "lilac",
	"meadow", "moon", "mist",
	"nest", "night", "nova",
	"oak", "oasis", "orchard",
	"pine", "pond", "prairie",
	"quail", "quartz", "quill",
	"river", "ridge", "rain",
	"stone", "stream", "star",
	"thicket", "trail", "thrush",
	"upland", "umbra", "union",
	"vale", "vine", "vista",
	"willow", "wood", "waterfall",
	"xylem", "xerus", "xenon",
	"year", "yonder", "yew",
	"zephyr", "zinnia", "zenith",
}

var loadedNames = loadNames("/etc/gone/assets", "names.txt")

// defaultNames returns the Cartesian product of adjectives
// and nouns with the first letter capitalized.
func defaultNames() []string {
	names := make([]string, 0, len(adjectives)*len(nouns))

	for _, a := range adjectives {
		for _, n := range nouns {
			names = append(names, a+upperFirst(n))
		}
	}

	return names
}

// loadNames returns trimmed names from a file or
// the default names list.
func loadNames(dir, filename string) []string {
	f, err := os.OpenInRoot(dir, filename)
	if err != nil {
		return defaultNames()
	}
	defer func() { _ = f.Close() }()

	return loadNamesFromReader(f)
}

func loadNamesFromReader(r io.Reader) []string {
	scanner := bufio.NewScanner(r)

	var names []string

	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name != "" {
			names = append(names, name)
		}
	}

	if scanner.Err() != nil {
		return defaultNames()
	}

	if len(names) == 0 {
		return defaultNames()
	}

	return names
}
