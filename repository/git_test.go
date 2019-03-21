package repository

import (
	"os/exec"
	"path"
	"strings"
	"testing"
)

const (
	emptyRepo  = "https://github.com/yusufRahmatullah/gamesave_empty_test.git"
	normalRepo = "https://github.com/yusufRahmatullah/gamesave_test.git"
	wrongRepo  = "https://a:a@github.com/yusufRahmatullah/wrong_and_inexist.git"
)

func TestCheckout(t *testing.T) {
	t.Run("checkout on existing condition", func(t *testing.T) {
		gitRepo := GitRepository{}
		ensureCloned(t, normalRepo)
		err := gitRepo.Checkout("game_1")
		assertNotError(t, err)
		currentBranch := gitCurrentBranchName(t)
		assertEqual(t, currentBranch, "game_1")
	})

	t.Run("checkout on new branch", func(t *testing.T) {
		gitRepo := GitRepository{}
		ensureCloned(t, normalRepo)
		ensureOnBranch(t, "master")
		// delete temp branch
		cmd := exec.Command("git", "branch", "-d", "deleted_game")
		cmd.Run()
		err := gitRepo.Checkout("deleted_game")
		assertNotError(t, err)
		currentBranch := gitCurrentBranchName(t)
		assertEqual(t, currentBranch, "deleted_game")
		// re-delete temp branch
		cmd = exec.Command("git", "branch", "-d", "deleted_game")
		cmd.Run()
	})

	t.Run("checkout on empty repo", func(t *testing.T) {
		gitRepo := GitRepository{}
		ensureCloned(t, emptyRepo)
		err := gitRepo.Checkout("game_1")
		assertNotError(t, err)
	})

	t.Run("checkout repo not set", func(t *testing.T) {
		deleteLocalRepo(t)
		gitRepo := GitRepository{}
		err := gitRepo.Checkout("game_1")
		assertError(t, err)
	})
}

func TestCommit(t *testing.T) {
	t.Run("commit on normal condition", func(t *testing.T) {
		gitRepo := GitRepository{}
		ensureCloned(t, normalRepo)
		createDummyFile(t, path.Join(GameSaveRoot, "new_game.save"))
		commitMsg := "Add dummy file"
		err := gitRepo.Commit(commitMsg)
		assertNotError(t, err)
		// check last commit
		cmd := exec.Command("git", "log", "-1", "--pretty=%B")
		cmd.Dir = GameSaveRoot
		output, err := cmd.Output()
		if err != nil {
			t.Errorf("[Helper-TestCommit] Error: %v", err)
		}
		assertEqual(t, strings.TrimSpace(string(output)), commitMsg)
	})

	t.Run("commit repo not set", func(t *testing.T) {
		gitRepo := GitRepository{}
		cleanLocalRepo(t)
		cmd := exec.Command("mkdir", GameSaveRoot)
		err := cmd.Run()
		if err != nil {
			t.Errorf("[Helper-TestCommit] Error: %v", err)
		}
		createDummyFile(t, path.Join(GameSaveRoot, "new_game.save"))
		commitMsg := "Add dummy file"
		err = gitRepo.Commit(commitMsg)
		assertError(t, err)
	})

	t.Run("commit without changes", func(t *testing.T) {
		gitRepo := GitRepository{}
		ensureCloned(t, normalRepo)
		commitMsg := "Add dummy file"
		err := gitRepo.Commit(commitMsg)
		assertError(t, err)
	})
}

func TestClone(t *testing.T) {
	t.Run("clone on normal condition", func(t *testing.T) {
		gitRepo := GitRepository{}
		cleanLocalRepo(t)
		err := gitRepo.Clone(normalRepo)
		assertNotError(t, err)
		assertRemoteSame(t, normalRepo)
	})

	t.Run("clone on existing repo", func(t *testing.T) {
		gitRepo := GitRepository{}
		ensureCloned(t, normalRepo)
		err := gitRepo.Clone(normalRepo)
		assertError(t, err)
		assertRemoteSame(t, normalRepo)
	})

	t.Run("clone on wrong URL", func(t *testing.T) {
		gitRepo := GitRepository{}
		cleanLocalRepo(t)
		err := gitRepo.Clone(wrongRepo)
		assertError(t, err)
	})
}

func TestFetchBranch(t *testing.T) {
	t.Run("fetch correct branch", func(t *testing.T) {
		gitRepo := GitRepository{}
		ensureCloned(t, normalRepo)
		err := gitRepo.FetchBranch("game_1")
		assertNotError(t, err)
		// able to checkout
		cmd := exec.Command("git", "checkout", "game_1")
		cmd.Dir = GameSaveRoot
		err = cmd.Run()
		assertNotError(t, err)
	})

	t.Run("fetch existing branch", func(t *testing.T) {
		gitRepo := GitRepository{}
		ensureCloned(t, normalRepo)
		err := gitRepo.FetchBranch("master")
		assertError(t, err)
	})

	t.Run("fetch inexists branch", func(t *testing.T) {
		gitRepo := GitRepository{}
		ensureCloned(t, normalRepo)
		err := gitRepo.FetchBranch("wrong_branch")
		assertError(t, err)
	})
}

