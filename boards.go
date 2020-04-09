package main

import (
	"strings"
)

func IsBoard(file string) bool {
	return strings.HasSuffix(file, ".board.toml")
}

func ReadBoards(dir string) []string {
	boards := []string{}

	for _, file := range FindFiles(dir) {
		if !IsBoard(file) {
			continue
		}
		board := strings.TrimPrefix(file, dir)
		board = strings.TrimSuffix(board, ".board.toml")
		boards = append(boards, "@"+board)
	}
	return boards
}
