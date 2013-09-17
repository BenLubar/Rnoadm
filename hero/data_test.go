package hero

import (
	"github.com/BenLubar/Rnoadm/resource"
	"regexp"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

var skinToneRegexp = regexp.MustCompile("^#[0-9a-f]{3}([0-9a-f]{3})?$")

func TestRacialData(t *testing.T) {
	if int(raceCount) != len(raceInfo) {
		t.Fatalf("UniverseImplosionException: there are %d races but there are %d races", raceCount, len(raceInfo))
	}
	names := make(map[string]Race, raceCount)
	for r := Race(0); r < raceCount; r++ {
		name := r.Name()
		if name == "" {
			t.Errorf("Race %d is missing a name", r)
		}
		firstLetter, _ := utf8.DecodeRuneInString(name)
		if !unicode.IsLetter(firstLetter) {
			t.Errorf("Race %d has name %q, which does not start with a letter", r, name)
		}
		if lower := strings.ToLower(name); lower != name {
			t.Errorf("Race %d has name %q, which is not fully lowercase (should be %q)", r, name, lower)
		}
		if _, ok := names[name]; ok {
			t.Errorf("Race %d has the same name as another race (%q)", r, name)
		}
		names[name] = r

		w, h := r.SpriteSize()
		if w < 32 || h < 32 {
			t.Errorf("Race %d has an invalid sprite size (%d by %d)", r, w, h)
		}

		health := r.BaseHealth()
		if health == nil {
			t.Errorf("Race %d has nil base health.", r)
		} else if health.Sign() <= 0 {
			t.Errorf("Race %d has %v base health.", r, health)
		}

		sprite := r.Sprite()
		if sprite == "" {
			t.Errorf("Race %d is missing a sprite.", r)
		} else if _, ok := resource.Resource[sprite+".png"]; !ok {
			t.Errorf("Race %d has a nonexistent sprite (%q)", r, sprite+".png")
		}

		genders := r.Genders()
		if len(genders) == 0 {
			t.Errorf("Race %d has no genders.", r)
		}

		skin := r.SkinTones()
		if len(skin) == 0 {
			t.Errorf("Race %d has no skin tones.", r)
		} else {
			for _, s := range skin {
				if !skinToneRegexp.MatchString(s) {
					t.Errorf("Race %d has skin invalid tone %q", r, s)
				}
			}
		}
	}
}
