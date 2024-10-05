package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mike-plivo/ipfilter"
)

func getRedisURL() string {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}
	return redisURL
}

func main() {
	redisAddr := getRedisURL()
	filter := ipfilter.NewIPFilter(redisAddr)

	// Example: Add rules
	rules := []ipfilter.Rule{
		{Action: "allow", Target: "1.1.1.1"},
		{Action: "deny", Target: "1.1.1.2"},
		{Action: "allow", Target: "1.1.1.3"},
		{Action: "allow", Target: "2.2.2.0/24"},
		{Action: "allow", Target: "2001:db8::/32"},
		{Action: "deny", Target: "2001:db8:1234::/48"},
		{Action: "deny", Target: "all"},
	}

	// Add rules and handle errors
	for _, rule := range rules {
		addedRule, err := filter.AppendRule(rule)
		if err != nil {
			log.Printf("Error adding rule %v: %v\n", rule, err)
		} else {
			log.Printf("Added rule: %+v\n", addedRule)
		}
	}

	// Example: Test IPs
	testIPs := []string{
		"1.1.1.1", "1.1.1.2", "1.1.1.3", "2.2.2.10",
		"3.3.3.3", "2001:db8::1", "2001:db8:1234::1", "2001:db8:5678::1",
	}

	// Test IPs and print results
	for _, ip := range testIPs {
		allowed, err := filter.IsAllowedIP(ip)
		if err != nil {
			log.Printf("Error checking IP %s: %v\n", ip, err)
			continue
		}
		fmt.Printf("IP %s is %s\n", ip, map[bool]string{true: "allowed", false: "denied"}[allowed])
	}
}
