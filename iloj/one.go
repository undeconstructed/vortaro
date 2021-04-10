package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var LETTERS = []rune("'abcĉdefgĝhĥijĵklmnoprsŝtuŭvz")

// OldEntry is
type OldEntry struct {
	Radiko string   `json:"radiko"`
	Fakoj  []string `json:"fakoj,omitempty"`
	En     string   `json:"en"`
	Vidu   string   `json:"vidu"`
}

// Entry is
type Entry struct {
	Radiko   string            `json:"radiko"`
	Speco    string            `json:"speco"`
	Fonto    string            `json:"fonto"`
	Bona     bool              `json:"bona"`
	Nivelo   int               `json:"nivelo"`
	Fakoj    []string          `json:"fakoj,omitempty"`
	Vidu     string            `json:"vidu,omitempty"`
	Tradukoj map[string]string `json:"tradukoj"`
}

func readOldJSON(in io.Reader) ([]Entry, error) {
	j := json.NewDecoder(in)
	var l0 []OldEntry
	err := j.Decode(&l0)
	if err != nil {
		return nil, err
	}
	var l []Entry
	for _, e0 := range l0 {
		tradukoj := map[string]string{
			"en": e0.En,
		}
		e := Entry{Radiko: e0.Radiko, Fakoj: e0.Fakoj, Vidu: e0.Vidu, Tradukoj: tradukoj}
		l = append(l, e)
	}
	return l, nil
}

