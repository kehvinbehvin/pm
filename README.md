# pm
Project manager for solo developers

## Summary
Record your requirements with your code and not compromise on the complexity of the relationships
between your requirements. Use hotkeys and a TUI to quickly break down and manage all your tasks
within your terminal.

All your issues are created as MD files in your project root so you wont lose anything. Track
the entire .pm directory in git so that you can always revert back to previous version and
share with your teammates.

No need to create a DB or docker containers. Just run download the binary and run pm on your terminal
to get started

## Features

![](https://github.com/kehvinbehvin/pm/blob/main/example.gif)

### File Structure
- .pm
    - blobs (Directory for all your issues)
    - dag (Storing relationship between your files)
    - fileTypes (Storing your file types)

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

## Building
- go build .
- sudo mv ./pm /usr/local/bin/pm