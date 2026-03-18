package builds

import (
	"context"
	"errors"
	"flag"
	"strings"
	"testing"
)

func TestBuildsCountCommand_MissingApp(t *testing.T) {
	isolateBuildsAuthEnv(t)
	t.Setenv("ASC_APP_ID", "")

	cmd := BuildsCountCommand()

	if err := cmd.FlagSet.Parse([]string{}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}

	err := cmd.Exec(context.Background(), []string{})
	if !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("expected flag.ErrHelp when --app is missing, got %v", err)
	}
}

func TestBuildsCountCommand_InvalidProcessingState(t *testing.T) {
	isolateBuildsAuthEnv(t)

	cmd := BuildsCountCommand()

	if err := cmd.FlagSet.Parse([]string{"--app", "123456789", "--processing-state", "INVALID_STATE"}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}

	err := cmd.Exec(context.Background(), []string{})
	if err == nil {
		t.Fatal("expected error for invalid --processing-state, got nil")
	}
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if !strings.Contains(err.Error(), "processing-state") {
		t.Fatalf("expected error to mention processing-state, got %v", err)
	}
}

func TestBuildsCountCommand_InvalidPlatform(t *testing.T) {
	isolateBuildsAuthEnv(t)

	cmd := BuildsCountCommand()

	if err := cmd.FlagSet.Parse([]string{"--app", "123456789", "--platform", "NOT_A_PLATFORM"}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}

	err := cmd.Exec(context.Background(), []string{})
	if err == nil {
		t.Fatal("expected error for invalid --platform, got nil")
	}
}

func TestBuildsCountCommand_UsesAppIDEnv(t *testing.T) {
	isolateBuildsAuthEnv(t)
	t.Setenv("ASC_APP_ID", "env-app-id")

	cmd := BuildsCountCommand()

	if err := cmd.FlagSet.Parse([]string{}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}

	err := cmd.Exec(context.Background(), []string{})
	if errors.Is(err, flag.ErrHelp) {
		t.Fatalf("expected ASC_APP_ID env fallback, got %v", err)
	}
}

func TestBuildsCountCommand_FlagsExist(t *testing.T) {
	cmd := BuildsCountCommand()

	required := []string{"app", "version", "build-number", "platform", "processing-state", "output"}
	for _, name := range required {
		if cmd.FlagSet.Lookup(name) == nil {
			t.Errorf("expected --%s flag to be defined", name)
		}
	}
}

func TestBuildsCountCommand_NoLimitFlag(t *testing.T) {
	cmd := BuildsCountCommand()

	if cmd.FlagSet.Lookup("limit") != nil {
		t.Error("builds count should not expose --limit (limit is set internally)")
	}
}

func TestBuildsCountCommand_NoPaginateFlag(t *testing.T) {
	cmd := BuildsCountCommand()

	if cmd.FlagSet.Lookup("paginate") != nil {
		t.Error("builds count should not expose --paginate (counting uses meta.paging.total)")
	}
}

func TestBuildsCountCommand_ShortHelp(t *testing.T) {
	cmd := BuildsCountCommand()

	if cmd.ShortHelp == "" {
		t.Error("expected non-empty ShortHelp")
	}
	if !strings.Contains(strings.ToLower(cmd.ShortHelp), "count") && !strings.Contains(strings.ToLower(cmd.ShortHelp), "total") {
		t.Errorf("expected ShortHelp to mention count or total, got %q", cmd.ShortHelp)
	}
}

func TestBuildsCountCommand_LongHelpHasExamples(t *testing.T) {
	cmd := BuildsCountCommand()

	if !strings.Contains(cmd.LongHelp, "asc builds count") {
		t.Errorf("expected LongHelp to contain example invocation, got %q", cmd.LongHelp)
	}
}
