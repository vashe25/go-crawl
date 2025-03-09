package crawler

import "regexp"

type rule struct {
	rule *regexp.Regexp
}

func compileRule(pattern string) (*rule, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &rule{re}, err
}

func (_this *rule) match(text string) bool {
	return _this.rule.MatchString(text)
}
