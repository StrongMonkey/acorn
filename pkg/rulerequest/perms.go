package rulerequest

import (
	"strings"

	v1 "github.com/acorn-io/acorn/pkg/apis/internal.acorn.io/v1"
)

type RuleRequest struct {
	Service      string
	Scope        string
	Verbs        string
	Resource     string
	ResourceName string
}

func ToRuleRequests(perms []v1.Permissions) (result []RuleRequest) {
	for _, perm := range perms {
		result = append(result, rulesToRequests(perm.ServiceName, perm.GetRules())...)
	}
	return
}

func rulesToRequests(serviceName string, rules []v1.PolicyRule) (result []RuleRequest) {
	for _, rule := range rules {
		if rule.IsAccountScoped() {
			result = append(result, ruleToRequests(serviceName, rule, "account")...)
		}
		if rule.IsProjectScoped() {
			result = append(result, ruleToRequests(serviceName, rule, "project")...)
		}
		namespaces := rule.Namespaces()
		if len(namespaces) > 0 {
			result = append(result, ruleToRequests(serviceName, rule, "namespaces:"+strings.Join(namespaces, ","))...)
		}
	}
	return
}

func ruleToRequests(serviceName string, rule v1.PolicyRule, scope string) (result []RuleRequest) {
	verbs := strings.Join(rule.Verbs, ",")

	if len(rule.NonResourceURLs) > 0 {
		for _, url := range rule.NonResourceURLs {
			result = append(result, RuleRequest{
				Scope:    "cluster",
				Resource: url,
				Verbs:    verbs,
			})
		}
		return
	}

	for _, apiGroup := range rule.APIGroups {
		for _, resource := range rule.Resources {
			if apiGroup != "" {
				resource += "." + apiGroup
			}

			if len(rule.ResourceNames) == 0 {
				result = append(result, RuleRequest{
					Service:  serviceName,
					Scope:    scope,
					Resource: resource,
					Verbs:    verbs,
				})
			} else {
				for _, resourceName := range rule.ResourceNames {
					result = append(result, RuleRequest{
						Service:  serviceName,
						Scope:    scope,
						Resource: resource + "/" + resourceName,
						Verbs:    verbs,
					})
				}
			}
		}
	}

	return
}
