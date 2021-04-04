package common

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

const sampleTableRows = `
<html>
  <table>
  <thead data-reactid="45">
    <tr>
      <th>
        <span data-reactid="48">Holder</span>
      </th>
      <th>
        <span data-reactid="50">Shares</span>
      </th>
      <th>
        <span data-reactid="52">Date Reported</span>
      </th>
      <th>
        <span data-reactid="54">% Out</span>
      </th>
      <th>
        <span data-reactid="56">Value</span>
      </th>
    </tr>
  </thead>
  <tbody data-reactid="57">
    <tr>
      <td>FMR, LLC</td>
      <td>9,276,087</td>
      <td>
        <span data-reactid="62">Dec 30, 2020</span>
      </td>
      <td>13.26%</td>
      <td>174,761,479</td>
    </tr>
    <tr>
      <td>Blackrock Inc.</td>
      <td>9,217,335</td>
      <td>
        <span data-reactid="69">Dec 30, 2020</span>
      </td>
      <td>13.18%</td>
      <td>173,654,591</td>
    </tr>
  </tbody>
</table>
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

func TestTableRows(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rsp := sampleTableRows
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
	require.Equal(t, 3, len(tables[0].Rows))
	require.Equal(t, 5, len(tables[0].Rows[0]))
	require.Equal(t, 5, len(tables[0].Rows[1]))
	require.Equal(t, 5, len(tables[0].Rows[2]))

	require.Equal(t, []string{"Holder", "Shares", "Date Reported", "% Out", "Value"}, tables[0].Rows[0])
	require.Equal(t, []string{"FMR, LLC", "9,276,087", "Dec 30, 2020", "13.26%", "174,761,479"}, tables[0].Rows[1])
	require.Equal(t, []string{"Blackrock Inc.", "9,217,335", "Dec 30, 2020", "13.18%", "173,654,591"}, tables[0].Rows[2])
}
