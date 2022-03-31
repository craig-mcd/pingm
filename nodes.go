package main

import (
	"net"
)

func cleanNodes(dirtyNodes []string) ([]string, []string) {

	nodes := []string{}
	invalidNodes := []string{}

	for _, node := range dirtyNodes {

		_, err := net.LookupIP(node)

		if err == nil {
			nodes = append(nodes, node)
		} else {
			invalidNodes = append(invalidNodes, node)
		}

	}

	return nodes, invalidNodes
}
