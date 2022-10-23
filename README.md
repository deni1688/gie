[![CI](https://github.com/deni1688/gie/actions/workflows/ci.yml/badge.svg)](https://github.com/deni1688/gie/actions/workflows/ci.yml)
# gie (git issue extractor)
cli tool to create multiple issues for a git provider

### Mission
To make it simple to extract issues to a configured git hosting service like GitHub or GitLab from your code base.

### Vision
Making it easier to extract issues will hopefully encourage developers to submit smaller commits to resolve
those specific issues. The tool itslef can be used locally or inside of a CI pipeline. A ports and adapters like 
architecture also makes it easy to extend the tool quite easily with your own adapters for the git host and notifier. 

### How it works

By prefixing comments with a specified tag inside your code and then running `gie` on that file or directory 
the comment(s) are extracted, formatted, and submitted to your git hosting service. The resulting issue(s) are then
appended to the commented line in order to have a direct reference to it. Each issue
will also have a reference to the file where it was extracted from in your git hosting service.

### Example

#### Before

```go
    // Issue: Make it possible to do XYZ in someFunctionThatShouldBerefactored
    func someFunctionThatShouldBerefactored() {
        // ...        
    }
```

#### Run
```bash
gie -path .
```


#### After
```go
    // Issue: Make it possible to do XYZ in someFunctionThatShouldBerefactored -> closes https://github.com/owner/project/issues/12
    func someFunctionThatShouldBerefactored() {
        // ...        
    }
```



