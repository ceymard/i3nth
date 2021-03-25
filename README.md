# i3nth, a small i3 companion

I3nth is a tool written in golang that provides the following commands for use with `bindsym`

- `nop i3nth nth <number>` : focus the nth window starting from the leftmost one to the rightmost one
    of all the currently displayed workspaces on all current screens
- `nop i3nth change-group [groupname]` : switch the current workspace group to another, focusing the workspaces
    of the new group if there were some previously associated with it. If no groupname is provided, `rofi`
    will be launched
- `nop i3nth rename-group [newname]` : rename the currently active group to `newname`. If no name is provided,
    `rofi` is brought up as an input prompt

# Installation

At this time there are no pre-built binaries available. This could change if there is demand for it.

For now, you will need golang installed on your system and type somewhere not in a go project the following

`go get github.com/ceymard/i3nth`

This will install `i3nth` in `~/go/bin`, so make sure this is in your path.

You may then create bindings in your i3 config file like the following ;

```
# Switch to the nth window on screen(s)
bindsym $mod+1 nop i3nth nth 1
bindsym $mod+2 nop i3nth nth 2
# ...

# You may tag group manually
bindsym $mod+Shift+F1 nop i3nth change first
bindsym $mod+Shift+F2 nop i3nth change second
bindsym $mod+Shift+F3 nop i3nth change third

# This will bring up rofi with a menu to select an existing group or create a new one
bindsym $mod+P nop i3nth change
# This brings up rofi to rename the current group to something else
bindsym $mod+Y nop i3nth rename
```

# Why rofi ?

Because I have it installed, I like it, and it was convenient to do so to have something working fast.
As this tool was mostly intended for me, I did not push it as far as allowing other methods, although this
could change if there is interest.
