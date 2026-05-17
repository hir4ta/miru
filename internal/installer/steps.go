package installer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var errCouldNotDetectShell = errors.New("could not detect shell rc file")

// installBinary copies the running executable to cfg.InstallDir/miru.
// If the running binary already lives at the destination, it's a no-op.
func installBinary(cfg Config) (string, error) {
	src, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("locate executable: %w", err)
	}
	srcAbs, err := filepath.Abs(src)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(cfg.InstallDir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", cfg.InstallDir, err)
	}
	dst := filepath.Join(cfg.InstallDir, "miru")
	dstAbs, _ := filepath.Abs(dst)
	if srcAbs == dstAbs {
		return dst, nil
	}
	if err := copyFile(srcAbs, dst); err != nil {
		return "", fmt.Errorf("copy: %w", err)
	}
	if err := os.Chmod(dst, 0o755); err != nil {
		return "", fmt.Errorf("chmod: %w", err)
	}
	return dst, nil
}

func copyFile(src, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	return err
}

type pathAction int

const (
	pathUpdated pathAction = iota
	pathAlreadyInPath
	pathAlreadyConfigured
)

func configurePath(cfg Config) (string, pathAction, error) {
	if isDirInPATH(cfg.InstallDir) {
		return "", pathAlreadyInPath, nil
	}
	rc, line, err := detectShellRC(cfg.InstallDir)
	if err != nil {
		return "", 0, err
	}
	if rc == "" {
		return "", 0, errCouldNotDetectShell
	}
	if has, _ := rcContainsLine(rc, line); has {
		return rc, pathAlreadyConfigured, nil
	}
	if err := appendRC(rc, line); err != nil {
		return "", 0, err
	}
	return rc, pathUpdated, nil
}

func isDirInPATH(dir string) bool {
	for _, p := range filepath.SplitList(os.Getenv("PATH")) {
		if p == dir {
			return true
		}
	}
	return false
}

func detectShellRC(installDir string) (string, string, error) {
	shell := filepath.Base(os.Getenv("SHELL"))
	home, _ := os.UserHomeDir()
	switch shell {
	case "zsh":
		rc := os.Getenv("ZDOTDIR")
		if rc == "" {
			rc = home
		}
		return filepath.Join(rc, ".zshrc"), fmt.Sprintf("export PATH=%q:$PATH", installDir), nil
	case "bash":
		var rc string
		if runtime.GOOS == "darwin" {
			rc = filepath.Join(home, ".bash_profile")
		} else {
			rc = filepath.Join(home, ".bashrc")
		}
		return rc, fmt.Sprintf("export PATH=%q:$PATH", installDir), nil
	case "fish":
		cfgHome := os.Getenv("XDG_CONFIG_HOME")
		if cfgHome == "" {
			cfgHome = filepath.Join(home, ".config")
		}
		return filepath.Join(cfgHome, "fish", "config.fish"), fmt.Sprintf("fish_add_path %q", installDir), nil
	}
	return "", "", nil
}

func rcContainsLine(rc, line string) (bool, error) {
	data, err := os.ReadFile(rc)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	want := strings.TrimSpace(line)
	for _, l := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(l) == want {
			return true, nil
		}
	}
	return false, nil
}

func appendRC(rc, line string) error {
	if err := os.MkdirAll(filepath.Dir(rc), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(rc, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "\n# added by miru installer\n%s\n", line)
	return err
}

func verifyInstall(cfg Config) error {
	dst := filepath.Join(cfg.InstallDir, "miru")
	st, err := os.Stat(dst)
	if err != nil {
		return err
	}
	if st.Mode()&0o111 == 0 {
		return errors.New("not executable")
	}
	return nil
}

func friendlyPath(p string) string {
	home, err := os.UserHomeDir()
	if err == nil && home != "" && strings.HasPrefix(p, home) {
		return "~" + strings.TrimPrefix(p, home)
	}
	return p
}
