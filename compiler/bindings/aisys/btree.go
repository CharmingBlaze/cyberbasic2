package aisys

import (
	"encoding/json"
	"strings"
)

type btreeStatus string

const (
	btreeSuccess btreeStatus = "success"
	btreeFailure btreeStatus = "failure"
	btreeRunning btreeStatus = "running"
)

// TickJSON evaluates a minimal data-driven behavior tree (sequence, selector, invert, always_success, always_fail).
// Leaf nodes that would call game code are not supported; this integrates Phase 12 as a testable runner skeleton.
func TickJSON(jsonStr string) (string, error) {
	var root map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &root); err != nil {
		return "", err
	}
	return string(evalNode(root)), nil
}

func childMaps(n map[string]interface{}) []map[string]interface{} {
	raw, ok := n["children"].([]interface{})
	if !ok {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(raw))
	for _, x := range raw {
		if m, ok := x.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

func evalNode(n map[string]interface{}) btreeStatus {
	if n == nil {
		return btreeSuccess
	}
	t, _ := n["type"].(string)
	t = strings.ToLower(strings.TrimSpace(t))
	switch t {
	case "always_success", "success", "succeed":
		return btreeSuccess
	case "always_fail", "failure", "fail":
		return btreeFailure
	case "running":
		return btreeRunning
	case "sequence", "seq":
		for _, c := range childMaps(n) {
			st := evalNode(c)
			if st != btreeSuccess {
				return st
			}
		}
		return btreeSuccess
	case "selector", "fallback", "sel":
		for _, c := range childMaps(n) {
			st := evalNode(c)
			if st != btreeFailure {
				return st
			}
		}
		return btreeFailure
	case "invert", "inverter":
		ch := childMaps(n)
		if len(ch) == 0 {
			return btreeSuccess
		}
		switch evalNode(ch[0]) {
		case btreeSuccess:
			return btreeFailure
		case btreeFailure:
			return btreeSuccess
		default:
			return btreeRunning
		}
	default:
		return btreeFailure
	}
}