func TestGetCurrentBranch(t *testing.T) {
	t.Run("get branch on normal condition", func(t *testing.T) {
		gitRepo := GitRepository{}
		ensureCloned(t, normalRepo)
		branch, err := gitRepo.GetCurrentBranch()
		assertNotError(t, err)
		assertEqual(t, branch, gitCurrentBranchName(t))
	})

	t.Run("get branch repo not set", func(t *testing.T) {
		gitRepo := GitRepository{}
		cleanLocalRepo(t)
		branch, err := gitRepo.GetCurrentBranch()
		assertError(t, err)
		assertEqual(t, branch, "")
	})
}

func TestGetRepoURL(t *testing.T) {
	t.Run("get repo url on normal condition", func(t *testing.T) {
		gitRepo := GitRepository{}
		ensureCloned(t, normalRepo)
		repo, err := gitRepo.GetRepoURL()
		assertNotError(t, err)
		assertEqual(t, repo, gitCurrentRepoURL(t))
	})

	t.Run("get repo url not set", func(t *testing.T) {
		gitRepo := GitRepository{}
		cleanLocalRepo(t)
		repo, err := gitRepo.GetRepoURL()
		assertError(t, err)
		assertEqual(t, repo, "")
	})
}

func TestPull(t *testing.T) {
	t.Run("pull on normal condition", func(t *testing.T) {
		gitRepo := GitRepository{}
		cleanLocalRepo(t)
		cmd := exec.Command("mkdir", GameSaveRoot)
		err := cmd.Run()
		if err != nil {
			t.Errorf("[Helper-TestPull] Error: %v", err)
		}
		gitAddRepoURL(t, normalRepo)
		err = gitRepo.Pull("game_1")
		assertNotError(t, err)
		assertEqual(t, gitCurrentBranchName(t), "game_1")
	})

	t.Run("pull repo not set", func(t *testing.T) {
		gitRepo := GitRepository{}
		cleanLocalRepo(t)
		cmd := exec.Command("mkdir", GameSaveRoot)
		err := cmd.Run()
		if err != nil {
			t.Errorf("[Helper-TestPull] Error: %v", err)
		}
		err = gitRepo.Pull("game_1")
		assertError(t, err)
	})

	t.Run("pull wrong branch", func(t *testing.T) {
		gitRepo := GitRepository{}
		cleanLocalRepo(t)
		cmd := exec.Command("mkdir", GameSaveRoot)
		err := cmd.Run()
		if err != nil {
			t.Errorf("[Helper-TestPull] Error: %v", err)
		}
		gitAddRepoURL(t, normalRepo)
		err = gitRepo.Pull("wrong_branch")
		assertError(t, err)
	})
}

func TestPush(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	t.Run("push on normal condition", func(t *testing.T) {
		gitRepo := GitRepository{}
		ensureCloned(t, normalRepo)
		gitDeleteRemoteBranch(t, "game_push")
		createDummyFile(t, path.Join(GameSaveRoot, "new_game.save"))
		gitAddAndCommit(t)
		err := gitRepo.Push("game_push")
		assertNotError(t, err)
		cmd := exec.Command("git", "show-branch", "remotes/origin/game_push")
		cmd.Dir = GameSaveRoot
		err = cmd.Run()
		if err != nil {
			t.Errorf("[Helper-TestPush] Error: %v", err)
		}
	})

	t.Run("push repo not set", func(t *testing.T) {
		gitRepo := GitRepository{}
		ensureCloned(t, normalRepo)
		createDummyFile(t, path.Join(GameSaveRoot, "new_game.save"))
		gitAddAndCommit(t)
		cmd := exec.Command("git", "remote", "remove", "origin")
		cmd.Dir = GameSaveRoot
		err := cmd.Run()
		if err != nil {
			t.Errorf("[Helper-TestPush] Error: %v", err)
		}
		err = gitRepo.Push("game_push")
		assertError(t, err)
	})
}

func TestSetRepoURL(t *testing.T) {
	t.Run("set repo url on normal condition", func(t *testing.T) {
		gitRepo := GitRepository{}
		ensureCloned(t, normalRepo)
		err := gitRepo.SetRepoURL(emptyRepo)
		assertNotError(t, err)
		assertRemoteSame(t, emptyRepo)
	})

	t.Run("set empty repo url", func(t *testing.T) {
		gitRepo := GitRepository{}
		cleanLocalRepo(t)
		err := gitRepo.SetRepoURL(normalRepo)
		assertError(t, err) // .git inexsist
	})
}
