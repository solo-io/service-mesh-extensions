package registry

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/solo-io/go-utils/vfsutils"
	"github.com/spf13/afero"

	"github.com/pkg/errors"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/protoutils"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
	"go.uber.org/zap"
)

type SpecReader interface {
	GetSpecs() ([]*v1.ApplicationSpec, error)
}

const (
	githubToken         = "GITHUB_TOKEN"
	specFilename        = "spec.yaml"
	descriptionFilename = "description.md"
)

var (
	FailedToDownloadAppSpecsError = func(err error) error {
		return errors.Wrap(err, "Failed to download application specs")
	}

	FailedToParseApplicationSpecsError = func(err error) error {
		return errors.Wrap(err, "Failed to parse application specs")
	}

	FailedToGetSpecsFromGithubError = func(err error) error {
		return errors.Wrap(err, "Failed to get application specs from github")
	}

	FailedToGetLocalSpecsError = func(err error) error {
		return errors.Wrap(err, "Failed to get local application specs")
	}
)

type RemoteSpecReader struct {
	ctx context.Context
	url string
}

func (r *RemoteSpecReader) GetSpecs() ([]*v1.ApplicationSpec, error) {
	specYaml, err := r.getBytesYaml()
	if err != nil {
		wrapped := FailedToDownloadAppSpecsError(err)
		return nil, wrapped
	}

	specs, err := getSpecsFromBytes(specYaml)
	if err != nil {
		wrapped := FailedToParseApplicationSpecsError(err)
		return nil, wrapped
	}

	return specs, nil
}

func (r *RemoteSpecReader) getBytesYaml() ([]byte, error) {
	response, err := http.Get(r.url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

var _ SpecReader = &RemoteSpecReader{}

func NewRemoteSpecReader(ctx context.Context, url string) *RemoteSpecReader {
	contextutils.LoggerFrom(ctx).Infow("Initializing reader for remote application spec registry",
		zap.String("url", url))

	return &RemoteSpecReader{
		ctx: ctx,
		url: url,
	}
}

func getSpecsFromBytes(bytes []byte) ([]*v1.ApplicationSpec, error) {
	var specsMessage v1.ApplicationSpecs
	if err := protoutils.UnmarshalYaml(bytes, &specsMessage); err != nil {
		return nil, err
	}
	return specsMessage.Specs, nil
}

type GithubSpecReader struct {
	ctx      context.Context
	location v1.GithubRepositoryLocation
}

func (r *GithubSpecReader) GetSpecs() ([]*v1.ApplicationSpec, error) {
	contextutils.LoggerFrom(r.ctx).Infow("getting all application specs from github directory",
		zap.Any("location", r.location))

	// Initialize github client
	var client *github.Client
	token, found := os.LookupEnv(githubToken)
	if !found {
		client = github.NewClient(nil)
		contextutils.LoggerFrom(r.ctx).Warnw(fmt.Sprintf("Could not find %s in environment. The hub will fail to load applications from private registries, and may be rate limited.", githubToken))
	} else {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(r.ctx, ts)
		client = github.NewClient(tc)
	}

	// Mount github repo
	fs := afero.NewMemMapFs()
	codeDir, err := vfsutils.MountCode(fs, r.ctx, client, r.location.Org, r.location.Repo, r.location.Ref)
	if err != nil {
		wrapped := FailedToGetSpecsFromGithubError(err)
		contextutils.LoggerFrom(r.ctx).Errorw(wrapped.Error(), zap.Error(err))
		return nil, wrapped
	}
	defer fs.Remove(codeDir)

	// Read spec parent directory
	specParent := filepath.Join(codeDir, r.location.Directory)
	subdirs, err := afero.ReadDir(fs, specParent)
	if err != nil {
		wrapped := FailedToGetSpecsFromGithubError(err)
		contextutils.LoggerFrom(r.ctx).Errorw(wrapped.Error(), zap.Error(err))
		return nil, wrapped
	}

	return getSpecsFromDirectory(r.ctx, fs, subdirs, specParent)
}

var _ SpecReader = &GithubSpecReader{}

func NewGithubSpecReader(ctx context.Context, location v1.GithubRepositoryLocation) *GithubSpecReader {
	contextutils.LoggerFrom(ctx).Infow("Initializing reader for github spec registry",
		zap.Any("location", location))

	return &GithubSpecReader{
		ctx:      ctx,
		location: location,
	}
}

type LocalSpecReader struct {
	ctx  context.Context
	path string
}

func (r *LocalSpecReader) GetSpecs() ([]*v1.ApplicationSpec, error) {
	fs := afero.NewOsFs()
	subdirs, err := afero.ReadDir(fs, r.path)
	if err != nil {
		wrapped := FailedToGetLocalSpecsError(err)
		contextutils.LoggerFrom(r.ctx).Errorw(wrapped.Error(), zap.Error(err))
		return nil, wrapped
	}

	return getSpecsFromDirectory(r.ctx, fs, subdirs, r.path)
}

var _ SpecReader = &LocalSpecReader{}

func NewLocalSpecReader(ctx context.Context, path string) *LocalSpecReader {
	return &LocalSpecReader{
		ctx:  ctx,
		path: path,
	}
}

func getSpecsFromDirectory(ctx context.Context, fs afero.Fs, subdirs []os.FileInfo, specParent string) ([]*v1.ApplicationSpec, error) {
	// Create an application spec for every subdirectory
	var specs []*v1.ApplicationSpec
	for _, subdir := range subdirs {
		specPath := filepath.Join(specParent, subdir.Name(), specFilename)
		specBytes, err := afero.ReadFile(fs, specPath)
		if err != nil {
			contextutils.LoggerFrom(ctx).Errorw("Failed to read spec file", zap.Error(err), zap.String("file", specPath))
			continue
		}

		spec := &v1.ApplicationSpec{}
		if err := protoutils.UnmarshalYamlAllowUnknown(specBytes, spec); err != nil {
			contextutils.LoggerFrom(ctx).Errorw("Failed to unmarshal spec file", zap.Error(err), zap.String("file", specPath))
			continue
		}

		// If provided, render description.md to html and override the inline long description.
		// Else, render the long description to html as if it were markdown to simplify rendering on web.
		descriptionPath := filepath.Join(specParent, subdir.Name(), descriptionFilename)
		descriptionBytes, err := afero.ReadFile(fs, descriptionPath)
		var renderedBytes []byte
		if err != nil {
			info := fmt.Sprintf("%v not loaded for %v, falling back to inline long description", descriptionFilename, subdir.Name())
			contextutils.LoggerFrom(ctx).Infow(info,
				zap.Error(err),
				zap.String("file", specPath))
			renderer := blackfriday.HtmlRenderer(0, subdir.Name(), "")
			extensions := blackfriday.EXTENSION_HARD_LINE_BREAK
			renderedBytes = blackfriday.Markdown([]byte(spec.LongDescription), renderer, extensions)
		} else {
			options := blackfriday.HTML_HREF_TARGET_BLANK
			renderer := blackfriday.HtmlRenderer(options, subdir.Name(), "")
			extensions := blackfriday.EXTENSION_FENCED_CODE
			renderedBytes = blackfriday.Markdown(descriptionBytes, renderer, extensions)
		}

		// Sanitize the rendered bytes.
		sanitized := bluemonday.UGCPolicy().SanitizeBytes(renderedBytes)
		spec.LongDescription = string(sanitized)

		specs = append(specs, spec)
	}

	return specs, nil
}
