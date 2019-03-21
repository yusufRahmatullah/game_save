package repository

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

func addLocalConfig(t *testing.T, key, value string) {
	t.Helper()
	var config map[string]string
	data, err := ioutil.ReadFile(LocalConfig)
	err = json.Unmarshal(data, &config)
	if err != nil {
		t.Errorf("[Helper-addLocalConfig] error unmarshalling: %v", err)
		return
	}
	config[key] = value
	byt, err := json.MarshalIndent(config, "", "  ")
	err = ioutil.WriteFile(LocalConfig, byt, 0644)
	if err != nil {
		t.Errorf("[Helper-addLocalConfig] error write file: %v", err)
	}
}

func assertEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("Got '%s' expect '%s'", got, want)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("Should be thrown an error")
	}
}

func assertExist(t *testing.T, path string) {
	t.Helper()
	cmd := exec.Command("ls", path)
	err := cmd.Run()
	if err != nil {
		t.Errorf("'%s' is not exist", path)
	}
}

func assertNotError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Should be not error. Error: %v", err)
	}
}

func assertRemoteSame(t *testing.T, repoURL string) {
	cmd := exec.Command("git", "remote", "-v")
	cmd.Dir = GameSaveRoot
	output, err := cmd.Output()
	if err != nil {
		t.Errorf("Should not throw Error: %v", err)
	} else if !strings.Contains(string(output), repoURL) {
		t.Errorf(
			"Different remote, got: '%s', want: '%s'",
			strings.TrimSpace(string(output)),
			repoURL,
		)
	}
}

func assertSameContent(t *testing.T, p1, p2 string) {
	t.Helper()
	cmd := exec.Command("cmp", p1, p2)
	output, _ := cmd.CombinedOutput()
	if string(output) != "" {
		t.Error("File are different")
	}
}

func cleanLocalRepo(t *testing.T) {
	cmd := exec.Command("rm", "-rf", GameSaveRoot)
	err := cmd.Run()
	if err != nil {
		t.Errorf("[Helper-ensureClone] Error: %v", err)
	}
}

func createBlankFile(t *testing.T, path string) {
	t.Helper()
	cmd := exec.Command("echo", `""`, ">", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("[Helper-createBlankFile] error: %v, output: %s", err, string(output))
	}
}

func createDummyDirectory(t *testing.T, path string) {
	t.Helper()
	cmd := exec.Command("mkdir", "-p", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("[Helper-createDummyDirectory] error: %v, output: %s", err, string(output))
	}
}

func createDummyFile(t *testing.T, path string) {
	t.Helper()
	ctn := []byte("this is dummy file\n")
	err := ioutil.WriteFile(path, ctn, 0644)
	if err != nil {
		t.Errorf("[Helper-createDummyFile] Error: %v", err)
	}
}

func deleteLocalRepo(t *testing.T) {
	t.Helper()
	// delete local repo
	cmd := exec.Command("rm", "-rf", GameSaveRoot)
	err := cmd.Run()
	if err != nil {
		t.Errorf("[Helper-deleteLocalRepo] Error: %v", err)
	}
}

func ensureOnBranch(t *testing.T, branchName string) {
	t.Helper()
	currentBranch := gitCurrentBranchName(t)
	if currentBranch != branchName {
		cmd := exec.Command("git", "checkout", "-B", branchName)
		err := cmd.Run()
		if err != nil {
			t.Errorf("[Helper-ensureBranch] Error: %v", err)
		}
	}
}

func ensureCloned(t *testing.T, repoURL string) {
	t.Helper()
	cleanLocalRepo(t)
	shadowClone(t, repoURL)
	assertRemoteSame(t, repoURL)
}

func getLocalConfig(t *testing.T, key string) string {
	t.Helper()
	var config map[string]string
	data, err := ioutil.ReadFile(LocalConfig)
	err = json.Unmarshal(data, &config)
	if err != nil {
		t.Errorf("[Helper-getLocalConfig] error: %v", err)
		return ""
	}
	if val, ok := config[key]; ok {
		return val
	}
	return ""
}

func gitAddRepoURL(t *testing.T, repoURL string) {
	t.Helper()
	cmd := exec.Command("git", "init")
	cmd.Dir = GameSaveRoot
	err := cmd.Run()
	if err != nil {
		t.Errorf("[Helper-gitAddRepoURL] Error: %v", err)
	}
	cmd = exec.Command("git", "remote", "add", "origin", repoURL)
	cmd.Dir = GameSaveRoot
	err = cmd.Run()
	if err != nil {
		t.Errorf("[Helper-gitAddRepoURL] Error: %v", err)
	}
}

func gitAddAndCommit(t *testing.T) {
	t.Helper()
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = GameSaveRoot
	err := cmd.Run()
	if err != nil {
		t.Errorf("[Helper-TestPush] Error: %v", err)
	}
	cmd = exec.Command("git", "commit", "-m", "add dummy file")
	cmd.Dir = GameSaveRoot
	err = cmd.Run()
	if err != nil {
		t.Errorf("[Helper-TestPush] Error: %v", err)
	}
}

func gitCurrentBranchName(t *testing.T) string {
	t.Helper()
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = GameSaveRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("[Helper-gitCurrentBranch] Error: %v", err)
	}
	return strings.TrimSpace(string(output))
}

func gitCurrentRepoURL(t *testing.T) string {
	t.Helper()
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = GameSaveRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("[Helper-gitRepoURL] Error: %v", err)
	}
	return strings.TrimSpace(string(output))
}

func gitDeleteRemoteBranch(t *testing.T, branch string) {
	t.Helper()
	// checkout to master
	cmd := exec.Command("git", "checkout", "-B", "master")
	cmd.Dir = GameSaveRoot
	cmd.Run() // branch may not exist
	cmd = exec.Command("git", "branch", "-d", branch)
	cmd.Dir = GameSaveRoot
	cmd.Run() // branch may not exist
	cmd = exec.Command("git", "push", "origin", ":"+branch)
	cmd.Dir = GameSaveRoot
	cmd.Run() // branch may not exist
}

func initLocalConfig(t *testing.T) {
	t.Helper()
	err := ioutil.WriteFile(LocalConfig, []byte("{}"), 0644)
	if err != nil {
		t.Errorf("[Helper-initLocalConfig] error: %v", err)
	}
}

func removeDummyFile(t *testing.T, path string) {
	t.Helper()
	cmd := exec.Command("rm", path)
	err := cmd.Run()
	if err != nil {
		t.Errorf("[Helper-removeDummyFile] Error: %v", err)
	}
}

func removeLocalConfig(t *testing.T) {
	t.Helper()
	cmd := exec.Command("rm", LocalConfig)
	cmd.Run() // LocalConfig may uninitialized
}

func shadowClone(t *testing.T, repoURL string) {
	t.Helper()
	parts := strings.Split(repoURL, "/")
	repoName := parts[len(parts)-1]
	repoName = "." + repoName[0:len(repoName)-4]
	repoPath := path.Join(path.Dir(GameSaveRoot), repoName)
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", repoURL, repoPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("[Helper-shadowClone] Error cloning: %v, output: %v", err, string(output))
		}
	}
	cmd := exec.Command("cp", "-rf", repoPath, GameSaveRoot)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("[Helper-shadowClone] Error copying: %v, output: %v", err, string(output))
	}

}
