# pm
Project manager for solo developers

## Summary
As a Solo developer, i need to manage my own projects. One problem i face is having my requirements
kept separately from my source code. This approach will help me to check in on my project status 
like in a git-style.

## Features

### File Structure
- .pm 
  - config
  - epics
  - stories
  - tasks
  - objects
  - markdown

### Storage method
- Bytes or Compressed (after a certain size)
- Hashed content as the key 
  - Allows for a flatter content directory
  - Easier to represent relationships
  - Avoid file name collisions
  - Version control
  - Data integrity

### Autocomplete and suggestions
- Trie prefix structure 
  - prefix as key, content hash as value

### Edit
- Use name to look up content
- Use working file to display editable content
- Hash edited content and create new blob
- Update name value to new content hash

### CLI Tool
- Create Epic and associate user stories to them
- Create User stories and associate tasks to them 
- Update Epics, Stories, Tasks (Individually)
- Access to the same pm directory across different directories / users
- Reference git tags to indicate work completed
- Autocomplete and suggestions

## Commands
- pm init (Creates master directory)
- pm attach (Creates symlink to master directory)
  - [path to master]
- pm tree (Displays full table of epics, stories, tasks)
  - [-r nid] (Choose root node of tree) Eg. pm tree -r S2
- pm list (Display list of epics, stories, tasks)
  - [-e | -s | -t]
- pm add -n 5 -[e | s | t] (5 prompts to create epics/stories/task)
- pm add -n 5 -[e | s | t]  > [e | s | t] [name] (5 prompts to create epic/stories/tasks linked to epic/stories/tasks)

## Central repository
- pm can add remotes to pm directory via config
- pm can pull changes from remote via ssh
- pm can push changes to remote via ssh
- pm can deal with conflicts

