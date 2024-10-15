package main

func main() {
  epicTrie := NewTrie()
  epicTrie.addFile("epicA", "This is the description of epic A")
  epicTrie.addFile("epicB", "This is the description of epic B")

  storyTrie := NewTrie()
  storyTrie.addFile("Research A", "This is the description of story A")
  storyTrie.addFile("Implement B", "This is the description of story A")

  taskTrie := NewTrie()
  taskTrie.addFile("Update function B", "This is the description of task A")
  taskTrie.addFile("Remove variable A", "This is the description of task A")

  Save(epicTrie, "epic")
  Save(storyTrie, "story")
  Save(taskTrie, "task")
  Execute();
}
