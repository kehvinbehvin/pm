package main

func main() {
  epicTrie := NewTrie()
  epicTrie.addWord("epicA")
  epicTrie.addWord("epicB")

  storyTrie := NewTrie()
  storyTrie.addWord("Research A")
  storyTrie.addWord("Implement B")

  taskTrie := NewTrie()
  taskTrie.addWord("Update function B")
  taskTrie.addWord("Remove variable A")

  Save(epicTrie, "epic")
  Save(storyTrie, "story")
  Save(taskTrie, "task")
  Execute();
}