func readJSON(in io.Reader) ([]Entry, error) {
	j := json.NewDecoder(in)
	var l []Entry
	err := j.Decode(&l)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func writeJSON(out io.Writer, l []Entry) error {
	j := json.NewEncoder(out)
	j.SetIndent("", "  ")
	err := j.Encode(l)
	if err != nil {
		return err
	}
	return nil
}

func readCSV(in io.Reader) ([]Entry, error) {
	r := csv.NewReader(in)

	seenHeader := false
	var l = []Entry{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if !seenHeader {
			seenHeader = true
			continue
		}

		bona := record[2] == "j"
		nivelo, err := strconv.Atoi(record[4])
		if err != nil {
			nivelo = 2
		}
		fakoj := strings.Split(record[5], ";")
		tradukoj := map[string]string{
			"en": record[7],
		}

		e := Entry{
			Radiko:   record[0],
			Speco:    record[1],
			Bona:     bona,
			Fonto:    record[3],
			Nivelo:   nivelo,
			Fakoj:    fakoj,
			Vidu:     record[6],
			Tradukoj: tradukoj,
		}
		l = append(l, e)
	}

	return l, nil
}

func writeCSV(out io.Writer, l []Entry) error {
	w := csv.NewWriter(out)

	r := []string{"radiko", "speco", "bona", "fonto", "nivelo", "fakoj", "vidu", "en"}
	if err := w.Write(r); err != nil {
		return err
	}

	for _, e := range l {
		bona := "j"
		if !e.Bona {
			bona = "n"
		}
		nivelo := strconv.Itoa(e.Nivelo)
		r = []string{e.Radiko, e.Speco, bona, e.Fonto, nivelo, strings.Join(e.Fakoj, ";"), e.Vidu, e.Tradukoj["en"]}
		if err := w.Write(r); err != nil {
			return err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

func contains(l []string, s string) bool {
	for _, e := range l {
		if e == s {
			return true
		}
	}
	return false
}

func runeIndex(r rune) int {
	for i, x := range LETTERS {
		if x == r {
			return i
		}
	}
	return -1
}

func readFileLines(f string) []string {
	fd, err := os.Open(f)
	if err != nil {
		log.Fatalf("error opening file: %v\n", err)
	}
	fs := bufio.NewScanner(fd)
	var lines []string
	for fs.Scan() {
		lines = append(lines, fs.Text())
	}
	err = fs.Err()
	if err != nil {
		log.Fatalf("error readind file lines: %v\n", err)
	}
	return lines
}

func readFundamento() map[string]bool {
	m := map[string]bool{}
	for _, f := range readFileLines("fundamento") {
		m[f] = true
	}
	return m
}

func rootLess(a, b string) bool {
	aRunes, bRunes := []rune(a), []rune(b)
	aDash := false

	if aRunes[0] == '-' {
		aRunes = aRunes[1:]
		aDash = true
	}
	if bRunes[0] == '-' {
		bRunes = bRunes[1:]
	}

	aLen, bLen := len(aRunes), len(bRunes)

	for ci, cj := 0, 0; ci < aLen && cj < bLen; ci, cj = ci+1, cj+1 {
		li, lj := aRunes[ci], bRunes[cj]
		ni := runeIndex(li)
		nj := runeIndex(lj)
		if ni < nj {
			return true
		}
		if ni > nj {
			return false
		}
	}

	if aLen < bLen {
		return true
	}
	if aLen > bLen {
		return false
	}
	return aDash
}

func sortEntries(l []Entry) []Entry {
	sort.SliceStable(l, func(i, j int) bool {
		return rootLess(l[i].Radiko, l[j].Radiko)
	})
	return l
}

type entryProc func(Entry) Entry

func transform(l []Entry, p entryProc) []Entry {
	var l1 []Entry
	for _, e := range l {
		e1 := p(e)
		l1 = append(l1, e1)
	}

	return l1
}

func markFundamento(l []Entry) []Entry {
	fj := readFundamento()

	return transform(l, func(e Entry) Entry {
		r := strings.ReplaceAll(e.Radiko, "/", "'")
		if fj[r] {
			e.Fonto = "f"
		} else {
			e.Fonto = "a"
		}
		return e
	})
}

func changePunc(l []Entry) []Entry {
	return transform(l, func(e Entry) Entry {
		e.Radiko = strings.ReplaceAll(e.Radiko, "/", "'")
		e.Vidu = strings.ReplaceAll(e.Vidu, "/", "'")
		return e
	})
}

type entryPred func(Entry) bool

func filter(l []Entry, p entryPred) []Entry {
	var l1 = []Entry{}
	for _, e := range l {
		if p(e) {
			l1 = append(l1, e)
		}
	}

	return l1
}

func filterFako(l []Entry, fako string) []Entry {
	return filter(l, func(e Entry) bool {
		return contains(e.Fakoj, fako)
	})
}

func filterPrefix(l []Entry, prefix string) []Entry {
	return filter(l, func(e Entry) bool {
		return strings.HasPrefix(e.Radiko, prefix)
	})
}

func readEntries(spec string) ([]Entry, error) {
	fname := filepath.Base(spec)
	format := filepath.Ext(spec)
	fname = strings.TrimSuffix(fname, format)

	var in *os.File
	if fname == "-" {
		in = os.Stdin
	} else {
		var err error
		in, err = os.Open(spec)
		if err != nil {
			return nil, err
		}
	}

	switch format {
	case ".csv":
		return readCSV(in)
	case ".json":
		return readJSON(in)
	default:
		return nil, errors.New("no understand " + spec)
	}
}

func writeEntries(spec string, l []Entry) error {
	fname := filepath.Base(spec)
	format := filepath.Ext(spec)
	fname = strings.TrimSuffix(fname, format)

	var out *os.File
	if fname == "-" {
		out = os.Stdout
	} else {
		var err error
		out, err = os.OpenFile(spec, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
	}

	switch format {
	case ".csv":
		return writeCSV(out, l)
	case ".json":
		return writeJSON(out, l)
	default:
		return writeJSON(out, l)
		// return errors.New("no understand " + spec)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <in> <out>\n", os.Args[0])
		os.Exit(1)
	}

	inspec, outspec := os.Args[1], os.Args[2]

	l, err := readEntries(inspec)
	if err != nil {
		log.Fatalf("cannot read %s: %v\n", inspec, err)
	}
	fmt.Fprintf(os.Stderr, "in: %d\n", len(l))

	// l = markFundamento(l)
	// l = changePunc(l)
	// l = filterFako(l, fako)
	// l = filterPrefix(l, prefix)
	l = sortEntries(l)

	err = writeEntries(outspec, l)
	if err != nil {
		log.Fatalf("cannot read %s: %v\n", inspec, err)
	}
	fmt.Fprintf(os.Stderr, "out: %d\n", len(l))
}
