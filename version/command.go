// Package version the version command.
package version

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/alecthomas/kong"
)

// Run check that a given version matches a semantic version constraint.
func (o Options) Run(ctx *kong.Context) error {
	c, err := semver.NewConstraint(o.Constraint)
	if err != nil {
		return fmt.Errorf("could not parse range %s: %w", o.Constraint, err)
	}
	current := ctx.Model.Vars()["versionNumber"]
	v, err := semver.NewVersion(current)
	if err != nil {
		return fmt.Errorf("could not parse version %s: %w", current, err)
	}
	if !c.Check(v) {
		return fmt.Errorf("gum version %q is not within given range %q", current, o.Constraint)
	}
	return nil
}
