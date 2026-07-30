package authz

import (
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

const modelText = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`

// sub = user_id or role name (admin / accountant / viewer)
// obj = resource, optionally prefixed with tenant: "company:<id>:journal-entries"
// act = create|read|update|delete|approve|post|export
var (
	enforcer *casbin.Enforcer
	initOnce sync.Once
	initErr  error
)

func Init() error {
	initOnce.Do(func() {
		m, err := model.NewModelFromString(modelText)
		if err != nil {
			initErr = fmt.Errorf("casbin model parse: %w", err)
			return
		}
		// Use a threaded enforcer — safe for concurrent HTTP requests.
		enforcer, initErr = casbin.NewEnforcer(m)
		if initErr != nil {
			return
		}
		enforcer.EnableLog(false)
	})
	return initErr
}

func MustInit() {
	if err := Init(); err != nil {
		panic(err)
	}
}

func E() *casbin.Enforcer {
	if enforcer == nil {
		panic("authz not initialized — call authz.Init() first")
	}
	return enforcer
}

// ─── Role helpers ────────────────────────────────────────────────────────────

// GrantRole binds user → role (RBAC grouping).
func GrantRole(userID, role string) error {
	_, err := E().AddRoleForUser(userID, role)
	return err
}

// RevokeRole removes user → role binding.
func RevokeRole(userID, role string) error {
	_, err := E().DeleteRoleForUser(userID, role)
	return err
}

// SetRoles replaces all role bindings for a user.
func SetRoles(userID string, roles []string) error {
	e := E()
	if _, err := e.DeleteRolesForUser(userID); err != nil {
		return err
	}
	for _, r := range roles {
		if _, err := e.AddRoleForUser(userID, r); err != nil {
			return err
		}
	}
	return nil
}

// ListRoles returns all roles bound to user.
func ListRoles(userID string) ([]string, error) {
	return E().GetRolesForUser(userID)
}

// ─── Policy helpers ──────────────────────────────────────────────────────────

// GrantPerm adds a policy rule. sub = role name or user ID.
func GrantPerm(sub, obj, act string) error {
	_, err := E().AddPolicy(sub, obj, act)
	return err
}

// RevokePerm removes a matching policy rule.
func RevokePerm(sub, obj, act string) error {
	_, err := E().RemovePolicy(sub, obj, act)
	return err
}

// HasPerm checks: can (user_or_role) perform act on obj?
func HasPerm(sub, obj, act string) (bool, error) {
	return E().Enforce(sub, obj, act)
}

// Enforce is shorthand for HasPerm.
func Enforce(sub, obj, act string) bool {
	ok, _ := E().Enforce(sub, obj, act)
	return ok
}

// ─── Bulk / seed ─────────────────────────────────────────────────────────────

// ImportPolicies replaces all g (grouping) + p (policy) rules atomically.
func ImportPolicies(gRules, pRules [][3]string) error {
	e := E()
	e.ClearPolicy()
	for _, g := range gRules {
		if _, err := e.AddRoleForUser(g[0], g[1]); err != nil {
			return fmt.Errorf("import g %v: %w", g, err)
		}
	}
	for _, p := range pRules {
		if _, err := e.AddPolicy(p[0], p[1], p[2]); err != nil {
			return fmt.Errorf("import p %v: %w", p, err)
		}
	}
	return nil
}

// DumpPolicy returns current policy rules as readable JSON.
func DumpPolicy() ([]map[string]interface{}, error) {
	e := E()
	pp := e.GetPolicy()
	out := make([]map[string]interface{}, 0, len(pp))
	for _, r := range pp {
		if len(r) == 3 {
			out = append(out, map[string]interface{}{
				"sub": r[0], "obj": r[1], "act": r[2],
			})
		}
	}
	return out, nil
}
