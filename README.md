# pm
Project manager for solo developers

## Summary
As a Solo developer, i need to manage my own projects. One problem i face is having my requirements
kept separately from my source code. This approach will help me to check in on my project status
like in a git-style.

## Features

![](https://github.com/kehvinbehvin/pm/blob/main/example.gif)

### File Structure
- .pm
    - config (txt)
    - blobs (addressable content)
    - prefix (binary)
    - dag (binary)

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

### Implementing the data structures
- https://intranet.icar.cnr.it/wp-content/uploads/2018/12/RT-ICAR-PA-2018-06.pdf
- Build DAG/Trie in memory, perform binary serialisation to store it on disk
- Use mmmap to fetch Dag/trie

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

## Adding shell completion
- Zsh
    - Add 'autoload -U compinit; compinit' to zshrc
    - add completion file in $fpath
    - Add 'source <path to completion function>'
        - source /Users/kevin/.oh-my-zsh/completions/_pm
    - source zshrc

## History
- Blobs
  - Blobs will not have a remove feature so it can always be referenced

- Tries delta will be snapshot on each operation (add/remove)
- Upon checking out snapshot in history, trie can be built forward or backwards
  Trie Delta 
  - Delta structure for tries will store the change in trie using operations as bytes and the filename as the path

- Dag delta will be snapshot on every operation
- Dag Delta
  - if the change is an edge, the delta would be the ID of parent vertex and the ID of the child vertex + the operation in bytes
  - if the change is a vertex, the delta would be the binary of the vertex + the operation in bytes

## Dealing with conflicts
- Local ahead of Remote
  - Diff the local against the remote
  - Returns extra deltas that local has 
  - Push extra deltas to remote
- Local behind of Remote
  - Diff the local against the remote
  - Returns extra deltas that remote has
  - Pull extra deltas from remote
- Local and Remote have deviated
  - Allows non-conflicting events to be applied
  - Diff the local against the remote
  - Identify all deltas that operate on the same vertex or add the same Edge(parent, child)
  - Group all deltas by events and end state
  - Ask user to choose which end state

- Deviation case
  - Given all the deltas after LCS on both trees
  - Build 2 trees that represent the end state of applying the deltas in sequence 

DeltaState {
  Vertex
  Opp: Create | Destroy | nil {make sure to store reference to delta here}
  Opp: Add Child | Remove Child | nil {make sure to store reference to delta here}
}
  - Construct DeltaStates from each node
  - Compare all DeltaStates, if any have conflicting, get User input
  - Else continue

- Add Vertex, Remove Vertex
- Add Vertex, Remove Vertex. Add Vertex

- Add Vertex 1, Add Edge 1-2, Remove Vertex 1
- Add Vertex 1, Add Edge 1-2, Remove Edge 1-2
- A V1, A V2, A E(1,2), R E(1,2) 

Case 1
Vertex 1
- Remove
Edge 1-2
- Add

Case 2
Edge 1-2
- Remove
Vertex 1
- Add

## Building
- go build .
- sudo mv ./pm /usr/local/bin/pm


## Revamp and Rescoping
- How to launch MVP asap

## MVP Definition
- 

