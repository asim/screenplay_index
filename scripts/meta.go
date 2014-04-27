package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

/*

cat u | while read id data ; do if [ -s "tmp/$id.meta" ] && [ `egrep '^\s+$' "tmp/$id.meta"|wc -l` -ne 1 ]; then d=`cat "tmp/$id.meta"`; while [ "$d" ]; do last=`echo $d | awk '{print $NF}' | tr '[:upper:]' '[:lower:]'`; if (egrep -q "^${last}$" words)&>/dev/null; then echo -e "$d..." >> "m/$id"; break ; else d=`echo $d | awk '{$(NF--)=""; print}'`; fi; done; fi; done



*/

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Require work dir")
		return
	}

	dir := os.Args[1]

	urls, err := ioutil.ReadFile(filepath.Join(dir, "urls"))
	if err != nil {
		fmt.Println(err)
		return
	}

	words, err := ioutil.ReadFile(filepath.Join(dir, "words"))
	if err != nil {
		fmt.Println(err)
		return
	}

	wordM := make(map[string]bool)

	buf := bytes.NewBuffer(words)

	for {
		w, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}

		wordM[w[:len(w)-2]] = true
	}

	fmt.Printf("%#v", wordM)

	ubuf := bytes.NewBuffer(urls)
	for {
		u, err := ubuf.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}

		w := strings.Split(u, "\t")
		if len(w) < 1 {
			continue
		}

		//		id := u[:len(u)-1]
		id := w[0]
		meta, err := ioutil.ReadFile(filepath.Join(dir, "meta", id+".meta"))
		if err != nil {
			fmt.Println(err)
			continue
		}

		m := strings.Split(string(meta), " ")

		for {
			if len(m) == 0 {
				break
			}

			if w := strings.ToLower(m[len(m)-1]); wordM[w] {
				break
			} else {
				m = m[:len(m)-1]
			}
		}

		if len(m) > 0 {
			ioutil.WriteFile(filepath.Join(dir, "m", id), []byte(strings.Join(m, " ")+"..."), 0666)
		}
	}
}
