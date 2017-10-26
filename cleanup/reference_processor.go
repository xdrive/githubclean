package cleanup

import (
	"log"
	"time"

	"sync"

	"github.com/xdrive/githubclean/github_api"
	"fmt"
)

type ReferencesProcessor struct {
	repoService *github_api.RepositoryService
	refService *github_api.ReferenceService
}

func NewReferencesProcessor(repoService *github_api.RepositoryService, refSrevice *github_api.ReferenceService) *ReferencesProcessor {
	return &ReferencesProcessor{
		repoService,
		refSrevice,
	}
}

// Requests github api (at concurrency level 10) for references commits.
// As a result it will return all references which have commit date older
// than filter.OlderThan value. If filter.OlderThan is not specified
// no requests are made and all references are returned.
func (rp *ReferencesProcessor) FilterReference(refs []*github_api.Reference, filter *FilterOptions) []*github_api.Reference {
	if filter == nil || filter.OlderThan == (time.Time{}) {

		return refs
	}

	var wg sync.WaitGroup
	jobs := make(chan *github_api.Reference, len(refs))
	results := make(chan *github_api.Reference, len(refs))
	for w := 1; w <= 10; w++ {
		wg.Add(1)
		go rp.getCommitWorker(jobs, results, &wg)
	}

	for k, ref := range refs {
		jobs <- ref

		// TODO: for testing only
		if k >= 5 {
			break
		}
	}
	close(jobs)

	wg.Wait()
	close(results)
	var oldRefs []*github_api.Reference
	for ref := range results {
		if ref.CommitDate.Before(filter.OlderThan) {
			oldRefs = append(oldRefs, ref)
		}
	}

	return oldRefs
}

func (rp *ReferencesProcessor) DeleteReferences(refs []*github_api.Reference) {
	var wg sync.WaitGroup
	jobs := make(chan *github_api.Reference, len(refs))
	for w := 1; w <= 10; w++ {
		wg.Add(1)
		go rp.refCleanupWorker(jobs, &wg)
	}

	for _, ref := range refs {
		jobs <- ref
	}
	close(jobs)

	wg.Wait()
}

func (rp *ReferencesProcessor) deleteReference(reference *github_api.Reference) error {
	fmt.Println(reference.Ref)
	return nil
}

func (rp *ReferencesProcessor) getCommitWorker(jobs <-chan *github_api.Reference, results chan<- *github_api.Reference, wg *sync.WaitGroup) {
	defer wg.Done()
	for ref := range jobs {
		commit, err := rp.repoService.GetCommit(ref.Object.Sha)
		if err != nil {
			log.Printf("could not get commit for ref %s: %v", ref.URL, err)
			continue
		}

		ref.CommitDate = commit.Commit.Author.Date
		results <- ref
	}
}

func (rp *ReferencesProcessor) refCleanupWorker(jobs <-chan *github_api.Reference, wg *sync.WaitGroup) {
	defer wg.Done()
	for ref := range jobs {
		err := rp.deleteReference(ref)
		if err != nil {
			log.Printf("could not get commit for ref %s: %v", ref.URL, err)
			continue
		}
	}
}

type FilterOptions struct {
	OlderThan time.Time
}

func FindRefsWithoutPulls(refs []*github_api.Reference, pulls []*github_api.PullRequest) []*github_api.Reference {
	var found, initRefsLen int
	initRefsLen = len(refs)
	for _, pull := range pulls {
		for i, ref := range refs {
			if ref.Object.Sha == pull.Head.Sha {
				//fmt.Println("Found match", pull.Head.Label)
				refs = append(refs[:i], refs[i+1:]...)
				found++
				break
			}

			//fmt.Println("!!!!! No match:", ref.Ref)
		}
	}

	fmt.Printf("Found pull requests: %d, Refs without pull requests: %d, Total # of refs: %d\n", found, len(refs), initRefsLen)

	return refs
}

