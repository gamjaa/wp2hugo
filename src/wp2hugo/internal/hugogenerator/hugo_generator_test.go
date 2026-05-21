package hugogenerator

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/ashishb/wp2hugo/src/wp2hugo/internal/wpparser"
	"github.com/mmcdole/gofeed/rss"
	"github.com/stretchr/testify/require"
)

func TestFootnote(t *testing.T) {
	t.Parallel()
	file, err := os.Open("./testdata/testcase.WordPress.2024-07-01.xml")
	require.NoError(t, err)

	parser := wpparser.NewParser()
	websiteInfo, err := parser.Parse(file, nil, nil)
	require.NoError(t, err)
	require.Len(t, websiteInfo.Posts(), 1)

	post := websiteInfo.Posts()[0]
	require.Len(t, post.Footnotes, 1)
	require.NotNil(t, post.CommonFields)

	generator := NewGenerator("/tmp", "", nil, false, false, false, true, *websiteInfo, Options{})
	url1, err := url.Parse(post.GUID.Value)
	require.NoError(t, err)
	hugoPage, err := generator.newHugoPage(url1, post.CommonFields)
	require.NoError(t, err)

	const expectedMarkdown = "Some text[^1] with a footnote\n\n[^1]: Here we are: the footnote."
	require.Contains(t, hugoPage.Markdown(), expectedMarkdown)
}

func TestPost(t *testing.T) {
	t.Parallel()
	file, err := os.Open("./testdata/testcase.WordPress_2.xml")
	require.NoError(t, err)

	parser := wpparser.NewParser()
	websiteInfo, err := parser.Parse(file, nil, nil)
	require.NoError(t, err)
	require.Len(t, websiteInfo.Posts(), 1)

	post := websiteInfo.Posts()[0]
	require.NotNil(t, post.CommonFields)
	require.Equal(t, "Kurz angemerkt zum Tag der Schachtels��tze", post.Title)
	require.Len(t, post.Categories, 1)
	require.Equal(t, "netzfundst��cke", post.Categories[0])
	require.Len(t, post.Content, 1276)
}

func TestWritePageCreatesParentDirectory(t *testing.T) {
	t.Parallel()

	info := wpparser.WebsiteInfo{}
	generator := NewGenerator(t.TempDir(), "", nil, false, false, false, false, info, Options{
		MigrateComments: false,
	})

	bundleDir := postBundleDirName(wpparser.PostInfo{
		CommonFields: wpparser.CommonFields{Title: "중년탐정 김정일을 보고나서..."},
	})
	require.Equal(t, "undated 중년탐정 김정일을 보고나서", bundleDir)

	pagePath := filepath.Join(t.TempDir(), "posts", bundleDir, "index.md")
	page := wpparser.CommonFields{
		PostID:        "1",
		Title:         "중년탐정 김정일을 보고나서...",
		Link:          "https://example.net/posts/2009-08-06",
		PublishStatus: wpparser.PublishStatusPublish,
		GUID:          &rss.GUID{Value: "https://example.net/?p=1"},
		Content:       "test content",
	}

	require.NoError(t, generator.writePage(context.Background(), t.TempDir(), pagePath, page, info))
	require.FileExists(t, pagePath)
}
