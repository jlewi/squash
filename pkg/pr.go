package pkg

import (
	"fmt"
	"github.com/jlewi/hydros/pkg/github"
)

func GetPR(org string, repo string, prNumber int) {
	//// PullRequestForBranch returns the PR for the given branch if it exists and nil if no PR exists.
	//// TODO(jeremy): Can we change this to api.PullRequest?
	//
	//baseBranch := h.BaseBranch
	//headBranch := h.BranchName
	//type response struct {
	//	Repository struct {
	//		PullRequests struct {
	//			ID    githubv4.ID
	//			Nodes []PullRequest
	//		}
	//	}
	//}
	//
	//query := `
	//query($owner: String!, $Repo: String!, $headRefName: String!) {
	//	repository(owner: $owner, name: $Repo) {
	//		pullRequests(headRefName: $headRefName, states: OPEN, first: 30) {
	//			nodes {
	//				id
	//				number
	//				title
	//				state
	//				body
	//				mergeable
	//				author {
	//					login
	//				}
	//				commits {
	//					totalCount
	//				}
	//				url
	//				baseRefName
	//				headRefName
	//			}
	//		}
	//	}
	//}`
	//
	//branchWithoutOwner := headBranch
	//if idx := strings.Index(headBranch, ":"); idx >= 0 {
	//	branchWithoutOwner = headBranch[idx+1:]
	//}
	//
	//variables := map[string]interface{}{
	//	"owner":       h.baseRepo.RepoOwner(),
	//	"Repo":        h.baseRepo.RepoName(),
	//	"headRefName": branchWithoutOwner,
	//}
	//
	//var resp response
	//err := h.client.GraphQL(h.baseRepo.RepoHost(), query, variables, &resp)
	//if err != nil {
	//	return nil, err
	//}
	//
	//for _, pr := range resp.Repository.PullRequests.Nodes {
	//	h.log.Info("found", "pr", pr)
	//	if pr.HeadLabel() == headBranch {
	//		if baseBranch != "" {
	//			if pr.BaseRefName != baseBranch {
	//				continue
	//			}
	//		}
	//		return &pr, nil
	//	}
	//}
	//
	//return nil, nil
}

func GetPRBranchPoint(h *github.RepoHelper, prNumber int) error {
	pr, err := h.PullRequestForBranch()

	if err != nil {
		return err
	}

	fmt.Printf("PR: %+v", pr)
	return nil
}
