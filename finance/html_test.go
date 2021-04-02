package finance

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const sampleHtmlResponse = `
<html>
<table>
  <tr>
    <td>A</td>
    <td>B</td>
  </tr>
  <tr>
    <td>C</td>
    <td>D</td> 
  </tr>
</table>
</html>
`

const sampleHtmlNoTable = `
<html>
</html>
`

const sampleHtmlEmptyTable = `
<html>
<table>
</table>
</html>
`

const sampleTables = `
<html>
<table>
  <tr>
    <td>One</td>
  </tr>
</table>
<table>
  <tr>
    <td>Two</td>
  </tr>
</table>
</html>
`

func TestHtmlResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rsp := sampleHtmlResponse
		w.Header()["Content-Type"] = []string{"text/html"}
		fmt.Fprintln(w, rsp)
	}))
	defer ts.Close()

	context := context.Background()
	client := &http.Client{
		Transport: &http.Transport{},
		Timeout:   time.Second,
	}

	tables, err := ReadHtml(context, client, ts.URL)
	require.Equal(t, io.EOF, err)
	require.NotNil(t, tables)
	require.Equal(t, 1, len(tables))
	require.Equal(t, 2, len(tables[0].Rows))
}

func TestNoTable(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rsp := sampleHtmlNoTable
		w.Header()["Content-Type"] = []string{"text/html"}
		fmt.Fprintln(w, rsp)
	}))
	defer ts.Close()

	context := context.Background()
	client := &http.Client{
		Transport: &http.Transport{},
		Timeout:   time.Second,
	}

	tables, err := ReadHtml(context, client, ts.URL)
	require.Equal(t, io.EOF, err)
	require.NotNil(t, tables)
	require.Empty(t, tables)
}

func TestEmptyTable(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rsp := sampleHtmlEmptyTable
		w.Header()["Content-Type"] = []string{"text/html"}
		fmt.Fprintln(w, rsp)
	}))
	defer ts.Close()

	context := context.Background()
	client := &http.Client{
		Transport: &http.Transport{},
		Timeout:   time.Second,
	}

	tables, err := ReadHtml(context, client, ts.URL)
	require.Equal(t, io.EOF, err)
	require.Equal(t, 1, len(tables))
	require.Equal(t, 0, len(tables[0].Rows))
}

func TestTables(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rsp := sampleTables
		w.Header()["Content-Type"] = []string{"text/html"}
		fmt.Fprintln(w, rsp)
	}))
	defer ts.Close()

	context := context.Background()
	client := &http.Client{
		Transport: &http.Transport{},
		Timeout:   time.Second,
	}

	tables, err := ReadHtml(context, client, ts.URL)
	require.Equal(t, io.EOF, err)
	require.NotNil(t, tables)
	require.Equal(t, 2, len(tables))
	require.Equal(t, 1, len(tables[0].Rows))
	require.Equal(t, 1, len(tables[1].Rows))
}
