package setup

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
)

const polkitRulePath = "/etc/polkit-1/rules.d/80-deckjoy.rules"

// InstallPolkitRuleAsRoot writes the polkit rule directly. Must be called from
// a process already running as root (e.g. the daemon).
func InstallPolkitRuleAsRoot() error {
	exePath, userName, err := getRuleInputs()
	if err != nil {
		return err
	}

	ruleContent := buildPolkitRule(exePath, userName)

	if existing, err := os.ReadFile(polkitRulePath); err == nil {
		if string(existing) == ruleContent {
			return nil
		}
	}

	return os.WriteFile(polkitRulePath, []byte(ruleContent), 0644)
}

func getRuleInputs() (exePath string, userName string, err error) {
	exePath, err = os.Readlink("/proc/self/exe")
	if err != nil {
		return "", "", err
	}

	userName, err = getExpectedSubjectUser()
	if err != nil {
		return "", "", err
	}

	return exePath, userName, nil
}

func getExpectedSubjectUser() (string, error) {
	// Default to deck for SteamOS environments.
	userName := "deck"

	// When called from the root daemon (started via pkexec), current user is root.
	// In that case, prefer the original caller from PKEXEC_UID/SUDO_USER.
	if os.Geteuid() == 0 {
		if pkexecUID := os.Getenv("PKEXEC_UID"); pkexecUID != "" {
			if pkexecUser, lookupErr := user.LookupId(pkexecUID); lookupErr == nil && pkexecUser.Username != "" {
				userName = pkexecUser.Username
				return userName, nil
			}
		}
		if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
			userName = sudoUser
			return userName, nil
		}
		return userName, nil
	}

	if currentUser, currentUserErr := user.Current(); currentUserErr == nil && currentUser.Username != "" {
		userName = currentUser.Username
	}

	return userName, nil
}

func buildPolkitRule(exePath string, userName string) string {
	exePathQuoted := strconv.Quote(exePath)
	userNameQuoted := strconv.Quote(userName)

	return strings.TrimSpace(fmt.Sprintf(`
// Allows DeckJoy GUI to start its own daemon with pkexec without repeated password prompts.
polkit.addRule(function(action, subject) {
    if (action.id == "org.freedesktop.policykit.exec" &&
        action.lookup("program") == %s &&
        subject.user == %s &&
        subject.local &&
        subject.active) {
        return polkit.Result.YES;
    }
});
`, exePathQuoted, userNameQuoted)) + "\n"
}
