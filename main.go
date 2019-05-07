package main

import (
	"fmt"
	"hack-assember/code"
	"hack-assember/parser"
	"hack-assember/symbol"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	p := parser.New(f)
	var line int
	var lineRom int64
	for p.HasMoreCommand() {
		err = p.Advance(&line)
		if err != nil {
			fmt.Println(err)
		}
		if p.CommandType() == parser.CCommand || p.CommandType() == parser.ACommand {
			lineRom++
		}

		if p.CommandType() == parser.LCommand {
			if _, b := symbol.Table[p.Symbol()]; !b {
				symbol.Table[p.Symbol()] = lineRom
			}
		}
	}
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Println(err)
		return
	}
	p = parser.New(f)
	line = 0
	var newAddress int64 = 16
	for p.HasMoreCommand() {
		err = p.Advance(&line)
		if err != nil {
			fmt.Println(err)
			return
		}
		if p.CommandType() == parser.LCommand {
			continue
		}
		if p.CommandType() == parser.ACommand {
			v := p.Symbol()
			var vInt int64
			vInt, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				var b bool
				vInt, b = symbol.Table[p.Symbol()]
				if !b {
					vInt = newAddress
					symbol.Table[p.Symbol()] = newAddress
					newAddress++
				}
			}
			command := strconv.FormatInt(vInt, 2)
			leading := 16 - len(command)
			command = strings.Repeat("0", leading) + command
			fmt.Println(command)
			continue
		}
		if p.CommandType() == parser.CCommand {
			d := p.Dest()
			c := p.Comp()
			j := p.Jump()
			command := "111" + code.Comp[c] + code.Dest[d] + code.Jump[j]
			fmt.Println(command)
		}
	}
}
