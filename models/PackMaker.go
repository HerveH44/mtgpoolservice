package models

import (
	"fmt"
	wr "mtgpoolservice/weighted"
)

func (s *Set) GetDefaultBoosterRule() (*BoosterRule, error) {
	for _, rule := range s.Booster {
		if rule.Name == "default" {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("packmaker: did not find any default booster rule for set %s", s.Code)
}

func (r *BoosterRule) GetRandomConfiguration() (*PackConfiguration, error) {
	configurations := r.Boosters
	if len(configurations) == 0 {
		return nil, fmt.Errorf("packmaker: Did not find any booster rule for %s-%s", r.SetID, r.Name)
	}

	choices := make([]wr.Choice, 0)
	for _, conf := range configurations {
		choices = append(choices, wr.Choice{
			Item:   conf,
			Weight: conf.Weight,
		})
	}

	chooser := wr.NewChooser(choices...)
	pick := chooser.Pick().(PackConfiguration)

	return &pick, nil
}
