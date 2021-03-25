package main

import (
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"go.i3wm.org/i3/v4"
)

///////////////////////////////////////////////////////////
type Set struct {
	mp map[string]struct{}
}

func NewSet() *Set {
	return &Set{mp: make(map[string]struct{})}
}
func (s *Set) Add(str string) {
	s.mp[str] = struct{}{}
}

func (s *Set) Has(str string) bool {
	if _, ok := s.mp[str]; ok {
		return true
	}
	return false
}

///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
type nodesByX []*i3.Node

func (s nodesByX) Len() int {
	return len(s)
}
func (s nodesByX) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s nodesByX) Less(i, j int) bool {
	return s[i].Rect.X < s[j].Rect.X
}

///////////////////////////////////////////////////////////

// filter a Tree on a condition
func filterTree(node *i3.Node, fn func(node *i3.Node) bool) []*i3.Node {
	var res []*i3.Node
	for _, n := range node.Nodes {
		res = append(res, filterTree(n, fn)...)
	}
	if fn(node) {
		res = append(res, node)
	}
	return res
}

// gotoNth sends a message to i3 to go to the nth window starting on the left of the screen.
func gotoNth(nth int) {
	var (
		wrk            []i3.Workspace
		visibleWrk     = NewSet()
		workspaceNodes []*i3.Node
		clients        []*i3.Node
		tree           i3.Tree

		// clients     []i3.Tree
		err error
	)

	if wrk, err = i3.GetWorkspaces(); err != nil {
		return
	}

	// Filter the workspaces to include only those currently on a screen
	// If there is only one monitor, the list will only have one workspaces
	for _, w := range wrk {
		if w.Visible {
			// pp.Print(w)
			visibleWrk.Add(w.Name)
		}
	}

	if tree, err = i3.GetTree(); err != nil {
		return
	}

	// Get the nodes that are workspaces
	workspaceNodes = filterTree(tree.Root, func(n *i3.Node) bool {
		return n.Type == "workspace" && visibleWrk.Has(n.Name)
	})

	for _, w := range workspaceNodes {
		clients = append(clients, filterTree(w, func(n *i3.Node) bool {
			return len(n.Nodes) == 0 && n.Type == "con"
		})...)
	}

	// Now that we have the clients, we order them by ascending x coordinates
	sort.Sort(nodesByX(clients))

	if nth > 0 && nth <= len(clients) {
		i3.RunCommand("[con_id=" + strconv.FormatInt(int64(clients[nth-1].ID), 10) + "] focus")
	}
}

var reNth = regexp.MustCompile(`nop wg nth (\d+)`)
var reChangeGroup = regexp.MustCompile(`nop wg change(?: +(.+))?`)
var reRenameGroup = regexp.MustCompile(`nop wg rename(?: +(.+))?`)

func handleBinding(v *i3.BindingEvent) error {
	if !strings.HasPrefix(v.Binding.Command, "nop wg") {
		return nil
	}
	var cmd = v.Binding.Command
	log.Print(cmd)

	// Handle focus nth command
	if match := reNth.FindStringSubmatch(cmd); match != nil {
		nth, _ := strconv.Atoi(match[1])
		gotoNth(nth)
		return nil
	}

	if match := reChangeGroup.FindStringSubmatch(cmd); match != nil {
		log.Print("received change")
		if len(match[1]) > 0 {
			activateGroup(match[1])
		} else {
			trySwitchToGroup()
		}
		return nil
	}

	if match := reRenameGroup.FindStringSubmatch(cmd); match != nil {
		if len(match[1]) > 0 {
			renameCurrentGroup(match[1])
		} else {
			tryRenameCurrentGroup()
		}
		return nil
	}

	// Handle

	return nil
}

func main() {
	rec := i3.Subscribe("binding")
	// i3.RunCommand()

	for rec.Next() {
		evt := rec.Event()

		switch v := evt.(type) {

		case *i3.BindingEvent:
			handleBinding(v)
		}
	}
}
