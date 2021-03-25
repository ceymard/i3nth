# i3nth, a small i3 companion

I3nth is a tool written in golang that provides the following commands for use with `bindsym`

- `nop i3nth nth <number>` : focus the nth window starting from the leftmost one to the rightmost one
    of all the currently displayed workspaces on all current screens
- `nop i3nth change-group [groupname]` : switch the current workspace group to another, focusing the workspaces
    of the new group if there were some previously associated with it. If no groupname is provided, `rofi`
    will be launched
- `nop i3nth rename-group [newname]` : rename the currently active group to `newname`. If no name is provided,
    `rofi` is brought up as an input prompt

When changing "groups", i3nth takes all the currently active workspaces and renames them to include
the previous group name. It also scans all the workspaces that exist and renames those that match
the new current group to their old name and displays them back if they existed.

If there was no workspace associated to the new group on a given screen, then i3nth will try to switch to a new
workspace with the same name as the one that was there before, this displaying a new empty workspace.

For now, there is no way to send a window to another workspace to another group, although this will probably added in the future.

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

Also, do not forget to run it as a daemon !

```
exec_always --no-startup-id $HOME/go/bin/i3nth
```

# Why rofi ?

Because I have it installed, I like it, and it was convenient to do so to have something working fast.
As this tool was mostly intended for me, I did not push it as far as allowing other methods, although this
could change if there is interest.

# Quirks

At the time of this writing, the workspaces that are renamed do not appear anymore, probably due to the
way i3-rocks interprets pango strings that have extraneous attributes. This suits me, but it it might not suit you.