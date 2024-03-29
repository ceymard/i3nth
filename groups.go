package main

/**
  Whenever the group changes, the following is performed ;
   - all current workspaces are renamed as "<group>::<workspace_name>"
   - all workspaces matching "<newgroup>::..." are renamed as the ... part.
	Workspaces stay on their assigned displays.
*/

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"go.i3wm.org/i3/v4"
)

var currentGroup = "default"
var previousGroup = "default"

// map[groupname]map[workspacename]outputname
var reGroupName = regexp.MustCompile(`^<span group='([^']+)'( visible='')?[^>]*>❱</span>(.*)`)

func tryRenameCurrentGroup() {
	var cmd = exec.Command("rofi", "-dmenu", "-mesg", fmt.Sprintf("Current group is '%s' Enter a new group name", currentGroup))
	var stdout, err = cmd.Output()
	if err != nil {
		return
	}
	var out = strings.TrimSpace(string(stdout))
	if out != "" {
		renameCurrentGroup(out)
	}
	log.Print(out)
}

// Find returns the smallest index i at which x == a[i],
// or len(a) if there is no such index.
func find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -1
}

func trySwitchToGroup() {
	var groups = make([]string, 0)
	var wk, err1 = i3.GetWorkspaces()
	if err1 != nil {
		log.Print(err1)
		return
	}
	for _, w := range wk {
		var match = reGroupName.FindStringSubmatch(w.Name)
		if match == nil {
			continue
		}
		if find(groups, match[1]) == -1 {
			groups = append(groups, match[1])
		}
	}

	var cmd = exec.Command("rofi", "-dmenu", "-mesg", "Switch to another group")
	var stdin, err2 = cmd.StdinPipe()
	if err2 != nil {
		log.Print(err2)
		return
	}
	sort.Strings(groups)
	var _g = []string{}
	if previousGroup != currentGroup {
		_g = append(_g, previousGroup)
	}

	for _, g := range groups {
		if g != previousGroup && g != currentGroup {
			_g = append(_g, g)
		}
	}

	groups = _g

	for i, g := range groups {
		if i > 0 {
			_, _ = stdin.Write([]byte{'\n'})
		}
		_, _ = stdin.Write([]byte(g))
	}
	_ = stdin.Close()
	var stdout, err = cmd.Output()
	if err != nil {
		return
	}
	var out = strings.TrimSpace(string(stdout))
	if out != "" {
		activateGroup(out)
	}
}

func renameCurrentGroup(newgroup string) {
	if currentGroup == previousGroup {
		previousGroup = newgroup
	}
	currentGroup = newgroup
}

func activateGroup(newgroup string) {
	// how do we know the current group name ?
	if newgroup == currentGroup {
		return
	}

	previousGroup = currentGroup
	var current = currentGroup
	var wk, err = i3.GetWorkspaces()
	if err != nil {
		return
	}

	var cmds = make([]string, 0)
	var activated = make(map[string]struct{})
	var previously_activated = make(map[string]string)

	for _, w := range wk {
		var match = reGroupName.FindStringSubmatch(w.Name)
		// first rename workspace to its group version
		if match == nil {
			if current == newgroup {
				continue
			}
			var visible = ""
			if w.Visible {
				visible = ` visible=''`
				previously_activated[w.Output] = w.Name
			}
			// If they don't match the regexp, they're then assigned to current.
			cmds = append(cmds, fmt.Sprintf(
				"rename workspace \"%s\" to \"%s\"",
				w.Name,
				fmt.Sprintf(`<span group='%s'%s>❱</span>%s`, current, visible, w.Name),
			))
		}
	}

	for _, w := range wk {
		// for all workspaces that we have, we run the regexp to find out which group they're part of.
		var match = reGroupName.FindStringSubmatch(w.Name)

		// If the group of the workspace is not the one we're trying to activate, don't activate it.
		if match == nil || match[1] != newgroup {
			continue
		}

		var visible = match[2] != ""
		var oldname = match[3]

		// Otherwise, rename it to its "ungrouped" version.
		cmds = append(cmds, fmt.Sprintf(
			"rename workspace \"%s\" to \"%s\"",
			w.Name,
			oldname,
		))

		// And focus it if it's the first time we encounter it for a given X-Y position
		var pos = w.Output
		if _, ok := activated[pos]; !ok && visible {
			cmds = append(cmds, fmt.Sprintf("workspace \"%s\"", oldname))
			activated[pos] = struct{}{}
		}

	}

	for output, name := range previously_activated {
		if _, ok := activated[output]; !ok {
			cmds = append(cmds, fmt.Sprintf(`focus output %s`, output))
			cmds = append(cmds, fmt.Sprintf(`workspace "%s"`, name))
			cmds = append(cmds, fmt.Sprintf(`move workspace to output %s`, output))
		}
	}

	for _, c := range cmds {
		log.Print(c)
		var cmdres, err = i3.RunCommand(c)
		if err != nil {
			log.Print("error: ", err)
		}
		for _, r := range cmdres {
			log.Print(r)
		}
	}
	currentGroup = newgroup
	// if len(cmds) > 0 {
	// i3.RunCommand(strings.Join(cmds, ";"))
	// }
}
